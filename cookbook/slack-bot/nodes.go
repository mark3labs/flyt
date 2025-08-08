package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/mark3labs/flyt"
)

// ParseMessageNode prepares the message for processing
type ParseMessageNode struct {
	*flyt.BaseNode
}

func (n *ParseMessageNode) Prep(ctx context.Context, shared *flyt.SharedStore) (any, error) {
	message, ok := shared.Get("message")
	if !ok {
		return nil, fmt.Errorf("no message found in shared store")
	}

	messageStr, ok := message.(string)
	if !ok {
		return nil, fmt.Errorf("message is not a string")
	}

	// Clean up the message (remove bot mentions, extra spaces, etc.)
	cleanedMessage := strings.TrimSpace(messageStr)
	cleanedMessage = strings.ReplaceAll(cleanedMessage, "  ", " ")

	// Remove bot mention if present (format: <@BOTID>)
	if strings.Contains(cleanedMessage, "<@") {
		parts := strings.Split(cleanedMessage, ">")
		if len(parts) > 1 {
			cleanedMessage = strings.TrimSpace(strings.Join(parts[1:], ">"))
		}
	}

	return cleanedMessage, nil
}

func (n *ParseMessageNode) Exec(ctx context.Context, prepResult any) (any, error) {
	// Message is already cleaned in Prep
	return prepResult, nil
}

func (n *ParseMessageNode) Post(ctx context.Context, shared *flyt.SharedStore, prepResult, execResult any) (flyt.Action, error) {
	shared.Set("cleaned_message", execResult)
	log.Printf("Parsed message: %v", execResult)
	return flyt.DefaultAction, nil
}

// LLMNode processes the message through OpenAI with function calling
type LLMNode struct {
	*flyt.BaseNode
	llm *LLMService
}

func (n *LLMNode) Prep(ctx context.Context, shared *flyt.SharedStore) (any, error) {
	// Check if we're processing tool responses
	if toolResponses, ok := shared.Get("tool_responses"); ok {
		return map[string]interface{}{
			"type":           "tool_response",
			"tool_responses": toolResponses,
		}, nil
	}

	// Otherwise, process user message
	message, ok := shared.Get("cleaned_message")
	if !ok {
		return nil, fmt.Errorf("no cleaned message found")
	}

	return map[string]interface{}{
		"type":    "user_message",
		"message": message,
	}, nil
}
func (n *LLMNode) Exec(ctx context.Context, prepResult any) (any, error) {
	data := prepResult.(map[string]interface{})

	if data["type"] == "tool_response" {
		// Process tool responses
		toolResponses := data["tool_responses"].(map[string]string)
		response, err := n.llm.ProcessToolResponses(ctx, toolResponses)
		if err != nil {
			return nil, fmt.Errorf("failed to process tool responses: %w", err)
		}
		return map[string]interface{}{
			"type":     "final_response",
			"response": response,
		}, nil
	}

	// Process user message
	message := data["message"].(string)
	response, toolCalls, err := n.llm.ProcessMessage(ctx, message)
	if err != nil {
		return nil, fmt.Errorf("failed to process with LLM: %w", err)
	}

	if len(toolCalls) > 0 {
		return map[string]interface{}{
			"type":       "tool_calls",
			"tool_calls": toolCalls,
		}, nil
	}

	return map[string]interface{}{
		"type":     "response",
		"response": response,
	}, nil
}
func (n *LLMNode) Post(ctx context.Context, shared *flyt.SharedStore, prepResult, execResult any) (flyt.Action, error) {
	result := execResult.(map[string]interface{})

	switch result["type"] {
	case "tool_calls":
		// Need to execute tools
		shared.Set("tool_calls", result["tool_calls"])
		log.Println("LLM requested tool calls")
		return "tool_call", nil
	case "response", "final_response":
		// Got a text response
		shared.Set("response", result["response"])
		log.Printf("LLM response: %v", result["response"])
		return "response", nil
	default:
		return flyt.DefaultAction, nil
	}
}

// ToolExecutorNode executes tool calls requested by the LLM
type ToolExecutorNode struct {
	*flyt.BaseNode
}

func (n *ToolExecutorNode) Prep(ctx context.Context, shared *flyt.SharedStore) (any, error) {
	toolCalls, ok := shared.Get("tool_calls")
	if !ok {
		return nil, fmt.Errorf("no tool calls found")
	}
	return toolCalls, nil
}

func (n *ToolExecutorNode) Exec(ctx context.Context, prepResult any) (any, error) {
	toolCalls, ok := prepResult.([]ToolCall)
	if !ok {
		return nil, fmt.Errorf("invalid tool calls format")
	}

	log.Printf("Executing %d tool calls", len(toolCalls))

	// Execute all tool calls
	results, err := ExecuteToolCalls(toolCalls)
	if err != nil {
		return nil, fmt.Errorf("failed to execute tools: %w", err)
	}

	return results, nil
}

func (n *ToolExecutorNode) Post(ctx context.Context, shared *flyt.SharedStore, prepResult, execResult any) (flyt.Action, error) {
	// Store tool responses and go back to LLM
	shared.Set("tool_responses", execResult)
	log.Printf("Tool execution completed with %d results", len(execResult.(map[string]string)))
	return flyt.DefaultAction, nil // Goes back to LLM node
}

// FormatResponseNode formats the final response for Slack
type FormatResponseNode struct {
	*flyt.BaseNode
}

func (n *FormatResponseNode) Prep(ctx context.Context, shared *flyt.SharedStore) (any, error) {
	response, ok := shared.Get("response")
	if !ok {
		return nil, fmt.Errorf("no response found")
	}
	return response, nil
}

func (n *FormatResponseNode) Exec(ctx context.Context, prepResult any) (any, error) {
	response := prepResult.(string)

	// Format the response for Slack (could add formatting, emojis, etc.)
	formattedResponse := response

	// Add some Slack-specific formatting if needed
	// For example, convert markdown bold to Slack bold
	formattedResponse = strings.ReplaceAll(formattedResponse, "**", "*")

	// Ensure response isn't too long for Slack
	maxLength := 3000
	if len(formattedResponse) > maxLength {
		formattedResponse = formattedResponse[:maxLength-3] + "..."
	}

	return formattedResponse, nil
}

func (n *FormatResponseNode) Post(ctx context.Context, shared *flyt.SharedStore, prepResult, execResult any) (flyt.Action, error) {
	shared.Set("response", execResult)
	log.Println("Response formatted for Slack")
	return flyt.DefaultAction, nil
}
