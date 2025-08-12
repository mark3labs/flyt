package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/mark3labs/flyt"
)

// CreateGreetingNode creates a node that generates a greeting message
func CreateGreetingNode(tracer *Tracer) flyt.Node {
	// Create a parent span for the entire node that will be closed when node completes
	var nodeSpan *Span

	return flyt.NewNode(
		flyt.WithPrepFunc(func(ctx context.Context, shared *flyt.SharedStore) (any, error) {
			name, ok := shared.Get("name")
			if !ok {
				name = "World"
			}

			// Create node-level span if not exists
			if nodeSpan == nil {
				nodeSpan = tracer.StartSpan("GreetingNode", map[string]any{
					"node": "GreetingNode",
					"type": "greeting",
				}, map[string]any{
					"initial_name": name,
				})
			}

			// Start prep span as child of node span
			span := tracer.StartChildSpan(nodeSpan, "GreetingNode.prep", map[string]any{
				"phase": "prep",
			}, map[string]any{
				"name_from_store": name,
			})

			// End with output
			defer span.EndWithOutput(map[string]any{
				"prepared_name": name,
			}, nil)

			return name, nil
		}),
		flyt.WithExecFunc(func(ctx context.Context, prepResult any) (any, error) {
			name := prepResult.(string)

			// Start exec span as child of node span
			span := tracer.StartChildSpan(nodeSpan, "GreetingNode.exec", map[string]any{
				"phase": "exec",
			}, map[string]any{
				"input_name": name,
			})

			greeting := fmt.Sprintf("Hello, %s!", name)

			// Simulate some work
			time.Sleep(100 * time.Millisecond)

			// End with output
			defer span.EndWithOutput(map[string]any{
				"greeting": greeting,
			}, nil)

			return greeting, nil
		}),
		flyt.WithPostFunc(func(ctx context.Context, shared *flyt.SharedStore, prepResult, execResult any) (flyt.Action, error) {
			greeting := execResult.(string)

			// Start post span as child of node span
			span := tracer.StartChildSpan(nodeSpan, "GreetingNode.post", map[string]any{
				"phase": "post",
			}, map[string]any{
				"greeting_to_store": greeting,
			})

			shared.Set("greeting", greeting)

			// End with output
			defer span.EndWithOutput(map[string]any{
				"next_action":  "uppercase",
				"stored_value": greeting,
			}, nil)

			// End the node span when post completes
			defer nodeSpan.EndWithOutput(map[string]any{
				"final_greeting": greeting,
				"next_action":    "uppercase",
			}, nil)

			return "uppercase", nil
		}),
	)
}

// CreateUppercaseNode creates a node that converts text to uppercase
func CreateUppercaseNode(tracer *Tracer) flyt.Node {
	// Create a parent span for the entire node
	var nodeSpan *Span

	return flyt.NewNode(
		flyt.WithPrepFunc(func(ctx context.Context, shared *flyt.SharedStore) (any, error) {
			greeting, ok := shared.Get("greeting")
			if !ok {
				return "", fmt.Errorf("greeting not found")
			}

			// Create node-level span if not exists
			if nodeSpan == nil {
				nodeSpan = tracer.StartSpan("UppercaseNode", map[string]any{
					"node": "UppercaseNode",
					"type": "transformation",
				}, map[string]any{
					"input_text": greeting,
				})
			}

			// Start prep span as child of node span
			span := tracer.StartChildSpan(nodeSpan, "UppercaseNode.prep", map[string]any{
				"phase": "prep",
			}, map[string]any{
				"greeting_from_store": greeting,
			})

			// End with output
			defer span.EndWithOutput(map[string]any{
				"prepared_greeting": greeting,
			}, nil)

			return greeting, nil
		}),
		flyt.WithExecFunc(func(ctx context.Context, prepResult any) (any, error) {
			greeting := prepResult.(string)

			// Start exec span as child of node span
			span := tracer.StartChildSpan(nodeSpan, "UppercaseNode.exec", map[string]any{
				"phase": "exec",
			}, map[string]any{
				"input_greeting": greeting,
			})

			uppercase := ""

			// Simulate character-by-character processing
			for _, ch := range greeting {
				uppercase += string(ch)
				time.Sleep(10 * time.Millisecond) // Simulate processing
			}

			// Convert to uppercase
			uppercase = string([]rune(uppercase))
			for i, ch := range greeting {
				if ch >= 'a' && ch <= 'z' {
					uppercase = uppercase[:i] + string(ch-32) + uppercase[i+1:]
				}
			}

			// End with output
			defer span.EndWithOutput(map[string]any{
				"uppercase_result": uppercase,
			}, nil)

			return uppercase, nil
		}),
		flyt.WithPostFunc(func(ctx context.Context, shared *flyt.SharedStore, prepResult, execResult any) (flyt.Action, error) {
			uppercase := execResult.(string)

			// Start post span as child of node span
			span := tracer.StartChildSpan(nodeSpan, "UppercaseNode.post", map[string]any{
				"phase": "post",
			}, map[string]any{
				"uppercase_to_store": uppercase,
			})

			shared.Set("uppercase_greeting", uppercase)

			// End with output
			defer span.EndWithOutput(map[string]any{
				"next_action":  "reverse",
				"stored_value": uppercase,
			}, nil)

			// End the node span when post completes
			defer nodeSpan.EndWithOutput(map[string]any{
				"final_uppercase": uppercase,
				"next_action":     "reverse",
			}, nil)

			return "reverse", nil
		}),
	)
}

