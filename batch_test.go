package flyt

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"
	"testing"
	"time"
)

func TestSimpleBatchNode(t *testing.T) {
	// Test the new simplified batch node
	var processedCount int32

	node := NewBatchNode().
		WithPrepFunc(func(ctx context.Context, shared *SharedStore) ([]Result, error) {
			// Return items to process
			return []Result{
				NewResult("item1"),
				NewResult("item2"),
				NewResult("item3"),
			}, nil
		}).
		WithExecFunc(func(ctx context.Context, item Result) (Result, error) {
			atomic.AddInt32(&processedCount, 1)
			// Process individual item
			value := item.Value().(string)
			return NewResult("processed_" + value), nil
		}).
		WithPostFunc(func(ctx context.Context, shared *SharedStore, prep, exec []Result) (Action, error) {
			// Handle aggregated results
			var results []any
			for _, r := range exec {
				if !r.IsError() {
					results = append(results, r.Value())
				}
			}
			shared.Set("results", results)
			return DefaultAction, nil
		})

	shared := NewSharedStore()
	ctx := context.Background()

	action, err := Run(ctx, node, shared)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if action != DefaultAction {
		t.Errorf("expected DefaultAction, got %v", action)
	}

	if atomic.LoadInt32(&processedCount) != 3 {
		t.Errorf("expected 3 items processed, got %d", processedCount)
	}

	results, _ := shared.Get("results")
	if len(results.([]any)) != 3 {
		t.Errorf("expected 3 results, got %d", len(results.([]any)))
	}
}

func TestBatchNodeWithErrors(t *testing.T) {
	node := NewBatchNode().
		WithPrepFunc(func(ctx context.Context, shared *SharedStore) ([]Result, error) {
			return []Result{
				NewResult(1),
				NewResult(2),
				NewResult(3),
			}, nil
		}).
		WithExecFunc(func(ctx context.Context, item Result) (Result, error) {
			val := item.Value().(int)
			if val == 2 {
				return Result{}, errors.New("error processing 2")
			}
			return NewResult(val * 10), nil
		}).
		WithPostFunc(func(ctx context.Context, shared *SharedStore, prep, exec []Result) (Action, error) {
			successCount := 0
			errorCount := 0

			for _, r := range exec {
				if r.IsError() {
					errorCount++
				} else {
					successCount++
				}
			}

			shared.Set("success_count", successCount)
			shared.Set("error_count", errorCount)

			if errorCount > 0 {
				return "partial_failure", nil
			}
			return DefaultAction, nil
		}).
		WithBatchErrorHandling(true) // Continue on errors

	shared := NewSharedStore()
	ctx := context.Background()

	action, err := Run(ctx, node, shared)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if action != "partial_failure" {
		t.Errorf("expected partial_failure action, got %v", action)
	}

	successCount, _ := shared.Get("success_count")
	if successCount != 2 {
		t.Errorf("expected 2 successes, got %v", successCount)
	}

	errorCount, _ := shared.Get("error_count")
	if errorCount != 1 {
		t.Errorf("expected 1 error, got %v", errorCount)
	}
}

func TestBatchNodeConcurrent(t *testing.T) {
	var maxConcurrent int32
	var currentConcurrent int32

	node := NewBatchNode().
		WithPrepFunc(func(ctx context.Context, shared *SharedStore) ([]Result, error) {
			items := make([]Result, 10)
			for i := range items {
				items[i] = NewResult(i)
			}
			return items, nil
		}).
		WithExecFunc(func(ctx context.Context, item Result) (Result, error) {
			current := atomic.AddInt32(&currentConcurrent, 1)
			defer atomic.AddInt32(&currentConcurrent, -1)

			// Track max concurrent
			for {
				max := atomic.LoadInt32(&maxConcurrent)
				if current <= max || atomic.CompareAndSwapInt32(&maxConcurrent, max, current) {
					break
				}
			}

			time.Sleep(10 * time.Millisecond)
			return item, nil
		}).
		WithPostFunc(func(ctx context.Context, shared *SharedStore, prep, exec []Result) (Action, error) {
			shared.Set("processed", len(exec))
			return DefaultAction, nil
		}).
		WithBatchConcurrency(3)

	shared := NewSharedStore()
	ctx := context.Background()

	_, err := Run(ctx, node, shared)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if maxConcurrent > 3 {
		t.Errorf("max concurrency exceeded limit: %d > 3", maxConcurrent)
	}

	processed, _ := shared.Get("processed")
	if processed != 10 {
		t.Errorf("expected 10 items processed, got %v", processed)
	}
}

