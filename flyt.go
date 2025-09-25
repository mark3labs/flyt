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
	"encoding/json"
	"fmt"
	"reflect"
	"sync"
	"time"
)

// SharedStore provides thread-safe access to shared data across nodes in a flow.
// It acts as a key-value store that can be safely accessed by multiple goroutines
// during concurrent batch processing.
//
// Example:
//
//	shared := flyt.NewSharedStore()
//	shared.Set("user_id", 123)
//	shared.Set("config", map[string]any{"timeout": 30})
//
//	if val, ok := shared.Get("user_id"); ok {
//	    userID := val.(int)
//	    // Use userID
//	}
type SharedStore struct {
	mu   sync.RWMutex
	data map[string]any
}

// NewSharedStore creates a new thread-safe shared store.
// The store is initialized empty and ready for use.
func NewSharedStore() *SharedStore {
	return &SharedStore{
		data: make(map[string]any),
	}
}

// Get retrieves a value from the store by key.
// Returns the value and true if the key exists, or nil and false if not found.
// This method is safe for concurrent access.
func (s *SharedStore) Get(key string) (any, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, ok := s.data[key]
	return val, ok
}

// Set stores a value in the store with the given key.
// If the key already exists, its value is overwritten.
// This method is safe for concurrent access.
func (s *SharedStore) Set(key string, value any) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = value
}

// GetAll returns a copy of all data in the store.
// The returned map is a shallow copy and can be safely modified
// without affecting the store's internal data.
// This method is safe for concurrent access.
func (s *SharedStore) GetAll() map[string]any {
	s.mu.RLock()
	defer s.mu.RUnlock()
	copy := make(map[string]any, len(s.data))
	for k, v := range s.data {
		copy[k] = v
	}
	return copy
}

// Merge merges another map into the store.
// Existing keys are overwritten with values from the provided map.
// If the provided map is nil, this method does nothing.
// This method is safe for concurrent access.
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

// Has checks if a key exists in the store.
// This method is safe for concurrent access.
func (s *SharedStore) Has(key string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.data[key]
	return ok
}

// Delete removes a key from the store.
// This method is safe for concurrent access.
func (s *SharedStore) Delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, key)
}

// Clear removes all keys from the store.
// This method is safe for concurrent access.
func (s *SharedStore) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data = make(map[string]any)
}

// Keys returns all keys in the store.
// The returned slice is a snapshot and can be safely modified
// without affecting the store's internal data.
// This method is safe for concurrent access.
func (s *SharedStore) Keys() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	keys := make([]string, 0, len(s.data))
	for k := range s.data {
		keys = append(keys, k)
	}
	return keys
}

// Len returns the number of items in the store.
// This method is safe for concurrent access.
func (s *SharedStore) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.data)
}

// GetString retrieves a string value from the store.
// Returns empty string if the key doesn't exist or the value is not a string.
// This method is safe for concurrent access.
func (s *SharedStore) GetString(key string) string {
	val, ok := s.Get(key)
	if !ok {
		return ""
	}
	str, _ := val.(string)
	return str
}

// GetStringOr retrieves a string value from the store.
// Returns the provided default value if the key doesn't exist or the value is not a string.
// This method is safe for concurrent access.
func (s *SharedStore) GetStringOr(key string, defaultVal string) string {
	val, ok := s.Get(key)
	if !ok {
		return defaultVal
	}
	str, ok := val.(string)
	if !ok {
		return defaultVal
	}
	return str
}

// GetInt retrieves an int value from the store.
// Returns 0 if the key doesn't exist or the value cannot be converted to int.
// Supports conversion from int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, and float types.
// This method is safe for concurrent access.
func (s *SharedStore) GetInt(key string) int {
	return s.GetIntOr(key, 0)
}

// GetIntOr retrieves an int value from the store.
// Returns the provided default value if the key doesn't exist or the value cannot be converted to int.
// Supports conversion from various numeric types.
// This method is safe for concurrent access.
func (s *SharedStore) GetIntOr(key string, defaultVal int) int {
	val, ok := s.Get(key)
	if !ok {
		return defaultVal
	}

	switch v := val.(type) {
	case int:
		return v
	case int8:
		return int(v)
	case int16:
		return int(v)
	case int32:
		return int(v)
	case int64:
		return int(v)
	case uint:
		return int(v)
	case uint8:
		return int(v)
	case uint16:
		return int(v)
	case uint32:
		return int(v)
	case uint64:
		return int(v)
	case float32:
		return int(v)
	case float64:
		return int(v)
	default:
		return defaultVal
	}
}

