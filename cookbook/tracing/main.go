package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/mark3labs/flyt"
)

// CreateGreetingNode creates a node that generates a greeting message
func CreateGreetingNode(tracer *Tracer) flyt.Node {
	return flyt.NewNode(
		flyt.WithPrepFunc(func(ctx context.Context, shared *flyt.SharedStore) (any, error) {
			name, ok := shared.Get("name")
			if !ok {
				name = "World"
			}

			// Start prep span with input
			span := tracer.StartSpan("GreetingNode.prep", map[string]any{
				"node":  "GreetingNode",
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

			// Start exec span with input
			span := tracer.StartSpan("GreetingNode.exec", map[string]any{
				"node":  "GreetingNode",
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

			// Start post span with input
			span := tracer.StartSpan("GreetingNode.post", map[string]any{
				"node":  "GreetingNode",
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

			return "uppercase", nil
		}),
	)
}

// CreateUppercaseNode creates a node that converts text to uppercase
func CreateUppercaseNode(tracer *Tracer) flyt.Node {
	return flyt.NewNode(
		flyt.WithPrepFunc(func(ctx context.Context, shared *flyt.SharedStore) (any, error) {
			greeting, ok := shared.Get("greeting")
			if !ok {
				return "", fmt.Errorf("greeting not found")
			}

			// Start prep span with input
			span := tracer.StartSpan("UppercaseNode.prep", map[string]any{
				"node":  "UppercaseNode",
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

			// Start exec span with input
			span := tracer.StartSpan("UppercaseNode.exec", map[string]any{
				"node":  "UppercaseNode",
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

			// Start post span with input
			span := tracer.StartSpan("UppercaseNode.post", map[string]any{
				"node":  "UppercaseNode",
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

			return "reverse", nil
		}),
	)
}

// CreateReverseNode creates a node that reverses text
func CreateReverseNode(tracer *Tracer) flyt.Node {
	return flyt.NewNode(
		flyt.WithPrepFunc(func(ctx context.Context, shared *flyt.SharedStore) (any, error) {
			uppercase, ok := shared.Get("uppercase_greeting")
			if !ok {
				return "", fmt.Errorf("uppercase_greeting not found")
			}

			// Start prep span with input
			span := tracer.StartSpan("ReverseNode.prep", map[string]any{
				"node":  "ReverseNode",
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

			// Start exec span with input
			span := tracer.StartSpan("ReverseNode.exec", map[string]any{
				"node":  "ReverseNode",
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

			// Start post span with input
			span := tracer.StartSpan("ReverseNode.post", map[string]any{
				"node":  "ReverseNode",
				"phase": "post",
			}, map[string]any{
				"reversed_to_store": reversed,
			})

			shared.Set("reversed_greeting", reversed)

			// End with output
			defer span.EndWithOutput(map[string]any{
				"final_result": reversed,
				"action":       "default",
			}, nil)

			return flyt.DefaultAction, nil
		}),
	)
}

// CreateTracedFlow creates a flow with tracing enabled
func CreateTracedFlow(tracer *Tracer) *flyt.Flow {
	// Create nodes with tracer injected
	greetingNode := CreateGreetingNode(tracer)
	uppercaseNode := CreateUppercaseNode(tracer)
	reverseNode := CreateReverseNode(tracer)

	// Create flow
	flow := flyt.NewFlow(greetingNode)
	flow.Connect(greetingNode, "uppercase", uppercaseNode)
	flow.Connect(uppercaseNode, "reverse", reverseNode)

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

	fmt.Printf("📤 Output:\n")
	fmt.Printf("   Greeting: %v\n", greeting)
	fmt.Printf("   Uppercase: %v\n", uppercase)
	fmt.Printf("   Reversed: %v\n", reversed)
	fmt.Println("✅ Flow completed successfully!")

	if host != "" {
		fmt.Printf("\n📊 Check your Langfuse dashboard to see the trace!\n")
		fmt.Printf("   Dashboard URL: %s\n", host)
	}
}
