package flyt_test

import (
	"context"
	"fmt"
	"log"

	"github.com/mark3labs/flyt"
)

// ExampleNewNode demonstrates creating a simple node using the NewNode helper
func ExampleNewNode() {
	// Create a node with custom exec function
	node := flyt.NewNode(
		flyt.WithExecFunc(func(ctx context.Context, prepResult any) (any, error) {
			fmt.Println("Hello from custom node!")
			return "success", nil
		}),
	)

	// Run the node
	shared := flyt.NewSharedStore()
	ctx := context.Background()

	action, err := flyt.Run(ctx, node, shared)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Action: %s\n", action)
	// Output:
	// Hello from custom node!
	// Action: default
}

// ExampleNewNode_withAllFunctions demonstrates using all three functions
func ExampleNewNode_withAllFunctions() {
	// Create a node that processes data through all phases
	node := flyt.NewNode(
		flyt.WithPrepFunc(func(ctx context.Context, shared *flyt.SharedStore) (any, error) {
			// Read input from shared store
			input, _ := shared.Get("input")
			fmt.Printf("Prep: reading input '%v'\n", input)
			return input, nil
		}),
		flyt.WithExecFunc(func(ctx context.Context, prepResult any) (any, error) {
			// Process the input
			result := fmt.Sprintf("processed: %v", prepResult)
			fmt.Printf("Exec: %s\n", result)
			return result, nil
		}),
		flyt.WithPostFunc(func(ctx context.Context, shared *flyt.SharedStore, prepResult, execResult any) (flyt.Action, error) {
			// Store the result
			shared.Set("output", execResult)
			fmt.Printf("Post: stored result\n")
			return flyt.DefaultAction, nil
		}),
	)

	// Setup and run
	shared := flyt.NewSharedStore()
	shared.Set("input", "hello world")
	ctx := context.Background()

	_, err := flyt.Run(ctx, node, shared)
	if err != nil {
		log.Fatal(err)
	}

	output, _ := shared.Get("output")
	fmt.Printf("Final output: %v\n", output)
	// Output:
	// Prep: reading input 'hello world'
	// Exec: processed: hello world
	// Post: stored result
	// Final output: processed: hello world
}

// ExampleNewNode_withRetries demonstrates combining with BaseNode options
func ExampleNewNode_withRetries() {
	attempts := 0
	node := flyt.NewNode(
		flyt.WithExecFunc(func(ctx context.Context, prepResult any) (any, error) {
			attempts++
			fmt.Printf("Attempt %d\n", attempts)
			if attempts < 3 {
				return nil, fmt.Errorf("temporary failure")
			}
			return "success after retries", nil
		}),
		flyt.WithMaxRetries(3),
	)

	shared := flyt.NewSharedStore()
	ctx := context.Background()

	_, err := flyt.Run(ctx, node, shared)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Success!")
	// Output:
	// Attempt 1
	// Attempt 2
	// Attempt 3
	// Success!
}