// GetFloat64 retrieves a float64 value from the store.
// Returns 0.0 if the key doesn't exist or the value cannot be converted to float64.
// Supports conversion from int, float32, and other numeric types.
// This method is safe for concurrent access.
func (s *SharedStore) GetFloat64(key string) float64 {
	return s.GetFloat64Or(key, 0.0)
}

// GetFloat64Or retrieves a float64 value from the store.
// Returns the provided default value if the key doesn't exist or the value cannot be converted to float64.
// Supports conversion from various numeric types.
// This method is safe for concurrent access.
func (s *SharedStore) GetFloat64Or(key string, defaultVal float64) float64 {
	val, ok := s.Get(key)
	if !ok {
		return defaultVal
	}

	switch v := val.(type) {
	case float64:
		return v
	case float32:
		return float64(v)
	case int:
		return float64(v)
	case int8:
		return float64(v)
	case int16:
		return float64(v)
	case int32:
		return float64(v)
	case int64:
		return float64(v)
	case uint:
		return float64(v)
	case uint8:
		return float64(v)
	case uint16:
		return float64(v)
	case uint32:
		return float64(v)
	case uint64:
		return float64(v)
	default:
		return defaultVal
	}
}

// GetBool retrieves a bool value from the store.
// Returns false if the key doesn't exist or the value is not a bool.
// This method is safe for concurrent access.
func (s *SharedStore) GetBool(key string) bool {
	return s.GetBoolOr(key, false)
}

// GetBoolOr retrieves a bool value from the store.
// Returns the provided default value if the key doesn't exist or the value is not a bool.
// This method is safe for concurrent access.
func (s *SharedStore) GetBoolOr(key string, defaultVal bool) bool {
	val, ok := s.Get(key)
	if !ok {
		return defaultVal
	}
	b, ok := val.(bool)
	if !ok {
		return defaultVal
	}
	return b
}

// GetSlice retrieves a []any slice from the store.
// Returns nil if the key doesn't exist or the value is not a slice.
// This method is safe for concurrent access.
func (s *SharedStore) GetSlice(key string) []any {
	return s.GetSliceOr(key, nil)
}

// GetSliceOr retrieves a []any slice from the store.
// Returns the provided default value if the key doesn't exist or the value is not a slice.
// Uses ToSlice to convert various slice types to []any.
// This method is safe for concurrent access.
func (s *SharedStore) GetSliceOr(key string, defaultVal []any) []any {
	val, ok := s.Get(key)
	if !ok {
		return defaultVal
	}

	// Use ToSlice for conversion which handles various slice types
	if val == nil {
		return defaultVal
	}

	// Check if it's already []any
	if slice, ok := val.([]any); ok {
		return slice
	}

	// Try to convert using reflection
	result := ToSlice(val)
	if len(result) == 1 && result[0] == val {
		// ToSlice wrapped a non-slice value, so this wasn't actually a slice
		return defaultVal
	}
	return result
}

// GetMap retrieves a map[string]any from the store.
// Returns nil if the key doesn't exist or the value is not a map[string]any.
// This method is safe for concurrent access.
func (s *SharedStore) GetMap(key string) map[string]any {
	return s.GetMapOr(key, nil)
}

// GetMapOr retrieves a map[string]any from the store.
// Returns the provided default value if the key doesn't exist or the value is not a map[string]any.
// This method is safe for concurrent access.
func (s *SharedStore) GetMapOr(key string, defaultVal map[string]any) map[string]any {
	val, ok := s.Get(key)
	if !ok {
		return defaultVal
	}
	m, ok := val.(map[string]any)
	if !ok {
		return defaultVal
	}
	return m
}

