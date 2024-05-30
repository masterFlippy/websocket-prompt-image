package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, request events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
    fmt.Println("Connect Event:", request)
    return events.APIGatewayProxyResponse{StatusCode: 200, Body: "Connected"}, nil
}

func main() {
    lambda.Start(handler)
}
