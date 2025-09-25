package flyt

import (
	"context"
	"time"
)

// NodeBuilder provides a fluent interface for creating and configuring nodes.
// It implements the Node interface while also providing chainable methods
// for configuration. This allows both styles:
//   - Traditional: flyt.NewNode(WithExecFunc(...), WithMaxRetries(3))
//   - Builder: flyt.NewNode().WithExecFunc(...).WithMaxRetries(3)
type NodeBuilder struct {
	*CustomNode
}

// Prep implements Node.Prep by delegating to the embedded CustomNode
func (b *NodeBuilder) Prep(ctx context.Context, shared *SharedStore) (any, error) {
	return b.CustomNode.Prep(ctx, shared)
}

// Exec implements Node.Exec by delegating to the embedded CustomNode
func (b *NodeBuilder) Exec(ctx context.Context, prepResult any) (any, error) {
	return b.CustomNode.Exec(ctx, prepResult)
}

// Post implements Node.Post by delegating to the embedded CustomNode
func (b *NodeBuilder) Post(ctx context.Context, shared *SharedStore, prepResult, execResult any) (Action, error) {
	return b.CustomNode.Post(ctx, shared, prepResult, execResult)
}

// ExecFallback implements FallbackNode.ExecFallback by delegating to the embedded CustomNode
func (b *NodeBuilder) ExecFallback(prepResult any, err error) (any, error) {
	return b.CustomNode.ExecFallback(prepResult, err)
}

// GetMaxRetries implements RetryableNode.GetMaxRetries by delegating to the embedded BaseNode
func (b *NodeBuilder) GetMaxRetries() int {
	return b.CustomNode.GetMaxRetries()
}

// GetWait implements RetryableNode.GetWait by delegating to the embedded BaseNode
func (b *NodeBuilder) GetWait() time.Duration {
	return b.CustomNode.GetWait()
}

// WithMaxRetries sets the maximum number of retries for the node's Exec phase.
// Returns the builder for method chaining.
func (b *NodeBuilder) WithMaxRetries(retries int) *NodeBuilder {
	WithMaxRetries(retries)(b.BaseNode)
	return b
}

// WithWait sets the wait duration between retries.
// Returns the builder for method chaining.
func (b *NodeBuilder) WithWait(wait time.Duration) *NodeBuilder {
	WithWait(wait)(b.BaseNode)
	return b
}

// WithPrepFunc sets a custom Prep implementation using Result types.
// Returns the builder for method chaining.
func (b *NodeBuilder) WithPrepFunc(fn func(context.Context, *SharedStore) (Result, error)) *NodeBuilder {
	b.prepFunc = fn
	return b
}

// WithExecFunc sets a custom Exec implementation using Result types.
// Returns the builder for method chaining.
func (b *NodeBuilder) WithExecFunc(fn func(context.Context, Result) (Result, error)) *NodeBuilder {
	b.execFunc = fn
	return b
}

// WithPostFunc sets a custom Post implementation using Result types.
// Returns the builder for method chaining.
func (b *NodeBuilder) WithPostFunc(fn func(context.Context, *SharedStore, Result, Result) (Action, error)) *NodeBuilder {
	b.postFunc = fn
	return b
}

// WithExecFallbackFunc sets a custom ExecFallback implementation.
// Returns the builder for method chaining.
func (b *NodeBuilder) WithExecFallbackFunc(fn func(any, error) (any, error)) *NodeBuilder {
	b.execFallbackFunc = fn
	return b
}

// WithPrepFuncAny sets a custom Prep implementation using any types.
// Returns the builder for method chaining.
func (b *NodeBuilder) WithPrepFuncAny(fn func(context.Context, *SharedStore) (any, error)) *NodeBuilder {
	b.prepFunc = func(ctx context.Context, shared *SharedStore) (Result, error) {
		val, err := fn(ctx, shared)
		if err != nil {
			return Result{}, err
		}
		return NewResult(val), nil
	}
	return b
}

// WithExecFuncAny sets a custom Exec implementation using any types.
// Returns the builder for method chaining.
func (b *NodeBuilder) WithExecFuncAny(fn func(context.Context, any) (any, error)) *NodeBuilder {
	b.execFunc = func(ctx context.Context, prepResult Result) (Result, error) {
		val, err := fn(ctx, prepResult.Value())
		if err != nil {
			return Result{}, err
		}
		return NewResult(val), nil
	}
	return b
}

// WithPostFuncAny sets a custom Post implementation using any types.
// Returns the builder for method chaining.
func (b *NodeBuilder) WithPostFuncAny(fn func(context.Context, *SharedStore, any, any) (Action, error)) *NodeBuilder {
	b.postFunc = func(ctx context.Context, shared *SharedStore, prepResult, execResult Result) (Action, error) {
		return fn(ctx, shared, prepResult.Value(), execResult.Value())
	}
	return b
}

// WithBatchConcurrency sets the concurrency level for batch processing.
// Returns the builder for method chaining.
func (b *NodeBuilder) WithBatchConcurrency(n int) *NodeBuilder {
	b.batchConcurrency = n
	return b
}

// WithBatchErrorHandling sets the error handling strategy for batch processing.
// Returns the builder for method chaining.
func (b *NodeBuilder) WithBatchErrorHandling(continueOnError bool) *NodeBuilder {
	if continueOnError {
		b.batchErrorHandling = "continue"
	} else {
		b.batchErrorHandling = "stop"
	}
	return b
}