// CreateReverseNode creates a node that reverses text
func CreateReverseNode(tracer *Tracer) flyt.Node {
	// Create a parent span for the entire node
	var nodeSpan *Span

	return flyt.NewNode(
		flyt.WithPrepFunc(func(ctx context.Context, shared *flyt.SharedStore) (any, error) {
			uppercase, ok := shared.Get("uppercase_greeting")
			if !ok {
				return "", fmt.Errorf("uppercase_greeting not found")
			}

			// Create node-level span if not exists
			if nodeSpan == nil {
				nodeSpan = tracer.StartSpan("ReverseNode", map[string]any{
					"node": "ReverseNode",
					"type": "transformation",
				}, map[string]any{
					"input_text": uppercase,
				})
			}

			// Start prep span as child of node span
			span := tracer.StartChildSpan(nodeSpan, "ReverseNode.prep", map[string]any{
				"phase": "prep",
			}, map[string]any{
				"uppercase_from_store": uppercase,
			})

			// End with output
			defer span.EndWithOutput(map[string]any{
				"prepared_text": uppercase,
			}, nil)

			return uppercase, nil
		}),
		flyt.WithExecFunc(func(ctx context.Context, prepResult any) (any, error) {
			text := prepResult.(string)

			// Start exec span as child of node span
			span := tracer.StartChildSpan(nodeSpan, "ReverseNode.exec", map[string]any{
				"phase": "exec",
			}, map[string]any{
				"input_text": text,
			})

			runes := []rune(text)

			// Reverse the runes
			for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
				runes[i], runes[j] = runes[j], runes[i]
				time.Sleep(5 * time.Millisecond) // Simulate processing
			}

			reversed := string(runes)

			// End with output
			defer span.EndWithOutput(map[string]any{
				"reversed_result": reversed,
			}, nil)

			return reversed, nil
		}),
		flyt.WithPostFunc(func(ctx context.Context, shared *flyt.SharedStore, prepResult, execResult any) (flyt.Action, error) {
			reversed := execResult.(string)

			// Start post span as child of node span
			span := tracer.StartChildSpan(nodeSpan, "ReverseNode.post", map[string]any{
				"phase": "post",
			}, map[string]any{
				"reversed_to_store": reversed,
			})

			shared.Set("reversed_greeting", reversed)

			// End with output
			defer span.EndWithOutput(map[string]any{
				"next_action":  "summarize",
				"stored_value": reversed,
			}, nil)

			// End the node span when post completes
			defer nodeSpan.EndWithOutput(map[string]any{
				"final_reversed": reversed,
				"next_action":    "summarize",
			}, nil)

			return "summarize", nil
		}),
	)
}

