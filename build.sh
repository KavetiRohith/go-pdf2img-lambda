#!/bin/sh
# https://awsteele.com/blog/2021/10/17/cgo-for-arm64-lambda-functions.html

set -eux
rm -f bootstrap lambda-handler.zip
export GOOS=linux
export CGO_ENABLED=1
export CC=$(pwd)/zcc.sh
export CXX=$(pwd)/zxx.sh

GOARCH=amd64 \
ZTARGET=x86_64-linux-musl \
go build -tags musl -ldflags="-linkmode external" -o bootstrap
zip lambda-handler.zip bootstrap