func TestBatchNodeStopOnError(t *testing.T) {
	var processedCount int32

	node := NewBatchNode().
		WithPrepFunc(func(ctx context.Context, shared *SharedStore) ([]Result, error) {
			items := make([]Result, 5)
			for i := range items {
				items[i] = NewResult(i + 1)
			}
			return items, nil
		}).
		WithExecFunc(func(ctx context.Context, item Result) (Result, error) {
			val := item.Value().(int)
			atomic.AddInt32(&processedCount, 1)

			// Fail on item 3
			if val == 3 {
				return Result{}, errors.New("error on item 3")
			}
			return NewResult(val * 10), nil
		}).
		WithPostFunc(func(ctx context.Context, shared *SharedStore, prep, exec []Result) (Action, error) {
			var errorCount int
			for _, r := range exec {
				if r.IsError() {
					errorCount++
				}
			}
			shared.Set("error_count", errorCount)
			return DefaultAction, nil
		}).
		WithBatchErrorHandling(false) // Stop on error

	shared := NewSharedStore()
	ctx := context.Background()

	_, err := Run(ctx, node, shared)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should have processed items 1, 2, and 3 (stopped after 3 failed)
	if processedCount := atomic.LoadInt32(&processedCount); processedCount != 3 {
		t.Errorf("expected 3 items processed before stopping, got %d", processedCount)
	}
}

func TestBatchNodeEmptyBatch(t *testing.T) {
	node := NewBatchNode().
		WithPrepFunc(func(ctx context.Context, shared *SharedStore) ([]Result, error) {
			// Return empty batch
			return []Result{}, nil
		}).
		WithExecFunc(func(ctx context.Context, item Result) (Result, error) {
			t.Error("Exec should not be called for empty batch")
			return Result{}, nil
		}).
		WithPostFunc(func(ctx context.Context, shared *SharedStore, prep, exec []Result) (Action, error) {
			if len(prep) != 0 || len(exec) != 0 {
				t.Errorf("expected empty results, got prep=%d, exec=%d", len(prep), len(exec))
			}
			shared.Set("empty_handled", true)
			return "empty", nil
		})

	shared := NewSharedStore()
	ctx := context.Background()

	action, err := Run(ctx, node, shared)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if action != "empty" {
		t.Errorf("expected 'empty' action, got %v", action)
	}

	handled, _ := shared.Get("empty_handled")
	if handled != true {
		t.Error("expected empty batch to be handled")
	}
}

func TestBatchNodeWithRetries(t *testing.T) {
	var attemptCounts = make(map[int]int)

	node := NewBatchNode().
		WithPrepFunc(func(ctx context.Context, shared *SharedStore) ([]Result, error) {
			return []Result{
				NewResult(1),
				NewResult(2),
			}, nil
		}).
		WithExecFunc(func(ctx context.Context, item Result) (Result, error) {
			val := item.Value().(int)
			attemptCounts[val]++

			// Fail first attempt for item 1
			if val == 1 && attemptCounts[val] < 2 {
				return Result{}, errors.New("transient error")
			}
			return NewResult(val * 10), nil
		}).
		WithPostFunc(func(ctx context.Context, shared *SharedStore, prep, exec []Result) (Action, error) {
			var results []any
			for _, r := range exec {
				if !r.IsError() {
					results = append(results, r.Value())
				}
			}
			shared.Set("results", results)
			return DefaultAction, nil
		}).
		WithMaxRetries(3).
		WithBatchErrorHandling(true)

	shared := NewSharedStore()
	ctx := context.Background()

	_, err := Run(ctx, node, shared)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Item 1 should have been retried
	if attemptCounts[1] != 2 {
		t.Errorf("expected 2 attempts for item 1, got %d", attemptCounts[1])
	}

	// Item 2 should have succeeded on first try
	if attemptCounts[2] != 1 {
		t.Errorf("expected 1 attempt for item 2, got %d", attemptCounts[2])
	}

	results, _ := shared.Get("results")
	if len(results.([]any)) != 2 {
		t.Errorf("expected 2 successful results, got %d", len(results.([]any)))
	}
}

