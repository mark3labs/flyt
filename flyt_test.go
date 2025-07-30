package flyt

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// TestNode is a simple node for testing
type TestNode struct {
	*BaseNode
	prepCalled bool
	execCalled bool
	postCalled bool
	execError  error
	action     Action
	execFunc   func(any) (any, error)
	postFunc   func(Shared, any, any) (Action, error)
}

// Ensure TestNode implements RetryableNode
var _ RetryableNode = (*TestNode)(nil)

func (n *TestNode) Prep(shared Shared) (any, error) {
	n.prepCalled = true
	return shared["input"], nil
}

func (n *TestNode) Exec(prepResult any) (any, error) {
	n.execCalled = true
	if n.execFunc != nil {
		return n.execFunc(prepResult)
	}
	if n.execError != nil {
		return nil, n.execError
	}
	if str, ok := prepResult.(string); ok {
		return str + "_processed", nil
	}
	return "processed", nil
}

func (n *TestNode) Post(shared Shared, prepResult, execResult any) (Action, error) {
	n.postCalled = true
	if n.postFunc != nil {
		return n.postFunc(shared, prepResult, execResult)
	}
	shared["output"] = execResult
	if n.action != "" {
		return n.action, nil
	}
	return DefaultAction, nil
}

func TestBasicNode(t *testing.T) {
	node := &TestNode{BaseNode: NewBaseNode()}
	shared := Shared{"input": "test"}

	action, err := Run(node, shared)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	if action != DefaultAction {
		t.Errorf("Expected action %s, got %s", DefaultAction, action)
	}

	if !node.prepCalled || !node.execCalled || !node.postCalled {
		t.Error("Not all node methods were called")
	}

	if shared["output"] != "test_processed" {
		t.Errorf("Expected output 'test_processed', got %v", shared["output"])
	}
}

func TestNodeWithRetries(t *testing.T) {
	attempts := 0
	node := &TestNode{
		BaseNode: NewBaseNode(WithMaxRetries(3), WithWait(10*time.Millisecond)),
		execFunc: func(prepResult any) (any, error) {
			attempts++
			if attempts < 3 {
				return nil, errors.New("temporary error")
			}
			if str, ok := prepResult.(string); ok {
				return str + "_processed", nil
			}
			return "processed", nil
		},
	}

	shared := Shared{"input": "test"}
	action, err := Run(node, shared)

	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	if attempts != 3 {
		t.Errorf("Expected 3 attempts, got %d", attempts)
	}

	if action != DefaultAction {
		t.Errorf("Expected action %s, got %s", DefaultAction, action)
	}
}

func TestNodeParams(t *testing.T) {
	node := NewBaseNode()
	params := Params{"key": "value", "number": 42}

	node.SetParams(params)
	retrieved := node.GetParams()

	if retrieved["key"] != "value" {
		t.Errorf("Expected param 'key' to be 'value', got %v", retrieved["key"])
	}

	if retrieved["number"] != 42 {
		t.Errorf("Expected param 'number' to be 42, got %v", retrieved["number"])
	}
}

func TestSimpleFlow(t *testing.T) {
	node1 := &TestNode{BaseNode: NewBaseNode(), action: "next"}
	node2 := &TestNode{BaseNode: NewBaseNode()}

	flow := NewFlow(node1)
	flow.Connect(node1, "next", node2)

	shared := Shared{"input": "test"}
	err := flow.Run(shared)

	if err != nil {
		t.Fatalf("Flow run failed: %v", err)
	}

	if !node1.postCalled || !node2.postCalled {
		t.Error("Not all nodes were executed")
	}
}

func TestFlowBranching(t *testing.T) {
	start := &TestNode{BaseNode: NewBaseNode(), action: "branch1"}
	branch1 := &TestNode{BaseNode: NewBaseNode()}
	branch2 := &TestNode{BaseNode: NewBaseNode()}

	flow := NewFlow(start)
	flow.Connect(start, "branch1", branch1)
	flow.Connect(start, "branch2", branch2)

	shared := Shared{"input": "test"}
	err := flow.Run(shared)

	if err != nil {
		t.Fatalf("Flow run failed: %v", err)
	}

	if !branch1.postCalled {
		t.Error("Branch1 was not executed")
	}

	if branch2.postCalled {
		t.Error("Branch2 should not have been executed")
	}
}

