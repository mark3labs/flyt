package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"
	"sync/atomic"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/mark3labs/flyt"
	"github.com/sst/sst/v3/sdk/golang/resource"
)

// ExtractedResponse represents the expected JSON response from Claude
type ExtractedResponse struct {
	Content string `json:"content"`
	Title   string `json:"title"`
}

// ProgressTracker tracks batch processing progress
type ProgressTracker struct {
	total     int
	completed int32
	failed    int32
}

// Simple node that downloads image from S3
func createDownloadNode() flyt.Node {
	return flyt.NewNode(
		flyt.WithPrepFuncAny(func(ctx context.Context, shared *flyt.SharedStore) (any, error) {
			bucket, _ := shared.Get("bucket")
			key, _ := shared.Get("key")
			return map[string]any{
				"bucket": bucket,
				"key":    key,
			}, nil
		}),
		flyt.WithExecFuncAny(func(ctx context.Context, prepResult any) (any, error) {
			data := prepResult.(map[string]any)
			bucket := data["bucket"].(string)
			key := data["key"].(string)

			// Create S3 client
			cfg, err := config.LoadDefaultConfig(ctx)
			if err != nil {
				log.Printf("Failed to load AWS config: %v", err)
				return nil, fmt.Errorf("failed to load AWS config: %w", err)
			}
			s3Client := s3.NewFromConfig(cfg)

			log.Printf("Downloading image from S3: bucket=%s, key=%s", bucket, key)
			result, err := s3Client.GetObject(ctx, &s3.GetObjectInput{
				Bucket: &bucket,
				Key:    &key,
			})
			if err != nil {
				log.Printf("Failed to download from S3: %v", err)
				return nil, fmt.Errorf("failed to download from S3: %w", err)
			}

			imageData, err := io.ReadAll(result.Body)
			if err != nil {
				log.Printf("Failed to read image data: %v", err)
				return nil, fmt.Errorf("failed to read image data: %w", err)
			}

			// Detect MIME type from file extension
			mimeType := "image/jpeg" // default
			lowerKey := strings.ToLower(key)
			if strings.HasSuffix(lowerKey, ".png") {
				mimeType = "image/png"
			} else if strings.HasSuffix(lowerKey, ".jpg") || strings.HasSuffix(lowerKey, ".jpeg") {
				mimeType = "image/jpeg"
			} else if strings.HasSuffix(lowerKey, ".webp") {
				mimeType = "image/webp"
			}

			log.Printf("Successfully downloaded image: %d bytes, MIME type: %s", len(imageData), mimeType)
			return map[string]any{
				"imageData": imageData,
				"mimeType":  mimeType,
			}, nil
		}),
		flyt.WithPostFuncAny(func(ctx context.Context, shared *flyt.SharedStore, prepResult, execResult any) (flyt.Action, error) {
			// Store the image data and mime type for the next node
			result := execResult.(map[string]any)
			shared.Set("imageData", result["imageData"])
			shared.Set("mimeType", result["mimeType"])
			log.Printf("Download node completed, MIME type: %s", result["mimeType"])
			return flyt.DefaultAction, nil
		}),
	)
}

