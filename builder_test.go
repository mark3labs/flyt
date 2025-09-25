package flyt_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/mark3labs/flyt"
)

// TestNodeBuilderBackwardsCompatibility verifies that existing code continues to work
func TestNodeBuilderBackwardsCompatibility(t *testing.T) {
	tests := []struct {
		name string
		node func() flyt.Node
	}{
		{
			name: "traditional_with_exec_func",
			node: func() flyt.Node {
				return flyt.NewNode(
					flyt.WithExecFunc(func(ctx context.Context, prepResult flyt.Result) (flyt.Result, error) {
						return flyt.R("executed"), nil
					}),
				)
			},
		},
		{
			name: "traditional_with_multiple_options",
			node: func() flyt.Node {
				return flyt.NewNode(
					flyt.WithMaxRetries(3),
					flyt.WithWait(time.Second),
					flyt.WithExecFuncAny(func(ctx context.Context, prepResult any) (any, error) {
						return "executed", nil
					}),
				)
			},
		},
		{
			name: "traditional_with_all_phases",
			node: func() flyt.Node {
				return flyt.NewNode(
					flyt.WithPrepFunc(func(ctx context.Context, shared *flyt.SharedStore) (flyt.Result, error) {
						return flyt.R("prepped"), nil
					}),
					flyt.WithExecFunc(func(ctx context.Context, prepResult flyt.Result) (flyt.Result, error) {
						return flyt.R("executed"), nil
					}),
					flyt.WithPostFunc(func(ctx context.Context, shared *flyt.SharedStore, prepResult, execResult flyt.Result) (flyt.Action, error) {
						return flyt.DefaultAction, nil
					}),
				)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := tt.node()
			if node == nil {
				t.Fatal("expected non-nil node")
			}

			// Verify it implements Node interface
			ctx := context.Background()
			shared := flyt.NewSharedStore()

			action, err := flyt.Run(ctx, node, shared)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if action == "" {
				t.Error("expected non-empty action")
			}
		})
	}
}

