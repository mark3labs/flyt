package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/mark3labs/flyt"
)

// CreateSummarizeNode creates a node that summarizes text using an LLM
// This version uses Result types for better type safety
func CreateSummarizeNode(apiKey string) flyt.Node {
	return flyt.NewNode(
		flyt.WithPrepFunc(func(ctx context.Context, shared *flyt.SharedStore) (flyt.Result, error) {
			// Read text from shared store with type-safe access
			text := shared.GetStringOr("text", "")
			if text == "" {
				return flyt.Result{}, fmt.Errorf("no text found in shared store")
			}
			return flyt.NewResult(text), nil
		}),
		flyt.WithExecFunc(func(ctx context.Context, prepResult flyt.Result) (flyt.Result, error) {
			// Type-safe string extraction
			text := prepResult.AsStringOr("")
			if text == "" {
				return flyt.NewResult("Empty text"), nil
			}

			// Create prompt for summarization
			prompt := fmt.Sprintf("Summarize this text in 10 words or less: %s", text)

			// Call LLM
			summary, err := CallLLM(apiKey, prompt)
			if err != nil {
				// Simulate retry behavior - the framework will retry based on node options
				return flyt.Result{}, fmt.Errorf("LLM call failed: %w", err)
			}

			return flyt.NewResult(summary), nil
		}),
		flyt.WithPostFunc(func(ctx context.Context, shared *flyt.SharedStore, prepResult, execResult flyt.Result) (flyt.Action, error) {
			// Store the summary in shared store with type-safe access
			summary := execResult.AsStringOr("Failed to generate summary")
			shared.Set("summary", summary)

			// Log the result
			fmt.Printf("✅ Summary generated: %s\n", summary)

			return flyt.DefaultAction, nil
		}),
		// Configure retries and wait time
		flyt.WithMaxRetries(3),
		flyt.WithWait(time.Second),
	)
}

// CreateSummarizeNodeWithFallback demonstrates custom error handling
// This version uses the Any variants for simpler code when type safety isn't critical
func CreateSummarizeNodeWithFallback(apiKey string) flyt.Node {
	attempts := 0

	return flyt.NewNode(
		flyt.WithPrepFuncAny(func(ctx context.Context, shared *flyt.SharedStore) (any, error) {
			text, ok := shared.Get("text")
			if !ok {
				return "", fmt.Errorf("no text found in shared store")
			}
			return text, nil
		}),
		flyt.WithExecFuncAny(func(ctx context.Context, prepResult any) (any, error) {
			attempts++
			text := prepResult.(string)

			// Simulate occasional failures for demonstration
			if attempts < 2 && os.Getenv("SIMULATE_FAILURE") == "true" {
				return "", fmt.Errorf("simulated LLM failure (attempt %d)", attempts)
			}

			if text == "" {
				return "Empty text", nil
			}

			prompt := fmt.Sprintf("Summarize this text in 10 words or less: %s", text)
			summary, err := CallLLM(apiKey, prompt)
			if err != nil {
				return "", err
			}

			return summary, nil
		}),
		flyt.WithPostFuncAny(func(ctx context.Context, shared *flyt.SharedStore, prepResult, execResult any) (flyt.Action, error) {
			summary := execResult.(string)

			// If we got here after retries, note that
			if attempts > 1 {
				fmt.Printf("⚡ Succeeded after %d attempts\n", attempts)
			}

			shared.Set("summary", summary)
			shared.Set("attempts", attempts)

			fmt.Printf("✅ Summary: %s\n", summary)

			return flyt.DefaultAction, nil
		}),
		flyt.WithMaxRetries(3),
		flyt.WithWait(500*time.Millisecond),
	)
}

func main() {
	// Example text to summarize
	text := `PocketFlow is a minimalist LLM framework that models workflows as a Nested Directed Graph.
Nodes handle simple LLM tasks, connecting through Actions for Agents.
Flows orchestrate these nodes for Task Decomposition, and can be nested.
It also supports Batch processing and Async execution.`

	// Get API key
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("Please set OPENAI_API_KEY environment variable")
	}

	// Initialize shared store
	shared := flyt.NewSharedStore()
	shared.Set("text", text)

	// Create the node
	var summarizeNode flyt.Node
	if os.Getenv("USE_FALLBACK_DEMO") == "true" {
		fmt.Println("🔄 Using node with retry demonstration")
		summarizeNode = CreateSummarizeNodeWithFallback(apiKey)
	} else {
		fmt.Println("📝 Using standard summarize node")
		summarizeNode = CreateSummarizeNode(apiKey)
	}

	// Create and run the flow
	flow := flyt.NewFlow(summarizeNode)
	ctx := context.Background()

	fmt.Println("\n📄 Input text:")
	fmt.Println(text)
	fmt.Println("\n🤖 Processing...")

	// Run the flow
	if err := flow.Run(ctx, shared); err != nil {
		// If all retries failed, provide a fallback
		fmt.Printf("❌ Error: %v\n", err)
		shared.Set("summary", "There was an error processing your request.")
	}

	// Print results
	summary, _ := shared.Get("summary")
	fmt.Printf("\n📋 Final summary: %s\n", summary)

	if attempts, ok := shared.Get("attempts"); ok {
		fmt.Printf("📊 Total attempts: %d\n", attempts)
	}
}
