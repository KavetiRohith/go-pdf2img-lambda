#!/usr/bin/env bash

# cleanup and create the build/ directory
rm -r build 2> /dev/null && mkdir -p build

# build the image
docker build -t go-pdf2img-lambda-zig -f Dockerfile.zig .

# Run the container
docker run --name go-pdf2img-lambda-zig go-pdf2img-lambda-zig

# Copy the lambda-handler.zip file from the container to the host
docker cp go-pdf2img-lambda-zig:/app/lambda-handler.zip build/lambda-handler.zip

# cleanup
docker container rm go-pdf2img-lambda-zig && docker image rm go-pdf2img-lambda-zig