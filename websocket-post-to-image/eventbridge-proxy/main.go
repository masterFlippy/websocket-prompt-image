package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/eventbridge"
)

type Payload struct {
    Text    string `json:"text"`
    S3Key   string `json:"s3Key"`
    Bedrock bool `json:"bedrock"`
}

type EventDetail struct {
    ConnectionId string `json:"connectionId"`
    Text         string `json:"text"`
    S3Key        string `json:"s3Key"`
    Bedrock      bool   `json:"bedrock"`
}

var svc *eventbridge.EventBridge

func init() {
    sesh := session.Must(session.NewSession())
    svc = eventbridge.New(sesh)
}


func handler(ctx context.Context, request events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
    eventBusName := os.Getenv("EVENT_BUS_NAME")
	if eventBusName == "" {
		log.Println("Error: EVENT_BUS_NAME environment variable is not set")
		return events.APIGatewayProxyResponse{StatusCode: 400}, fmt.Errorf("missing EVENT_BUS_NAME")
	}

    connectionId := request.RequestContext.ConnectionID

    var payload Payload
    err := json.Unmarshal([]byte(request.Body), &payload)
    if err != nil {
        log.Printf("Error parsing JSON: %s\n", err)
        return events.APIGatewayProxyResponse{StatusCode: 400}, nil
    }

    eventDetail := EventDetail{
        ConnectionId: connectionId,
        Text:         payload.Text,
        S3Key:        payload.S3Key,
        Bedrock:      payload.Bedrock,
    }

    eventDetailJSON, err := json.Marshal(eventDetail)
    if err != nil {
        log.Printf("Error marshaling event detail: %s\n", err)
        return events.APIGatewayProxyResponse{StatusCode: 500}, nil
    }

	input := &eventbridge.PutEventsInput{
        Entries: []*eventbridge.PutEventsRequestEntry{
            {
                Detail:     aws.String(string(eventDetailJSON)),
                DetailType: aws.String("PreparePrompt"),
                Source:     aws.String("api-gateway"),
                EventBusName: aws.String(eventBusName),
            },
        },
    }

	result, err := svc.PutEvents(input)
    if err != nil {
        log.Printf("Error putting event: %s\n", err)
        return events.APIGatewayProxyResponse{StatusCode: 500}, nil
    }

	log.Printf("PutEvents result: %v\n", result)


    return events.APIGatewayProxyResponse{StatusCode: 200}, nil
}

func main() {
    lambda.Start(handler)
}
