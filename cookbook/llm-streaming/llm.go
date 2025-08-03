package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Message represents a chat message
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// LLM handles all LLM interactions including streaming
type LLM struct {
	apiKey string
	model  string
	client *http.Client
}

// NewLLM creates a new LLM instance with streaming support
func NewLLM(apiKey string) *LLM {
	return &LLM{
		apiKey: apiKey,
		model:  "gpt-4o",
		client: &http.Client{
			Timeout: 5 * time.Minute, // Longer timeout for streaming
		},
	}
}

// StreamChat sends messages to the LLM and returns a channel for streaming responses
func (l *LLM) StreamChat(ctx context.Context, messages []Message) (<-chan string, error) {
	url := "https://api.openai.com/v1/chat/completions"

	// Convert messages to the format expected by the API
	apiMessages := make([]map[string]string, len(messages))
	for i, msg := range messages {
		apiMessages[i] = map[string]string{
			"role":    msg.Role,
			"content": msg.Content,
		}
	}

	reqBody := map[string]any{
		"model":       l.model,
		"messages":    apiMessages,
		"temperature": 0.7,
		"stream":      true, // Enable streaming
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+l.apiKey)
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")

	resp, err := l.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Create channel for streaming chunks
	chunkChan := make(chan string, 100)

	// Start goroutine to read SSE stream
	go func() {
		defer close(chunkChan)
		defer resp.Body.Close()

		reader := bufio.NewReader(resp.Body)

		for {
			// Check if context is cancelled
			select {
			case <-ctx.Done():
				return
			default:
			}

			line, err := reader.ReadString('\n')
			if err != nil {
				if err != io.EOF {
					fmt.Printf("Error reading stream: %v\n", err)
				}
				return
			}

			line = strings.TrimSpace(line)

			// Skip empty lines
			if line == "" {
				continue
			}

			// Check for SSE data prefix
			if !strings.HasPrefix(line, "data: ") {
				continue
			}

			// Extract JSON data
			data := strings.TrimPrefix(line, "data: ")

			// Check for end of stream
			if data == "[DONE]" {
				return
			}

			// Parse the JSON response
			var streamResp struct {
				Choices []struct {
					Delta struct {
						Content string `json:"content"`
					} `json:"delta"`
				} `json:"choices"`
			}

			if err := json.Unmarshal([]byte(data), &streamResp); err != nil {
				// Skip malformed JSON
				continue
			}

			// Extract content from the response
			if len(streamResp.Choices) > 0 && streamResp.Choices[0].Delta.Content != "" {
				select {
				case chunkChan <- streamResp.Choices[0].Delta.Content:
				case <-ctx.Done():
					return
				}
			}
		}
	}()

	return chunkChan, nil
}

// Stream sends a prompt to the LLM and returns a channel for streaming responses
func (l *LLM) Stream(ctx context.Context, prompt string) (<-chan string, error) {
	// Convert single prompt to message format
	messages := []Message{
		{Role: "user", Content: prompt},
	}
	return l.StreamChat(ctx, messages)
}