// Simple node that calls Claude to extract text
func createExtractTextNode() flyt.Node {
	return flyt.NewNode(
		flyt.WithMaxRetries(3), // Add retry capability
		flyt.WithPrepFuncAny(func(ctx context.Context, store *flyt.SharedStore) (any, error) {
			// Get the image data and mime type from the previous node
			imageData, _ := store.Get("imageData")
			mimeType, _ := store.Get("mimeType")
			return map[string]any{
				"imageData": imageData,
				"mimeType":  mimeType,
			}, nil
		}),
		flyt.WithPostFuncAny(func(ctx context.Context, store *flyt.SharedStore, prepResult, execResult any) (flyt.Action, error) {
			// Check if extraction failed
			if execResult == nil {
				log.Println("Extract text failed after retries, moving to cleanup")
				return "skip", nil
			}

			// Check for unsupported type
			if str, ok := execResult.(string); ok && strings.HasPrefix(str, "Unsupported image type:") {
				log.Printf("Skipping save for unsupported type: %s", str)
				return "skip", nil
			}

			// Store the extracted response for the next node
			response := execResult.(*ExtractedResponse)
			store.Set("extractedResponse", response)
			log.Printf("Extract text node completed, title: %s, content length: %d", response.Title, len(response.Content))

			return flyt.DefaultAction, nil
		}),
		flyt.WithExecFuncAny(func(ctx context.Context, prepResult any) (any, error) {
			log.Println("Extract text node executing...")
			data := prepResult.(map[string]any)
			imageData := data["imageData"].([]byte)
			mimeType := data["mimeType"].(string)

			// Check if we support this image type
			if mimeType != "image/png" && mimeType != "image/jpeg" {
				log.Printf("Unsupported image type: %s, skipping text extraction", mimeType)
				return "Unsupported image type: " + mimeType, nil
			}

			// Get API key from SST v3 secret
			apiKey, err := resource.Get("AnthropicApiKey", "value")
			if err != nil {
				log.Printf("Failed to get API key from SST: %v", err)
				return nil, fmt.Errorf("failed to get API key: %w", err)
			}

			// Create Anthropic client with API key
			client := anthropic.NewClient(
				option.WithAPIKey(apiKey.(string)),
			)

			// Encode image to base64
			imageEncoded := base64.StdEncoding.EncodeToString(imageData)
			log.Printf("Encoded image for Claude API: %d bytes -> %d base64 chars", len(imageData), len(imageEncoded))

			// Prepare the prompt
			prompt := `Follow these instructions.
- Transcribe the handwritten note. Format the transcription with markdown.
- Based on the contents of the note, come up with a descriptive but concise title.
- Return the transcription as "content" and title as "title" in a JSON object.`

			// Call Claude API to extract text
			log.Println("Calling Claude API to extract text from image...")
			message, err := client.Messages.New(ctx, anthropic.MessageNewParams{
				MaxTokens: 1024,
				Messages: []anthropic.MessageParam{
					anthropic.NewUserMessage(
						anthropic.NewTextBlock(prompt),
						anthropic.NewImageBlockBase64(mimeType, imageEncoded),
					),
				},
				Model: "claude-sonnet-4-20250514",
			})
			if err != nil {
				log.Printf("Failed to call Claude API: %v", err)
				return nil, fmt.Errorf("failed to call Claude API: %w", err)
			}

			// Extract text from response
			if len(message.Content) == 0 {
				log.Println("No content in Claude API response")
				return nil, fmt.Errorf("no content in Claude API response")
			}

			responseText := message.Content[0].Text
			log.Printf("Claude response: %s", responseText)

			// Strip JSON code block markers if present
			responseText = strings.TrimSpace(responseText)
			if strings.HasPrefix(responseText, "```json") {
				responseText = strings.TrimPrefix(responseText, "```json")
				responseText = strings.TrimSuffix(responseText, "```")
				responseText = strings.TrimSpace(responseText)
			}

			// Parse JSON response
			var extracted ExtractedResponse
			if err := json.Unmarshal([]byte(responseText), &extracted); err != nil {
				log.Printf("Failed to parse JSON response: %v", err)
				return nil, fmt.Errorf("failed to parse JSON response: %w", err)
			}
			log.Printf("Successfully parsed response - Title: %s, Content length: %d", extracted.Title, len(extracted.Content))
			return &extracted, nil
		}),
		flyt.WithExecFallbackFunc(func(prepResult any, err error) (any, error) {
			log.Printf("Extract text failed after all retries: %v", err)
			// Return nil to indicate failure, which will trigger skip action in Post
			return nil, nil
		}),
	)
}

