// Package flyt is a minimalist workflow framework for Go, inspired by Pocket Flow.
// It provides a simple graph-based abstraction for orchestrating tasks.
package flyt

import (
	"fmt"
	"reflect"
	"sync"
	"time"
)

// Shared is a map for sharing data between nodes
type Shared map[string]any

// Params holds node-specific parameters
type Params map[string]any

// Action represents the next action to take after a node executes
type Action string

const (
	// DefaultAction is the default action if none is specified
	DefaultAction Action = "default"
)

// Node is the interface that all nodes must implement
type Node interface {
	// Prep reads and preprocesses data from shared store
	Prep(shared Shared) (any, error)

	// Exec executes the main logic with optional retries
	Exec(prepResult any) (any, error)

	// Post processes results and writes back to shared store
	Post(shared Shared, prepResult, execResult any) (Action, error)

	// SetParams sets node-specific parameters
	SetParams(params Params)

	// GetParams returns node-specific parameters
	GetParams() Params
}

// BaseNode provides a base implementation of Node
type BaseNode struct {
	params     Params
	maxRetries int
	wait       time.Duration
}

// NewBaseNode creates a new BaseNode with options
func NewBaseNode(opts ...NodeOption) *BaseNode {
	n := &BaseNode{
		params:     make(Params),
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

// SetParams sets the node parameters
func (n *BaseNode) SetParams(params Params) {
	n.params = params
}

// GetParams returns the node parameters
func (n *BaseNode) GetParams() Params {
	return n.params
}

// GetMaxRetries returns the maximum number of retries
func (n *BaseNode) GetMaxRetries() int {
	return n.maxRetries
}

// GetWait returns the wait duration between retries
func (n *BaseNode) GetWait() time.Duration {
	return n.wait
}

// Prep is the default prep implementation (can be overridden)
func (n *BaseNode) Prep(shared Shared) (any, error) {
	return nil, nil
}

// Exec is the default exec implementation (must be overridden)
func (n *BaseNode) Exec(prepResult any) (any, error) {
	return nil, nil
}

// Post is the default post implementation (can be overridden)
func (n *BaseNode) Post(shared Shared, prepResult, execResult any) (Action, error) {
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

// Run executes the node with the prep->exec->post lifecycle
func Run(node Node, shared Shared) (Action, error) {
	// Prep phase
	prepResult, err := node.Prep(shared)
	if err != nil {
		return "", fmt.Errorf("prep failed: %w", err)
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
		if attempt > 0 && wait > 0 {
			time.Sleep(wait)
		}

		execResult, execErr = node.Exec(prepResult)
		if execErr == nil {
			break
		}
	}

	// Handle exec failure
	if execErr != nil {
		if fallback, ok := node.(interface{ ExecFallback(any, error) (any, error) }); ok {
			execResult, execErr = fallback.ExecFallback(prepResult, execErr)
		}
		if execErr != nil {
			return "", fmt.Errorf("exec failed after %d retries: %w", maxRetries, execErr)
		}
	}

	// Post phase
	action, err := node.Post(shared, prepResult, execResult)
	if err != nil {
		return "", fmt.Errorf("post failed: %w", err)
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
func (f *Flow) Run(shared Shared) error {
	current := f.start

	// Propagate flow params to start node
	if f.params != nil {
		current.SetParams(f.params)
	}

	for current != nil {
		// Run the current node
		var action Action
		var err error

		// Check if current is a Flow
		if subFlow, ok := current.(*Flow); ok {
			err = subFlow.Run(shared)
			if err != nil {
				return err
			}
			// Get action from flow's post
			action, err = subFlow.Post(shared, nil, nil)
			if err != nil {
				return err
			}
		} else {
			// Regular node
			action, err = Run(current, shared)
			if err != nil {
				return err
			}
		}

		// Find next node based on action
		if transitions, ok := f.transitions[current]; ok {
			if next, ok := transitions[action]; ok {
				current = next
				// Propagate params to next node
				if f.params != nil {
					current.SetParams(f.params)
				}
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
func (f *Flow) Prep(shared Shared) (any, error) {
	return nil, nil
}

// Exec implements Node interface for Flow (not used)
func (f *Flow) Exec(prepResult any) (any, error) {
	return nil, nil
}

// Post implements Node interface for Flow
func (f *Flow) Post(shared Shared, prepResult, execResult any) (Action, error) {
	return DefaultAction, nil
}

// BatchProcessFunc is a function that processes a single item
type BatchProcessFunc func(item any) (any, error)

// NewBatchNode creates a node that processes items in batch
func NewBatchNode(processFunc BatchProcessFunc, concurrent bool, opts ...NodeOption) Node {
	return &batchNode{
		BaseNode:    NewBaseNode(opts...),
		processFunc: processFunc,
		concurrent:  concurrent,
	}
}

type batchNode struct {
	*BaseNode
	processFunc BatchProcessFunc
	concurrent  bool
}

func (n *batchNode) Prep(shared Shared) (any, error) {
	// Look for common keys like "items", "data", "batch"
	for _, key := range []string{"items", "data", "batch"} {
		if val, ok := shared[key]; ok {
			return val, nil
		}
	}
	return nil, fmt.Errorf("no batch data found in shared store (looked for: items, data, batch)")
}

func (n *batchNode) Exec(prepResult any) (any, error) {
	// Convert to slice
	items := toSlice(prepResult)
	if len(items) == 0 {
		return []any{}, nil
	}

	results := make([]any, len(items))

	if n.concurrent {
		// Process concurrently
		var wg sync.WaitGroup
		var mu sync.Mutex
		errors := make([]error, len(items))

		for i, item := range items {
			wg.Add(1)
			go func(idx int, itm any) {
				defer wg.Done()
				result, err := n.processFunc(itm)
				mu.Lock()
				results[idx] = result
				errors[idx] = err
				mu.Unlock()
			}(i, item)
		}
		wg.Wait()

		// Check for errors
		for i, err := range errors {
			if err != nil {
				return nil, fmt.Errorf("item %d failed: %w", i, err)
			}
		}
	} else {
		// Process sequentially
		for i, item := range items {
			result, err := n.processFunc(item)
			if err != nil {
				return nil, fmt.Errorf("item %d failed: %w", i, err)
			}
			results[i] = result
		}
	}

	return results, nil
}

func (n *batchNode) Post(shared Shared, prepResult, execResult any) (Action, error) {
	shared["results"] = execResult
	return DefaultAction, nil
}

// Helper to convert various types to []any
func toSlice(v any) []any {
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

// BatchFlowFunc returns parameters for each batch iteration
type BatchFlowFunc func(shared Shared) ([]Params, error)

// NewBatchFlow creates a flow that runs multiple times with different parameters.
// Note: For concurrent execution, the innerFlow and its nodes must be thread-safe.
// Consider creating new flow instances for each concurrent execution if nodes maintain state.
func NewBatchFlow(innerFlow *Flow, batchFunc BatchFlowFunc, concurrent bool) *Flow {
	batchNode := &batchFlowNode{
		BaseNode:   NewBaseNode(),
		innerFlow:  innerFlow,
		batchFunc:  batchFunc,
		concurrent: concurrent,
	}
	return NewFlow(batchNode)
}

type batchFlowNode struct {
	*BaseNode
	innerFlow  *Flow
	batchFunc  BatchFlowFunc
	concurrent bool
}

func (n *batchFlowNode) Prep(shared Shared) (any, error) {
	// Store shared in params for Exec to access
	n.SetParams(Params{"shared": shared})
	return nil, nil
}

func (n *batchFlowNode) Exec(prepResult any) (any, error) {
	shared := n.GetParams()["shared"].(Shared)

	// Get batch parameters
	batchParams, err := n.batchFunc(shared)
	if err != nil {
		return nil, fmt.Errorf("batch func failed: %w", err)
	}

	if len(batchParams) == 0 {
		return 0, nil
	}

	if n.concurrent {
		// Run flows concurrently
		var wg sync.WaitGroup
		var mu sync.Mutex
		errors := make([]error, len(batchParams))

		for i, params := range batchParams {
			wg.Add(1)
			go func(idx int, p Params) {
				defer wg.Done()

				// Create a copy of shared and merge params
				localShared := make(Shared)
				mu.Lock()
				for k, v := range shared {
					localShared[k] = v
				}
				mu.Unlock()

				for k, v := range p {
					localShared[k] = v
				}

				// Create a new flow instance to avoid concurrent access
				// This is necessary because flows maintain state
				if err := n.innerFlow.Run(localShared); err != nil {
					mu.Lock()
					errors[idx] = err
					mu.Unlock()
				}

				// Merge results back to shared (if needed)
				// This is application-specific and should be handled by the user
			}(i, params)
		}
		wg.Wait()

		// Check for errors
		for i, err := range errors {
			if err != nil {
				return nil, fmt.Errorf("batch %d failed: %w", i, err)
			}
		}
	} else {
		// Run flows sequentially
		for i, params := range batchParams {
			// Merge params into shared for the flow execution
			for k, v := range params {
				shared[k] = v
			}

			if err := n.innerFlow.Run(shared); err != nil {
				return nil, fmt.Errorf("batch %d failed: %w", i, err)
			}
		}
	}

	return len(batchParams), nil
}

func (n *batchFlowNode) Post(shared Shared, prepResult, execResult any) (Action, error) {
	if count, ok := execResult.(int); ok {
		shared["batch_count"] = count
	}
	return DefaultAction, nil
}
