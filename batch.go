package flyt

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// BatchError aggregates multiple errors from batch operations.
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

// BatchNode is a marker type that indicates batch processing.
// It embeds CustomNode and works exactly like a regular node,
// but the framework handles it specially.
type BatchNode struct {
	*CustomNode
	batchPrepFunc func(context.Context, *SharedStore) ([]Result, error)
	batchPostFunc func(context.Context, *SharedStore, []Result, []Result) (Action, error)
}

// Prep delegates to batchPrepFunc if set
func (n *BatchNode) Prep(ctx context.Context, shared *SharedStore) (any, error) {
	if n.batchPrepFunc != nil {
		return n.batchPrepFunc(ctx, shared)
	}
	// Fallback to CustomNode's Prep
	return n.CustomNode.Prep(ctx, shared)
}

// Post handles batch results
func (n *BatchNode) Post(ctx context.Context, shared *SharedStore, prepResult, execResult any) (Action, error) {
	if n.batchPostFunc != nil {
		prep := prepResult.([]Result)
		exec := execResult.([]Result)
		return n.batchPostFunc(ctx, shared, prep, exec)
	}
	return DefaultAction, nil
}

// BatchNodeBuilder provides fluent interface for BatchNode
type BatchNodeBuilder struct {
	*BatchNode
}

// NewBatchNode creates a new batch node with fluent interface
// Usage: flyt.NewBatchNode() from outside the package
func NewBatchNode(opts ...any) *BatchNodeBuilder {
	customNode := &CustomNode{
		BaseNode: NewBaseNode(),
	}

	// Process options (similar to NewNode)
	var baseOpts []NodeOption
	for _, opt := range opts {
		switch o := opt.(type) {
		case NodeOption:
			baseOpts = append(baseOpts, o)
		case func(*BaseNode):
			baseOpts = append(baseOpts, NodeOption(o))
		}
	}

	for _, opt := range baseOpts {
		opt(customNode.BaseNode)
	}

	return &BatchNodeBuilder{
		BatchNode: &BatchNode{CustomNode: customNode},
	}
}

// Implement Node interface
func (b *BatchNodeBuilder) Prep(ctx context.Context, shared *SharedStore) (any, error) {
	return b.BatchNode.Prep(ctx, shared)
}

func (b *BatchNodeBuilder) Exec(ctx context.Context, prepResult any) (any, error) {
	return b.BatchNode.Exec(ctx, prepResult)
}

func (b *BatchNodeBuilder) Post(ctx context.Context, shared *SharedStore, prepResult, execResult any) (Action, error) {
	return b.BatchNode.Post(ctx, shared, prepResult, execResult)
}

// Fluent configuration methods
func (b *BatchNodeBuilder) WithMaxRetries(retries int) *BatchNodeBuilder {
	WithMaxRetries(retries)(b.BatchNode.CustomNode.BaseNode)
	return b
}

func (b *BatchNodeBuilder) WithWait(wait time.Duration) *BatchNodeBuilder {
	WithWait(wait)(b.BatchNode.CustomNode.BaseNode)
	return b
}

func (b *BatchNodeBuilder) WithBatchConcurrency(n int) *BatchNodeBuilder {
	b.BatchNode.CustomNode.BaseNode.batchConcurrency = n
	return b
}

func (b *BatchNodeBuilder) WithBatchErrorHandling(continueOnError bool) *BatchNodeBuilder {
	if continueOnError {
		b.BatchNode.CustomNode.BaseNode.batchErrorHandling = "continue"
	} else {
		b.BatchNode.CustomNode.BaseNode.batchErrorHandling = "stop"
	}
	return b
}

// Note: When used from outside the package, these will be:
// - func(context.Context, *flyt.SharedStore) ([]flyt.Result, error)
// - func(context.Context, flyt.Result) (flyt.Result, error)
// - func(context.Context, *flyt.SharedStore, []flyt.Result, []flyt.Result) (flyt.Action, error)

func (b *BatchNodeBuilder) WithPrepFunc(fn func(context.Context, *SharedStore) ([]Result, error)) *BatchNodeBuilder {
	b.BatchNode.batchPrepFunc = fn
	return b
}

func (b *BatchNodeBuilder) WithExecFunc(fn func(context.Context, Result) (Result, error)) *BatchNodeBuilder {
	b.BatchNode.CustomNode.execFunc = fn
	return b
}

func (b *BatchNodeBuilder) WithPostFunc(fn func(context.Context, *SharedStore, []Result, []Result) (Action, error)) *BatchNodeBuilder {
	b.BatchNode.batchPostFunc = fn
	return b
}

// Alternative: WithExecFuncAny for compatibility
func (b *BatchNodeBuilder) WithExecFuncAny(fn func(context.Context, any) (any, error)) *BatchNodeBuilder {
	b.BatchNode.CustomNode.execFunc = func(ctx context.Context, prepResult Result) (Result, error) {
		val, err := fn(ctx, prepResult.Value())
		if err != nil {
			return Result{}, err
		}
		return NewResult(val), nil
	}
	return b
}