// CreateSummarizeNode creates a node that summarizes text using an LLM
func CreateSummarizeNode(tracer *Tracer) flyt.Node {
	// Create a parent span for the entire node
	var nodeSpan *Span

	return flyt.NewNode(
		flyt.WithPrepFunc(func(ctx context.Context, shared *flyt.SharedStore) (any, error) {
			reversed, ok := shared.Get("reversed_greeting")
			if !ok {
				return "", fmt.Errorf("reversed_greeting not found")
			}

			// Create node-level span if not exists
			if nodeSpan == nil {
				nodeSpan = tracer.StartSpan("SummarizeNode", map[string]any{
					"node": "SummarizeNode",
					"type": "llm",
				}, map[string]any{
					"input_text": reversed,
				})
			}

			// Start prep span as child of node span
			span := tracer.StartChildSpan(nodeSpan, "SummarizeNode.prep", map[string]any{
				"phase": "prep",
			}, map[string]any{
				"text_from_store": reversed,
			})

			// End with output
			defer span.EndWithOutput(map[string]any{
				"prepared_text": reversed,
			}, nil)

			return reversed, nil
		}),
		flyt.WithExecFunc(func(ctx context.Context, prepResult any) (any, error) {
			text := prepResult.(string)

			// Start exec span as child of node span
			span := tracer.StartChildSpan(nodeSpan, "SummarizeNode.exec", map[string]any{
				"phase": "exec",
			}, map[string]any{
				"input_text": text,
			})
			defer func() {
				span.EndWithOutput(map[string]any{
					"llm_called": true,
				}, nil)
			}()

			// Get API key from environment
			apiKey := os.Getenv("OPENAI_API_KEY")
			if apiKey == "" {
				// If no API key, return a mock response
				mockResponse := fmt.Sprintf("Summary of '%s': This is a mock LLM response. Set OPENAI_API_KEY to use real API.", text)
				return mockResponse, nil
			}

			// Prepare messages for the LLM
			messages := []map[string]string{
				{
					"role":    "system",
					"content": "You are a helpful assistant that creates playful, one-sentence summaries.",
				},
				{
					"role":    "user",
					"content": fmt.Sprintf("Create a fun, creative one-sentence summary of this text: '%s'", text),
				},
			}

			// Start generation tracking as child of exec span
			generation := tracer.StartGeneration(span, "gpt-4.1-summary", "gpt-4.1", messages, map[string]any{
				"temperature": 0.7,
				"max_tokens":  100,
			})

			// Call the LLM
			response, usage, err := CallLLM(apiKey, messages)

			// End generation with response
			generation.EndWithResponse(response, usage, err)

			if err != nil {
				return "", fmt.Errorf("LLM call failed: %v", err)
			}

			return response, nil
		}),
		flyt.WithPostFunc(func(ctx context.Context, shared *flyt.SharedStore, prepResult, execResult any) (flyt.Action, error) {
			summary := execResult.(string)

			// Start post span as child of node span
			span := tracer.StartChildSpan(nodeSpan, "SummarizeNode.post", map[string]any{
				"phase": "post",
			}, map[string]any{
				"summary_to_store": summary,
			})

			shared.Set("llm_summary", summary)

			// End with output
			defer span.EndWithOutput(map[string]any{
				"final_summary": summary,
				"action":        "default",
			}, nil)

			// End the node span when post completes
			defer nodeSpan.EndWithOutput(map[string]any{
				"final_summary": summary,
				"action":        "default",
			}, nil)

			return flyt.DefaultAction, nil
		}),
	)
}

