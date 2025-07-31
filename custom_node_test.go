package flyt

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"
)

// TestNewNodeWithExecFunc tests creating a node with only Exec function
func TestNewNodeWithExecFunc(t *testing.T) {
	executed := false
	node := NewNode(
		WithExecFunc(func(ctx context.Context, prepResult any) (any, error) {
			executed = true
			return "exec result", nil
		}),
	)

	shared := NewSharedStore()
	ctx := context.Background()

	action, err := Run(ctx, node, shared)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !executed {
		t.Error("exec function was not called")
	}

	if action != DefaultAction {
		t.Errorf("expected default action, got %v", action)
	}
}

// TestNewNodeWithAllFuncs tests creating a node with all three functions
func TestNewNodeWithAllFuncs(t *testing.T) {
	var prepCalled, execCalled, postCalled bool

	node := NewNode(
		WithPrepFunc(func(ctx context.Context, shared *SharedStore) (any, error) {
			prepCalled = true
			return "prep data", nil
		}),
		WithExecFunc(func(ctx context.Context, prepResult any) (any, error) {
			execCalled = true
			if prepResult != "prep data" {
				t.Errorf("expected prep data, got %v", prepResult)
			}
			return "exec result", nil
		}),
		WithPostFunc(func(ctx context.Context, shared *SharedStore, prepResult, execResult any) (Action, error) {
			postCalled = true
			if prepResult != "prep data" {
				t.Errorf("expected prep data in post, got %v", prepResult)
			}
			if execResult != "exec result" {
				t.Errorf("expected exec result in post, got %v", execResult)
			}
			shared.Set("final", "done")
			return "next", nil
		}),
	)

	shared := NewSharedStore()
	ctx := context.Background()

	action, err := Run(ctx, node, shared)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !prepCalled || !execCalled || !postCalled {
		t.Error("not all functions were called")
	}

	if action != "next" {
		t.Errorf("expected 'next' action, got %v", action)
	}

	val, _ := shared.Get("final")
	if val != "done" {
		t.Errorf("expected 'done' in shared store, got %v", val)
	}
}

// TestNewNodeErrorHandling tests error propagation from custom functions
func TestNewNodeErrorHandling(t *testing.T) {
	tests := []struct {
		name          string
		node          Node
		expectedError string
	}{
		{
			name: "prep error",
			node: NewNode(
				WithPrepFunc(func(ctx context.Context, shared *SharedStore) (any, error) {
					return nil, errors.New("prep failed")
				}),
			),
			expectedError: "prep failed",
		},
		{
			name: "exec error",
			node: NewNode(
				WithExecFunc(func(ctx context.Context, prepResult any) (any, error) {
					return nil, errors.New("exec failed")
				}),
			),
			expectedError: "exec failed",
		},
		{
			name: "post error",
			node: NewNode(
				WithPostFunc(func(ctx context.Context, shared *SharedStore, prepResult, execResult any) (Action, error) {
					return "", errors.New("post failed")
				}),
			),
			expectedError: "post failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shared := NewSharedStore()
			ctx := context.Background()

			_, err := Run(ctx, tt.node, shared)
			if err == nil {
				t.Fatal("expected error but got none")
			}

			if !strings.Contains(err.Error(), tt.expectedError) {
				t.Errorf("expected error containing %q, got %v", tt.expectedError, err)
			}
		})
	}
}

// TestNewNodeWithBaseNodeOptions tests combining custom functions with BaseNode options
func TestNewNodeWithBaseNodeOptions(t *testing.T) {
	attemptCount := 0
	node := NewNode(
		WithExecFunc(func(ctx context.Context, prepResult any) (any, error) {
			attemptCount++
			if attemptCount < 3 {
				return nil, errors.New("temporary failure")
			}
			return "success", nil
		}),
		WithMaxRetries(3),
		WithWait(10*time.Millisecond),
	)

	shared := NewSharedStore()
	ctx := context.Background()

	start := time.Now()
	action, err := Run(ctx, node, shared)
	duration := time.Since(start)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if attemptCount != 3 {
		t.Errorf("expected 3 attempts, got %d", attemptCount)
	}

	// Should have waited at least 20ms (2 retries * 10ms)
	if duration < 20*time.Millisecond {
		t.Errorf("expected wait time of at least 20ms, got %v", duration)
	}

	if action != DefaultAction {
		t.Errorf("expected default action, got %v", action)
	}
}

// TestNewNodeContextCancellation tests that custom functions respect context cancellation
func TestNewNodeContextCancellation(t *testing.T) {
	node := NewNode(
		WithExecFunc(func(ctx context.Context, prepResult any) (any, error) {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(100 * time.Millisecond):
				return "should not reach here", nil
			}
		}),
	)

	shared := NewSharedStore()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	_, err := Run(ctx, node, shared)
	if err == nil {
		t.Fatal("expected context cancellation error")
	}

	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("expected deadline exceeded error, got %v", err)
	}
}