// runBatch handles the execution of batch nodes
func runBatch(ctx context.Context, node Node, shared *SharedStore) (Action, error) {
	// Prep phase - returns []Result
	prepResult, err := node.Prep(ctx, shared)
	if err != nil {
		return "", fmt.Errorf("run: prep failed: %w", err)
	}

	// Convert to []Result
	var items []Result
	switch v := prepResult.(type) {
	case []Result:
		items = v
	case []any:
		items = make([]Result, len(v))
		for i, item := range v {
			items[i] = NewResult(item)
		}
	default:
		// Try to convert using ToSlice
		slice := ToSlice(prepResult)
		items = make([]Result, len(slice))
		for i, item := range slice {
			items[i] = NewResult(item)
		}
	}

	if len(items) == 0 {
		// No items to process
		action, err := node.Post(ctx, shared, []Result{}, []Result{})
		if err != nil {
			return "", fmt.Errorf("run: post failed: %w", err)
		}
		return action, nil
	}

	// Get batch configuration
	var concurrency int
	var errorHandling string = "continue"

	if baseNode, ok := node.(*BaseNode); ok {
		concurrency = baseNode.GetBatchConcurrency()
		errorHandling = baseNode.GetBatchErrorHandling()
	} else if customNode, ok := node.(*CustomNode); ok {
		concurrency = customNode.BaseNode.GetBatchConcurrency()
		errorHandling = customNode.BaseNode.GetBatchErrorHandling()
	} else if batchNode, ok := node.(*BatchNode); ok {
		concurrency = batchNode.CustomNode.BaseNode.GetBatchConcurrency()
		errorHandling = batchNode.CustomNode.BaseNode.GetBatchErrorHandling()
	} else if batchBuilder, ok := node.(*BatchNodeBuilder); ok {
		concurrency = batchBuilder.BatchNode.CustomNode.BaseNode.GetBatchConcurrency()
		errorHandling = batchBuilder.BatchNode.CustomNode.BaseNode.GetBatchErrorHandling()
	}

	// Execute items
	results := make([]Result, len(items))

	if concurrency > 0 {
		runBatchConcurrent(ctx, node, items, results, concurrency, errorHandling)
	} else {
		runBatchSequential(ctx, node, items, results, errorHandling)
	}

	// Post phase - called once with all results
	action, err := node.Post(ctx, shared, items, results)
	if err != nil {
		return "", fmt.Errorf("run: post failed: %w", err)
	}

	if action == "" {
		action = DefaultAction
	}

	return action, nil
}

func runBatchSequential(ctx context.Context, node Node, items []Result, results []Result, errorHandling string) {
	for i, item := range items {
		if ctx.Err() != nil {
			results[i] = NewErrorResult(fmt.Errorf("context cancelled"))
			if errorHandling == "stop" {
				break
			}
			continue
		}

		execResult, err := runExecWithRetries(ctx, node, item)
		if err != nil {
			results[i] = NewErrorResult(err)
			if errorHandling == "stop" {
				break
			}
		} else {
			if r, ok := execResult.(Result); ok {
				results[i] = r
			} else {
				results[i] = NewResult(execResult)
			}
		}
	}
}

func runBatchConcurrent(ctx context.Context, node Node, items []Result, results []Result, concurrency int, errorHandling string) {
	pool := NewWorkerPool(concurrency)
	defer pool.Close()

	var mu sync.Mutex
	shouldStop := false

	for i, item := range items {
		idx := i
		itm := item

		pool.Submit(func() {
			mu.Lock()
			if shouldStop && errorHandling == "stop" {
				results[idx] = NewErrorResult(fmt.Errorf("batch stopped due to error"))
				mu.Unlock()
				return
			}
			mu.Unlock()

			if ctx.Err() != nil {
				results[idx] = NewErrorResult(fmt.Errorf("context cancelled"))
				return
			}

			execResult, err := runExecWithRetries(ctx, node, itm)

			mu.Lock()
			if err != nil {
				results[idx] = NewErrorResult(err)
				if errorHandling == "stop" {
					shouldStop = true
				}
			} else {
				if r, ok := execResult.(Result); ok {
					results[idx] = r
				} else {
					results[idx] = NewResult(execResult)
				}
			}
			mu.Unlock()
		})
	}

	pool.Wait()
}

func runExecWithRetries(ctx context.Context, node Node, item Result) (any, error) {
	// Get retry settings
	var maxRetries int = 1
	var wait time.Duration = 0

	if retryable, ok := node.(RetryableNode); ok {
		maxRetries = retryable.GetMaxRetries()
		wait = retryable.GetWait()
	}

	var execResult any
	var execErr error

	for attempt := 0; attempt < maxRetries; attempt++ {
		if ctx.Err() != nil {
			return nil, fmt.Errorf("context cancelled during retry: %w", ctx.Err())
		}

		if attempt > 0 && wait > 0 {
			select {
			case <-time.After(wait):
			case <-ctx.Done():
				return nil, fmt.Errorf("context cancelled during wait: %w", ctx.Err())
			}
		}

		execResult, execErr = node.Exec(ctx, item)
		if execErr == nil {
			break
		}
	}

	if execErr != nil {
		if fallback, ok := node.(FallbackNode); ok {
			return fallback.ExecFallback(item, execErr)
		}
		return nil, execErr
	}

	return execResult, nil
}
