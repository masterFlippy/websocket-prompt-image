package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi"
	"github.com/aws/aws-sdk-go/service/s3"
)

type InputEvent struct {
	ConnectionId string `json:"connectionId"`
	Prompt string `json:"prompt"`
	S3Key  string `json:"s3Key"`
}


type OutputEvent struct {
	Url   string `json:"url"`
}

// Titan specific types for requests towards bedrock
type TitanInputTextToImageInput struct {
	TaskType              string                                                `json:"taskType"`
	ImageGenerationConfig TitanInputTextToImageConfig `json:"imageGenerationConfig"`
	TextToImageParams     TitanInputTextToImageParams     `json:"textToImageParams"`
}

type TitanInputTextToImageParams struct {
	Text         string `json:"text"`
	NegativeText string `json:"negativeText,omitempty"`
}

type TitanInputTextToImageConfig struct {
	NumberOfImages int     `json:"numberOfImages,omitempty"`
	Height         int     `json:"height,omitempty"`
	Width          int     `json:"width,omitempty"`
	Scale          float64 `json:"cfgScale,omitempty"`
	Seed           int     `json:"seed,omitempty"`
}

type TitanInputTextToImageOutput struct {
	Images []string `json:"images"`
	Error  string   `json:"error"`
}

// Function to decode base64 image
func decodeImage(base64Image string) ([]byte, error) {
	decoded, err := base64.StdEncoding.DecodeString(base64Image)
	if err != nil {
		return nil, err
	}
	return decoded, nil
}

func handler(ctx context.Context, event InputEvent)  error {
	bucketName := os.Getenv("BUCKET_NAME")
	if bucketName == "" {
		log.Println("Error: BUCKET_NAME environment variable is not set")
		return fmt.Errorf("missing BUCKET_NAME")
	}

	websocketEndpoint := os.Getenv("WEBSOCKET_API_ENDPOINT")
    if websocketEndpoint == "" {
        log.Println("Error: WEBSOCKET_API_ENDPOINT environment variable is not set")
        return fmt.Errorf("missing WEBSOCKET_API_ENDPOINT")
    }

	region := os.Getenv("REGION")
	if region == "" {
		log.Println("Error: REGION environment variable is not set")
		return fmt.Errorf("missing REGION")
	}

	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion("us-east-1"))
	if err != nil {
		return err
	}

	runtime := bedrockruntime.NewFromConfig(cfg)

	payload := TitanInputTextToImageInput{
		TaskType: "TEXT_IMAGE",
		TextToImageParams: TitanInputTextToImageParams{
			Text: event.Prompt,
		},
		ImageGenerationConfig: TitanInputTextToImageConfig{
			NumberOfImages: 1,
			Scale: 8.0,
			Height: 1024.0,
			Width: 1024.0,

		},
	}

	payloadString, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("unable to marshal body: %v", err)
	}


	accept := "*/*"
	contentType := "application/json"
	model := "amazon.titan-image-generator-v1"

	resp, err := runtime.InvokeModel(context.TODO(), &bedrockruntime.
	InvokeModelInput{
		Accept:      &accept,
		ModelId:     &model,
		ContentType: &contentType,
		Body:        payloadString,
	})

	if err != nil {
		return fmt.Errorf("error from Bedrock, %v", err)
	}

	var output TitanInputTextToImageOutput

	err = json.Unmarshal(resp.Body, &output)
	if err != nil {
		return fmt.Errorf("unable to unmarshal response from Bedrock: %v", err)
	}
	decoded, err := decodeImage(output.Images[0])
	if err != nil {
		return fmt.Errorf("unable to decode image: %v", err)
	}

	sesh := session.Must(session.NewSession())

	s3Client := s3.New(sesh)

	objectKey := event.S3Key
	
	_, err = s3Client.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(objectKey),
		Body:        bytes.NewReader(decoded),
		ContentType: aws.String(http.DetectContentType(decoded)),
	})
	if err != nil {
		log.Println("Error uploading image to S3:", err)
		return err
	}

	log.Println("Successfully uploaded image to S3:", objectKey)

    url := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", bucketName, region, objectKey)


    apiClient := apigatewaymanagementapi.New(sesh, aws.NewConfig().WithEndpoint(websocketEndpoint))

    postToConnectionInput := &apigatewaymanagementapi.PostToConnectionInput{
        ConnectionId: aws.String(event.ConnectionId),
        Data:         []byte(fmt.Sprintf(`{"url": "%s"}`, url)),
    }

	_, err = apiClient.PostToConnection(postToConnectionInput)
    if err != nil {
        log.Println("Error sending response to WebSocket client:", err)
        return err
    }

    log.Println("Successfully sent response to WebSocket client:", event.ConnectionId)

	return  nil
}

func main() {
	lambda.Start(handler)
}
