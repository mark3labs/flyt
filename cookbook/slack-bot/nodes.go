package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/mark3labs/flyt"
)

// NewParseMessageNode creates a node that prepares the message for processing
func NewParseMessageNode(slackService *SlackService) flyt.Node {
	return flyt.NewNode().
		WithPrepFuncAny(func(ctx context.Context, shared *flyt.SharedStore) (any, error) {
			messageStr := shared.GetString("message")
			if messageStr == "" {
				return flyt.R(nil), fmt.Errorf("no message found in shared store")
			}

			// Clean up the message (remove bot mentions, extra spaces, etc.)
			cleanedMessage := strings.TrimSpace(messageStr)
			cleanedMessage = strings.ReplaceAll(cleanedMessage, "  ", " ")

			// Remove bot mention if present (format: <@BOTID>)
			if strings.Contains(cleanedMessage, "<@") {
				// If we have SlackService, we can be more precise about bot ID
				if slackService != nil {
					botID := slackService.GetBotUserID()
					mentionPattern := fmt.Sprintf("<@%s>", botID)
					cleanedMessage = strings.ReplaceAll(cleanedMessage, mentionPattern, "")
				} else {
					// Fallback to generic mention removal
					parts := strings.Split(cleanedMessage, ">")
					if len(parts) > 1 {
						cleanedMessage = strings.TrimSpace(strings.Join(parts[1:], ">"))
					}
				}
			}

			return flyt.R(strings.TrimSpace(cleanedMessage)), nil
		}).
		WithExecFuncAny(func(ctx context.Context, prepResult any) (any, error) {
			// Message is already cleaned in Prep
			return prepResult, nil
		}).
		WithPostFuncAny(func(ctx context.Context, shared *flyt.SharedStore, prepResult, execResult any) (flyt.Action, error) {
			shared.Set("cleaned_message", execResult)
			log.Printf("Parsed message: %v", execResult)
			return flyt.DefaultAction, nil
		}).
		Build()
}

// NewLLMNode creates a node that processes the message through OpenAI with function calling
func NewLLMNode(llmService *LLMService) flyt.Node {
	return flyt.NewNode().
		WithPrepFuncAny(func(ctx context.Context, shared *flyt.SharedStore) (any, error) {
			// Check if we're processing tool responses
			if shared.Has("tool_responses") {
				toolResponses, _ := shared.Get("tool_responses")
				return flyt.R(map[string]interface{}{
					"type":           "tool_response",
					"tool_responses": toolResponses,
				}), nil
			}

			// Otherwise, process user message
			message := shared.GetString("cleaned_message")
			if message == "" {
				return flyt.R(nil), fmt.Errorf("no cleaned message found")
			}

			// Get conversation history if available
			var history []map[string]string
			if h, ok := shared.Get("history"); ok {
				if hist, ok := h.([]map[string]string); ok {
					history = hist
				}
			}

			return flyt.R(map[string]interface{}{
				"type":    "user_message",
				"message": message,
				"history": history,
			}), nil
		}).
		WithExecFuncAny(func(ctx context.Context, prepResult any) (any, error) {
			prepRes, ok := prepResult.(flyt.Result)
			if !ok {
				prepRes = flyt.R(prepResult)
			}
			data := prepRes.MustMap()

			if data["type"] == "tool_response" {
				// Process tool responses
				toolResponses := data["tool_responses"].(map[string]string)
				response, err := llmService.ProcessToolResponses(ctx, toolResponses)
				if err != nil {
					return flyt.R(nil), fmt.Errorf("failed to process tool responses: %w", err)
				}
				return flyt.R(map[string]interface{}{
					"type":     "final_response",
					"response": response,
				}), nil
			}

			// Process user message with optional history
			message := data["message"].(string)

			var response string
			var toolCalls []ToolCall
			var err error

			// Check if we have history
			if history, ok := data["history"].([]map[string]string); ok && len(history) > 0 {
				response, toolCalls, err = llmService.ProcessMessageWithHistory(ctx, message, history)
			} else {
				response, toolCalls, err = llmService.ProcessMessage(ctx, message)
			}

			if err != nil {
				return flyt.R(nil), fmt.Errorf("failed to process with LLM: %w", err)
			}

			if len(toolCalls) > 0 {
				return flyt.R(map[string]interface{}{
					"type":       "tool_calls",
					"tool_calls": toolCalls,
				}), nil
			}

			return flyt.R(map[string]interface{}{
				"type":     "response",
				"response": response,
			}), nil
		}).
		WithPostFuncAny(func(ctx context.Context, shared *flyt.SharedStore, prepResult, execResult any) (flyt.Action, error) {
			execRes, ok := execResult.(flyt.Result)
			if !ok {
				execRes = flyt.R(execResult)
			}
			result := execRes.MustMap()

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
		}).
		Build()
}

// NewToolExecutorNode creates a node that executes tool calls requested by the LLM
func NewToolExecutorNode() flyt.Node {
	return flyt.NewNode().
		WithPrepFuncAny(func(ctx context.Context, shared *flyt.SharedStore) (any, error) {
			toolCalls, ok := shared.Get("tool_calls")
			if !ok {
				return flyt.R(nil), fmt.Errorf("no tool calls found")
			}
			return flyt.R(toolCalls), nil
		}).
		WithExecFuncAny(func(ctx context.Context, prepResult any) (any, error) {
			prepRes, ok := prepResult.(flyt.Result)
			if !ok {
				prepRes = flyt.R(prepResult)
			}

			toolCalls, ok := flyt.As[[]ToolCall](prepRes)
			if !ok {
				return flyt.R(nil), fmt.Errorf("invalid tool calls format")
			}

			log.Printf("Executing %d tool calls", len(toolCalls))

			// Execute all tool calls
			results, err := ExecuteToolCalls(toolCalls)
			if err != nil {
				return flyt.R(nil), fmt.Errorf("failed to execute tools: %w", err)
			}

			return flyt.R(results), nil
		}).
		WithPostFuncAny(func(ctx context.Context, shared *flyt.SharedStore, prepResult, execResult any) (flyt.Action, error) {
			// Store tool responses and go back to LLM
			shared.Set("tool_responses", execResult)
			log.Printf("Tool execution completed with %d results", len(execResult.(map[string]string)))
			return flyt.DefaultAction, nil // Goes back to LLM node
		}).
		Build()
}

// NewFormatResponseNode creates a node that formats the final response for Slack
func NewFormatResponseNode() flyt.Node {
	return flyt.NewNode().
		WithPrepFuncAny(func(ctx context.Context, shared *flyt.SharedStore) (any, error) {
			response := shared.GetString("response")
			if response == "" {
				return flyt.R(nil), fmt.Errorf("no response found")
			}
			return flyt.R(response), nil
		}).
		WithExecFuncAny(func(ctx context.Context, prepResult any) (any, error) {
			prepRes, ok := prepResult.(flyt.Result)
			if !ok {
				prepRes = flyt.R(prepResult)
			}
			response := prepRes.MustString()

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

			return flyt.R(formattedResponse), nil
		}).
		WithPostFuncAny(func(ctx context.Context, shared *flyt.SharedStore, prepResult, execResult any) (flyt.Action, error) {
			shared.Set("response", execResult)
			log.Println("Response formatted for Slack")
			return flyt.DefaultAction, nil
		}).
		Build()
}
