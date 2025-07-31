package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/mark3labs/flyt"
)

// Message represents a chat message
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// CreateChatNode creates an interactive chat node using the NewNode helper
func CreateChatNode(apiKey string) flyt.Node {
	return flyt.NewNode(
		flyt.WithPrepFunc(func(ctx context.Context, shared *flyt.SharedStore) (any, error) {
			// Initialize messages if this is the first run
			messages, ok := shared.Get("messages")
			if !ok {
				messages = []Message{}
				shared.Set("messages", messages)
				fmt.Println("Welcome to the chat! Type 'exit' to end the conversation.")
			}

			// Get user input
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("\nYou: ")
			userInput, err := reader.ReadString('\n')
			if err != nil {
				return nil, err
			}
			userInput = strings.TrimSpace(userInput)

			// Check if user wants to exit
			if strings.ToLower(userInput) == "exit" {
				return nil, nil
			}

			// Add user message to history
			messageList := messages.([]Message)
			messageList = append(messageList, Message{
				Role:    "user",
				Content: userInput,
			})
			shared.Set("messages", messageList)

			return messageList, nil
		}),
		flyt.WithExecFunc(func(ctx context.Context, prepResult any) (any, error) {
			if prepResult == nil {
				return nil, nil
			}

			messages := prepResult.([]Message)

			// Call LLM with the entire conversation history
			response, err := CallLLM(apiKey, messages)
			if err != nil {
				return nil, fmt.Errorf("LLM call failed: %w", err)
			}

			return response, nil
		}),
		flyt.WithPostFunc(func(ctx context.Context, shared *flyt.SharedStore, prepResult, execResult any) (flyt.Action, error) {
			if prepResult == nil || execResult == nil {
				fmt.Println("\nGoodbye!")
				return "end", nil
			}

			// Print and store the response
			response := execResult.(string)
			fmt.Printf("\nAssistant: %s\n", response)

			messages, _ := shared.Get("messages")
			messageList := messages.([]Message)
			messageList = append(messageList, Message{
				Role:    "assistant",
				Content: response,
			})
			shared.Set("messages", messageList)

			return "continue", nil
		}),
	)
}

func main() {
	// Get API key from environment
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("Please set OPENAI_API_KEY environment variable")
	}

	// Create chat node
	chatNode := CreateChatNode(apiKey)

	// Create the flow with self-loop
	flow := flyt.NewFlow(chatNode)
	flow.Connect(chatNode, "continue", chatNode) // Loop back to continue conversation

	// Start the chat
	shared := flyt.NewSharedStore()
	ctx := context.Background()

	if err := flow.Run(ctx, shared); err != nil {
		log.Fatal(err)
	}
}
