package flyt

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"
)

// TestSharedStoreConcurrency tests concurrent access to SharedStore
func TestSharedStoreConcurrency(t *testing.T) {
	store := NewSharedStore()
	var wg sync.WaitGroup

	// Concurrent writes
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			store.Set(fmt.Sprintf("key%d", n), n)
		}(i)
	}

	// Concurrent reads
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			store.Get(fmt.Sprintf("key%d", n))
		}(i)
	}

	wg.Wait()

	// Verify all values
	for i := 0; i < 100; i++ {
		val, ok := store.Get(fmt.Sprintf("key%d", i))
		if !ok {
			t.Errorf("key%d not found", i)
		}
		if val != i {
			t.Errorf("key%d: expected %d, got %v", i, i, val)
		}
	}
}

// TestSharedStoreGetAll tests GetAll returns a copy
func TestSharedStoreGetAll(t *testing.T) {
	store := NewSharedStore()
	store.Set("key1", "value1")
	store.Set("key2", "value2")

	// Get all data
	data := store.GetAll()

	// Modify the returned map
	data["key1"] = "modified"
	data["key3"] = "new"

	// Original store should be unchanged
	val, _ := store.Get("key1")
	if val != "value1" {
		t.Errorf("store was modified: expected value1, got %v", val)
	}

	_, ok := store.Get("key3")
	if ok {
		t.Error("store has key3 which should not exist")
	}
}

// TestSharedStoreMerge tests the Merge method
func TestSharedStoreMerge(t *testing.T) {
	store := NewSharedStore()
	store.Set("key1", "value1")
	store.Set("key2", "value2")

	// Merge new data
	newData := map[string]any{
		"key2": "updated",
		"key3": "new",
	}
	store.Merge(newData)

	// Check results
	val1, _ := store.Get("key1")
	if val1 != "value1" {
		t.Errorf("key1: expected value1, got %v", val1)
	}

	val2, _ := store.Get("key2")
	if val2 != "updated" {
		t.Errorf("key2: expected updated, got %v", val2)
	}

	val3, _ := store.Get("key3")
	if val3 != "new" {
		t.Errorf("key3: expected new, got %v", val3)
	}
}

// SlowNode is a test node that takes time to execute
type SlowNode struct {
	*BaseNode
	delay time.Duration
}