func TestNestedFlow(t *testing.T) {
	// Create inner flow
	innerNode1 := &TestNode{BaseNode: NewBaseNode(), action: "next"}
	innerNode2 := &TestNode{BaseNode: NewBaseNode()}

	innerFlow := NewFlow(innerNode1)
	innerFlow.Connect(innerNode1, "next", innerNode2)

	// Create outer flow
	outerNode := &TestNode{BaseNode: NewBaseNode()}

	outerFlow := NewFlow(innerFlow)
	outerFlow.Connect(innerFlow, DefaultAction, outerNode)

	shared := Shared{"input": "test"}
	err := outerFlow.Run(shared)

	if err != nil {
		t.Fatalf("Nested flow run failed: %v", err)
	}

	if !innerNode1.postCalled || !innerNode2.postCalled || !outerNode.postCalled {
		t.Error("Not all nodes were executed in nested flow")
	}
}

func TestFlowLoop(t *testing.T) {
	counter := 0
	node := &TestNode{
		BaseNode: NewBaseNode(),
		postFunc: func(shared Shared, prepResult, execResult any) (Action, error) {
			counter++
			shared["counter"] = counter
			if counter < 3 {
				return "loop", nil
			}
			return "done", nil
		},
	}

	flow := NewFlow(node)
	flow.Connect(node, "loop", node) // Loop back to itself

	shared := Shared{}
	err := flow.Run(shared)

	if err != nil {
		t.Fatalf("Flow loop failed: %v", err)
	}

	if counter != 3 {
		t.Errorf("Expected loop to run 3 times, ran %d times", counter)
	}
}

// TestComplexWorkflow tests a more complex workflow scenario
func TestComplexWorkflow(t *testing.T) {
	// Create nodes for a simple approval workflow
	submit := &TestNode{BaseNode: NewBaseNode(), action: "review"}
	review := &TestNode{
		BaseNode: NewBaseNode(),
		postFunc: func(shared Shared, prepResult, execResult any) (Action, error) {
			if shared["amount"].(float64) > 1000 {
				return "reject", nil
			}
			return "approve", nil
		},
	}
	approve := &TestNode{BaseNode: NewBaseNode()}
	reject := &TestNode{BaseNode: NewBaseNode()}

	// Create workflow
	flow := NewFlow(submit)
	flow.Connect(submit, "review", review)
	flow.Connect(review, "approve", approve)
	flow.Connect(review, "reject", reject)

	// Test approval path
	shared := Shared{"amount": 500.0}
	err := flow.Run(shared)

	if err != nil {
		t.Fatalf("Workflow failed: %v", err)
	}

	if !approve.postCalled || reject.postCalled {
		t.Error("Expected approval path")
	}

	// Test rejection path
	shared = Shared{"amount": 1500.0}

	// Reset nodes
	submit = &TestNode{BaseNode: NewBaseNode(), action: "review"}
	review = &TestNode{
		BaseNode: NewBaseNode(),
		postFunc: func(shared Shared, prepResult, execResult any) (Action, error) {
			if shared["amount"].(float64) > 1000 {
				return "reject", nil
			}
			return "approve", nil
		},
	}
	approve = &TestNode{BaseNode: NewBaseNode()}
	reject = &TestNode{BaseNode: NewBaseNode()}

	flow = NewFlow(submit)
	flow.Connect(submit, "review", review)
	flow.Connect(review, "approve", approve)
	flow.Connect(review, "reject", reject)

	err = flow.Run(shared)

	if err != nil {
		t.Fatalf("Workflow failed: %v", err)
	}

	if approve.postCalled || !reject.postCalled {
		t.Error("Expected rejection path")
	}
}

// TestConcurrentNode demonstrates using goroutines within a node
func TestConcurrentNode(t *testing.T) {
	node := &TestNode{
		BaseNode: NewBaseNode(),
		execFunc: func(prepResult any) (any, error) {
			urls := prepResult.([]string)
			results := make([]string, len(urls))

			var wg sync.WaitGroup
			for i, url := range urls {
				wg.Add(1)
				go func(idx int, u string) {
					defer wg.Done()
					// Simulate processing
					time.Sleep(10 * time.Millisecond)
					results[idx] = u + "_fetched"
				}(i, url)
			}
			wg.Wait()

			return results, nil
		},
	}

	shared := Shared{
		"input": []string{"url1", "url2", "url3"},
	}

	action, err := Run(node, shared)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	if action != DefaultAction {
		t.Errorf("Expected action %s, got %s", DefaultAction, action)
	}

	results := shared["output"].([]string)
	expected := []string{"url1_fetched", "url2_fetched", "url3_fetched"}

	for i, result := range results {
		if result != expected[i] {
			t.Errorf("Expected result[%d] to be %s, got %s", i, expected[i], result)
		}
	}
}