// Simple node that saves text to S3
func createSaveTextNode() flyt.Node {
	return flyt.NewNode(
		flyt.WithPrepFuncAny(func(ctx context.Context, store *flyt.SharedStore) (any, error) {
			// Get the extracted response
			response, _ := store.Get("extractedResponse")
			return response, nil
		}),
		flyt.WithExecFuncAny(func(ctx context.Context, prepResult any) (any, error) {
			response := prepResult.(*ExtractedResponse)

			// Create filename from title (replace spaces with underscores)
			filename := strings.ReplaceAll(response.Title, " ", "_") + ".md"

			// Create S3 client
			cfg, err := config.LoadDefaultConfig(ctx)
			if err != nil {
				log.Printf("Failed to load AWS config: %v", err)
				return nil, fmt.Errorf("failed to load AWS config: %w", err)
			}
			s3Client := s3.NewFromConfig(cfg)

			// Get destination bucket name from SST
			destBucketName, err := resource.Get("ExtractedText", "name")
			if err != nil {
				log.Printf("Failed to get destination bucket name: %v", err)
				return nil, fmt.Errorf("failed to get destination bucket name: %w", err)
			}

			bucketName := destBucketName.(string)
			log.Printf("Saving extracted text to S3: bucket=%s, key=%s, size=%d bytes", bucketName, filename, len(response.Content))

			_, err = s3Client.PutObject(ctx, &s3.PutObjectInput{
				Bucket: &bucketName,
				Key:    &filename,
				Body:   strings.NewReader(response.Content),
			})

			if err != nil {
				log.Printf("Failed to save text to S3: %v", err)
				return nil, fmt.Errorf("failed to save text to S3: %w", err)
			}

			log.Printf("Successfully saved extracted text to: s3://%s/%s", bucketName, filename)
			return filename, nil
		}),
	)
}

// Simple node that cleans up the original file from S3
func createCleanupNode() flyt.Node {
	return flyt.NewNode(
		flyt.WithPrepFuncAny(func(ctx context.Context, store *flyt.SharedStore) (any, error) {
			bucket, _ := store.Get("bucket")
			key, _ := store.Get("key")
			return map[string]any{
				"bucket": bucket,
				"key":    key,
			}, nil
		}),
		flyt.WithExecFuncAny(func(ctx context.Context, prepResult any) (any, error) {
			data := prepResult.(map[string]any)
			bucket := data["bucket"].(string)
			key := data["key"].(string)

			// Create S3 client
			cfg, err := config.LoadDefaultConfig(ctx)
			if err != nil {
				log.Printf("Failed to load AWS config for cleanup: %v", err)
				return nil, fmt.Errorf("failed to load AWS config: %w", err)
			}
			s3Client := s3.NewFromConfig(cfg)

			log.Printf("Deleting original file from S3: bucket=%s, key=%s", bucket, key)
			_, err = s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
				Bucket: &bucket,
				Key:    &key,
			})
			if err != nil {
				log.Printf("Failed to delete from S3: %v", err)
				return nil, fmt.Errorf("failed to delete from S3: %w", err)
			}

			log.Printf("Successfully deleted original file: s3://%s/%s", bucket, key)
			return fmt.Sprintf("Deleted s3://%s/%s", bucket, key), nil
		}),
	)
}

