package flyt

import (
	"context"
	"errors"
	"fmt"
	"reflect"
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

// TestSharedStoreTypeGetters tests all the type-specific getter methods
func TestSharedStoreTypeGetters(t *testing.T) {
	store := NewSharedStore()

	// Set various types
	store.Set("string", "hello")
	store.Set("int", 42)
	store.Set("int64", int64(100))
	store.Set("float32", float32(3.14))
	store.Set("float64", 2.718)
	store.Set("bool", true)
	store.Set("slice", []string{"a", "b", "c"})
	store.Set("map", map[string]any{"name": "test", "value": 123})

	// Test GetString
	if got := store.GetString("string"); got != "hello" {
		t.Errorf("GetString: expected hello, got %v", got)
	}
	if got := store.GetString("notfound"); got != "" {
		t.Errorf("GetString (not found): expected empty string, got %v", got)
	}
	if got := store.GetString("int"); got != "" {
		t.Errorf("GetString (wrong type): expected empty string, got %v", got)
	}

	// Test GetStringOr
	if got := store.GetStringOr("string", "default"); got != "hello" {
		t.Errorf("GetStringOr: expected hello, got %v", got)
	}
	if got := store.GetStringOr("notfound", "default"); got != "default" {
		t.Errorf("GetStringOr (not found): expected default, got %v", got)
	}

	// Test GetInt with numeric conversions
	if got := store.GetInt("int"); got != 42 {
		t.Errorf("GetInt: expected 42, got %v", got)
	}
	if got := store.GetInt("int64"); got != 100 {
		t.Errorf("GetInt (int64): expected 100, got %v", got)
	}
	if got := store.GetInt("float32"); got != 3 {
		t.Errorf("GetInt (float32): expected 3, got %v", got)
	}
	if got := store.GetInt("notfound"); got != 0 {
		t.Errorf("GetInt (not found): expected 0, got %v", got)
	}

	// Test GetIntOr
	if got := store.GetIntOr("int", -1); got != 42 {
		t.Errorf("GetIntOr: expected 42, got %v", got)
	}
	if got := store.GetIntOr("notfound", -1); got != -1 {
		t.Errorf("GetIntOr (not found): expected -1, got %v", got)
	}

	// Test GetFloat64 with numeric conversions
	if got := store.GetFloat64("float64"); got != 2.718 {
		t.Errorf("GetFloat64: expected 2.718, got %v", got)
	}
	if got := store.GetFloat64("float32"); got != float64(float32(3.14)) {
		t.Errorf("GetFloat64 (float32): expected %v, got %v", float64(float32(3.14)), got)
	}
	if got := store.GetFloat64("int"); got != 42.0 {
		t.Errorf("GetFloat64 (int): expected 42.0, got %v", got)
	}
	if got := store.GetFloat64("notfound"); got != 0.0 {
		t.Errorf("GetFloat64 (not found): expected 0.0, got %v", got)
	}

	// Test GetFloat64Or
	if got := store.GetFloat64Or("float64", -1.0); got != 2.718 {
		t.Errorf("GetFloat64Or: expected 2.718, got %v", got)
	}
	if got := store.GetFloat64Or("notfound", -1.0); got != -1.0 {
		t.Errorf("GetFloat64Or (not found): expected -1.0, got %v", got)
	}

	// Test GetBool
	if got := store.GetBool("bool"); got != true {
		t.Errorf("GetBool: expected true, got %v", got)
	}
	if got := store.GetBool("notfound"); got != false {
		t.Errorf("GetBool (not found): expected false, got %v", got)
	}
	if got := store.GetBool("string"); got != false {
		t.Errorf("GetBool (wrong type): expected false, got %v", got)
	}

	// Test GetBoolOr
	if got := store.GetBoolOr("bool", false); got != true {
		t.Errorf("GetBoolOr: expected true, got %v", got)
	}
	if got := store.GetBoolOr("notfound", true); got != true {
		t.Errorf("GetBoolOr (not found): expected true, got %v", got)
	}

	// Test GetSlice
	if got := store.GetSlice("slice"); got == nil || len(got) != 3 {
		t.Errorf("GetSlice: expected slice of length 3, got %v", got)
	}
	if got := store.GetSlice("notfound"); got != nil {
		t.Errorf("GetSlice (not found): expected nil, got %v", got)
	}

	// Test GetSliceOr
	defaultSlice := []any{"default"}
	if got := store.GetSliceOr("slice", defaultSlice); got == nil || len(got) != 3 {
		t.Errorf("GetSliceOr: expected slice of length 3, got %v", got)
	}
	if got := store.GetSliceOr("notfound", defaultSlice); !reflect.DeepEqual(got, defaultSlice) {
		t.Errorf("GetSliceOr (not found): expected %v, got %v", defaultSlice, got)
	}

	// Test GetMap
	if got := store.GetMap("map"); got == nil || got["name"] != "test" {
		t.Errorf("GetMap: expected map with name=test, got %v", got)
	}
	if got := store.GetMap("notfound"); got != nil {
		t.Errorf("GetMap (not found): expected nil, got %v", got)
	}

	// Test GetMapOr
	defaultMap := map[string]any{"default": true}
	if got := store.GetMapOr("map", defaultMap); got == nil || got["name"] != "test" {
		t.Errorf("GetMapOr: expected map with name=test, got %v", got)
	}
	if got := store.GetMapOr("notfound", defaultMap); !reflect.DeepEqual(got, defaultMap) {
		t.Errorf("GetMapOr (not found): expected %v, got %v", defaultMap, got)
	}
}

// TestSharedStoreUtilityMethods tests Has, Delete, Clear, Keys, Len
func TestSharedStoreUtilityMethods(t *testing.T) {
	store := NewSharedStore()

	// Test Has
	store.Set("key1", "value1")
	if !store.Has("key1") {
		t.Error("Has: expected true for existing key")
	}
	if store.Has("notfound") {
		t.Error("Has: expected false for non-existing key")
	}

	// Test Len
	store.Set("key2", "value2")
	store.Set("key3", "value3")
	if got := store.Len(); got != 3 {
		t.Errorf("Len: expected 3, got %v", got)
	}

	// Test Keys
	keys := store.Keys()
	if len(keys) != 3 {
		t.Errorf("Keys: expected 3 keys, got %v", len(keys))
	}
	keyMap := make(map[string]bool)
	for _, k := range keys {
		keyMap[k] = true
	}
	if !keyMap["key1"] || !keyMap["key2"] || !keyMap["key3"] {
		t.Errorf("Keys: missing expected keys, got %v", keys)
	}

	// Test Delete
	store.Delete("key2")
	if store.Has("key2") {
		t.Error("Delete: key2 should be deleted")
	}
	if got := store.Len(); got != 2 {
		t.Errorf("Len after delete: expected 2, got %v", got)
	}

	// Test Clear
	store.Clear()
	if got := store.Len(); got != 0 {
		t.Errorf("Len after clear: expected 0, got %v", got)
	}
	if store.Has("key1") || store.Has("key3") {
		t.Error("Clear: store should be empty")
	}
}

// TestSharedStoreBind tests the Bind and MustBind methods
func TestSharedStoreBind(t *testing.T) {
	store := NewSharedStore()

	// Test struct binding
	type User struct {
		ID   int      `json:"id"`
		Name string   `json:"name"`
		Tags []string `json:"tags"`
	}

	// Store as map
	store.Set("user", map[string]any{
		"id":   123,
		"name": "Alice",
		"tags": []string{"admin", "developer"},
	})

	var user User
	if err := store.Bind("user", &user); err != nil {
		t.Errorf("Bind: unexpected error: %v", err)
	}
	if user.ID != 123 || user.Name != "Alice" || len(user.Tags) != 2 {
		t.Errorf("Bind: incorrect values: %+v", user)
	}

	// Test binding to non-pointer
	var notPointer User
	if err := store.Bind("user", notPointer); err == nil {
		t.Error("Bind: should error on non-pointer")
	}

	// Test binding non-existent key
	var missing User
	if err := store.Bind("missing", &missing); err == nil {
		t.Error("Bind: should error on missing key")
	}

	// Test direct type match
	type Config struct {
		Timeout int
		Debug   bool
	}
	config := Config{Timeout: 30, Debug: true}
	store.Set("config", config)

	var boundConfig Config
	if err := store.Bind("config", &boundConfig); err != nil {
		t.Errorf("Bind (direct type): unexpected error: %v", err)
	}
	if boundConfig.Timeout != 30 || boundConfig.Debug != true {
		t.Errorf("Bind (direct type): incorrect values: %+v", boundConfig)
	}

	// Test MustBind success
	var user2 User
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("MustBind: unexpected panic: %v", r)
			}
		}()
		store.MustBind("user", &user2)
	}()

	if user2.ID != 123 {
		t.Error("MustBind: binding failed")
	}

	// Test MustBind panic
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Error("MustBind: expected panic for missing key")
			}
		}()
		var missing User
		store.MustBind("notfound", &missing)
	}()
}

