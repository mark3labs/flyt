package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// OpenAI API structures for function calling
type ChatCompletionRequest struct {
	Model       string        `json:"model"`
	Messages    []ChatMessage `json:"messages"`
	Tools       []Tool        `json:"tools,omitempty"`
	ToolChoice  interface{}   `json:"tool_choice,omitempty"`
	Temperature float64       `json:"temperature"`
}

type ChatMessage struct {
	Role       string     `json:"role"`
	Content    string     `json:"content,omitempty"`
	ToolCalls  []ToolCall `json:"tool_calls,omitempty"`
	ToolCallID string     `json:"tool_call_id,omitempty"`
}

type Tool struct {
	Type     string      `json:"type"`
	Function FunctionDef `json:"function"`
}

type FunctionDef struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
}

type ToolCall struct {
	ID       string       `json:"id"`
	Type     string       `json:"type"`
	Function FunctionCall `json:"function"`
}

type FunctionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type ChatCompletionResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

type Choice struct {
	Index        int         `json:"index"`
	Message      ChatMessage `json:"message"`
	FinishReason string      `json:"finish_reason"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// Tool definitions
func getToolDefinitions() []Tool {
	return []Tool{
		{
			Type: "function",
			Function: FunctionDef{
				Name:        "calculator",
				Description: "Perform mathematical calculations",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"expression": map[string]interface{}{
							"type":        "string",
							"description": "The mathematical expression to evaluate (e.g., '2 + 2', '10 * 5', 'sqrt(16)')",
						},
					},
					"required": []string{"expression"},
				},
			},
		},
		{
			Type: "function",
			Function: FunctionDef{
				Name:        "chuck_norris_fact",
				Description: "Get a random Chuck Norris fact",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"category": map[string]interface{}{
							"type":        "string",
							"description": "Optional category for the fact (e.g., 'dev', 'movie', 'food', 'sport')",
							"enum":        []string{"dev", "movie", "food", "sport", "random"},
						},
					},
				},
			},
		},
	}
}

// LLMClient handles OpenAI API interactions
type LLMClient struct {
	apiKey     string
	httpClient *http.Client
}

func NewLLMClient(apiKey string) *LLMClient {
	return &LLMClient{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *LLMClient) CreateChatCompletion(ctx context.Context, messages []ChatMessage) (*ChatCompletionResponse, error) {
	request := ChatCompletionRequest{
		Model:       "gpt-4.1",
		Messages:    messages,
		Tools:       getToolDefinitions(),
		ToolChoice:  "auto",
		Temperature: 0.7,
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var completion ChatCompletionResponse
	if err := json.Unmarshal(body, &completion); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &completion, nil
}

// LLMService encapsulates the LLM client and conversation management
type LLMService struct {
	client       *LLMClient
	conversation *ConversationManager
}

// NewLLMService creates a new LLM service with its own conversation manager
func NewLLMService(apiKey string) *LLMService {
	return &LLMService{
		client:       NewLLMClient(apiKey),
		conversation: NewConversationManager(),
	}
}

// ProcessMessage processes a user message through the LLM
func (s *LLMService) ProcessMessage(ctx context.Context, message string) (string, []ToolCall, error) {
	// Add user message to conversation
	s.conversation.AddUserMessage(message)

	// Get completion from OpenAI
	response, err := s.client.CreateChatCompletion(ctx, s.conversation.GetMessages())
	if err != nil {
		return "", nil, fmt.Errorf("failed to get completion: %w", err)
	}

	if len(response.Choices) == 0 {
		return "", nil, fmt.Errorf("no choices in response")
	}

	choice := response.Choices[0]
	assistantMessage := choice.Message

	// Add assistant message to conversation
	s.conversation.AddAssistantMessage(assistantMessage)

	// Check if there are tool calls
	if len(assistantMessage.ToolCalls) > 0 {
		return "", assistantMessage.ToolCalls, nil
	}

	// Return the text response
	return assistantMessage.Content, nil, nil
}

// ProcessMessageWithHistory processes a message with conversation history context
func (s *LLMService) ProcessMessageWithHistory(ctx context.Context, message string, history []map[string]string) (string, []ToolCall, error) {
	// If we have history and this is a fresh conversation, add the history as context
	if len(history) > 0 && len(s.conversation.messages) <= 2 {
		// Build conversation history in a natural format
		for _, msg := range history {
			// Add each historical message as a user message
			// This gives the LLM the full context of the conversation
			s.conversation.messages = append(s.conversation.messages, ChatMessage{
				Role:    "user",
				Content: msg["text"],
			})
			// Add a placeholder assistant response to maintain conversation flow
			// (In a real implementation, you'd store and retrieve actual bot responses)
			s.conversation.messages = append(s.conversation.messages, ChatMessage{
				Role:    "assistant",
				Content: "[Previous response in thread]",
			})
		}
	}

	// Process the current message
	return s.ProcessMessage(ctx, message)
}

// ProcessToolResponses processes tool responses and gets final answer from LLM
func (s *LLMService) ProcessToolResponses(ctx context.Context, toolResponses map[string]string) (string, error) {
	// Add tool responses to conversation
	for toolCallID, response := range toolResponses {
		s.conversation.AddToolResponse(toolCallID, response)
	}

	// Get final response from LLM
	response, err := s.client.CreateChatCompletion(ctx, s.conversation.GetMessages())
	if err != nil {
		return "", fmt.Errorf("failed to get final completion: %w", err)
	}

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("no choices in final response")
	}

	finalMessage := response.Choices[0].Message
	s.conversation.AddAssistantMessage(finalMessage)

	return finalMessage.Content, nil
}

// ConversationManager manages the conversation history
type ConversationManager struct {
	messages []ChatMessage
	maxSize  int
}

func NewConversationManager() *ConversationManager {
	return &ConversationManager{
		messages: []ChatMessage{
			{
				Role: "system",
				Content: `You are a helpful Slack bot assistant with access to two tools:
1. A calculator for mathematical computations
2. A Chuck Norris fact generator for entertainment

You should use these tools when appropriate to help users. Be friendly, concise, and helpful in your responses.
When users ask for calculations, use the calculator tool.
When users want entertainment or Chuck Norris facts, use the chuck_norris_fact tool.

You are operating in Slack, so:
- Keep responses concise and well-formatted for Slack
- Use *bold* for emphasis (not **bold**)
- Be aware you may be in a thread with conversation history
- In channels, you only respond when directly mentioned
- Be professional but friendly`,
			},
		},
		maxSize: 20, // Keep last 20 messages
	}
}

func (cm *ConversationManager) AddUserMessage(content string) {
	cm.messages = append(cm.messages, ChatMessage{
		Role:    "user",
		Content: content,
	})
	cm.trimMessages()
}

func (cm *ConversationManager) AddAssistantMessage(message ChatMessage) {
	cm.messages = append(cm.messages, message)
	cm.trimMessages()
}

func (cm *ConversationManager) AddToolResponse(toolCallID, content string) {
	cm.messages = append(cm.messages, ChatMessage{
		Role:       "tool",
		Content:    content,
		ToolCallID: toolCallID,
	})
	cm.trimMessages()
}

func (cm *ConversationManager) GetMessages() []ChatMessage {
	return cm.messages
}

func (cm *ConversationManager) trimMessages() {
	// Keep system message and last N messages
	if len(cm.messages) > cm.maxSize {
		systemMsg := cm.messages[0]
		cm.messages = append([]ChatMessage{systemMsg}, cm.messages[len(cm.messages)-cm.maxSize+1:]...)
	}
}
