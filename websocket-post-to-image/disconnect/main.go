package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, request events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
    fmt.Println("Disconnect Event:", request)
    return events.APIGatewayProxyResponse{StatusCode: 200, Body: "Disconnected"}, nil
}

func main() {
    lambda.Start(handler)
}