// TestSharedStoreNumericConversions tests numeric type conversions
func TestSharedStoreNumericConversions(t *testing.T) {
	store := NewSharedStore()

	// Test all numeric types to int conversion
	store.Set("int8", int8(8))
	store.Set("int16", int16(16))
	store.Set("int32", int32(32))
	store.Set("int64", int64(64))
	store.Set("uint", uint(10))
	store.Set("uint8", uint8(8))
	store.Set("uint16", uint16(16))
	store.Set("uint32", uint32(32))
	store.Set("uint64", uint64(64))

	tests := []struct {
		key      string
		expected int
	}{
		{"int8", 8},
		{"int16", 16},
		{"int32", 32},
		{"int64", 64},
		{"uint", 10},
		{"uint8", 8},
		{"uint16", 16},
		{"uint32", 32},
		{"uint64", 64},
	}

	for _, test := range tests {
		if got := store.GetInt(test.key); got != test.expected {
			t.Errorf("GetInt(%s): expected %d, got %d", test.key, test.expected, got)
		}
	}

	// Test all numeric types to float64 conversion
	for _, test := range tests {
		if got := store.GetFloat64(test.key); got != float64(test.expected) {
			t.Errorf("GetFloat64(%s): expected %f, got %f", test.key, float64(test.expected), got)
		}
	}
}

