// Package flyt provides batch processing capabilities for the flyt workflow framework.
//
// The batch package includes utilities for processing collections of items either
// sequentially or concurrently, with configurable batch sizes and concurrency limits.
//
// Key Features:
//   - Batch processing nodes for item collections
//   - Batch flow execution with parameter variations
//   - Concurrent and sequential processing modes
//   - Configurable batch sizes and concurrency limits
//   - Error aggregation for batch operations
//
// Example (Batch Processing):
//
//	// Process items concurrently
//	processFunc := func(ctx context.Context, item any) (any, error) {
//	    // Process each item
//	    return fmt.Sprintf("processed: %v", item), nil
//	}
//
//	node := flyt.NewBatchNode(processFunc, true) // true for concurrent
//	shared := flyt.NewSharedStore()
//	shared.Set(flyt.KeyItems, []string{"item1", "item2", "item3"})
//
//	action, err := flyt.Run(ctx, node, shared)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	results, _ := shared.Get(flyt.KeyResults)
//	fmt.Println(results) // ["processed: item1", "processed: item2", "processed: item3"]
//
// Example (Batch Flow):
//
//	// Run a flow multiple times with different parameters
//	flowFactory := func() *flyt.Flow {
//	    // Create your flow here
//	    return flyt.NewFlow(startNode)
//	}
//
//	batchFunc := func(ctx context.Context, shared *flyt.SharedStore) ([]flyt.FlowInputs, error) {
//	    return []flyt.FlowInputs{
//	        {"user_id": 1, "action": "process"},
//	        {"user_id": 2, "action": "process"},
//	    }, nil
//	}
//
//	batchFlow := flyt.NewBatchFlow(flowFactory, batchFunc, true)
//	err := batchFlow.Run(ctx, shared)
package flyt

import (
	"context"
	"fmt"
	"sync"
)

// BatchError aggregates multiple errors from batch operations.
// It implements the error interface and provides detailed information
// about all errors that occurred during batch processing.
type BatchError struct {
	Errors []error
}

func (e *BatchError) Error() string {
	if len(e.Errors) == 0 {
		return "batch: no errors recorded"
	}
	if len(e.Errors) == 1 {
		return fmt.Sprintf("batch: %v", e.Errors[0])
	}
	return fmt.Sprintf("batch: %d errors occurred, first: %v", len(e.Errors), e.Errors[0])
}

// BatchConfig holds configuration for batch operations.
// It allows customization of batch processing behavior including
// size limits, concurrency settings, and storage keys.
type BatchConfig struct {
	// MaxBatchSize is the maximum number of items to process in a single batch.
	// If a batch exceeds this size, an error will be returned.
	// Default: 1000
	MaxBatchSize int

	// MaxConcurrency is the maximum number of concurrent workers for parallel processing.
	// This is only used when concurrent processing is enabled.
	// Default: 10
	MaxConcurrency int

	// ItemsKey is the SharedStore key to read items from.
	// Default: "items"
	ItemsKey string

	// ResultsKey is the SharedStore key to write results to.
	// Default: "results"
	ResultsKey string

	// BatchCountKey is the SharedStore key to write the batch count for BatchFlow.
	// Default: "batch_count"
	BatchCountKey string
}

// DefaultBatchConfig returns a BatchConfig with sensible defaults.
// The defaults are:
//   - MaxBatchSize: 1000
//   - MaxConcurrency: 10
//   - ItemsKey: "items"
//   - ResultsKey: "results"
//   - BatchCountKey: "batch_count"
func DefaultBatchConfig() *BatchConfig {
	return &BatchConfig{
		MaxBatchSize:   1000,
		MaxConcurrency: 10,
		ItemsKey:       KeyItems,
		ResultsKey:     KeyResults,
		BatchCountKey:  KeyBatchCount,
	}
}

// BatchProcessFunc is a function that processes a single item in a batch.
// It receives a context and an item, and returns the processed result or an error.
// The function should be thread-safe if concurrent processing is enabled.
type BatchProcessFunc func(ctx context.Context, item any) (any, error)

// NewBatchNode creates a node that processes items in batch.
// The node reads items from the SharedStore (using the key "items" by default),
// processes each item using the provided function, and stores results back
// to the SharedStore (using the key "results" by default).
//
// Parameters:
//   - processFunc: Function to process each item
//   - concurrent: If true, items are processed concurrently using a worker pool
//   - opts: Additional node options (e.g., WithMaxRetries, WithWait)
//
// Example:
//
//	processFunc := func(ctx context.Context, item any) (any, error) {
//	    // Process the item
//	    return processedItem, nil
//	}
//	node := flyt.NewBatchNode(processFunc, true)
func NewBatchNode(processFunc BatchProcessFunc, concurrent bool, opts ...NodeOption) Node {
	baseOpts := append([]NodeOption{}, opts...)
	return &batchNode{
		BaseNode:    NewBaseNode(baseOpts...),
		processFunc: processFunc,
		concurrent:  concurrent,
		config:      DefaultBatchConfig(),
	}
}

