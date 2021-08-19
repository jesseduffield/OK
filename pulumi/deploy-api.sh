#!sh

set -e

cd api/handler
go get github.com/aws/aws-lambda-go/cmd/build-lambda-zip
GOOS=linux GOARCH=amd64 go build -o handler main.go
build-lambda-zip -output handler.zip handler
cd ..
pulumi up
cd ..
