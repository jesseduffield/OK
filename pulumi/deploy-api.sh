#!sh

set -e

cd api/handler
# direct so that we're not waiting for the proxy to receive our latest changes
GOPROXY=direct go get github.com/jesseduffield/OK/ok
go get github.com/aws/aws-lambda-go/cmd/build-lambda-zip
go mod tidy
GOOS=linux GOARCH=amd64 go build -o handler main.go
build-lambda-zip -output handler.zip handler
cd ..
pulumi up
cd ..
