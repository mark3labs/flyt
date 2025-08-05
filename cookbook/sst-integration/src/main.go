package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/mark3labs/flyt"
	"github.com/sst/sst/v3/sdk/golang/resource"
)

// Simple node that downloads image from S3
func createDownloadNode() flyt.Node {
	return flyt.NewNode(
		flyt.WithExecFunc(func(ctx context.Context, prepResult any) (any, error) {
			shared := prepResult.(*flyt.SharedStore)
			bucket, _ := shared.Get("bucket")
			key, _ := shared.Get("key")

			// Create S3 client
			cfg, err := config.LoadDefaultConfig(ctx)
			if err != nil {
				return nil, err
			}
			s3Client := s3.NewFromConfig(cfg)

			bucketStr := bucket.(string)
			keyStr := key.(string)

			result, err := s3Client.GetObject(ctx, &s3.GetObjectInput{
				Bucket: &bucketStr,
				Key:    &keyStr,
			})
			if err != nil {
				return nil, err
			}

			imageData, _ := io.ReadAll(result.Body)
			return imageData, nil
		}),
	)
}

// Simple node that calls Claude to extract text
func createExtractTextNode() flyt.Node {
	return flyt.NewNode(
		flyt.WithExecFunc(func(ctx context.Context, prepResult any) (any, error) {
			imageData := prepResult.([]byte)

			// Get API key from SST v3 secret
			apiKey, err := resource.Get("AnthropicApiKey", "value")
			if err != nil {
				return nil, fmt.Errorf("failed to get API key: %w", err)
			}

			// Create Anthropic client with API key
			client := anthropic.NewClient(
				option.WithAPIKey(apiKey.(string)),
			)

			// Encode image to base64
			imageEncoded := base64.StdEncoding.EncodeToString(imageData)

			// Call Claude API to extract text
			message, err := client.Messages.New(ctx, anthropic.MessageNewParams{
				MaxTokens: 1024,
				Messages: []anthropic.MessageParam{
					anthropic.NewUserMessage(
						anthropic.NewTextBlock("Extract all text from this image. If there is no text, describe what you see."),
						anthropic.NewImageBlockBase64("image/jpeg", imageEncoded),
					),
				},
				Model: anthropic.ModelClaude_3_Sonnet_20240229,
			})
			if err != nil {
				return nil, fmt.Errorf("failed to call Claude API: %w", err)
			}

			// Extract text from response
			if len(message.Content) > 0 {
				return message.Content[0].Text, nil
			}

			return "No text extracted", nil
		}),
	)
}

// Simple node that saves text to S3
func createSaveTextNode() flyt.Node {
	return flyt.NewNode(
		flyt.WithExecFunc(func(ctx context.Context, prepResult any) (any, error) {
			text := prepResult.(string)
			shared := ctx.Value("shared").(*flyt.SharedStore)
			originalKey, _ := shared.Get("key")

			// Create S3 client
			cfg, err := config.LoadDefaultConfig(ctx)
			if err != nil {
				return nil, err
			}
			s3Client := s3.NewFromConfig(cfg)

			// Get destination bucket name from SST
			destBucketName, _ := resource.Get("ExtractedText", "name")
			// Save as .txt file with same name
			textKey := strings.Replace(originalKey.(string), filepath.Ext(originalKey.(string)), ".txt", 1)

			bucketName := destBucketName.(string)
			_, err = s3Client.PutObject(ctx, &s3.PutObjectInput{
				Bucket: &bucketName,
				Key:    &textKey,
				Body:   strings.NewReader(text),
			})

			return textKey, err
		}),
	)
}

func HandleS3Event(ctx context.Context, s3Event events.S3Event) error {
	// Create nodes
	downloadNode := createDownloadNode()
	extractNode := createExtractTextNode()
	saveNode := createSaveTextNode()

	// Create flow
	flow := flyt.NewFlow(downloadNode)
	flow.Connect(downloadNode, "", extractNode)
	flow.Connect(extractNode, "", saveNode)

	// Process each S3 record
	for _, record := range s3Event.Records {
		shared := flyt.NewSharedStore()
		shared.Set("bucket", record.S3.Bucket.Name)
		shared.Set("key", record.S3.Object.Key)

		if err := flow.Run(ctx, shared); err != nil {
			return err
		}
	}

	return nil
}

func main() {
	lambda.Start(HandleS3Event)
}
