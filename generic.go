package flyt

import "context"

// NodeG is a generic version of Node interface for type-safe implementations
type NodeG[P, E any] interface {
	// PrepG reads and preprocesses data from shared store with typed result
	PrepG(ctx context.Context, shared *SharedStore) (P, error)

	// ExecG executes the main logic with typed input and output
	ExecG(ctx context.Context, prepResult P) (E, error)

	// PostG processes results with typed inputs
	PostG(ctx context.Context, shared *SharedStore, prepResult P, execResult E) (Action, error)
}

// BaseNodeG provides a generic base implementation
type BaseNodeG[P, E any] struct {
	*BaseNode
}

// NewBaseNodeG creates a new generic BaseNode
func NewBaseNodeG[P, E any](opts ...NodeOption) *BaseNodeG[P, E] {
	return &BaseNodeG[P, E]{
		BaseNode: NewBaseNode(opts...),
	}
}

// PrepG default implementation
func (n *BaseNodeG[P, E]) PrepG(ctx context.Context, shared *SharedStore) (P, error) {
	var zero P
	return zero, nil
}

// ExecG default implementation
func (n *BaseNodeG[P, E]) ExecG(ctx context.Context, prepResult P) (E, error) {
	var zero E
	return zero, nil
}

// PostG default implementation
func (n *BaseNodeG[P, E]) PostG(ctx context.Context, shared *SharedStore, prepResult P, execResult E) (Action, error) {
	return DefaultAction, nil
}

// NodeAdapter adapts a generic node to the standard Node interface
type NodeAdapter[P, E any] struct {
	node NodeG[P, E]
}

// NewNodeAdapter creates an adapter from generic to standard node
func NewNodeAdapter[P, E any](node NodeG[P, E]) Node {
	return &NodeAdapter[P, E]{node: node}
}

func (a *NodeAdapter[P, E]) Prep(ctx context.Context, shared *SharedStore) (any, error) {
	return a.node.PrepG(ctx, shared)
}

func (a *NodeAdapter[P, E]) Exec(ctx context.Context, prepResult any) (any, error) {
	// Safe type assertion with fallback
	if typed, ok := prepResult.(P); ok {
		return a.node.ExecG(ctx, typed)
	}
	// Attempt zero value if nil
	if prepResult == nil {
		var zero P
		return a.node.ExecG(ctx, zero)
	}
	var zero E
	return zero, nil
}

func (a *NodeAdapter[P, E]) Post(ctx context.Context, shared *SharedStore, prepResult, execResult any) (Action, error) {
	var p P
	var e E

	if typed, ok := prepResult.(P); ok {
		p = typed
	}
	if typed, ok := execResult.(E); ok {
		e = typed
	}

	return a.node.PostG(ctx, shared, p, e)
}

// RunG executes a generic node with type safety
func RunG[P, E any](ctx context.Context, node NodeG[P, E], shared *SharedStore) (Action, error) {
	// Use the adapter to run through standard flow
	return Run(ctx, NewNodeAdapter(node), shared)
}