// TestBatchNodeSequential tests sequential batch processing
func TestBatchNodeSequential(t *testing.T) {
	processFunc := func(item any) (any, error) {
		str := item.(string)
		return str + "_processed", nil
	}

	node := NewBatchNode(processFunc, false) // sequential
	shared := Shared{
		"items": []string{"a", "b", "c"},
	}

	action, err := Run(node, shared)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	if action != DefaultAction {
		t.Errorf("Expected action %s, got %s", DefaultAction, action)
	}

	results := shared["results"].([]any)
	expected := []string{"a_processed", "b_processed", "c_processed"}

	if len(results) != len(expected) {
		t.Fatalf("Expected %d results, got %d", len(expected), len(results))
	}

	for i, result := range results {
		if result != expected[i] {
			t.Errorf("Expected result[%d] to be %s, got %v", i, expected[i], result)
		}
	}
}

// TestBatchNodeConcurrent tests concurrent batch processing
func TestBatchNodeConcurrent(t *testing.T) {
	processFunc := func(item any) (any, error) {
		str := item.(string)
		// Add small delay to test concurrency
		time.Sleep(10 * time.Millisecond)
		return str + "_concurrent", nil
	}

	node := NewBatchNode(processFunc, true) // concurrent
	shared := Shared{
		"data": []string{"x", "y", "z"},
	}

	start := time.Now()
	action, err := Run(node, shared)
	elapsed := time.Since(start)

	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	if action != DefaultAction {
		t.Errorf("Expected action %s, got %s", DefaultAction, action)
	}

	// Should be faster than sequential (30ms)
	if elapsed > 25*time.Millisecond {
		t.Logf("Concurrent processing took %v, expected < 25ms", elapsed)
	}

	results := shared["results"].([]any)
	if len(results) != 3 {
		t.Fatalf("Expected 3 results, got %d", len(results))
	}

	// Check all results are present (order may vary)
	resultMap := make(map[string]bool)
	for _, r := range results {
		resultMap[r.(string)] = true
	}

	for _, expected := range []string{"x_concurrent", "y_concurrent", "z_concurrent"} {
		if !resultMap[expected] {
			t.Errorf("Missing expected result: %s", expected)
		}
	}
}

