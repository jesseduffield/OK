Steps:

## Deploy Playground API (lambda)

```sh
cd ok/pulumi/api/handler
get -u github.com/aws/aws-lambda-go/cmd/build-lambda-zip
GOOS=linux GOARCH=amd64 go build -o handler main.go
build-lambda-zip -output handler.zip handler
pulumi up
```

## Deploy Playground Site

```sh
cd ok/site
npm run build
cd ../pulumi/site
pulumi up
```
