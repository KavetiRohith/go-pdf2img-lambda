package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image/jpeg"
	"log"
	"strings"
	"sync"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/gen2brain/go-fitz"
)

func DownloadS3File(ctx context.Context, downloader *manager.Downloader, objectKey string, bucket string) ([]byte, error) {
	buffer := manager.NewWriteAtBuffer([]byte{})

	numBytes, err := downloader.Download(ctx, buffer, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		return nil, err
	}

	if numBytes < 1 {
		return nil, errors.New("zero bytes written to memory")
	}

	return buffer.Bytes(), nil
}

func UploadToS3(ctx context.Context, uploader *manager.Uploader, image *bytes.Buffer, objectKey string, bucket string) error {

	_, err := uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(objectKey),
		Body:   image,
	})

	return err
}

func ConvertToPdf(ctx context.Context, downloader *manager.Downloader, uploader *manager.Uploader, objectKey string, bucket string) error {
	pdfBytes, err := DownloadS3File(ctx, downloader, objectKey, bucket)
	if err != nil {
		return err
	}

	doc, err := fitz.NewFromMemory(pdfBytes)
	if err != nil {
		return err
	}

	defer doc.Close()

	imgPrefix := strings.TrimSuffix(objectKey, ".pdf")

	var wg sync.WaitGroup

	// limit max concurrency to 4 images
	limitter := make(chan struct{}, 4)

	for n := 0; n < doc.NumPage(); n++ {
		n := n

		wg.Add(1)
		limitter <- struct{}{}

		go func() {
			defer func() {
				<-limitter
			}()
			defer wg.Done()

			img, err := doc.Image(n)
			if err != nil {
				log.Printf("Unable to extract image %s-%v.jpeg err: %v\n", imgPrefix, n+1, err)
				return
			}

			imgBuffer := bytes.NewBuffer([]byte{})

			err = jpeg.Encode(imgBuffer, img, &jpeg.Options{Quality: jpeg.DefaultQuality})
			if err != nil {
				log.Printf("Unable to Encode image %s-%v.jpeg err: %v\n", imgPrefix, n+1, err)
				return
			}

			err = UploadToS3(ctx, uploader, imgBuffer, fmt.Sprintf("%s-%v.jpeg", imgPrefix, n+1), bucket)
			if err != nil {
				log.Printf("Unable to Upload image %s-%v.jpeg to s3 err: %v\n", imgPrefix, n+1, err)
				return
			}
		}()
	}
	wg.Wait()

	return nil
}

func handler(ctx context.Context, s3Event events.S3Event) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatal(err)
	}

	s3Client := s3.NewFromConfig(cfg)

	downloader := manager.NewDownloader(s3Client)
	uploader := manager.NewUploader(s3Client)

	for _, record := range s3Event.Records {
		s3 := record.S3
		err := ConvertToPdf(ctx, downloader, uploader, s3.Object.Key, s3.Bucket.Name)
		if err != nil {
			log.Printf("Unable to convert pdf %s to image err: %s\n", s3.Object.Key, err)
		}
	}
}

func main() {
	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(handler)
}