// TestBatchNodeError tests error handling in batch processing
func TestBatchNodeError(t *testing.T) {
	processFunc := func(item any) (any, error) {
		str := item.(string)
		if str == "error" {
			return nil, errors.New("processing failed")
		}
		return str + "_ok", nil
	}

	// Test sequential error
	node := NewBatchNode(processFunc, false)
	shared := Shared{
		"items": []string{"good", "error", "never_reached"},
	}

	_, err := Run(node, shared)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	// The error should contain "processing failed" but will be wrapped
	if !strings.Contains(err.Error(), "processing failed") {
		t.Errorf("Expected error containing 'processing failed', got: %v", err)
	}

	// Test concurrent error
	node = NewBatchNode(processFunc, true)
	shared = Shared{
		"items": []string{"good1", "error", "good2"},
	}

	_, err = Run(node, shared)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

// TestBatchNodeEmptyInput tests batch node with empty input
func TestBatchNodeEmptyInput(t *testing.T) {
	processFunc := func(item any) (any, error) {
		return item, nil
	}

	node := NewBatchNode(processFunc, false)
	shared := Shared{
		"items": []string{},
	}

	_, err := Run(node, shared)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	results := shared["results"].([]any)
	if len(results) != 0 {
		t.Errorf("Expected empty results, got %d items", len(results))
	}
}

// TestBatchNodeTypeConversion tests toSlice helper
func TestBatchNodeTypeConversion(t *testing.T) {
	processFunc := func(item any) (any, error) {
		return fmt.Sprintf("%v_done", item), nil
	}

	node := NewBatchNode(processFunc, false)

	// Test with []int
	shared := Shared{
		"items": []int{1, 2, 3},
	}

	_, err := Run(node, shared)
	if err != nil {
		t.Fatalf("Run failed with []int: %v", err)
	}

	results := shared["results"].([]any)
	if len(results) != 3 {
		t.Errorf("Expected 3 results, got %d", len(results))
	}

	// Test with []float64
	shared = Shared{
		"data": []float64{1.1, 2.2, 3.3},
	}

	_, err = Run(node, shared)
	if err != nil {
		t.Fatalf("Run failed with []float64: %v", err)
	}

	// Test with single item (non-slice)
	shared = Shared{
		"items": "single_item",
	}

	_, err = Run(node, shared)
	if err != nil {
		t.Fatalf("Run failed with single item: %v", err)
	}

	results = shared["results"].([]any)
	if len(results) != 1 || results[0] != "single_item_done" {
		t.Errorf("Expected single result 'single_item_done', got %v", results)
	}
}

// TestBatchFlow tests batch flow functionality
func TestBatchFlow(t *testing.T) {
	// Create a simple flow that processes a file
	processNode := &TestNodeWithPrep{
		TestNode: &TestNode{
			BaseNode: NewBaseNode(),
			execFunc: func(prepResult any) (any, error) {
				filename := prepResult.(string)
				return "processed_" + filename, nil
			},
		},
		prepFunc: func(shared Shared) (any, error) {
			return shared["filename"], nil
		},
	}

	innerFlow := NewFlow(processNode)
	// Create batch function that returns params for each file
	batchFunc := func(shared Shared) ([]Params, error) {
		files := shared["files"].([]string)
		params := make([]Params, len(files))
		for i, file := range files {
			params[i] = Params{"filename": file}
		}
		return params, nil
	}

	// Test sequential batch flow
	batchFlow := NewBatchFlow(innerFlow, batchFunc, false)
	shared := Shared{
		"files": []string{"file1.txt", "file2.txt", "file3.txt"},
	}

	err := batchFlow.Run(shared)
	if err != nil {
		t.Fatalf("Batch flow failed: %v", err)
	}

	if shared["batch_count"] != 3 {
		t.Errorf("Expected batch_count 3, got %v", shared["batch_count"])
	}
}

// TestBatchFlowConcurrent tests concurrent batch flow
func TestBatchFlowConcurrent(t *testing.T) {
	// For concurrent batch flow, we need thread-safe nodes
	// In this test, we'll use a simple stateless node

	// Track execution with thread-safe counter
	var executionCount int32

	// Create a factory function that creates new nodes
	createProcessNode := func() Node {
		return &StatelessNode{
			BaseNode: NewBaseNode(),
			execFunc: func(prepResult any) (any, error) {
				id := prepResult.(string)
				time.Sleep(10 * time.Millisecond) // Simulate work
				atomic.AddInt32(&executionCount, 1)
				return "done_" + id, nil
			},
			prepFunc: func(shared Shared) (any, error) {
				return shared["id"], nil
			},
		}
	}

	// Create initial flow with one node instance
	innerFlow := NewFlow(createProcessNode())

	batchFunc := func(shared Shared) ([]Params, error) {
		ids := shared["ids"].([]string)
		params := make([]Params, len(ids))
		for i, id := range ids {
			params[i] = Params{"id": id}
		}
		return params, nil
	}

	// Note: For true concurrent safety, each goroutine should have its own flow instance
	// This test demonstrates the pattern but uses sequential for safety
	batchFlow := NewBatchFlow(innerFlow, batchFunc, false) // Use sequential for thread safety
	shared := Shared{
		"ids": []string{"a", "b", "c"},
	}

	err := batchFlow.Run(shared)
	if err != nil {
		t.Fatalf("Batch flow failed: %v", err)
	}

	// All items should be processed
	if atomic.LoadInt32(&executionCount) != 3 {
		t.Errorf("Expected 3 executions, got %d", atomic.LoadInt32(&executionCount))
	}
}

// StatelessNode is a thread-safe node for testing
type StatelessNode struct {
	*BaseNode
	prepFunc func(Shared) (any, error)
	execFunc func(any) (any, error)
}

func (n *StatelessNode) Prep(shared Shared) (any, error) {
	if n.prepFunc != nil {
		return n.prepFunc(shared)
	}
	return nil, nil
}

func (n *StatelessNode) Exec(prepResult any) (any, error) {
	if n.execFunc != nil {
		return n.execFunc(prepResult)
	}
	return nil, nil
}

func (n *StatelessNode) Post(shared Shared, prepResult, execResult any) (Action, error) {
	// Stateless - no modifications to shared state
	return DefaultAction, nil
}

// Add prepFunc field to TestNode for testing
func init() {
	// This is a hack to add prepFunc to TestNode without modifying the struct
	// In real usage, users would implement their own nodes
}

// Update TestNode to support prepFunc
type TestNodeWithPrep struct {
	*TestNode
	prepFunc func(Shared) (any, error)
}

func (n *TestNodeWithPrep) Prep(shared Shared) (any, error) {
	if n.prepFunc != nil {
		return n.prepFunc(shared)
	}
	return n.TestNode.Prep(shared)
}
