package main

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/mark3labs/flyt"
)

// Example nodes for the flow that will process each item
// These nodes need to pass the shared store from Prep to Exec
type TransformNode struct {
	*flyt.BaseNode
	multiplier int
}

func (n *TransformNode) Prep(ctx context.Context, shared *flyt.SharedStore) (any, error) {
	// Pass shared store to Exec
	return shared, nil
}

func (n *TransformNode) Exec(ctx context.Context, prepResult any) (any, error) {
	shared := prepResult.(*flyt.SharedStore)
	input, _ := shared.Get("input")
	result := input.(int) * n.multiplier
	shared.Set("transformed", result)
	return nil, nil
}

type ValidateNode struct {
	*flyt.BaseNode
	maxValue int
}

func (n *ValidateNode) Prep(ctx context.Context, shared *flyt.SharedStore) (any, error) {
	// Pass shared store to Exec
	return shared, nil
}

func (n *ValidateNode) Exec(ctx context.Context, prepResult any) (any, error) {
	shared := prepResult.(*flyt.SharedStore)
	value, _ := shared.Get("transformed")
	if value.(int) > n.maxValue {
		// For invalid items, we won't set output so they get filtered
		return nil, nil
	}
	return nil, nil
}

func (n *ValidateNode) Post(ctx context.Context, shared *flyt.SharedStore, prepResult, execResult any) (flyt.Action, error) {
	value, _ := shared.Get("transformed")
	if value.(int) > n.maxValue {
		return "invalid", nil
	}
	return flyt.DefaultAction, nil
}

type SaveNode struct {
	*flyt.BaseNode
	mu      sync.Mutex
	results []int
}

func (n *SaveNode) Prep(ctx context.Context, shared *flyt.SharedStore) (any, error) {
	// Pass shared store to Exec
	return shared, nil
}

func (n *SaveNode) Exec(ctx context.Context, prepResult any) (any, error) {
	shared := prepResult.(*flyt.SharedStore)
	value, _ := shared.Get("transformed")
	n.mu.Lock()
	n.results = append(n.results, value.(int))
	n.mu.Unlock()
	shared.Set("output", value.(int))
	return nil, nil
}