// Bind binds a value from the store to a struct using JSON marshaling/unmarshaling.
// This allows for easy conversion of complex types stored in the SharedStore.
// The destination must be a pointer to the target struct.
// Returns an error if the key is not found or binding fails.
// This method is safe for concurrent access.
//
// Example:
//
//	type User struct {
//	    ID   int    `json:"id"`
//	    Name string `json:"name"`
//	}
//	var user User
//	err := shared.Bind("user", &user)
func (s *SharedStore) Bind(key string, dest any) error {
	val, ok := s.Get(key)
	if !ok {
		return fmt.Errorf("key %q not found in shared store", key)
	}

	// Check if dest is a pointer
	rv := reflect.ValueOf(dest)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return fmt.Errorf("destination must be a non-nil pointer")
	}

	// If val is already the correct type, assign directly
	valType := reflect.TypeOf(val)
	destType := rv.Type().Elem()
	if valType == destType {
		rv.Elem().Set(reflect.ValueOf(val))
		return nil
	}

	// Otherwise use JSON as intermediate format
	jsonBytes, err := json.Marshal(val)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	if err := json.Unmarshal(jsonBytes, dest); err != nil {
		return fmt.Errorf("failed to unmarshal to destination: %w", err)
	}

	return nil
}

// MustBind is like Bind but panics if binding fails.
// Use this only when binding failure should be considered a programming error.
// This method is safe for concurrent access.
//
// Example:
//
//	var config Config
//	shared.MustBind("config", &config)  // Panics if binding fails
func (s *SharedStore) MustBind(key string, dest any) {
	if err := s.Bind(key, dest); err != nil {
		panic(fmt.Sprintf("SharedStore.MustBind failed: %v", err))
	}
}

// Action represents the next action to take after a node executes.
// Actions are used to determine flow control in workflows, allowing
// nodes to direct execution to different paths based on their results.
//
// Example:
//
//	const (
//	    ActionSuccess Action = "success"
//	    ActionRetry   Action = "retry"
//	    ActionFail    Action = "fail"
//	)
type Action string