// NewBatchNodeWithKeys creates a batch node with custom keys for items and results.
// This allows you to specify which keys in the SharedStore to read items from
// and write results to, instead of using the default keys.
//
// Parameters:
//   - processFunc: Function to process each item
//   - concurrent: If true, items are processed concurrently
//   - itemsKey: SharedStore key to read items from
//   - resultsKey: SharedStore key to write results to
//   - opts: Additional node options
//
// Example:
//
//	node := flyt.NewBatchNodeWithKeys(processFunc, true, "input_data", "output_data")
func NewBatchNodeWithKeys(processFunc BatchProcessFunc, concurrent bool, itemsKey, resultsKey string, opts ...NodeOption) Node {
	config := DefaultBatchConfig()
	config.ItemsKey = itemsKey
	config.ResultsKey = resultsKey

	baseOpts := append([]NodeOption{}, opts...)
	return &batchNode{
		BaseNode:    NewBaseNode(baseOpts...),
		processFunc: processFunc,
		concurrent:  concurrent,
		config:      config,
	}
}

// NewBatchNodeWithConfig creates a batch node with custom configuration.
// This provides full control over all batch processing settings.
//
// Parameters:
//   - processFunc: Function to process each item
//   - concurrent: If true, items are processed concurrently
//   - config: Custom batch configuration
//   - opts: Additional node options
//
// Example:
//
//	config := &flyt.BatchConfig{
//	    MaxBatchSize:   500,
//	    MaxConcurrency: 20,
//	    ItemsKey:       "tasks",
//	    ResultsKey:     "completed_tasks",
//	}
//	node := flyt.NewBatchNodeWithConfig(processFunc, true, config)
func NewBatchNodeWithConfig(processFunc BatchProcessFunc, concurrent bool, config *BatchConfig, opts ...NodeOption) Node {
	baseOpts := append([]NodeOption{}, opts...)
	return &batchNode{
		BaseNode:    NewBaseNode(baseOpts...),
		processFunc: processFunc,
		concurrent:  concurrent,
		config:      config,
	}
}

type batchNode struct {
	*BaseNode
	processFunc BatchProcessFunc
	concurrent  bool
	config      *BatchConfig
}

func (n *batchNode) Prep(ctx context.Context, shared *SharedStore) (any, error) {
	// Use configured key or fall back to common keys
	keysToCheck := []string{}
	if n.config != nil && n.config.ItemsKey != "" {
		keysToCheck = []string{n.config.ItemsKey}
	} else {
		keysToCheck = []string{KeyItems, "data", "batch"}
	}

	for _, key := range keysToCheck {
		if val, ok := shared.Get(key); ok {
			return val, nil
		}
	}

	return nil, fmt.Errorf("batchNode: prep failed: no batch data found in shared store (looked for: %v)",
		keysToCheck)
}

func (n *batchNode) Exec(ctx context.Context, prepResult any) (any, error) {
	// Convert to slice
	items := ToSlice(prepResult)

	// Validate batch size
	if n.config != nil && n.config.MaxBatchSize > 0 && len(items) > n.config.MaxBatchSize {
		return nil, fmt.Errorf("batchNode: exec failed: batch size %d exceeds maximum %d", len(items), n.config.MaxBatchSize)
	}

	if len(items) == 0 {
		return []any{}, nil
	}

	results := make([]any, len(items))

	if n.concurrent {
		// Process concurrently with worker pool
		maxWorkers := 10
		if n.config != nil && n.config.MaxConcurrency > 0 {
			maxWorkers = n.config.MaxConcurrency
		}

		pool := NewWorkerPool(maxWorkers)
		defer pool.Close()

		errors := make([]error, len(items))
		var mu sync.Mutex

		for i, item := range items {
			idx, itm := i, item // Capture loop variables
			pool.Submit(func() {
				// Check context
				if ctx.Err() != nil {
					mu.Lock()
					errors[idx] = fmt.Errorf("batch item %d: context cancelled: %w", idx, ctx.Err())
					mu.Unlock()
					return
				}

				result, err := n.processFunc(ctx, itm)
				mu.Lock()
				results[idx] = result
				errors[idx] = err
				mu.Unlock()
			})
		}

		pool.Wait()

		// Collect all errors
		var batchErrors []error
		for _, err := range errors {
			if err != nil {
				batchErrors = append(batchErrors, err)
			}
		}

		if len(batchErrors) > 0 {
			return nil, &BatchError{Errors: batchErrors}
		}
	} else {
		// Process sequentially
		for i, item := range items {
			// Check context
			if err := ctx.Err(); err != nil {
				return nil, fmt.Errorf("batchNode: exec cancelled at item %d: %w", i, err)
			}

			result, err := n.processFunc(ctx, item)
			if err != nil {
				return nil, fmt.Errorf("batch item %d: process failed: %w", i, err)
			}
			results[i] = result
		}
	}

	return results, nil
}

