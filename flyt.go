// Package flyt is a minimalist workflow framework for Go, inspired by Pocket Flow.
// It provides a simple graph-based abstraction for orchestrating tasks.
//
// Thread Safety:
// When using concurrent batch operations, ensure that your Node implementations
// are thread-safe. The framework provides SharedStore for safe concurrent access
// to shared data.
//
// Example:
//
//	// Define a simple node
//	type PrintNode struct {
//	    *flyt.BaseNode
//	}
//
//	func (n *PrintNode) Exec(ctx context.Context, prepResult any) (any, error) {
//	    fmt.Println("Hello from node!")
//	    return nil, nil
//	}
//
//	// Create and run a flow
//	node := &PrintNode{BaseNode: flyt.NewBaseNode()}
//	shared := flyt.NewSharedStore()
//
//	ctx := context.Background()
//	action, err := flyt.Run(ctx, node, shared)
//	if err != nil {
//	    log.Fatal(err)
//	}
package flyt

import (
	"context"
	"fmt"
	"reflect"
	"sync"
	"time"
)

// SharedStore provides thread-safe access to shared data
type SharedStore struct {
	mu   sync.RWMutex
	data map[string]any
}

// NewSharedStore creates a new thread-safe shared store
func NewSharedStore() *SharedStore {
	return &SharedStore{
		data: make(map[string]any),
	}
}

// Get retrieves a value from the store
func (s *SharedStore) Get(key string) (any, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, ok := s.data[key]
	return val, ok
}

// Set stores a value in the store
func (s *SharedStore) Set(key string, value any) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = value
}

// GetAll returns a copy of all data
func (s *SharedStore) GetAll() map[string]any {
	s.mu.RLock()
	defer s.mu.RUnlock()
	copy := make(map[string]any, len(s.data))
	for k, v := range s.data {
		copy[k] = v
	}
	return copy
}

