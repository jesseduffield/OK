package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jesseduffield/OK/ok/interpreter"
)

type Body struct {
	Content string `json:"content"`
}

func HandleRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	headers := map[string]string{
		"Content-Type":                "text/plain",
		"Access-Control-Allow-Origin": "*",
	}

	b64String, _ := base64.StdEncoding.DecodeString(request.Body)
	rawIn := json.RawMessage(b64String)
	bodyBytes, err := rawIn.MarshalJSON()
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500, IsBase64Encoded: false, Headers: headers}, nil
	}

	data := Body{}
	err = json.Unmarshal(bodyBytes, &data)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500, IsBase64Encoded: false, Headers: headers}, nil
	}

	reader := strings.NewReader(data.Content)
	writer := new(strings.Builder)

	err = withTimeout(
		func() { interpreter.Interpret(reader, writer) },
	)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500, IsBase64Encoded: false, Headers: headers}, nil
	}

	return events.APIGatewayProxyResponse{Body: writer.String(), StatusCode: 200, IsBase64Encoded: false, Headers: headers}, nil
}

func main() {
	lambda.Start(HandleRequest)
}

// TODO: dry this up, currently also defined in ok/cmd/playground/main.go.
func withTimeout(f func()) error {
	c := make(chan struct{}, 1)
	go func() {
		f()
		c <- struct{}{}
	}()

	select {
	case <-c:
		return nil
	case <-time.After(5 * time.Second):
		return errors.New("Timed out (program must complete within 5 seconds")
	}
}