// Create a flow factory that returns a new flow instance for each S3 record
func createImageProcessingFlowFactory(tracker *ProgressTracker) func() *flyt.Flow {
	return func() *flyt.Flow {
		// Create nodes with progress tracking
		downloadNode := createDownloadNode()
		extractNode := createExtractTextNode()
		saveNode := createSaveTextNode()

		// Wrap cleanup node with progress tracking
		cleanupNode := flyt.NewNode(
			flyt.WithPrepFuncAny(func(ctx context.Context, store *flyt.SharedStore) (any, error) {
				bucket, _ := store.Get("bucket")
				key, _ := store.Get("key")
				return map[string]any{
					"bucket": bucket,
					"key":    key,
				}, nil
			}),
			flyt.WithExecFuncAny(func(ctx context.Context, prepResult any) (any, error) {
				data := prepResult.(map[string]any)
				bucket := data["bucket"].(string)
				key := data["key"].(string)

				// Create S3 client
				cfg, err := config.LoadDefaultConfig(ctx)
				if err != nil {
					log.Printf("Failed to load AWS config for cleanup: %v", err)
					return nil, fmt.Errorf("failed to load AWS config: %w", err)
				}
				s3Client := s3.NewFromConfig(cfg)

				log.Printf("Deleting original file from S3: bucket=%s, key=%s", bucket, key)
				_, err = s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
					Bucket: &bucket,
					Key:    &key,
				})
				if err != nil {
					log.Printf("Failed to delete from S3: %v", err)
					return nil, fmt.Errorf("failed to delete from S3: %w", err)
				}

				log.Printf("Successfully deleted original file: s3://%s/%s", bucket, key)
				return fmt.Sprintf("Deleted s3://%s/%s", bucket, key), nil
			}),
			flyt.WithPostFuncAny(func(ctx context.Context, shared *flyt.SharedStore, prepResult, execResult any) (flyt.Action, error) {
				// Update progress tracking
				if execResult != nil {
					atomic.AddInt32(&tracker.completed, 1)
				} else {
					atomic.AddInt32(&tracker.failed, 1)
				}

				progress := atomic.LoadInt32(&tracker.completed) + atomic.LoadInt32(&tracker.failed)
				percentage := float64(progress) / float64(tracker.total) * 100

				log.Printf("Batch progress: %.1f%% (%d/%d) - Completed: %d, Failed: %d",
					percentage, progress, tracker.total,
					atomic.LoadInt32(&tracker.completed),
					atomic.LoadInt32(&tracker.failed))

				return flyt.DefaultAction, nil
			}),
		)

		// Create flow
		flow := flyt.NewFlow(downloadNode)
		flow.Connect(downloadNode, flyt.DefaultAction, extractNode)
		flow.Connect(extractNode, flyt.DefaultAction, saveNode)
		flow.Connect(extractNode, "skip", cleanupNode) // Skip saving for unsupported types
		flow.Connect(saveNode, flyt.DefaultAction, cleanupNode)

		return flow
	}
}

func HandleS3Event(ctx context.Context, s3Event events.S3Event) error {
	log.Printf("Received S3 event with %d records", len(s3Event.Records))

	// Create progress tracker
	tracker := &ProgressTracker{
		total: len(s3Event.Records),
	}

	// Create batch function that generates inputs for each S3 record
	batchFunc := func(ctx context.Context, shared *flyt.SharedStore) ([]flyt.FlowInputs, error) {
		inputs := make([]flyt.FlowInputs, len(s3Event.Records))

		for i, record := range s3Event.Records {
			inputs[i] = flyt.FlowInputs{
				"bucket": record.S3.Bucket.Name,
				"key":    record.S3.Object.Key,
				"index":  i + 1,
			}
			log.Printf("Queued record %d/%d: s3://%s/%s",
				i+1, len(s3Event.Records),
				record.S3.Bucket.Name,
				record.S3.Object.Key)
		}

		return inputs, nil
	}

	// Create batch flow with concurrent processing
	batchFlow := flyt.NewBatchFlow(
		createImageProcessingFlowFactory(tracker),
		batchFunc,
		true, // Enable concurrent processing
	)

	// Run the batch flow
	shared := flyt.NewSharedStore()
	if err := batchFlow.Run(ctx, shared); err != nil {
		log.Printf("Batch flow failed: %v", err)
		return fmt.Errorf("batch flow failed: %w", err)
	}

	log.Printf("Successfully processed all %d records - Completed: %d, Failed: %d",
		tracker.total,
		atomic.LoadInt32(&tracker.completed),
		atomic.LoadInt32(&tracker.failed))

	// Return error if any records failed
	if failed := atomic.LoadInt32(&tracker.failed); failed > 0 {
		return fmt.Errorf("%d out of %d records failed processing", failed, tracker.total)
	}

	return nil
}
func main() {
	lambda.Start(HandleS3Event)
}