// CallLLM calls the OpenAI API with messages
func CallLLM(apiKey string, messages []map[string]string) (string, map[string]int, error) {
	// Create request body
	reqBody := map[string]any{
		"model":       "gpt-4", // Using gpt-4 as gpt-4.1 doesn't exist
		"messages":    messages,
		"temperature": 0.7,
		"max_tokens":  100,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", nil, err
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// Make request
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", nil, err
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return "", nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var result map[string]any
	if err := json.Unmarshal(body, &result); err != nil {
		return "", nil, err
	}

	// Extract the response
	choices, ok := result["choices"].([]any)
	if !ok || len(choices) == 0 {
		return "", nil, fmt.Errorf("no response from API")
	}

	choice := choices[0].(map[string]any)
	message := choice["message"].(map[string]any)
	content := message["content"].(string)

	// Extract usage statistics if available
	usage := make(map[string]int)
	if usageData, ok := result["usage"].(map[string]any); ok {
		if promptTokens, ok := usageData["prompt_tokens"].(float64); ok {
			usage["prompt_tokens"] = int(promptTokens)
		}
		if completionTokens, ok := usageData["completion_tokens"].(float64); ok {
			usage["completion_tokens"] = int(completionTokens)
		}
		if totalTokens, ok := usageData["total_tokens"].(float64); ok {
			usage["total_tokens"] = int(totalTokens)
		}
	}

	return content, usage, nil
}

// CreateTracedFlow creates a flow with tracing enabled
func CreateTracedFlow(tracer *Tracer) *flyt.Flow {
	// Create nodes with tracer injected
	greetingNode := CreateGreetingNode(tracer)
	uppercaseNode := CreateUppercaseNode(tracer)
	reverseNode := CreateReverseNode(tracer)
	summarizeNode := CreateSummarizeNode(tracer)

	// Create flow
	flow := flyt.NewFlow(greetingNode)
	flow.Connect(greetingNode, "uppercase", uppercaseNode)
	flow.Connect(uppercaseNode, "reverse", reverseNode)
	flow.Connect(reverseNode, "summarize", summarizeNode)

	return flow
}

// RunTracedFlow runs a flow with tracing
func RunTracedFlow(flow *flyt.Flow, tracer *Tracer, shared *flyt.SharedStore) error {
	// Start main trace
	trace := tracer.StartTrace("BasicGreetingFlow", map[string]any{
		"flow_type": "greeting_pipeline",
		"version":   "1.0.0",
	})
	defer func() {
		// End trace and flush
		trace.End(nil)
		tracer.Flush(context.Background())
	}()

	// Add input metadata
	trace.AddMetadata(map[string]any{
		"input": shared,
	})

	// Run the flow
	ctx := context.Background()
	err := flow.Run(ctx, shared)

	if err != nil {
		trace.End(err)
		return err
	}

	// Add output metadata
	trace.AddMetadata(map[string]any{
		"output": shared,
	})

	return nil
}

func main() {
	fmt.Println("🚀 Starting Flyt Tracing Example")
	fmt.Println("=" + string(make([]byte, 49)))

	// Check for Langfuse configuration
	host := os.Getenv("LANGFUSE_HOST")
	publicKey := os.Getenv("LANGFUSE_PUBLIC_KEY")
	secretKey := os.Getenv("LANGFUSE_SECRET_KEY")

	if host == "" || publicKey == "" || secretKey == "" {
		log.Println("⚠️  Warning: Langfuse environment variables not set")
		log.Println("   Set LANGFUSE_HOST, LANGFUSE_PUBLIC_KEY, and LANGFUSE_SECRET_KEY")
		log.Println("   Running in demo mode without actual tracing")
	}

	// Check for OpenAI API key
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Println("ℹ️  Note: OPENAI_API_KEY not set - LLM node will use mock responses")
	}

	// Create tracer
	tracer := NewTracer()

	// Create the traced flow
	flow := CreateTracedFlow(tracer)

	// Prepare shared data
	shared := flyt.NewSharedStore()
	shared.Set("name", "Flyt User")

	fmt.Printf("📥 Input: name=%v\n", "Flyt User")

	// Run the flow with tracing
	err := RunTracedFlow(flow, tracer, shared)
	if err != nil {
		fmt.Printf("❌ Flow failed with error: %v\n", err)
		log.Fatal(err)
	}

	// Print results
	greeting, _ := shared.Get("greeting")
	uppercase, _ := shared.Get("uppercase_greeting")
	reversed, _ := shared.Get("reversed_greeting")
	summary, _ := shared.Get("llm_summary")

	fmt.Printf("📤 Output:\n")
	fmt.Printf("   Greeting: %v\n", greeting)
	fmt.Printf("   Uppercase: %v\n", uppercase)
	fmt.Printf("   Reversed: %v\n", reversed)
	fmt.Printf("   Summary: %v\n", summary)
	fmt.Println("✅ Flow completed successfully!")

	if host != "" {
		fmt.Printf("\n📊 Check your Langfuse dashboard to see the trace!\n")
		fmt.Printf("   Dashboard URL: %s\n", host)
	}
}