func TestBatchNodeResultErrors(t *testing.T) {
	node := NewBatchNode().
		WithPrepFunc(func(ctx context.Context, shared *SharedStore) ([]Result, error) {
			return []Result{
				NewResult("ok1"),
				NewErrorResult(errors.New("prep error")),
				NewResult("ok2"),
			}, nil
		}).
		WithExecFunc(func(ctx context.Context, item Result) (Result, error) {
			if item.IsError() {
				// Pass through error Results
				return item, nil
			}
			return NewResult(fmt.Sprintf("processed_%v", item.Value())), nil
		}).
		WithPostFunc(func(ctx context.Context, shared *SharedStore, prep, exec []Result) (Action, error) {
			var prepErrors, execErrors int
			for _, r := range prep {
				if r.IsError() {
					prepErrors++
				}
			}
			for _, r := range exec {
				if r.IsError() {
					execErrors++
				}
			}
			shared.Set("prep_errors", prepErrors)
			shared.Set("exec_errors", execErrors)
			return DefaultAction, nil
		})

	shared := NewSharedStore()
	ctx := context.Background()

	_, err := Run(ctx, node, shared)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	prepErrors, _ := shared.Get("prep_errors")
	if prepErrors != 1 {
		t.Errorf("expected 1 prep error, got %v", prepErrors)
	}

	execErrors, _ := shared.Get("exec_errors")
	if execErrors != 1 {
		t.Errorf("expected 1 exec error, got %v", execErrors)
	}
}

func TestBatchError(t *testing.T) {
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

func TestBatchNodeInFlow(t *testing.T) {
	// Test using a BatchNode within a Flow
	batchNode := NewBatchNode().
		WithPrepFunc(func(ctx context.Context, shared *SharedStore) ([]Result, error) {
			items, _ := shared.Get("items")
			slice := ToSlice(items)
			results := make([]Result, len(slice))
			for i, item := range slice {
				results[i] = NewResult(item)
			}
			return results, nil
		}).
		WithExecFunc(func(ctx context.Context, item Result) (Result, error) {
			val := item.Value().(int)
			return NewResult(val * 2), nil
		}).
		WithPostFunc(func(ctx context.Context, shared *SharedStore, prep, exec []Result) (Action, error) {
			var results []any
			for _, r := range exec {
				if !r.IsError() {
					results = append(results, r.Value())
				}
			}
			shared.Set("doubled", results)
			return DefaultAction, nil
		})

	sumNode := NewNode().
		WithPrepFuncAny(func(ctx context.Context, shared *SharedStore) (any, error) {
			doubled, _ := shared.Get("doubled")
			return doubled, nil
		}).
		WithExecFuncAny(func(ctx context.Context, prepResult any) (any, error) {
			doubled := prepResult.([]any)
			var sum int
			for _, v := range doubled {
				sum += v.(int)
			}
			return sum, nil
		}).
		WithPostFuncAny(func(ctx context.Context, shared *SharedStore, prep, exec any) (Action, error) {
			shared.Set("sum", exec)
			return DefaultAction, nil
		})

	flow := NewFlow(batchNode)
	flow.Connect(batchNode, DefaultAction, sumNode)

	shared := NewSharedStore()
	shared.Set("items", []int{1, 2, 3, 4, 5})

	ctx := context.Background()
	err := flow.Run(ctx, shared)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	sumResult, _ := shared.Get("sum")
	// The sum is wrapped in a Result from CustomNode.Exec
	var sum int
	if r, ok := sumResult.(Result); ok {
		sum = r.Value().(int)
	} else {
		sum = sumResult.(int)
	}

	if sum != 30 { // (1+2+3+4+5) * 2 = 30
		t.Errorf("expected sum of 30, got %v", sum)
	}
}
