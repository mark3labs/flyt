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

// CreateStreamNode creates a node that streams LLM responses
func CreateStreamNode(llm *LLM) flyt.Node {
	return flyt.NewNode().
		WithPrepFuncAny(func(ctx context.Context, shared *flyt.SharedStore) (any, error) {
			// Get messages from shared store
			messages, ok := shared.Get("messages")
			if !ok {
				return nil, fmt.Errorf("messages not found in shared store")
			}

			return map[string]any{
				"messages": messages.([]Message),
			}, nil
		}).
		WithExecFuncAny(func(ctx context.Context, prepResult any) (any, error) {
			data := prepResult.(map[string]any)
			messages := data["messages"].([]Message)

			fmt.Printf("\n🤖 Assistant: ")

			// Stream the response with full conversation history
			streamChan, err := llm.StreamChat(ctx, messages)
			if err != nil {
				return nil, fmt.Errorf("failed to start streaming: %w", err)
			}

			response := strings.Builder{}

			for chunk := range streamChan {
				fmt.Print(chunk)
				response.WriteString(chunk)
			}

			fmt.Println() // New line after streaming

			return response.String(), nil
		}).
		WithPostFuncAny(func(ctx context.Context, shared *flyt.SharedStore, prepResult, execResult any) (flyt.Action, error) {
			if execResult != nil {
				// Add assistant's response to messages
				messages, _ := shared.Get("messages")
				messageList := messages.([]Message)
				messageList = append(messageList, Message{
					Role:    "assistant",
					Content: execResult.(string),
				})
				shared.Set("messages", messageList)
			}

			return flyt.DefaultAction, nil
		}).
		Build()
}

// CreateInteractiveChatFlow creates a flow that continuously prompts for input and streams responses
func CreateInteractiveChatFlow(llm *LLM) *flyt.Flow {
	// Create input node
	inputNode := flyt.NewNode().
		WithPrepFuncAny(func(ctx context.Context, shared *flyt.SharedStore) (any, error) {
			// Initialize messages if this is the first run
			messages, ok := shared.Get("messages")
			if !ok {
				messages = []Message{}
				shared.Set("messages", messages)
				fmt.Println("Welcome to the streaming chat! Type 'exit' to end the conversation.")
			}
			return nil, nil
		}).
		WithExecFuncAny(func(ctx context.Context, prepResult any) (any, error) {
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("\n💭 You: ")
			input, err := reader.ReadString('\n')
			if err != nil {
				return nil, err
			}
			input = strings.TrimSpace(input)

			if strings.ToLower(input) == "exit" {
				return nil, nil
			}

			return input, nil
		}).
		WithPostFuncAny(func(ctx context.Context, shared *flyt.SharedStore, prepResult, execResult any) (flyt.Action, error) {
			if execResult == nil {
				fmt.Println("\n👋 Goodbye!")
				return "exit", nil
			}

			// Add user message to history
			messages, _ := shared.Get("messages")
			messageList := messages.([]Message)
			messageList = append(messageList, Message{
				Role:    "user",
				Content: execResult.(string),
			})
			shared.Set("messages", messageList)

			return "stream", nil
		}).
		Build()

	// Create stream node
	streamNode := CreateStreamNode(llm)

	// Create flow
	return flyt.NewFlow(inputNode).
		Connect(inputNode, "stream", streamNode).
		Connect(streamNode, flyt.DefaultAction, inputNode) // Loop back for next input
}

func main() {
	// Get API key from environment or command line
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		if len(os.Args) > 1 {
			apiKey = os.Args[1]
		} else {
			log.Fatal("Please set OPENAI_API_KEY environment variable or provide as argument")
		}
	}

	// Create LLM instance with streaming support
	llm := NewLLM(apiKey)

	fmt.Println("🚀 LLM Streaming Demo")
	fmt.Println("Type 'exit' to quit")
	fmt.Println(strings.Repeat("-", 50))

	if len(os.Args) > 2 && os.Args[2] == "--single" {
		// Single prompt mode
		prompt := "What's the meaning of life?"
		if len(os.Args) > 3 {
			prompt = strings.Join(os.Args[3:], " ")
		}

		node := CreateStreamNode(llm)
		flow := flyt.NewFlow(node)

		shared := flyt.NewSharedStore()
		// Initialize with a single user message
		shared.Set("messages", []Message{
			{Role: "user", Content: prompt},
		})

		fmt.Printf("📝 Prompt: %s\n", prompt)

		ctx := context.Background()
		if err := flow.Run(ctx, shared); err != nil {
			log.Fatal(err)
		}

		messages, _ := shared.Get("messages")
		messageList := messages.([]Message)
		if len(messageList) > 1 {
			fmt.Printf("\n📊 Summary:\n")
			fmt.Printf("   Response length: %d characters\n", len(messageList[1].Content))
		}
	} else {
		// Interactive chat mode
		flow := CreateInteractiveChatFlow(llm)
		shared := flyt.NewSharedStore()

		ctx := context.Background()
		if err := flow.Run(ctx, shared); err != nil {
			log.Fatal(err)
		}
	}
}