func (n *SlowNode) Exec(ctx context.Context, prepResult any) (any, error) {
	select {
	case <-time.After(n.delay):
		return "completed", nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// TestContextCancellation tests context cancellation
func TestContextCancellation(t *testing.T) {
	node := &SlowNode{
		BaseNode: NewBaseNode(),
		delay:    500 * time.Millisecond,
	}
	shared := NewSharedStore()

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	_, err := Run(ctx, node, shared)
	if err == nil {
		t.Error("expected context cancellation error")
	}

	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("expected DeadlineExceeded, got %v", err)
	}
}

// TestContextCancellationDuringRetry tests cancellation during retry wait
func TestContextCancellationDuringRetry(t *testing.T) {
	attempts := 0
	node := &RetryNode{
		BaseNode: NewBaseNode(
			WithMaxRetries(3),
			WithWait(200*time.Millisecond),
		),
		onExec: func() error {
			attempts++
			return errors.New("retry me")
		},
	}

	shared := NewSharedStore()
	ctx, cancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
	defer cancel()

	_, err := Run(ctx, node, shared)
	if err == nil {
		t.Error("expected context cancellation error")
	}

	// Should have attempted once, then cancelled during wait
	if attempts != 1 {
		t.Errorf("expected 1 attempt, got %d", attempts)
	}
}

// RetryNode is a test node for retry logic
type RetryNode struct {
	*BaseNode
	attempts int
	onExec   func() error
}

func (n *RetryNode) Exec(ctx context.Context, prepResult any) (any, error) {
	n.attempts++
	if n.onExec != nil {
		return nil, n.onExec()
	}
	if n.attempts < 3 {
		return nil, errors.New("retry me")
	}
	return "success", nil
}

// TestRetryLogic tests retry functionality
func TestRetryLogic(t *testing.T) {
	node := &RetryNode{
		BaseNode: NewBaseNode(
			WithMaxRetries(3),
			WithWait(10*time.Millisecond),
		),
	}

	shared := NewSharedStore()
	ctx := context.Background()

	action, err := Run(ctx, node, shared)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if action != DefaultAction {
		t.Errorf("expected DefaultAction, got %v", action)
	}

	if node.attempts != 3 {
		t.Errorf("expected 3 attempts, got %d", node.attempts)
	}
}

// TestBatchErrorAggregation tests error aggregation in batch operations
func TestFlowExecution(t *testing.T) {
	// Create nodes
	node1 := &TestNode{
		BaseNode: NewBaseNode(),
		name:     "node1",
		action:   "next",
	}
	node2 := &TestNode{
		BaseNode: NewBaseNode(),
		name:     "node2",
		action:   DefaultAction,
	}

	// Create flow
	flow := NewFlow(node1)
	flow.Connect(node1, "next", node2)

	// Execute
	shared := NewSharedStore()
	ctx := context.Background()

	err := flow.Run(ctx, shared)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Check both nodes executed
	if !node1.executed {
		t.Error("node1 not executed")
	}
	if !node2.executed {
		t.Error("node2 not executed")
	}
}

// TestNode is a simple test node
type TestNode struct {
	*BaseNode
	name     string
	action   Action
	executed bool
	postFunc func(context.Context, *SharedStore, any, any) (Action, error)
}

func (n *TestNode) Exec(ctx context.Context, prepResult any) (any, error) {
	n.executed = true
	return n.name, nil
}

func (n *TestNode) Post(ctx context.Context, shared *SharedStore, prepResult, execResult any) (Action, error) {
	if n.postFunc != nil {
		return n.postFunc(ctx, shared, prepResult, execResult)
	}
	shared.Set(n.name+"_result", execResult)
	return n.action, nil
}


// TestBatchFlowCustomCountKey tests batch flow with custom count key
// TestFlowParameterPropagation tests parameter propagation in flows
func TestFlowParameterPropagation(t *testing.T) {
	node := &ParamCheckNode{
		BaseNode: NewBaseNode(),
	}

	flow := NewFlow(node)
	flow.SetParams(Params{"flow_param": "value"})

	shared := NewSharedStore()
	ctx := context.Background()

	err := flow.Run(ctx, shared)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Check that node received flow parameters
	nodeParams := node.GetParams()
	if nodeParams["flow_param"] != "value" {
		t.Error("flow parameters not propagated to node")
	}
}

// ParamCheckNode verifies it receives parameters
type ParamCheckNode struct {
	*BaseNode
}

func (n *ParamCheckNode) Exec(ctx context.Context, prepResult any) (any, error) {
	params := n.GetParams()
	if params == nil || len(params) == 0 {
		return nil, errors.New("no parameters received")
	}
	return nil, nil
}

// TestNestedFlowExecution tests flows within flows
func TestNestedFlowExecution(t *testing.T) {
	// Create inner flow
	innerNode := &TestNode{
		BaseNode: NewBaseNode(),
		name:     "inner",
		action:   DefaultAction,
	}
	innerFlow := NewFlow(innerNode)

	// Create outer flow
	outerNode := &TestNode{
		BaseNode: NewBaseNode(),
		name:     "outer",
		action:   "next",
	}
	outerFlow := NewFlow(outerNode)
	outerFlow.Connect(outerNode, "next", innerFlow)

	// Execute
	shared := NewSharedStore()
	ctx := context.Background()

	err := outerFlow.Run(ctx, shared)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Check both nodes executed
	if !outerNode.executed {
		t.Error("outer node not executed")
	}
	if !innerNode.executed {
		t.Error("inner node not executed")
	}
}

// BenchmarkSharedStoreConcurrentAccess benchmarks concurrent store access
func BenchmarkSharedStoreConcurrentAccess(b *testing.B) {
	store := NewSharedStore()

	// Pre-populate store
	for i := 0; i < 100; i++ {
		store.Set(fmt.Sprintf("key%d", i), i)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			if i%2 == 0 {
				store.Get(fmt.Sprintf("key%d", i%100))
			} else {
				store.Set(fmt.Sprintf("key%d", i%100), i)
			}
			i++
		}
	})
}

// BenchmarkBatchNodeProcessing benchmarks batch processing
