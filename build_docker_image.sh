#!/usr/bin/env sh
echo "==> Compiling Binary of project."
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o mysql-to-strict .

version=`cat VERSION`
echo "==> New version: $version"

docker build -t mysql-to-strict:${version} -f Dockerfile .

echo "==> Removing binary"
rm -rf ./mysql-to-strict
