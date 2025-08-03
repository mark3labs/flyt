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
			// Start prep span
			span := tracer.StartSpan("GreetingNode.prep", map[string]any{
				"node":  "GreetingNode",
				"phase": "prep",
			})
			defer span.End(nil)

			name, ok := shared.Get("name")
			if !ok {
				name = "World"
			}

			span.AddMetadata(map[string]any{
				"input_name": name,
			})

			return name, nil
		}),
		flyt.WithExecFunc(func(ctx context.Context, prepResult any) (any, error) {
			// Start exec span
			span := tracer.StartSpan("GreetingNode.exec", map[string]any{
				"node":  "GreetingNode",
				"phase": "exec",
			})
			defer span.End(nil)

			name := prepResult.(string)
			greeting := fmt.Sprintf("Hello, %s!", name)

			span.AddMetadata(map[string]any{
				"greeting": greeting,
			})

			// Simulate some work
			time.Sleep(100 * time.Millisecond)

			return greeting, nil
		}),
		flyt.WithPostFunc(func(ctx context.Context, shared *flyt.SharedStore, prepResult, execResult any) (flyt.Action, error) {
			// Start post span
			span := tracer.StartSpan("GreetingNode.post", map[string]any{
				"node":  "GreetingNode",
				"phase": "post",
			})
			defer span.End(nil)

			greeting := execResult.(string)
			shared.Set("greeting", greeting)

			span.AddMetadata(map[string]any{
				"stored_greeting": greeting,
			})

			return "uppercase", nil
		}),
	)
}

// CreateUppercaseNode creates a node that converts text to uppercase
func CreateUppercaseNode(tracer *Tracer) flyt.Node {
	return flyt.NewNode(
		flyt.WithPrepFunc(func(ctx context.Context, shared *flyt.SharedStore) (any, error) {
			// Start prep span
			span := tracer.StartSpan("UppercaseNode.prep", map[string]any{
				"node":  "UppercaseNode",
				"phase": "prep",
			})
			defer span.End(nil)

			greeting, ok := shared.Get("greeting")
			if !ok {
				return "", fmt.Errorf("greeting not found")
			}

			span.AddMetadata(map[string]any{
				"input_greeting": greeting,
			})

			return greeting, nil
		}),
		flyt.WithExecFunc(func(ctx context.Context, prepResult any) (any, error) {
			// Start exec span
			span := tracer.StartSpan("UppercaseNode.exec", map[string]any{
				"node":  "UppercaseNode",
				"phase": "exec",
			})
			defer span.End(nil)

			greeting := prepResult.(string)
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

			span.AddMetadata(map[string]any{
				"uppercase_result": uppercase,
			})

			return uppercase, nil
		}),
		flyt.WithPostFunc(func(ctx context.Context, shared *flyt.SharedStore, prepResult, execResult any) (flyt.Action, error) {
			// Start post span
			span := tracer.StartSpan("UppercaseNode.post", map[string]any{
				"node":  "UppercaseNode",
				"phase": "post",
			})
			defer span.End(nil)

			uppercase := execResult.(string)
			shared.Set("uppercase_greeting", uppercase)

			span.AddMetadata(map[string]any{
				"final_result": uppercase,
			})

			return "reverse", nil
		}),
	)
}

// CreateReverseNode creates a node that reverses text
func CreateReverseNode(tracer *Tracer) flyt.Node {
	return flyt.NewNode(
		flyt.WithPrepFunc(func(ctx context.Context, shared *flyt.SharedStore) (any, error) {
			// Start prep span
			span := tracer.StartSpan("ReverseNode.prep", map[string]any{
				"node":  "ReverseNode",
				"phase": "prep",
			})
			defer span.End(nil)

			uppercase, ok := shared.Get("uppercase_greeting")
			if !ok {
				return "", fmt.Errorf("uppercase_greeting not found")
			}

			span.AddMetadata(map[string]any{
				"input_text": uppercase,
			})

			return uppercase, nil
		}),
		flyt.WithExecFunc(func(ctx context.Context, prepResult any) (any, error) {
			// Start exec span
			span := tracer.StartSpan("ReverseNode.exec", map[string]any{
				"node":  "ReverseNode",
				"phase": "exec",
			})
			defer span.End(nil)

			text := prepResult.(string)
			runes := []rune(text)

			// Reverse the runes
			for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
				runes[i], runes[j] = runes[j], runes[i]
				time.Sleep(5 * time.Millisecond) // Simulate processing
			}

			reversed := string(runes)

			span.AddMetadata(map[string]any{
				"reversed_result": reversed,
			})

			return reversed, nil
		}),
		flyt.WithPostFunc(func(ctx context.Context, shared *flyt.SharedStore, prepResult, execResult any) (flyt.Action, error) {
			// Start post span
			span := tracer.StartSpan("ReverseNode.post", map[string]any{
				"node":  "ReverseNode",
				"phase": "post",
			})
			defer span.End(nil)

			reversed := execResult.(string)
			shared.Set("reversed_greeting", reversed)

			span.AddMetadata(map[string]any{
				"final_reversed": reversed,
			})

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
