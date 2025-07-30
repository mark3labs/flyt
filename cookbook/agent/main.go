// Package main implements an AI research agent using the Flyt workflow framework.
// The agent can answer questions by searching the web when needed and generating
// comprehensive answers using OpenAI's API.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	// Define command line flags
	question := flag.String("q", "Who won the Nobel Prize in Physics 2024?", "Question to ask the agent")
	apiKey := flag.String("key", "", "OpenAI API key (or set OPENAI_API_KEY env var)")
	braveKey := flag.String("brave", "", "Brave Search API key (or set BRAVE_API_KEY env var)")
	flag.Parse()

	key := *apiKey
	if key == "" {
		key = os.Getenv("OPENAI_API_KEY")
	}
	if key == "" {
		log.Fatal("Please provide OpenAI API key via -key flag or OPENAI_API_KEY environment variable")
	}

	brave := *braveKey
	if brave == "" {
		brave = os.Getenv("BRAVE_API_KEY")
	}

	agentFlow := CreateAgentFlow(key)

	shared := NewSharedState()
	shared.Set("question", *question)
	shared.Set("api_key", key)
	if brave != "" {
		shared.Set("brave_api_key", brave)
	}

	fmt.Printf("🤔 Processing question: %s\n", *question)

	ctx := context.Background()
	err := agentFlow.Run(ctx, shared)
	if err != nil {
		log.Fatalf("Error running agent: %v", err)
	}

	answerVal, _ := shared.Get("answer")
	answer, ok := answerVal.(string)
	if !ok || answer == "" {
		log.Fatal("No answer found")
	}
	fmt.Println("\n🎯 Final Answer:")
	fmt.Println(answer)
}
