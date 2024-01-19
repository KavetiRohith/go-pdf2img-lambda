#!/bin/sh
# https://awsteele.com/blog/2021/10/17/cgo-for-arm64-lambda-functions.html

set -eux
rm -f bootstrap lambda-handler.zip
export GOOS=linux
export CGO_ENABLED=1
export CC="zig cc -target x86_64-linux-musl"
export CXX="zig c++ -target x86_64-linux-musl"
export GOARCH=amd64
go build -tags musl -ldflags="-linkmode external" -o bootstrap
zip lambda-handler.zip bootstrap
rm bootstrap