// TestSharedStoreConcurrentTypeGetters tests concurrent access to type-specific getters
func TestSharedStoreConcurrentTypeGetters(t *testing.T) {
	store := NewSharedStore()
	store.Set("string", "test")
	store.Set("int", 42)
	store.Set("bool", true)

	var wg sync.WaitGroup

	// Concurrent type-specific reads
	for i := 0; i < 100; i++ {
		wg.Add(3)
		go func() {
			defer wg.Done()
			_ = store.GetString("string")
		}()
		go func() {
			defer wg.Done()
			_ = store.GetInt("int")
		}()
		go func() {
			defer wg.Done()
			_ = store.GetBool("bool")
		}()
	}

	// Concurrent writes and utility methods
	for i := 0; i < 50; i++ {
		wg.Add(4)
		go func(n int) {
			defer wg.Done()
			store.Set(fmt.Sprintf("concurrent_%d", n), n)
		}(i)
		go func(n int) {
			defer wg.Done()
			store.Has(fmt.Sprintf("concurrent_%d", n))
		}(i)
		go func() {
			defer wg.Done()
			_ = store.Keys()
		}()
		go func() {
			defer wg.Done()
			_ = store.Len()
		}()
	}

	wg.Wait()
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

// TestFlowConnectChaining tests that Connect method can be chained
func TestFlowConnectChaining(t *testing.T) {
	// Create nodes
	node1 := &TestNode{
		BaseNode: NewBaseNode(),
		name:     "node1",
		action:   "next",
	}

	node2 := &TestNode{
		BaseNode: NewBaseNode(),
		name:     "node2",
		action:   "final",
	}

	node3 := &TestNode{
		BaseNode: NewBaseNode(),
		name:     "node3",
		action:   DefaultAction,
	}

	// Test chaining
	flow := NewFlow(node1)
	flow.Connect(node1, "next", node2).
		Connect(node2, "final", node3)

	// Run the flow
	ctx := context.Background()
	shared := NewSharedStore()
	err := flow.Run(ctx, shared)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify execution order
	if !node1.executed {
		t.Error("node1 not executed")
	}
	if !node2.executed {
		t.Error("node2 not executed")
	}
	if !node3.executed {
		t.Error("node3 not executed")
	}
}

// TestFlowConnectChainingWithMultipleActions tests chaining with multiple actions from same node
func TestFlowConnectChainingWithMultipleActions(t *testing.T) {
	// Create nodes
	decisionNode := &TestNode{
		BaseNode: NewBaseNode(),
		name:     "decision",
		action:   "branch",
	}

	successNode := &TestNode{
		BaseNode: NewBaseNode(),
		name:     "success",
		action:   DefaultAction,
	}

	failNode := &TestNode{
		BaseNode: NewBaseNode(),
		name:     "fail",
		action:   DefaultAction,
	}

	retryNode := &TestNode{
		BaseNode: NewBaseNode(),
		name:     "retry",
		action:   DefaultAction,
	}

	// Test chaining with multiple actions from same node
	flow := NewFlow(decisionNode)
	flow.Connect(decisionNode, "branch", successNode).
		Connect(decisionNode, "fail", failNode).
		Connect(decisionNode, "retry", retryNode)

	// Run the flow
	ctx := context.Background()
	shared := NewSharedStore()
	err := flow.Run(ctx, shared)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify execution
	if !decisionNode.executed {
		t.Error("decisionNode not executed")
	}
	if !successNode.executed {
		t.Error("successNode not executed")
	}
	if failNode.executed {
		t.Error("failNode should not be executed")
	}
	if retryNode.executed {
		t.Error("retryNode should not be executed")
	}
}

// TestFlowConnectChainingBackwardsCompatibility tests that traditional usage still works
func TestFlowConnectChainingBackwardsCompatibility(t *testing.T) {
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

	// Test that traditional non-chaining usage still works
	flow := NewFlow(node1)
	flow.Connect(node1, "next", node2) // Traditional usage without chaining

	// Run the flow
	ctx := context.Background()
	shared := NewSharedStore()
	err := flow.Run(ctx, shared)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify execution
	if !node1.executed {
		t.Error("node1 not executed")
	}
	if !node2.executed {
		t.Error("node2 not executed")
	}
}

// BenchmarkBatchNodeProcessing benchmarks batch processing