// Merge merges another map into the store
func (s *SharedStore) Merge(data map[string]any) {
	if data == nil {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	for k, v := range data {
		s.data[k] = v
	}
}

// Action represents the next action to take after a node executes
type Action string

const (
	// DefaultAction is the default action if none is specified
	DefaultAction Action = "default"

	// KeyItems is the shared store key for items to be processed
	KeyItems = "items"
	// KeyData is the shared store key for generic data
	KeyData = "data"
	// KeyBatch is the shared store key for batch data
	KeyBatch = "batch"
	// KeyResults is the shared store key for processing results
	KeyResults = "results"
	// KeyBatchCount is the shared store key for batch count
	KeyBatchCount = "batch_count"
)

// Node is the interface that all nodes must implement.
//
// Important: Nodes should not be shared across concurrent flow executions.
// If you need to run the same logic concurrently, create separate node instances.
type Node interface {
	// Prep reads and preprocesses data from shared store
	Prep(ctx context.Context, shared *SharedStore) (any, error)

	// Exec executes the main logic with optional retries
	Exec(ctx context.Context, prepResult any) (any, error)

	// Post processes results and writes back to shared store
	Post(ctx context.Context, shared *SharedStore, prepResult, execResult any) (Action, error)
}

// BaseNode provides a base implementation of Node
type BaseNode struct {
	mu         sync.RWMutex
	maxRetries int
	wait       time.Duration
}

// NewBaseNode creates a new BaseNode with options
func NewBaseNode(opts ...NodeOption) *BaseNode {
	n := &BaseNode{
		maxRetries: 1,
		wait:       0,
	}

	for _, opt := range opts {
		opt(n)
	}

	return n
}

// NodeOption is a function that configures a BaseNode
type NodeOption func(*BaseNode)

// WithMaxRetries sets the maximum number of retries
func WithMaxRetries(retries int) NodeOption {
	return func(n *BaseNode) {
		n.maxRetries = retries
	}
}

// WithWait sets the wait duration between retries
func WithWait(wait time.Duration) NodeOption {
	return func(n *BaseNode) {
		n.wait = wait
	}
}

// GetMaxRetries returns the maximum number of retries
func (n *BaseNode) GetMaxRetries() int {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.maxRetries
}

// GetWait returns the wait duration between retries
func (n *BaseNode) GetWait() time.Duration {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.wait
}

// Prep is the default prep implementation (can be overridden)
func (n *BaseNode) Prep(ctx context.Context, shared *SharedStore) (any, error) {
	return nil, nil
}

// Exec is the default exec implementation (must be overridden)
func (n *BaseNode) Exec(ctx context.Context, prepResult any) (any, error) {
	return nil, nil
}

// Post is the default post implementation (can be overridden)
func (n *BaseNode) Post(ctx context.Context, shared *SharedStore, prepResult, execResult any) (Action, error) {
	return DefaultAction, nil
}

// ExecFallback handles errors after all retries are exhausted
func (n *BaseNode) ExecFallback(prepResult any, err error) (any, error) {
	return nil, err
}

// RetryableNode is a node that supports retries
type RetryableNode interface {
	Node
	GetMaxRetries() int
	GetWait() time.Duration
}

// FallbackNode is a node that supports fallback on error
type FallbackNode interface {
	ExecFallback(prepResult any, err error) (any, error)
}

// Run executes the node with the prep->exec->post lifecycle
func Run(ctx context.Context, node Node, shared *SharedStore) (Action, error) {
	// Check context before each phase
	if err := ctx.Err(); err != nil {
		return "", fmt.Errorf("run: context cancelled: %w", err)
	}

	// Prep phase
	prepResult, err := node.Prep(ctx, shared)
	if err != nil {
		return "", fmt.Errorf("run: prep failed: %w", err)
	}

	// Check context again
	if err := ctx.Err(); err != nil {
		return "", fmt.Errorf("run: context cancelled after prep: %w", err)
	}

	// Get retry settings if available
	var maxRetries int = 1
	var wait time.Duration = 0

	if retryable, ok := node.(RetryableNode); ok {
		maxRetries = retryable.GetMaxRetries()
		wait = retryable.GetWait()
	}

	// Exec phase with retries
	var execResult any
	var execErr error

	for attempt := 0; attempt < maxRetries; attempt++ {
		// Check context before retry
		if err := ctx.Err(); err != nil {
			return "", fmt.Errorf("run: context cancelled during retry: %w", err)
		}

		if attempt > 0 && wait > 0 {
			select {
			case <-time.After(wait):
				// Continue with retry
			case <-ctx.Done():
				return "", fmt.Errorf("run: context cancelled during wait: %w", ctx.Err())
			}
		}

		execResult, execErr = node.Exec(ctx, prepResult)
		if execErr == nil {
			break
		}
	}

	// Handle exec failure
	if execErr != nil {
		if fallback, ok := node.(FallbackNode); ok {
			execResult, execErr = fallback.ExecFallback(prepResult, execErr)
		}
		if execErr != nil {
			return "", fmt.Errorf("run: exec failed after %d retries: %w", maxRetries, execErr)
		}
	}

	// Post phase
	action, err := node.Post(ctx, shared, prepResult, execResult)
	if err != nil {
		return "", fmt.Errorf("run: post failed: %w", err)
	}

	if action == "" {
		action = DefaultAction
	}

	return action, nil
}

// Flow represents a workflow of connected nodes
type Flow struct {
	*BaseNode
	start       Node
	transitions map[Node]map[Action]Node
}

// NewFlow creates a new Flow with a start node
func NewFlow(start Node) *Flow {
	return &Flow{
		BaseNode:    NewBaseNode(),
		start:       start,
		transitions: make(map[Node]map[Action]Node),
	}
}

// Connect adds a transition from one node to another based on an action
func (f *Flow) Connect(from Node, action Action, to Node) {
	if f.transitions[from] == nil {
		f.transitions[from] = make(map[Action]Node)
	}
	f.transitions[from][action] = to
}

// Run executes the flow starting from the start node
func (f *Flow) Run(ctx context.Context, shared *SharedStore) error {
	if f.start == nil {
		return fmt.Errorf("flow: run failed: no start node configured")
	}
	current := f.start

	for current != nil {
		// Check context
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("flow: run cancelled: %w", err)
		}

		// Run the current node
		var action Action
		var err error

		// Check if current is a Flow
		if subFlow, ok := current.(*Flow); ok {
			err = subFlow.Run(ctx, shared)
			if err != nil {
				return err
			}
			// Get action from flow's post
			action, err = subFlow.Post(ctx, shared, nil, nil)
			if err != nil {
				return err
			}
		} else {
			// Regular node
			action, err = Run(ctx, current, shared)
			if err != nil {
				return err
			}
		}

		// Find next node based on action
		if transitions, ok := f.transitions[current]; ok {
			if next, ok := transitions[action]; ok {
				current = next

			} else {
				// No transition for this action, flow ends
				break
			}
		} else {
			// No transitions defined for this node, flow ends
			break
		}
	}

	return nil
}

// Prep implements Node interface for Flow
func (f *Flow) Prep(ctx context.Context, shared *SharedStore) (any, error) {
	return nil, nil
}

// Exec implements Node interface for Flow (not used)
func (f *Flow) Exec(ctx context.Context, prepResult any) (any, error) {
	return nil, nil
}

// Post implements Node interface for Flow
func (f *Flow) Post(ctx context.Context, shared *SharedStore, prepResult, execResult any) (Action, error) {
	return DefaultAction, nil
}

// WorkerPool manages concurrent task execution
type WorkerPool struct {
	workers int
	tasks   chan func()
	wg      sync.WaitGroup
	done    chan struct{}
}

// NewWorkerPool creates a new worker pool
func NewWorkerPool(workers int) *WorkerPool {
	if workers <= 0 {
		workers = 1
	}

	p := &WorkerPool{
		workers: workers,
		tasks:   make(chan func(), workers*2),
		done:    make(chan struct{}),
	}

	// Start workers
	for i := 0; i < workers; i++ {
		go p.worker()
	}

	return p
}