// TestNewNodeInFlow tests using CustomNode in a flow
func TestNewNodeInFlow(t *testing.T) {
	// Create nodes using NewNode
	startNode := NewNode(
		WithExecFunc(func(ctx context.Context, prepResult any) (any, error) {
			return "start", nil
		}),
		WithPostFunc(func(ctx context.Context, shared *SharedStore, prepResult, execResult any) (Action, error) {
			shared.Set("step", "start")
			return "continue", nil
		}),
	)

	endNode := NewNode(
		WithPrepFunc(func(ctx context.Context, shared *SharedStore) (any, error) {
			step, _ := shared.Get("step")
			return step, nil
		}),
		WithExecFunc(func(ctx context.Context, prepResult any) (any, error) {
			if prepResult != "start" {
				return nil, errors.New("unexpected prep result")
			}
			return "end", nil
		}),
		WithPostFunc(func(ctx context.Context, shared *SharedStore, prepResult, execResult any) (Action, error) {
			shared.Set("result", execResult)
			return DefaultAction, nil
		}),
	)

	// Create flow
	flow := NewFlow(startNode)
	flow.Connect(startNode, "continue", endNode)

	// Execute
	shared := NewSharedStore()
	ctx := context.Background()

	err := flow.Run(ctx, shared)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify result
	result, ok := shared.Get("result")
	if !ok || result != "end" {
		t.Errorf("expected 'end' result, got %v", result)
	}
}

// TestNewNodeWithNoFunctions tests that node works with no custom functions
func TestNewNodeWithNoFunctions(t *testing.T) {
	node := NewNode() // No functions provided

	shared := NewSharedStore()
	ctx := context.Background()

	action, err := Run(ctx, node, shared)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if action != DefaultAction {
		t.Errorf("expected default action, got %v", action)
	}
}

// TestNewNodeDataFlow tests data flow through all three phases
func TestNewNodeDataFlow(t *testing.T) {
	node := NewNode(
		WithPrepFunc(func(ctx context.Context, shared *SharedStore) (any, error) {
			data, _ := shared.Get("input")
			return data, nil
		}),
		WithExecFunc(func(ctx context.Context, prepResult any) (any, error) {
			if str, ok := prepResult.(string); ok {
				return str + " processed", nil
			}
			return nil, errors.New("unexpected input type")
		}),
		WithPostFunc(func(ctx context.Context, shared *SharedStore, prepResult, execResult any) (Action, error) {
			shared.Set("output", execResult)
			return DefaultAction, nil
		}),
	)

	shared := NewSharedStore()
	shared.Set("input", "test data")
	ctx := context.Background()

	action, err := Run(ctx, node, shared)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output, ok := shared.Get("output")
	if !ok {
		t.Fatal("output not found in shared store")
	}

	if output != "test data processed" {
		t.Errorf("expected 'test data processed', got %v", output)
	}

	if action != DefaultAction {
		t.Errorf("expected default action, got %v", action)
	}
}

// TestNewNodeConcurrentExecution tests that CustomNode can be used concurrently
func TestNewNodeConcurrentExecution(t *testing.T) {
	// Create a node that increments a counter
	var counter int
	var mu sync.Mutex

	createNode := func(id int) Node {
		return NewNode(
			WithExecFunc(func(ctx context.Context, prepResult any) (any, error) {
				mu.Lock()
				counter++
				mu.Unlock()
				return id, nil
			}),
			WithPostFunc(func(ctx context.Context, shared *SharedStore, prepResult, execResult any) (Action, error) {
				shared.Set(fmt.Sprintf("result_%d", execResult.(int)), true)
				return DefaultAction, nil
			}),
		)
	}

	// Run multiple nodes concurrently
	shared := NewSharedStore()
	ctx := context.Background()
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			node := createNode(id)
			_, err := Run(ctx, node, shared)
			if err != nil {
				t.Errorf("node %d failed: %v", id, err)
			}
		}(i)
	}

	wg.Wait()

	// Verify all nodes executed
	if counter != 10 {
		t.Errorf("expected counter to be 10, got %d", counter)
	}

	// Verify all results are in shared store
	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("result_%d", i)
		if _, ok := shared.Get(key); !ok {
			t.Errorf("result for node %d not found", i)
		}
	}
}

// TestNewNodeWithOnlyPrepFunc tests node with only prep function
func TestNewNodeWithOnlyPrepFunc(t *testing.T) {
	prepCalled := false
	node := NewNode(
		WithPrepFunc(func(ctx context.Context, shared *SharedStore) (any, error) {
			prepCalled = true
			shared.Set("prep_executed", true)
			return "prep result", nil
		}),
	)

	shared := NewSharedStore()
	ctx := context.Background()

	action, err := Run(ctx, node, shared)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !prepCalled {
		t.Error("prep function was not called")
	}

	val, ok := shared.Get("prep_executed")
	if !ok || val != true {
		t.Error("prep function did not set expected value")
	}

	if action != DefaultAction {
		t.Errorf("expected default action, got %v", action)
	}
}

// TestNewNodeWithOnlyPostFunc tests node with only post function
func TestNewNodeWithOnlyPostFunc(t *testing.T) {
	postCalled := false
	node := NewNode(
		WithPostFunc(func(ctx context.Context, shared *SharedStore, prepResult, execResult any) (Action, error) {
			postCalled = true
			shared.Set("post_executed", true)
			return "custom_action", nil
		}),
	)

	shared := NewSharedStore()
	ctx := context.Background()

	action, err := Run(ctx, node, shared)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !postCalled {
		t.Error("post function was not called")
	}

	val, ok := shared.Get("post_executed")
	if !ok || val != true {
		t.Error("post function did not set expected value")
	}

	if action != "custom_action" {
		t.Errorf("expected 'custom_action', got %v", action)
	}
}
