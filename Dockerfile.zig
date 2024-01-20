# Use Alpine Linux as base image
FROM alpine:latest

# Install necessary dependencies including xz
RUN apk update && \
    apk add --no-cache go tar wget xz zip && \
    wget -q https://ziglang.org/download/0.11.0/zig-linux-x86_64-0.11.0.tar.xz && \
    tar -xvf zig-linux-x86_64-0.11.0.tar.xz && \
    rm zig-linux-x86_64-0.11.0.tar.xz

# Set the working directory to /app
WORKDIR /app

# Copy all contents from the current directory to the container
COPY . .

# Build the Go application with Zig as the C/CXX compiler
RUN PATH="/zig-linux-x86_64-0.11.0:$PATH" CGO_ENABLED=1 GOOS=linux GOARCH=amd64 CC="zig cc -target x86_64-linux-musl" CXX="zig c++ -target x86_64-linux-musl" go build -tags musl -ldflags="-linkmode external" -o bootstrap && zip lambda-handler.zip bootstrap
