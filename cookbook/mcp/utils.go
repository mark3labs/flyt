package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
)

// CallLLMWithFunctions calls OpenAI with function calling
func CallLLMWithFunctions(apiKey string, messages []map[string]interface{}, tools []mcp.Tool) (map[string]interface{}, error) {
	url := "https://api.openai.com/v1/chat/completions"

	// Convert MCP tools to OpenAI function format
	functions := make([]map[string]interface{}, 0, len(tools))
	for _, tool := range tools {
		function := map[string]interface{}{
			"name":        tool.Name,
			"description": tool.Description,
			"parameters":  tool.InputSchema,
		}
		functions = append(functions, function)
	}

	requestBody := map[string]interface{}{
		"model":         "gpt-4o-mini",
		"messages":      messages,
		"functions":     functions,
		"function_call": "auto",
		"temperature":   0.3,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var response struct {
		ID      string `json:"id"`
		Choices []struct {
			Message struct {
				Role         string `json:"role"`
				Content      string `json:"content"`
				FunctionCall *struct {
					Name      string `json:"name"`
					Arguments string `json:"arguments"`
				} `json:"function_call"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(response.Choices) == 0 {
		return nil, fmt.Errorf("no choices in response")
	}

	choice := response.Choices[0]
	result := map[string]interface{}{
		"message": choice.Message,
	}

	// Check if there's a function call
	if choice.Message.FunctionCall != nil {
		// Parse the function arguments
		var args map[string]interface{}
		if err := json.Unmarshal([]byte(choice.Message.FunctionCall.Arguments), &args); err != nil {
			return nil, fmt.Errorf("failed to parse function arguments: %w", err)
		}

		result["function_call"] = map[string]interface{}{
			"type":      "function_call",
			"id":        response.ID,
			"call_id":   "call_" + response.ID,
			"name":      choice.Message.FunctionCall.Name,
			"arguments": args,
		}
	}

	return result, nil
}