func main() {
	// Define the flow factory that creates a processing pipeline
	flowFactory := func() *flyt.Flow {
		// Create new node instances for each flow
		transform := &TransformNode{
			BaseNode:   flyt.NewBaseNode(),
			multiplier: 2,
		}
		validate := &ValidateNode{
			BaseNode: flyt.NewBaseNode(),
			maxValue: 100,
		}
		save := &SaveNode{
			BaseNode: flyt.NewBaseNode(),
			results:  []int{},
		}

		// Build the flow
		flow := flyt.NewFlow(transform)
		flow.Connect(transform, flyt.DefaultAction, validate)
		flow.Connect(validate, flyt.DefaultAction, save)
		// Note: There's no explicit End node in the current API,
		// invalid items just won't reach save node

		return flow
	}

	// Create a batch node that processes items through the flow factory
	batchFlowNode := flyt.NewBatchNode().
		WithPrepFunc(func(ctx context.Context, shared *flyt.SharedStore) ([]flyt.Result, error) {
			// Extract items to process
			items, _ := shared.Get("items")
			itemSlice := items.([]int)
			results := make([]flyt.Result, len(itemSlice))
			for i, item := range itemSlice {
				results[i] = flyt.NewResult(item)
			}
			return results, nil
		}).
		WithExecFunc(func(ctx context.Context, item flyt.Result) (flyt.Result, error) {
			// Create a new flow instance for this item
			flow := flowFactory()

			// Create an isolated shared store for this flow
			flowShared := flyt.NewSharedStore()
			flowShared.Set("input", item.Value())

			// Run the flow
			_, err := flyt.Run(ctx, flow, flowShared)
			if err != nil {
				return flyt.NewErrorResult(err), nil
			}

			// Extract the output from the flow's shared store
			output, ok := flowShared.Get("output")
			if !ok {
				// Item was filtered out (e.g., invalid)
				return flyt.NewResult(nil), nil
			}
			return flyt.NewResult(output), nil
		}).
		WithPostFunc(func(ctx context.Context, shared *flyt.SharedStore, prep []flyt.Result, exec []flyt.Result) (flyt.Action, error) {
			// Aggregate results
			var outputs []int
			var errors []error

			for _, r := range exec {
				if r.IsError() {
					errors = append(errors, r.Error())
				} else if r.Value() != nil {
					outputs = append(outputs, r.Value().(int))
				}
			}

			shared.Set("outputs", outputs)
			shared.Set("errors", errors)
			return flyt.DefaultAction, nil
		}).
		WithBatchConcurrency(3)

	// Example 1: Using the BatchFlowNode
	fmt.Println("Example 1: BatchNode with FlowFactory")
	ctx := context.Background()
	shared1 := flyt.NewSharedStore()
	shared1.Set("items", []int{10, 20, 30, 40, 50, 60})

	_, err := flyt.Run(ctx, batchFlowNode, shared1)
	if err != nil {
		log.Fatal(err)
	}

	outputs, _ := shared1.Get("outputs")
	errors, _ := shared1.Get("errors")
	fmt.Printf("Outputs: %v\n", outputs)
	fmt.Printf("Errors: %v\n", errors)

	// Example 2: Inline lambda batch flow
	fmt.Println("\nExample 2: Inline BatchNode with FlowFactory")
	simpleBatchFlow := flyt.NewBatchNode().
		WithPrepFunc(func(ctx context.Context, shared *flyt.SharedStore) ([]flyt.Result, error) {
			// Generate items
			count, _ := shared.Get("count")
			n := count.(int)
			results := make([]flyt.Result, n)
			for i := 0; i < n; i++ {
				results[i] = flyt.NewResult(i + 1)
			}
			return results, nil
		}).
		WithExecFunc(func(ctx context.Context, item flyt.Result) (flyt.Result, error) {
			// Process through flow
			flow := flowFactory()
			flowShared := flyt.NewSharedStore()
			flowShared.Set("input", item.Value())

			_, err := flyt.Run(ctx, flow, flowShared)
			if err != nil {
				return flyt.NewErrorResult(err), nil
			}

			output, _ := flowShared.Get("output")
			return flyt.NewResult(output), nil
		}).
		WithPostFunc(func(ctx context.Context, shared *flyt.SharedStore, prep []flyt.Result, exec []flyt.Result) (flyt.Action, error) {
			var values []int
			for _, r := range exec {
				if !r.IsError() && r.Value() != nil {
					values = append(values, r.Value().(int))
				}
			}
			shared.Set("processed", values)
			return flyt.DefaultAction, nil
		}).
		WithBatchConcurrency(5)

	shared2 := flyt.NewSharedStore()
	shared2.Set("count", 10)

	_, err = flyt.Run(ctx, simpleBatchFlow, shared2)
	if err != nil {
		log.Fatal(err)
	}

	processed, _ := shared2.Get("processed")
	fmt.Printf("Processed: %v\n", processed)

	// Example 3: Advanced - BatchNode with conditional flow selection
	fmt.Println("\nExample 3: Conditional Flow Selection")

	// Helper function to create a priority processing flow
	createPriorityFlow := func() *flyt.Flow {
		// Priority flow processes with multiplier of 3 instead of 2
		transform := &TransformNode{
			BaseNode:   flyt.NewBaseNode(),
			multiplier: 3, // Higher multiplier for priority items
		}
		validate := &ValidateNode{
			BaseNode: flyt.NewBaseNode(),
			maxValue: 150, // Higher threshold for priority
		}
		save := &SaveNode{
			BaseNode: flyt.NewBaseNode(),
			results:  []int{},
		}

		flow := flyt.NewFlow(transform)
		flow.Connect(transform, flyt.DefaultAction, validate)
		flow.Connect(validate, flyt.DefaultAction, save)

		return flow
	}

	advancedBatchFlow := flyt.NewBatchNode().
		WithPrepFunc(func(ctx context.Context, shared *flyt.SharedStore) ([]flyt.Result, error) {
			data, _ := shared.Get("data")
			items := data.([]map[string]interface{})
			results := make([]flyt.Result, len(items))
			for i, item := range items {
				results[i] = flyt.NewResult(item)
			}
			return results, nil
		}).
		WithExecFunc(func(ctx context.Context, item flyt.Result) (flyt.Result, error) {
			data := item.Value().(map[string]interface{})
			itemType := data["type"].(string)

			// Select flow factory based on item type
			var flow *flyt.Flow
			if itemType == "priority" {
				flow = createPriorityFlow()
			} else {
				flow = flowFactory()
			}

			flowShared := flyt.NewSharedStore()
			flowShared.Set("input", data["value"])

			_, err := flyt.Run(ctx, flow, flowShared)
			if err != nil {
				return flyt.NewErrorResult(err), nil
			}

			output, ok := flowShared.Get("output")
			if !ok {
				return flyt.NewResult(nil), nil
			}

			return flyt.NewResult(map[string]interface{}{
				"type":   itemType,
				"result": output,
			}), nil
		}).
		WithPostFunc(func(ctx context.Context, shared *flyt.SharedStore, prep []flyt.Result, exec []flyt.Result) (flyt.Action, error) {
			var priority []int
			var standard []int

			for _, r := range exec {
				if !r.IsError() && r.Value() != nil {
					data := r.Value().(map[string]interface{})
					if data["type"] == "priority" {
						priority = append(priority, data["result"].(int))
					} else {
						standard = append(standard, data["result"].(int))
					}
				}
			}

			shared.Set("priority_results", priority)
			shared.Set("standard_results", standard)
			return flyt.DefaultAction, nil
		}).
		WithBatchConcurrency(4)

	shared3 := flyt.NewSharedStore()
	shared3.Set("data", []map[string]interface{}{
		{"type": "standard", "value": 5},
		{"type": "priority", "value": 10},
		{"type": "standard", "value": 15},
		{"type": "priority", "value": 20},
		{"type": "standard", "value": 55}, // Will be filtered as invalid
		{"type": "priority", "value": 60}, // Will pass with priority threshold
	})

	_, err = flyt.Run(ctx, advancedBatchFlow, shared3)
	if err != nil {
		log.Fatal(err)
	}

	priorityResults, _ := shared3.Get("priority_results")
	standardResults, _ := shared3.Get("standard_results")
	fmt.Printf("Priority results: %v\n", priorityResults)
	fmt.Printf("Standard results: %v\n", standardResults)

	// Example 4: Reusable BatchFlowNode wrapper
	fmt.Println("\nExample 4: Reusable BatchFlowNode Pattern")

	// Helper to create a reusable batch flow node
	createBatchFlowNode := func(name string, factory func() *flyt.Flow, concurrency int) flyt.Node {
		return flyt.NewBatchNode().
			WithPrepFunc(func(ctx context.Context, shared *flyt.SharedStore) ([]flyt.Result, error) {
				items, _ := shared.Get("items")
				// Handle different input types
				var results []flyt.Result
				switch v := items.(type) {
				case []int:
					results = make([]flyt.Result, len(v))
					for i, item := range v {
						results[i] = flyt.NewResult(item)
					}
				case []interface{}:
					results = make([]flyt.Result, len(v))
					for i, item := range v {
						results[i] = flyt.NewResult(item)
					}
				default:
					return nil, fmt.Errorf("unsupported items type: %T", items)
				}
				return results, nil
			}).
			WithExecFunc(func(ctx context.Context, item flyt.Result) (flyt.Result, error) {
				flow := factory()
				flowShared := flyt.NewSharedStore()
				flowShared.Set("input", item.Value())

				_, err := flyt.Run(ctx, flow, flowShared)
				if err != nil {
					return flyt.NewErrorResult(err), nil
				}

				output, _ := flowShared.Get("output")
				return flyt.NewResult(output), nil
			}).
			WithPostFunc(func(ctx context.Context, shared *flyt.SharedStore, prep []flyt.Result, exec []flyt.Result) (flyt.Action, error) {
				var results []interface{}
				for _, r := range exec {
					if !r.IsError() && r.Value() != nil {
						results = append(results, r.Value())
					}
				}
				shared.Set("results", results)
				return flyt.DefaultAction, nil
			}).
			WithBatchConcurrency(concurrency)
	}

	// Use the reusable pattern
	reusableNode := createBatchFlowNode("processor", flowFactory, 2)

	shared4 := flyt.NewSharedStore()
	shared4.Set("items", []int{1, 2, 3, 4, 5})

	_, err = flyt.Run(ctx, reusableNode, shared4)
	if err != nil {
		log.Fatal(err)
	}

	results, _ := shared4.Get("results")
	fmt.Printf("Reusable pattern results: %v\n", results)
}