const (
	// DefaultAction is the default action if none is specified.
	// Flows use this when a node doesn't explicitly return an action.
	DefaultAction Action = "default"

	// KeyItems is the shared store key for items to be processed.
	// Batch nodes look for items under this key by default.
	KeyItems = "items"

	// KeyResults is the shared store key for processing results.
	// Batch nodes store their results under this key by default.
	KeyResults = "results"

	// KeyBatchCount is the shared store key for batch count.
	// Batch flows store the number of iterations under this key.
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

// BaseNode provides a base implementation of the Node interface.
// It includes common functionality like retry configuration and default
// implementations of Prep, Exec, and Post methods that can be overridden.
//
// BaseNode is designed to be embedded in custom node implementations:
//
//	type MyNode struct {
//	    *flyt.BaseNode
//	    // custom fields
//	}
//
//	func (n *MyNode) Exec(ctx context.Context, prepResult any) (any, error) {
//	    // custom implementation
//	}
type BaseNode struct {
	mu         sync.RWMutex
	maxRetries int
	wait       time.Duration

	// Batch configuration
	batchConcurrency   int    // 0 = sequential, >0 = concurrent with limit
	batchErrorHandling string // "stop", "continue"
}

// NewBaseNode creates a new BaseNode with the provided options.
// By default, maxRetries is set to 1 (no retries) and wait is 0.
//
// Example:
//
//	node := flyt.NewBaseNode(
//	    flyt.WithMaxRetries(3),
//	    flyt.WithWait(time.Second),
//	)
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

// NodeOption is a function that configures a BaseNode.
// Options can be passed to NewBaseNode to customize its behavior.
type NodeOption func(*BaseNode)

// WithMaxRetries sets the maximum number of retries for the Exec phase.
// The default is 1 (no retries). Setting this to a value greater than 1
// enables automatic retry on Exec failures.
//
// Example:
//
//	node := flyt.NewBaseNode(flyt.WithMaxRetries(3))
func WithMaxRetries(retries int) NodeOption {
	return func(n *BaseNode) {
		n.maxRetries = retries
	}
}

// WithWait sets the wait duration between retries.
// This only applies when maxRetries is greater than 1.
// The default is 0 (no wait between retries).
//
// Example:
//
//	node := flyt.NewBaseNode(
//	    flyt.WithMaxRetries(3),
//	    flyt.WithWait(time.Second * 2),
//	)
func WithWait(wait time.Duration) NodeOption {
	return func(n *BaseNode) {
		n.wait = wait
	}
}

// WithBatchConcurrency sets the concurrency level for batch processing.
// 0 means sequential processing, >0 means concurrent with the specified limit.
//
// Example:
//
//	node := flyt.NewBatchNode(flyt.WithBatchConcurrency(10))
func WithBatchConcurrency(n int) NodeOption {
	return func(node *BaseNode) {
		node.batchConcurrency = n
	}
}

// WithBatchErrorHandling sets the error handling strategy for batch processing.
// If continueOnError is true, processing continues even if some items fail.
// If false, processing stops on the first error.
//
// Example:
//
//	node := flyt.NewBatchNode(flyt.WithBatchErrorHandling(true))
func WithBatchErrorHandling(continueOnError bool) NodeOption {
	return func(node *BaseNode) {
		if continueOnError {
			node.batchErrorHandling = "continue"
		} else {
			node.batchErrorHandling = "stop"
		}
	}
}

// GetMaxRetries returns the maximum number of retries configured for this node.
// This method is thread-safe.
func (n *BaseNode) GetMaxRetries() int {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.maxRetries
}

// GetWait returns the wait duration between retries configured for this node.
// This method is thread-safe.
func (n *BaseNode) GetWait() time.Duration {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.wait
}

// GetBatchConcurrency returns the batch concurrency level configured for this node.
// This method is thread-safe.
func (n *BaseNode) GetBatchConcurrency() int {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.batchConcurrency
}

// GetBatchErrorHandling returns the batch error handling strategy configured for this node.
// This method is thread-safe.
func (n *BaseNode) GetBatchErrorHandling() string {
	n.mu.RLock()
	defer n.mu.RUnlock()
	if n.batchErrorHandling == "" {
		return "continue" // default
	}
	return n.batchErrorHandling
}

// Prep is the default prep implementation that returns nil.
// Override this method in your node implementation to read and preprocess
// data from the SharedStore before execution.
func (n *BaseNode) Prep(ctx context.Context, shared *SharedStore) (any, error) {
	return nil, nil
}

// Exec is the default exec implementation that returns nil.
// This method should be overridden in your node implementation to provide
// the main processing logic.
func (n *BaseNode) Exec(ctx context.Context, prepResult any) (any, error) {
	return nil, nil
}

// Post is the default post implementation that returns DefaultAction.
// Override this method in your node implementation to process results
// and determine the next action in the flow.
func (n *BaseNode) Post(ctx context.Context, shared *SharedStore, prepResult, execResult any) (Action, error) {
	return DefaultAction, nil
}

// ExecFallback handles errors after all retries are exhausted.
// The default implementation simply returns the error.
// Override this method to provide custom fallback behavior.
func (n *BaseNode) ExecFallback(prepResult any, err error) (any, error) {
	return nil, err
}

// RetryableNode is a node that supports automatic retries on Exec failures.
// Nodes implementing this interface can specify retry behavior through
// GetMaxRetries and GetWait methods.
type RetryableNode interface {
	Node
	GetMaxRetries() int
	GetWait() time.Duration
}

// FallbackNode is a node that supports custom fallback behavior on error.
// When all retries are exhausted, the ExecFallback method is called to
// provide an alternative result or handle the error gracefully.
type FallbackNode interface {
	ExecFallback(prepResult any, err error) (any, error)
}

// Run executes a node with the standard prep->exec->post lifecycle.
// This is the main entry point for executing individual nodes.
//
// The execution flow is:
//  1. Prep: Read and preprocess data from SharedStore
//  2. Exec: Execute main logic with automatic retries if configured
//  3. Post: Process results and write back to SharedStore
//
// If the node implements RetryableNode, the Exec phase will automatically
// retry on failure according to the configured settings.
//
// If the node implements FallbackNode and all retries fail, ExecFallback
// is called to provide alternative handling.
//
// Parameters:
//   - ctx: Context for cancellation and timeouts
//   - node: The node to execute
//   - shared: SharedStore for data exchange
//
// Returns:
//   - Action: The next action to take in the flow
//   - error: Any error that occurred during execution
//
// Example:
//
//	node := &MyNode{BaseNode: flyt.NewBaseNode()}
//	shared := flyt.NewSharedStore()
//	shared.Set("input", "data")
//
//	action, err := flyt.Run(ctx, node, shared)
//	if err != nil {
//	    log.Fatal(err)
//	}
func Run(ctx context.Context, node Node, shared *SharedStore) (Action, error) {
	// Check if this is a BatchNode
	if _, ok := node.(*BatchNode); ok {
		return runBatch(ctx, node, shared)
	}
	if batchBuilder, ok := node.(*BatchNodeBuilder); ok {
		return runBatch(ctx, batchBuilder.BatchNode, shared)
	}

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

// Flow represents a workflow of connected nodes.
// A Flow is itself a Node, allowing flows to be nested within other flows.
// Nodes are connected via actions, creating a directed graph of execution.
//
// Example:
//
//	// Create nodes
//	validateNode := &ValidateNode{BaseNode: flyt.NewBaseNode()}
//	processNode := &ProcessNode{BaseNode: flyt.NewBaseNode()}
//	errorNode := &ErrorNode{BaseNode: flyt.NewBaseNode()}
//
//	// Create flow
//	flow := flyt.NewFlow(validateNode)
//	flow.Connect(validateNode, ActionSuccess, processNode)
//	flow.Connect(validateNode, ActionFail, errorNode)
//
//	// Or with chaining:
//	flow.Connect(validateNode, ActionSuccess, processNode)
//	     .Connect(validateNode, ActionFail, errorNode)
//
//	// Run flow
//	err := flow.Run(ctx, shared)
type Flow struct {
	*BaseNode
	start       Node
	transitions map[Node]map[Action]Node
}

// NewFlow creates a new Flow with a start node.
// The flow begins execution at the start node and follows
// transitions based on the actions returned by each node.
//
// Parameters:
//   - start: The first node to execute in the flow
//
// Example:
//
//	flow := flyt.NewFlow(startNode)
func NewFlow(start Node) *Flow {
	return &Flow{
		BaseNode:    NewBaseNode(),
		start:       start,
		transitions: make(map[Node]map[Action]Node),
	}
}

// Connect adds a transition from one node to another based on an action.
// When the 'from' node returns the specified action, execution continues
// with the 'to' node. Multiple actions can be connected from a single node.
// Returns the flow instance for method chaining.
//
// Parameters:
//   - from: The source node
//   - action: The action that triggers this transition
//   - to: The destination node
//
// Example:
//
//	flow.Connect(nodeA, "success", nodeB)
//	flow.Connect(nodeA, "retry", nodeA)  // Self-loop for retry
//	flow.Connect(nodeA, "fail", errorNode)
//
// Example with chaining:
//
//	flow.Connect(nodeA, "success", nodeB)
//	     .Connect(nodeB, "success", finalNode)
//	     .Connect(nodeB, "fail", errorNode)
func (f *Flow) Connect(from Node, action Action, to Node) *Flow {
	if f.transitions[from] == nil {
		f.transitions[from] = make(map[Action]Node)
	}
	f.transitions[from][action] = to
	return f
}

// Run executes the flow starting from the start node.
// This is a convenience method that wraps the standard Run function.
// The flow executes nodes in sequence based on their action transitions
// until no more transitions are available or an error occurs.
//
// Parameters:
//   - ctx: Context for cancellation and timeouts
//   - shared: SharedStore for data exchange between nodes
//
// Returns:
//   - error: Any error that occurred during flow execution
func (f *Flow) Run(ctx context.Context, shared *SharedStore) error {
	// Use the standard Run function to execute this flow as a node
	_, err := Run(ctx, f, shared)
	return err
}

// Prep implements Node interface for Flow
func (f *Flow) Prep(ctx context.Context, shared *SharedStore) (any, error) {
	// Pass the shared store to Exec
	return shared, nil
}

// Exec orchestrates the flow execution by running nodes in sequence
func (f *Flow) Exec(ctx context.Context, prepResult any) (any, error) {
	// Extract shared store from prepResult
	shared, ok := prepResult.(*SharedStore)
	if !ok {
		return nil, fmt.Errorf("flow: exec failed: invalid prepResult type %T, expected *SharedStore", prepResult)
	}

	if f.start == nil {
		return nil, fmt.Errorf("flow: exec failed: no start node configured")
	}

	current := f.start
	var lastAction Action

	for current != nil {
		// Check context
		if err := ctx.Err(); err != nil {
			return nil, fmt.Errorf("flow: exec cancelled: %w", err)
		}

		// Run the current node using the standard Run function
		action, err := Run(ctx, current, shared)
		if err != nil {
			return nil, err
		}

		lastAction = action

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

	// Return the last action as the result
	return lastAction, nil
}

// Post implements Node interface for Flow
func (f *Flow) Post(ctx context.Context, shared *SharedStore, prepResult, execResult any) (Action, error) {
	// Return the last action from the flow execution
	if action, ok := execResult.(Action); ok {
		return action, nil
	}
	return DefaultAction, nil
}

// WorkerPool manages concurrent task execution with a fixed number of workers.
// It provides a simple way to limit concurrency and execute tasks in parallel.
// WorkerPool is used internally by batch processing nodes but can also be
// used directly for custom concurrent operations.
//
// Example:
//
//	pool := flyt.NewWorkerPool(5)
//	defer pool.Close()
//
//	for _, item := range items {
//	    item := item // capture loop variable
//	    pool.Submit(func() {
//	        // Process item
//	    })
//	}
//
//	pool.Wait()
type WorkerPool struct {
	workers int
	tasks   chan func()
	wg      sync.WaitGroup
	done    chan struct{}
}

// NewWorkerPool creates a new worker pool with the specified number of workers.
// If workers is less than or equal to 0, it defaults to 1.
// The pool starts workers immediately and is ready to accept tasks.
//
// Parameters:
//   - workers: Number of concurrent workers
//
// Returns:
//   - *WorkerPool: A new worker pool ready for use
//
// Remember to call Close() when done to clean up resources.
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

// Submit submits a task to the pool for execution.
// The task will be executed by one of the available workers.
// This method blocks if all workers are busy and the task buffer is full.
//
// Parameters:
//   - task: Function to execute
func (p *WorkerPool) Submit(task func()) {
	p.wg.Add(1)
	p.tasks <- func() {
		defer p.wg.Done()
		task()
	}
}

// Wait waits for all submitted tasks to complete.
// This method blocks until all tasks have finished executing.
func (p *WorkerPool) Wait() {
	p.wg.Wait()
}

// Close closes the worker pool and waits for all workers to finish.
// After calling Close, no new tasks can be submitted.
// This method should be called when the pool is no longer needed
// to properly clean up resources.
func (p *WorkerPool) Close() {
	close(p.done)
	close(p.tasks)
}

// ToSlice converts various types to []any for batch processing.
// This utility function handles common slice types and single values,
// making it easier to work with batch operations.
//
// Supported types:
//   - []any: returned as-is
//   - []string, []int, []float64: converted to []any
//   - []map[string]any: converted to []any
//   - Other slice types: converted using reflection
//   - Single values: wrapped in a slice
//   - nil: returns empty slice
//
// Example:
//
//	items1 := flyt.ToSlice([]string{"a", "b", "c"})  // []any{"a", "b", "c"}
//	items2 := flyt.ToSlice("single")                 // []any{"single"}
//	items3 := flyt.ToSlice(nil)                      // []any{}
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

// FlowFactory creates new instances of a flow.
// This is used by batch flow operations to create isolated flow instances
// for concurrent execution. Each call should return a new flow instance
// to avoid race conditions.
//
// Example:
//
//	factory := func() *flyt.Flow {
//	    node1 := &ProcessNode{BaseNode: flyt.NewBaseNode()}
//	    node2 := &SaveNode{BaseNode: flyt.NewBaseNode()}
//
//	    flow := flyt.NewFlow(node1)
//	    flow.Connect(node1, flyt.DefaultAction, node2)
//	    return flow
//	}
type FlowFactory func() *Flow

// CustomNode is a node implementation that uses custom functions
// for Prep, Exec, and Post phases. This allows creating nodes
// without defining new types, useful for simple operations.
//
// CustomNode supports two styles of functions:
//
//  1. Result-based functions (WithPrepFunc, WithExecFunc, WithPostFunc):
//     These use the Result type which provides convenient type assertion methods.
//     Best for complex data processing where type safety and conversion helpers are valuable.
//
//  2. Any-based functions (WithPrepFuncAny, WithExecFuncAny, WithPostFuncAny):
//     These use standard any types matching the Node interface directly.
//     Best for simple operations or when migrating existing code.
//
// Example with Result types:
//
//	node := flyt.NewNode(
//	    flyt.WithExecFunc(func(ctx context.Context, prepResult flyt.Result) (flyt.Result, error) {
//	        data := prepResult.MustString()  // Type-safe access
//	        return flyt.NewResult(processData(data)), nil
//	    }),
//	)
//
// Example with any types:
//
//	node := flyt.NewNode(
//	    flyt.WithExecFuncAny(func(ctx context.Context, prepResult any) (any, error) {
//	        data := prepResult.(string)  // Manual type assertion
//	        return processData(data), nil
//	    }),
//	)
type CustomNode struct {
	*BaseNode
	prepFunc         func(context.Context, *SharedStore) (Result, error)
	execFunc         func(context.Context, Result) (Result, error)
	postFunc         func(context.Context, *SharedStore, Result, Result) (Action, error)
	execFallbackFunc func(any, error) (any, error)
}

// Prep implements Node.Prep by calling the custom prepFunc if provided
func (n *CustomNode) Prep(ctx context.Context, shared *SharedStore) (any, error) {
	if n.prepFunc != nil {
		result, err := n.prepFunc(ctx, shared)
		if err != nil {
			return nil, err
		}
		return result.Value(), nil
	}
	return n.BaseNode.Prep(ctx, shared)
}

// Exec implements Node.Exec by calling the custom execFunc if provided
func (n *CustomNode) Exec(ctx context.Context, prepResult any) (any, error) {
	if n.execFunc != nil {
		// Check if prepResult is already a Result (when called from batch processing)
		if r, ok := prepResult.(Result); ok {
			return n.execFunc(ctx, r)
		}
		return n.execFunc(ctx, NewResult(prepResult))
	}
	return n.BaseNode.Exec(ctx, prepResult)
}

// Post implements Node.Post by calling the custom postFunc if provided
func (n *CustomNode) Post(ctx context.Context, shared *SharedStore, prepResult, execResult any) (Action, error) {
	if n.postFunc != nil {
		return n.postFunc(ctx, shared, NewResult(prepResult), NewResult(execResult))
	}
	return n.BaseNode.Post(ctx, shared, prepResult, execResult)
}

// ExecFallback implements FallbackNode.ExecFallback by calling the custom execFallbackFunc if provided
func (n *CustomNode) ExecFallback(prepResult any, err error) (any, error) {
	if n.execFallbackFunc != nil {
		return n.execFallbackFunc(prepResult, err)
	}
	return n.BaseNode.ExecFallback(prepResult, err)
}

// CustomNodeOption is an option for configuring a CustomNode.
// It allows setting custom implementations for Prep, Exec, and Post methods.

// NewNode creates a new node with custom function implementations.
// Returns a NodeBuilder that implements Node and provides chainable configuration methods.
//
// Can be used in two styles:
//
// Traditional (backwards compatible):
//
//	node := flyt.NewNode(
//	    flyt.WithMaxRetries(3),
//	    flyt.WithExecFunc(execFn),
//	)
//
// Builder pattern:
//
//	node := flyt.NewNode().
//	    WithMaxRetries(3).
//	    WithExecFunc(execFn).
//	    Build()  // Build() is optional
func NewNode(opts ...any) *NodeBuilder {
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

	return &NodeBuilder{CustomNode: node}
}

// CustomNodeOption is an option for configuring a CustomNode.
// It allows setting custom implementations for Prep, Exec, and Post methods.
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

// WithPrepFunc sets a custom Prep implementation for a CustomNode.
// The provided function will be called during the Prep phase to read
// and preprocess data from the SharedStore.
//
// Example:
//
//	flyt.WithPrepFunc(func(ctx context.Context, shared *flyt.SharedStore) (flyt.Result, error) {
//	    data, _ := shared.Get("input")
//	    // Preprocess data
//	    return flyt.NewResult(preprocessedData), nil
//	})
func WithPrepFunc(fn func(context.Context, *SharedStore) (Result, error)) CustomNodeOption {
	return &customNodeOption{
		f: func(n *CustomNode) {
			n.prepFunc = fn
		},
	}
}

// WithExecFunc sets a custom Exec implementation for a CustomNode.
// The provided function will be called during the Exec phase with
// the result from Prep as input.
//
// Example:
//
//	flyt.WithExecFunc(func(ctx context.Context, prepResult flyt.Result) (flyt.Result, error) {
//	    // Process the data
//	    return flyt.NewResult(processedResult), nil
//	})
func WithExecFunc(fn func(context.Context, Result) (Result, error)) CustomNodeOption {
	return &customNodeOption{
		f: func(n *CustomNode) {
			n.execFunc = fn
		},
	}
}

// WithPostFunc sets a custom Post implementation for a CustomNode.
// The provided function will be called during the Post phase to process
// results and determine the next action.
//
// Example:
//
//	flyt.WithPostFunc(func(ctx context.Context, shared *flyt.SharedStore, prepResult, execResult flyt.Result) (flyt.Action, error) {
//	    shared.Set("output", execResult.Value())
//	    if success {
//	        return "success", nil
//	    }
//	    return "retry", nil
//	})
func WithPostFunc(fn func(context.Context, *SharedStore, Result, Result) (Action, error)) CustomNodeOption {
	return &customNodeOption{
		f: func(n *CustomNode) {
			n.postFunc = fn
		},
	}
}

// WithExecFallbackFunc sets a custom ExecFallback implementation for a CustomNode.
// The provided function will be called when Exec fails after all retries are exhausted.
// This allows for custom error handling or returning a default value.
//
// Example:
//
//	flyt.WithExecFallbackFunc(func(prepResult any, err error) (any, error) {
//	    log.Printf("Exec failed after retries: %v", err)
//	    // Return nil to indicate failure, which can be handled in Post
//	    return nil, nil
//	})
func WithExecFallbackFunc(fn func(any, error) (any, error)) CustomNodeOption {
	return &customNodeOption{
		f: func(n *CustomNode) {
			n.execFallbackFunc = fn
		},
	}
}

// WithPrepFuncAny sets a custom Prep implementation for a CustomNode using any types.
// This is an alternative to WithPrepFunc that doesn't require Result types,
// useful for simpler cases or when migrating existing code.
//
// Example:
//
//	flyt.WithPrepFuncAny(func(ctx context.Context, shared *flyt.SharedStore) (any, error) {
//	    data, _ := shared.Get("input")
//	    // Preprocess data
//	    return preprocessedData, nil
//	})
func WithPrepFuncAny(fn func(context.Context, *SharedStore) (any, error)) CustomNodeOption {
	return &customNodeOption{
		f: func(n *CustomNode) {
			n.prepFunc = func(ctx context.Context, shared *SharedStore) (Result, error) {
				val, err := fn(ctx, shared)
				if err != nil {
					return Result{}, err
				}
				return NewResult(val), nil
			}
		},
	}
}

// WithExecFuncAny sets a custom Exec implementation for a CustomNode using any types.
// This is an alternative to WithExecFunc that doesn't require Result types,
// useful for simpler cases or when you don't need the type assertion helpers.
//
// Example:
//
//	flyt.WithExecFuncAny(func(ctx context.Context, prepResult any) (any, error) {
//	    // Process the data
//	    return processedResult, nil
//	})
func WithExecFuncAny(fn func(context.Context, any) (any, error)) CustomNodeOption {
	return &customNodeOption{
		f: func(n *CustomNode) {
			n.execFunc = func(ctx context.Context, prepResult Result) (Result, error) {
				val, err := fn(ctx, prepResult.Value())
				if err != nil {
					return Result{}, err
				}
				return NewResult(val), nil
			}
		},
	}
}

// WithPostFuncAny sets a custom Post implementation for a CustomNode using any types.
// This is an alternative to WithPostFunc that doesn't require Result types,
// useful for simpler cases or backward compatibility.
//
// Example:
//
//	flyt.WithPostFuncAny(func(ctx context.Context, shared *flyt.SharedStore, prepResult, execResult any) (flyt.Action, error) {
//	    shared.Set("output", execResult)
//	    if success {
//	        return "success", nil
//	    }
//	    return "retry", nil
//	})
func WithPostFuncAny(fn func(context.Context, *SharedStore, any, any) (Action, error)) CustomNodeOption {
	return &customNodeOption{
		f: func(n *CustomNode) {
			n.postFunc = func(ctx context.Context, shared *SharedStore, prepResult, execResult Result) (Action, error) {
				return fn(ctx, shared, prepResult.Value(), execResult.Value())
			}
		},
	}
}