func (p *WorkerPool) worker() {
	for {
		select {
		case task, ok := <-p.tasks:
			if !ok {
				return
			}
			task()
		case <-p.done:
			return
		}
	}
}

// Submit submits a task to the pool
func (p *WorkerPool) Submit(task func()) {
	p.wg.Add(1)
	p.tasks <- func() {
		defer p.wg.Done()
		task()
	}
}

// Wait waits for all tasks to complete
func (p *WorkerPool) Wait() {
	p.wg.Wait()
}

// Close closes the worker pool and waits for all workers to finish
func (p *WorkerPool) Close() {
	close(p.done)
	close(p.tasks)
}

// ToSlice converts various types to []any (exported for testing)
func ToSlice(v any) []any {
	if v == nil {
		return []any{}
	}

	switch val := v.(type) {
	case []any:
		return val
	case []string:
		result := make([]any, len(val))
		for i, v := range val {
			result[i] = v
		}
		return result
	case []int:
		result := make([]any, len(val))
		for i, v := range val {
			result[i] = v
		}
		return result
	case []float64:
		result := make([]any, len(val))
		for i, v := range val {
			result[i] = v
		}
		return result
	case []map[string]any:
		result := make([]any, len(val))
		for i, v := range val {
			result[i] = v
		}
		return result
	default:
		// Try reflection for other slice types
		rv := reflect.ValueOf(v)
		if rv.Kind() == reflect.Slice {
			result := make([]any, rv.Len())
			for i := 0; i < rv.Len(); i++ {
				result[i] = rv.Index(i).Interface()
			}
			return result
		}
		return []any{v} // Single item
	}
}

// FlowFactory creates new instances of a flow
type FlowFactory func() *Flow

// CustomNode is a node implementation that uses custom functions
type CustomNode struct {
	*BaseNode
	prepFunc func(context.Context, *SharedStore) (any, error)
	execFunc func(context.Context, any) (any, error)
	postFunc func(context.Context, *SharedStore, any, any) (Action, error)
}

// Prep implements Node.Prep by calling the custom prepFunc if provided
func (n *CustomNode) Prep(ctx context.Context, shared *SharedStore) (any, error) {
	if n.prepFunc != nil {
		return n.prepFunc(ctx, shared)
	}
	return n.BaseNode.Prep(ctx, shared)
}

// Exec implements Node.Exec by calling the custom execFunc if provided
func (n *CustomNode) Exec(ctx context.Context, prepResult any) (any, error) {
	if n.execFunc != nil {
		return n.execFunc(ctx, prepResult)
	}
	return n.BaseNode.Exec(ctx, prepResult)
}

// Post implements Node.Post by calling the custom postFunc if provided
func (n *CustomNode) Post(ctx context.Context, shared *SharedStore, prepResult, execResult any) (Action, error) {
	if n.postFunc != nil {
		return n.postFunc(ctx, shared, prepResult, execResult)
	}
	return n.BaseNode.Post(ctx, shared, prepResult, execResult)
}

// NewNode creates a new node with custom function implementations
func NewNode(opts ...any) Node {
	node := &CustomNode{
		BaseNode: NewBaseNode(),
	}

	// Separate options by type
	var customOpts []CustomNodeOption
	var baseOpts []NodeOption

	for _, opt := range opts {
		switch o := opt.(type) {
		case CustomNodeOption:
			customOpts = append(customOpts, o)
		case NodeOption:
			baseOpts = append(baseOpts, o)
		case func(*BaseNode):
			baseOpts = append(baseOpts, NodeOption(o))
		default:
			// Ignore unknown option types
		}
	}

	// Apply base node options first
	for _, opt := range baseOpts {
		opt(node.BaseNode)
	}

	// Apply custom node options
	for _, opt := range customOpts {
		opt.apply(node)
	}

	return node
}

// CustomNodeOption is an option for configuring a CustomNode
type CustomNodeOption interface {
	apply(*CustomNode)
}

// customNodeOption is the internal implementation of CustomNodeOption
type customNodeOption struct {
	f func(*CustomNode)
}

func (o *customNodeOption) apply(n *CustomNode) {
	o.f(n)
}

// WithPrepFunc sets a custom Prep implementation
func WithPrepFunc(fn func(context.Context, *SharedStore) (any, error)) CustomNodeOption {
	return &customNodeOption{
		f: func(n *CustomNode) {
			n.prepFunc = fn
		},
	}
}

// WithExecFunc sets a custom Exec implementation
func WithExecFunc(fn func(context.Context, any) (any, error)) CustomNodeOption {
	return &customNodeOption{
		f: func(n *CustomNode) {
			n.execFunc = fn
		},
	}
}

// WithPostFunc sets a custom Post implementation
func WithPostFunc(fn func(context.Context, *SharedStore, any, any) (Action, error)) CustomNodeOption {
	return &customNodeOption{
		f: func(n *CustomNode) {
			n.postFunc = fn
		},
	}
}
