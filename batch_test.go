package flyt

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// TestBatchErrorAggregation tests the BatchError type
func TestBatchErrorAggregation(t *testing.T) {
	// Test empty errors
	be := &BatchError{Errors: []error{}}
	if be.Error() != "batch: no errors recorded" {
		t.Errorf("unexpected error message: %s", be.Error())
	}

	// Test single error
	be = &BatchError{Errors: []error{errors.New("test error")}}
	expected := "batch: test error"
	if be.Error() != expected {
		t.Errorf("expected %q, got %q", expected, be.Error())
	}

	// Test multiple errors
	be = &BatchError{
		Errors: []error{
			errors.New("error 1"),
			errors.New("error 2"),
			errors.New("error 3"),
		},
	}
	expected = "batch: 3 errors occurred, first: error 1"
	if be.Error() != expected {
		t.Errorf("expected %q, got %q", expected, be.Error())
	}
}

// TestBatchNodeValidation tests batch node validation
func TestBatchNodeValidation(t *testing.T) {
	processFunc := func(ctx context.Context, item any) (any, error) {
		return item, nil
	}

	config := &BatchConfig{
		MaxBatchSize: 5,
	}

	node := NewBatchNodeWithConfig(processFunc, false, config)
	shared := NewSharedStore()

	// Test exceeding max batch size
	items := make([]int, 10)
	for i := range items {
		items[i] = i
	}
	shared.Set("items", items)

	ctx := context.Background()
	_, err := Run(ctx, node, shared)
	if err == nil {
		t.Error("expected batch size validation error")
	}

	expectedErr := "batchNode: exec failed: batch size 10 exceeds maximum 5"
	if err.Error() != fmt.Sprintf("run: exec failed after 1 retries: %s", expectedErr) {
		t.Errorf("unexpected error: %v", err)
	}
}

// TestBatchNodeConcurrency tests concurrent batch processing
func TestBatchNodeConcurrency(t *testing.T) {
	var counter int32
	var maxConcurrent int32

	processFunc := func(ctx context.Context, item any) (any, error) {
		current := atomic.AddInt32(&counter, 1)
		defer atomic.AddInt32(&counter, -1)

		// Track max concurrent executions
		for {
			max := atomic.LoadInt32(&maxConcurrent)
			if current <= max || atomic.CompareAndSwapInt32(&maxConcurrent, max, current) {
				break
			}
		}

		// Simulate work
		time.Sleep(10 * time.Millisecond)
		return item, nil
	}

	config := &BatchConfig{
		MaxConcurrency: 3,
	}

	node := NewBatchNodeWithConfig(processFunc, true, config)
	shared := NewSharedStore()
	shared.Set("items", []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})

	ctx := context.Background()
	_, err := Run(ctx, node, shared)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Check max concurrency was respected
	if maxConcurrent > 3 {
		t.Errorf("max concurrency exceeded: %d > 3", maxConcurrent)
	}

	// Check results
	results, _ := shared.Get("results")
	if results == nil {
		t.Error("results not found")
	}
}

// TestBatchNodeSequential tests sequential batch processing
func TestBatchNodeSequential(t *testing.T) {
	var order []int
	var mu sync.Mutex

	processFunc := func(ctx context.Context, item any) (any, error) {
		mu.Lock()
		order = append(order, item.(int))
		mu.Unlock()
		return item, nil
	}

	node := NewBatchNode(processFunc, false) // sequential
	shared := NewSharedStore()
	shared.Set("items", []int{1, 2, 3, 4, 5})

	ctx := context.Background()
	_, err := Run(ctx, node, shared)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Check order is preserved
	for i, v := range order {
		if v != i+1 {
			t.Errorf("order[%d]: expected %d, got %d", i, i+1, v)
		}
	}
}

