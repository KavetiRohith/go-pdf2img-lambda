# go-pdf2img-lambda

This repository contains the source code for an AWS Lambda service that automatically splits a PDF file into images whenever a PDF is uploaded to an S3 bucket associated with the Lambda function. The service is implemented in Go and utilizes the [gen2brain/go-fitz](https://github.com/gen2brain/go-fitz) library, which is a Go wrapper for the libmupdf C library.

Using CGO (C Go) binaries in AWS Lambda can be challenging due to the serverless nature of Lambda and the constraints it imposes. Here are some common issues and potential solutions:

1. **Lambda Execution Environment:**
   - Lambda provides a specific execution environment, and the binaries compiled using CGO might depend on libraries or system calls not available in that environment.
   - Using docker and aws lambda go base image to compile a go binary.
   - In our case we can do and we would get a cgo binary that uses dynamic linking.
   ```bash
    CGO_ENABLED=1 go build  -o bootstrap
    ldd bootstrap
            linux-vdso.so.1 (0x00007ffd5a5ef000)
            libm.so.6 => /lib/x86_64-linux-gnu/libm.so.6 (0x00007f7d260ec000)
            libc.so.6 => /lib/x86_64-linux-gnu/libc.so.6 (0x00007f7d25ec4000)
            /lib64/ld-linux-x86-64.so.2 (0x00007f7d261db000)
   ```
   

2. **Static Linking:**
   - CGO binaries often dynamically link to C libraries, which may not be available in the Lambda environment.
   - Try statically linking the C libraries into your binary to reduce dependencies on external libraries.
   - In our case we can do and we would get a cgo binary that uses static linking.
   ```bash
    CGO_ENABLED=1 GOOS=linux GOARCH=amd64 CC="zig cc -target x86_64-linux-musl" CXX="zig c++ -target x86_64-linux-musl" go build -tags musl -ldflags="-linkmode external" -o bootstrap
    ldd bootstrap
        statically linked
   ```
   
    - Using Zig to compile CGO programs and targeting the Linux-musl architecture can indeed help alleviate some of the issues associated with using CGO binaries in AWS Lambda like cross compiling for arm64 target on an amd64 machine. Here are some advantages and considerations:

        1. **Static Binaries:**
        - Zig supports the creation of statically linked binaries, which helps reduce dependencies on external libraries. This is beneficial for deploying to environments like AWS Lambda, where the runtime environment may be limited.

        2. **Linux-musl Target:**
        - The Linux-musl target in Zig produces binaries that are compatible with the musl libc, which is known for its minimalistic and lightweight nature. This can be more suitable for serverless environments like AWS Lambda compared to binaries linked against glibc.

        3. **Cross-Compilation:**
        - Zig supports cross-compilation, allowing you to build binaries for different architectures and platforms. This is crucial for ensuring compatibility with the AWS Lambda execution environment.


    - While using Zig and the Linux-musl target can be a promising approach, it's essential to stay informed about updates and changes in the Zig ecosystem and AWS Lambda environment.

## Prerequisites

Before deploying and using this Lambda service, make sure you have the following prerequisites:

- An S3 bucket configured to trigger events when a new PDF file is uploaded.
- [Go](https://go.dev/dl/) and [Zig](https://ziglang.org/download/) installed on your local machine.

## Installation and Deployment

Follow these steps to deploy the Lambda service:

1. Clone this repository:

   ```bash
   git clone https://github.com/KavetiRohith/go-pdf2img-lambda.git
   cd go-pdf2img-lambda
   ```

2. Compile a static binary:

   ```bash
   CGO_ENABLED=1 GOOS=linux GOARCH=amd64 CC="zig cc -target x86_64-linux-musl" CXX="zig c++ -target x86_64-linux-musl" go build --tags musl -ldflags="-linkmode external" -o bootstrap
   ```

3. Package the Lambda function:

   ```bash
   zip lambda-handler.zip bootstrap
   ```
4. For deployment and additional information
    - https://docs.aws.amazon.com/lambda/latest/dg/golang-package.html#golang-package-mac-linux
    - https://docs.aws.amazon.com/lambda/latest/dg/with-s3-example.html
    
    
## Usage

Once the Lambda function is deployed, it will automatically trigger whenever a new PDF file is uploaded to the associated S3 bucket. The PDF file will be split into images, and the resulting images will be stored in the same S3 bucket.

## Dependencies

- [gen2brain/go-fitz](https://github.com/gen2brain/go-fitz): Go wrapper for the libmupdf C library.
- [Zig](https://ziglang.org/): Used to compile static binaries for the Lambda service.
- [aws-lambda-go sdk](github.com/aws/aws-lambda-go)
- [aws-sdk-go-v2](github.com/aws/aws-sdk-go-v2)

## Notes

- The Lambda service relies on the gen2brain/go-fitz library, which requires CGO to compile. To address compatibility issues with AWS Lambda, this repository uses Zig to compile static binaries that can be deployed as Lambda functions.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