func (n *batchNode) Post(ctx context.Context, shared *SharedStore, prepResult, execResult any) (Action, error) {
	resultsKey := KeyResults
	if n.config != nil && n.config.ResultsKey != "" {
		resultsKey = n.config.ResultsKey
	}
	shared.Set(resultsKey, execResult)
	return DefaultAction, nil
}

// FlowInputs holds input parameters for a flow iteration in batch processing.
// These parameters are merged into each flow's isolated SharedStore,
// allowing each flow instance to have its own set of input data.
//
// Example:
//
//	inputs := flyt.FlowInputs{
//	    "user_id": 123,
//	    "action": "process",
//	    "priority": "high",
//	}
type FlowInputs map[string]any

// BatchFlowFunc returns input parameters for each flow iteration in batch processing.
// This function is called during the Prep phase to determine how many flow instances
// to run and what parameters each should receive.
//
// The function receives the parent flow's SharedStore and should return a slice
// of FlowInputs, where each element represents parameters for one flow iteration.
type BatchFlowFunc func(ctx context.Context, shared *SharedStore) ([]FlowInputs, error)

// NewBatchFlow creates a flow that runs multiple times with different parameters.
// Each iteration gets its own isolated SharedStore with merged parameters.
//
// Important: The flowFactory must create new flow instances for concurrent
// execution to avoid race conditions. Do not reuse flow instances.
//
// Parameters:
//   - flowFactory: Function that creates new flow instances
//   - batchFunc: Function that returns parameters for each iteration
//   - concurrent: If true, flows run concurrently
//
// Example:
//
//	flowFactory := func() *flyt.Flow {
//	    // Create and return a new flow instance
//	    return flyt.NewFlow(startNode)
//	}
//
//	batchFunc := func(ctx context.Context, shared *flyt.SharedStore) ([]flyt.FlowInputs, error) {
//	    users, _ := shared.Get("users")
//	    var inputs []flyt.FlowInputs
//	    for _, user := range users.([]User) {
//	        inputs = append(inputs, flyt.FlowInputs{"user": user})
//	    }
//	    return inputs, nil
//	}
//
//	batchFlow := flyt.NewBatchFlow(flowFactory, batchFunc, true)
func NewBatchFlow(flowFactory FlowFactory, batchFunc BatchFlowFunc, concurrent bool) *Flow {
	batchNode := &batchFlowNode{
		BaseNode:    NewBaseNode(),
		flowFactory: flowFactory,
		batchFunc:   batchFunc,
		concurrent:  concurrent,
		config:      DefaultBatchConfig(),
	}
	return NewFlow(batchNode)
}

// NewBatchFlowWithCountKey creates a batch flow with a custom key for storing the batch count.
// After execution, the number of flows that were run is stored in the SharedStore
// using the specified key instead of the default "batch_count" key.
//
// Parameters:
//   - flowFactory: Function that creates new flow instances
//   - batchFunc: Function that returns parameters for each iteration
//   - concurrent: If true, flows run concurrently
//   - countKey: SharedStore key to store the batch count
//
// Example:
//
//	batchFlow := flyt.NewBatchFlowWithCountKey(flowFactory, batchFunc, true, "processed_count")
//	// After execution, retrieve count with: count, _ := shared.Get("processed_count")
func NewBatchFlowWithCountKey(flowFactory FlowFactory, batchFunc BatchFlowFunc, concurrent bool, countKey string) *Flow {
	config := DefaultBatchConfig()
	config.BatchCountKey = countKey

	batchNode := &batchFlowNode{
		BaseNode:    NewBaseNode(),
		flowFactory: flowFactory,
		batchFunc:   batchFunc,
		concurrent:  concurrent,
		config:      config,
	}
	return NewFlow(batchNode)
}

// NewBatchFlowWithConfig creates a batch flow with custom configuration.
// This provides full control over batch flow execution settings.
//
// Parameters:
//   - flowFactory: Function that creates new flow instances
//   - batchFunc: Function that returns parameters for each iteration
//   - concurrent: If true, flows run concurrently
//   - config: Custom batch configuration
//
// Example:
//
//	config := &flyt.BatchConfig{
//	    MaxBatchSize:   100,      // Limit to 100 flows
//	    MaxConcurrency: 5,        // Run max 5 flows concurrently
//	    BatchCountKey:  "total",  // Store count in "total" key
//	}
//	batchFlow := flyt.NewBatchFlowWithConfig(flowFactory, batchFunc, true, config)
func NewBatchFlowWithConfig(flowFactory FlowFactory, batchFunc BatchFlowFunc, concurrent bool, config *BatchConfig) *Flow {
	batchNode := &batchFlowNode{
		BaseNode:    NewBaseNode(),
		flowFactory: flowFactory,
		batchFunc:   batchFunc,
		concurrent:  concurrent,
		config:      config,
	}
	return NewFlow(batchNode)
}

