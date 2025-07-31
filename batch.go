package flyt

import (
	"context"
	"fmt"
	"sync"
)

// BatchError aggregates multiple errors from batch operations
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

// BatchConfig holds configuration for batch operations
type BatchConfig struct {
	MaxBatchSize   int
	MaxConcurrency int
	ItemsKey       string // Key to read items from SharedStore (defaults to KeyItems)
	ResultsKey     string // Key to write results to SharedStore (defaults to KeyResults)
	BatchCountKey  string // Key to write batch count for BatchFlow (defaults to KeyBatchCount)
}

// DefaultBatchConfig returns sensible defaults
func DefaultBatchConfig() *BatchConfig {
	return &BatchConfig{
		MaxBatchSize:   1000,
		MaxConcurrency: 10,
		ItemsKey:       KeyItems,
		ResultsKey:     KeyResults,
		BatchCountKey:  KeyBatchCount,
	}
}

// BatchProcessFunc is a function that processes a single item
type BatchProcessFunc func(ctx context.Context, item any) (any, error)

// NewBatchNode creates a node that processes items in batch
func NewBatchNode(processFunc BatchProcessFunc, concurrent bool, opts ...NodeOption) Node {
	baseOpts := append([]NodeOption{}, opts...)
	return &batchNode{
		BaseNode:    NewBaseNode(baseOpts...),
		processFunc: processFunc,
		concurrent:  concurrent,
		config:      DefaultBatchConfig(),
	}
}

// NewBatchNodeWithKeys creates a batch node with custom keys for items and results
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

// NewBatchNodeWithConfig creates a batch node with custom configuration
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

// BatchParams holds parameters for a batch iteration
type BatchParams map[string]any

// BatchFlowFunc returns parameters for each batch iteration
type BatchFlowFunc func(ctx context.Context, shared *SharedStore) ([]BatchParams, error)

// NewBatchFlow creates a flow that runs multiple times with different parameters.
// The flowFactory must create new flow instances for concurrent execution to avoid
// race conditions.
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

// NewBatchFlowWithCountKey creates a batch flow with a custom key for storing the batch count
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

// NewBatchFlowWithConfig creates a batch flow with custom configuration
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

	batchParams, ok := data["batchParams"].([]BatchParams)
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
		// Create a shared store for sequential execution
		seqShared := NewSharedStore()
		seqShared.Merge(sharedData)

		for i, params := range batchParams {
			// Check context
			if err := ctx.Err(); err != nil {
				return nil, fmt.Errorf("batchFlowNode: exec cancelled at flow %d: %w", i, err)
			}

			// Create a new flow instance for each iteration
			flow := n.flowFactory()

			// Merge params into shared for the flow execution
			seqShared.Merge(params)

			if err := flow.Run(ctx, seqShared); err != nil {
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
