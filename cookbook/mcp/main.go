package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/mark3labs/flyt"
)

func main() {
	// Parse command line flags
	var question string
	var apiKey string
	flag.StringVar(&question, "q", "What is 10 plus 20?", "Question to ask")
	flag.StringVar(&apiKey, "key", "", "OpenAI API key (or set OPENAI_API_KEY env var)")
	flag.Parse()

	// Get API key from flag or environment
	if apiKey == "" {
		apiKey = os.Getenv("OPENAI_API_KEY")
	}
	if apiKey == "" {
		log.Fatal("Please provide OpenAI API key via -key flag or OPENAI_API_KEY environment variable")
	}

	fmt.Printf("🤔 Processing question: %s\n", question)

	ctx := context.Background()

	// Create MCP client with in-process server
	mcpClient, err := NewMCPClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create MCP client: %v", err)
	}

	// Create LLM instance
	llm := NewLLM(apiKey)

	// Create nodes
	getToolsNode := NewGetToolsNode(mcpClient)
	decideNode := NewDecideToolNode(llm)
	executeNode := NewExecuteToolNode(mcpClient)

	// Create flow
	flow := flyt.NewFlow(getToolsNode).
		Connect(getToolsNode, "decide", decideNode).
		Connect(decideNode, "execute", executeNode).
		Connect(executeNode, "decide", decideNode) // Loop back for continued conversation

	// Create shared store and initialize with the user's question
	shared := flyt.NewSharedStore()

	// Initialize messages with system prompt and user question
	messages := []map[string]interface{}{
		{
			"role":    "system",
			"content": "You are a helpful assistant that performs mathematical calculations using the available functions.",
		},
		{
			"role":    "user",
			"content": question,
		},
	}
	shared.Set("messages", messages)

	// Run the flow
	if err := flow.Run(ctx, shared); err != nil {
		log.Fatalf("Flow failed: %v", err)
	}
}