// TestBatchNodeCustomKeys tests batch node with custom keys
func TestBatchNodeCustomKeys(t *testing.T) {
	processFunc := func(ctx context.Context, item any) (any, error) {
		return fmt.Sprintf("processed_%v", item), nil
	}

	// Create batch node with custom keys
	node := NewBatchNodeWithKeys(processFunc, false, "my_items", "my_results")

	shared := NewSharedStore()
	shared.Set("my_items", []string{"a", "b", "c"})

	ctx := context.Background()
	_, err := Run(ctx, node, shared)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Check results are stored with custom key
	results, ok := shared.Get("my_results")
	if !ok {
		t.Error("results not found with custom key")
	}

	resultSlice, ok := results.([]any)
	if !ok || len(resultSlice) != 3 {
		t.Error("unexpected results format")
	}

	// Verify default keys were not used
	if _, ok := shared.Get(KeyResults); ok {
		t.Error("results should not be stored with default key")
	}
}

// TestBatchFlowConcurrency tests concurrent batch flow execution
func TestBatchFlowConcurrency(t *testing.T) {
	var executionCount int32

	// Factory creates new flow instances
	flowFactory := func() *Flow {
		node := &CounterNode{
			BaseNode: NewBaseNode(),
			counter:  &executionCount,
		}
		return NewFlow(node)
	}

	// Batch function returns parameters for each iteration
	batchFunc := func(ctx context.Context, shared *SharedStore) ([]Params, error) {
		params := make([]Params, 5)
		for i := range params {
			params[i] = Params{"index": i}
		}
		return params, nil
	}

	config := &BatchConfig{
		MaxConcurrency: 3,
	}

	batchFlow := NewBatchFlowWithConfig(flowFactory, batchFunc, true, config)
	shared := NewSharedStore()
	ctx := context.Background()

	err := batchFlow.Run(ctx, shared)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Check all flows executed
	if atomic.LoadInt32(&executionCount) != 5 {
		t.Errorf("expected 5 executions, got %d", executionCount)
	}

	// Check batch count was stored
	count, ok := shared.Get(KeyBatchCount)
	if !ok {
		t.Error("batch_count not found in shared store")
	}
	if count != 5 {
		t.Errorf("expected batch_count 5, got %v", count)
	}
}

// TestBatchFlowCustomCountKey tests batch flow with custom count key
func TestBatchFlowCustomCountKey(t *testing.T) {
	// Factory creates new flow instances
	flowFactory := func() *Flow {
		node := &TestNode{
			BaseNode: NewBaseNode(),
			name:     "test",
			action:   DefaultAction,
		}
		return NewFlow(node)
	}

	// Batch function returns parameters
	batchFunc := func(ctx context.Context, shared *SharedStore) ([]Params, error) {
		return []Params{
			{"id": 1},
			{"id": 2},
			{"id": 3},
		}, nil
	}

	// Create batch flow with custom count key
	batchFlow := NewBatchFlowWithCountKey(flowFactory, batchFunc, false, "my_batch_count")

	shared := NewSharedStore()
	ctx := context.Background()

	err := batchFlow.Run(ctx, shared)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Check custom count key was used
	count, ok := shared.Get("my_batch_count")
	if !ok {
		t.Error("custom count key not found in shared store")
	}
	if count != 3 {
		t.Errorf("expected count 3, got %v", count)
	}

	// Verify default key was not used
	if _, ok := shared.Get(KeyBatchCount); ok {
		t.Error("count should not be stored with default key")
	}
}

// CounterNode counts executions
type CounterNode struct {
	*BaseNode
	counter *int32
}

func (n *CounterNode) Exec(ctx context.Context, prepResult any) (any, error) {
	atomic.AddInt32(n.counter, 1)
	return nil, nil
}

// BenchmarkBatchNodeProcessing benchmarks batch processing
func BenchmarkBatchNodeProcessing(b *testing.B) {
	processFunc := func(ctx context.Context, item any) (any, error) {
		// Simulate some work
		time.Sleep(time.Microsecond)
		return item, nil
	}

	node := NewBatchNode(processFunc, true)
	shared := NewSharedStore()

	// Create test data
	items := make([]int, 100)
	for i := range items {
		items[i] = i
	}
	shared.Set("items", items)

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Run(ctx, node, shared)
		if err != nil {
			b.Fatal(err)
		}
	}
}