type batchFlowNode struct {
	*BaseNode
	flowFactory FlowFactory
	batchFunc   BatchFlowFunc
	concurrent  bool
	config      *BatchConfig
}

func (n *batchFlowNode) Prep(ctx context.Context, shared *SharedStore) (any, error) {
	// Get batch parameters during prep phase
	batchParams, err := n.batchFunc(ctx, shared)
	if err != nil {
		return nil, fmt.Errorf("batch func failed: %w", err)
	}

	// Return both params and shared data as prepResult
	return map[string]any{
		"batchParams": batchParams,
		"sharedData":  shared.GetAll(),
	}, nil
}

func (n *batchFlowNode) Exec(ctx context.Context, prepResult any) (any, error) {
	// Extract data from prepResult
	data, ok := prepResult.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("batchFlowNode: exec failed: invalid prepResult type %T, expected map[string]any", prepResult)
	}

	batchParams, ok := data["batchParams"].([]FlowInputs)
	if !ok {
		return nil, fmt.Errorf("batchFlowNode: exec failed: invalid batchParams type in prepResult")
	}

	sharedData, ok := data["sharedData"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("batchFlowNode: exec failed: invalid sharedData type in prepResult")
	}

	// Validate batch size
	if n.config != nil && n.config.MaxBatchSize > 0 && len(batchParams) > n.config.MaxBatchSize {
		return nil, fmt.Errorf("batchFlowNode: exec failed: batch size %d exceeds maximum %d", len(batchParams), n.config.MaxBatchSize)
	}

	if len(batchParams) == 0 {
		return 0, nil
	}

	if n.concurrent {
		// Run flows concurrently with worker pool
		maxWorkers := 10
		if n.config != nil && n.config.MaxConcurrency > 0 {
			maxWorkers = n.config.MaxConcurrency
		}

		pool := NewWorkerPool(maxWorkers)
		defer pool.Close()

		errors := make([]error, len(batchParams))
		var errMu sync.Mutex
		hasError := false

		for i, params := range batchParams {
			idx, p := i, params // Capture loop variables
			sd := sharedData    // Capture shared data
			pool.Submit(func() {
				// Check context
				if ctx.Err() != nil {
					errMu.Lock()
					errors[idx] = fmt.Errorf("batchFlowNode: flow %d cancelled: %w", idx, ctx.Err())
					hasError = true
					errMu.Unlock()
					return
				}

				// Create a new flow instance for concurrent execution
				flow := n.flowFactory()

				// Create isolated shared store
				localShared := NewSharedStore()

				// Copy shared data
				localShared.Merge(sd)

				// Merge params
				localShared.Merge(p)

				if err := flow.Run(ctx, localShared); err != nil {
					errMu.Lock()
					errors[idx] = fmt.Errorf("batchFlowNode: flow %d failed: %w", idx, err)
					hasError = true
					errMu.Unlock()
				}
			})
		}

		pool.Wait()

		if hasError {
			// Collect non-nil errors
			var batchErrors []error
			for _, err := range errors {
				if err != nil {
					batchErrors = append(batchErrors, err)
				}
			}
			return nil, &BatchError{Errors: batchErrors}
		}
	} else {
		// Run flows sequentially
		for i, params := range batchParams {
			// Check context
			if err := ctx.Err(); err != nil {
				return nil, fmt.Errorf("batchFlowNode: exec cancelled at flow %d: %w", i, err)
			}

			// Create a new flow instance for each iteration
			flow := n.flowFactory()

			// Create isolated shared store for each iteration
			iterShared := NewSharedStore()
			iterShared.Merge(sharedData)
			iterShared.Merge(params)

			if err := flow.Run(ctx, iterShared); err != nil {
				return nil, fmt.Errorf("batchFlowNode: flow %d failed: %w", i, err)
			}
		}
	}

	return len(batchParams), nil
}

func (n *batchFlowNode) Post(ctx context.Context, shared *SharedStore, prepResult, execResult any) (Action, error) {
	if count, ok := execResult.(int); ok {
		countKey := KeyBatchCount
		if n.config != nil && n.config.BatchCountKey != "" {
			countKey = n.config.BatchCountKey
		}
		shared.Set(countKey, count)
	}
	return DefaultAction, nil
}
