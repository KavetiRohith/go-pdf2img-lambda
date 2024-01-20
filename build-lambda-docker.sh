#!/usr/bin/env bash

# cleanup and create the build/ directory
rm -r build 2> /dev/null && mkdir -p build

# build the image
docker build -t go-pdf2img-lambda -f Dockerfile.lambda .

# Run the container
docker run --name go-pdf2img-lambda go-pdf2img-lambda

# Copy the lambda-handler.zip file from the container to the host
docker cp go-pdf2img-lambda:/app/lambda-handler.zip build/lambda-handler.zip

# Cleanup
docker container rm go-pdf2img-lambda && docker image rm go-pdf2img-lambda