// TestNodeBuilderChaining tests the new builder pattern
func TestNodeBuilderChaining(t *testing.T) {
	tests := []struct {
		name     string
		buildFn  func() flyt.Node
		validate func(t *testing.T, node flyt.Node)
	}{
		{
			name: "simple_exec_chain",
			buildFn: func() flyt.Node {
				return flyt.NewNode().
					WithExecFunc(func(ctx context.Context, prepResult flyt.Result) (flyt.Result, error) {
						return flyt.R("chained"), nil
					})
			},
			validate: func(t *testing.T, node flyt.Node) {
				ctx := context.Background()
				result, err := node.Exec(ctx, nil)
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				r := result.(flyt.Result)
				if r.MustString() != "chained" {
					t.Errorf("expected 'chained', got %v", r.Value())
				}
			},
		},
		{
			name: "full_chain_with_retries",
			buildFn: func() flyt.Node {
				return flyt.NewNode().
					WithMaxRetries(5).
					WithWait(100 * time.Millisecond).
					WithPrepFunc(func(ctx context.Context, shared *flyt.SharedStore) (flyt.Result, error) {
						return flyt.R("prep_data"), nil
					}).
					WithExecFunc(func(ctx context.Context, prepResult flyt.Result) (flyt.Result, error) {
						return flyt.R("exec_result"), nil
					}).
					WithPostFunc(func(ctx context.Context, shared *flyt.SharedStore, prepResult, execResult flyt.Result) (flyt.Action, error) {
						shared.Set("result", execResult.Value())
						return "custom_action", nil
					})
			},
			validate: func(t *testing.T, node flyt.Node) {
				// Check if it's retryable
				if retryable, ok := node.(flyt.RetryableNode); ok {
					if retryable.GetMaxRetries() != 5 {
						t.Errorf("expected 5 retries, got %d", retryable.GetMaxRetries())
					}
					if retryable.GetWait() != 100*time.Millisecond {
						t.Errorf("expected 100ms wait, got %v", retryable.GetWait())
					}
				} else {
					t.Error("expected node to implement RetryableNode")
				}

				// Run the node
				ctx := context.Background()
				shared := flyt.NewSharedStore()
				action, err := flyt.Run(ctx, node, shared)
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if action != "custom_action" {
					t.Errorf("expected 'custom_action', got %v", action)
				}
				resultVal, exists := shared.Get("result")
				if !exists {
					t.Error("expected stored result to exist")
				} else if r, ok := resultVal.(flyt.Result); ok {
					if r.MustString() != "exec_result" {
						t.Errorf("expected stored result 'exec_result', got %v", r.MustString())
					}
				} else if resultVal != "exec_result" {
					t.Errorf("expected stored result 'exec_result', got %v", resultVal)
				}
			},
		},
		{
			name: "chain_with_any_functions",
			buildFn: func() flyt.Node {
				return flyt.NewNode().
					WithPrepFuncAny(func(ctx context.Context, shared *flyt.SharedStore) (any, error) {
						return map[string]string{"key": "value"}, nil
					}).
					WithExecFuncAny(func(ctx context.Context, prepResult any) (any, error) {
						m := prepResult.(map[string]string)
						return m["key"], nil
					}).
					WithPostFuncAny(func(ctx context.Context, shared *flyt.SharedStore, prepResult, execResult any) (flyt.Action, error) {
						// Extract value from Result if it's wrapped
						var resultValue any
						if r, ok := execResult.(flyt.Result); ok {
							resultValue = r.Value()
						} else {
							resultValue = execResult
						}

						if resultValue == "value" {
							return "success", nil
						}
						return "failure", nil
					})
			},
			validate: func(t *testing.T, node flyt.Node) {
				ctx := context.Background()
				shared := flyt.NewSharedStore()
				action, err := flyt.Run(ctx, node, shared)
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if action != "success" {
					t.Errorf("expected 'success', got %v", action)
				}
			},
		},
		{
			name: "chain_with_fallback",
			buildFn: func() flyt.Node {
				attempts := 0
				return flyt.NewNode().
					WithMaxRetries(2).
					WithExecFunc(func(ctx context.Context, prepResult flyt.Result) (flyt.Result, error) {
						attempts++
						return flyt.Result{}, errors.New("exec failed")
					}).
					WithExecFallbackFunc(func(prepResult any, err error) (any, error) {
						return "fallback_value", nil
					})
			},
			validate: func(t *testing.T, node flyt.Node) {
				ctx := context.Background()
				shared := flyt.NewSharedStore()
				action, err := flyt.Run(ctx, node, shared)
				if err != nil {
					t.Errorf("expected fallback to handle error, got: %v", err)
				}
				if action != flyt.DefaultAction {
					t.Errorf("expected DefaultAction, got %v", action)
				}
			},
		},
		{
			name: "chain_with_build_method",
			buildFn: func() flyt.Node {
				return flyt.NewNode().
					WithExecFunc(func(ctx context.Context, prepResult flyt.Result) (flyt.Result, error) {
						return flyt.R("built"), nil
					}).
					Build() // Explicitly call Build()
			},
			validate: func(t *testing.T, node flyt.Node) {
				ctx := context.Background()
				result, err := node.Exec(ctx, nil)
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				r := result.(flyt.Result)
				if r.MustString() != "built" {
					t.Errorf("expected 'built', got %v", r.Value())
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := tt.buildFn()
			if node == nil {
				t.Fatal("expected non-nil node")
			}
			tt.validate(t, node)
		})
	}
}

// TestNodeBuilderMixedStyle tests mixing traditional and builder patterns
func TestNodeBuilderMixedStyle(t *testing.T) {
	// Create node with traditional options
	builder := flyt.NewNode(
		flyt.WithMaxRetries(3),
		flyt.WithWait(time.Second),
	)

	// Continue configuring with builder pattern
	var node flyt.Node = builder.
		WithExecFunc(func(ctx context.Context, prepResult flyt.Result) (flyt.Result, error) {
			return flyt.R("mixed_result"), nil
		}).
		WithPostFunc(func(ctx context.Context, shared *flyt.SharedStore, prepResult, execResult flyt.Result) (flyt.Action, error) {
			return "mixed_action", nil
		})

	// Validate
	if retryable, ok := node.(flyt.RetryableNode); ok {
		if retryable.GetMaxRetries() != 3 {
			t.Errorf("expected 3 retries, got %d", retryable.GetMaxRetries())
		}
		if retryable.GetWait() != time.Second {
			t.Errorf("expected 1s wait, got %v", retryable.GetWait())
		}
	} else {
		t.Error("expected node to implement RetryableNode")
	}

	ctx := context.Background()
	shared := flyt.NewSharedStore()
	action, err := flyt.Run(ctx, node, shared)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if action != "mixed_action" {
		t.Errorf("expected 'mixed_action', got %v", action)
	}
}

// TestNodeBuilderInterfaceCompliance ensures NodeBuilder properly implements all interfaces
func TestNodeBuilderInterfaceCompliance(t *testing.T) {
	node := flyt.NewNode().
		WithMaxRetries(3).
		WithWait(time.Second)

	// Test Node interface
	var _ flyt.Node = node

	// Test RetryableNode interface
	var _ flyt.RetryableNode = node

	// Test FallbackNode interface
	var _ flyt.FallbackNode = node

	// All interfaces should be satisfied
	t.Log("NodeBuilder successfully implements all required interfaces")
}

// BenchmarkNodeCreationTraditional benchmarks traditional node creation
func BenchmarkNodeCreationTraditional(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = flyt.NewNode(
			flyt.WithMaxRetries(3),
			flyt.WithWait(time.Second),
			flyt.WithExecFunc(func(ctx context.Context, prepResult flyt.Result) (flyt.Result, error) {
				return flyt.R("result"), nil
			}),
		)
	}
}

// BenchmarkNodeCreationBuilder benchmarks builder pattern node creation
func BenchmarkNodeCreationBuilder(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = flyt.NewNode().
			WithMaxRetries(3).
			WithWait(time.Second).
			WithExecFunc(func(ctx context.Context, prepResult flyt.Result) (flyt.Result, error) {
				return flyt.R("result"), nil
			})
	}
}
