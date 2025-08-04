![Flyt Logo](flyt-logo.png)

# Flyt

*Norwegian for "flow" • Pronounced "fleet"*

A minimalist workflow framework for Go with zero dependencies inspired by [Pocket Flow](https://github.com/The-Pocket/PocketFlow).

## Table of Contents

- [Installation](#installation)
- [Getting Started](#getting-started)
  - [Using the Project Template (Recommended)](#using-the-project-template-recommended)
  - [Manual Setup](#manual-setup)
- [Core Concepts](#core-concepts)
  - [Nodes](#nodes)
  - [Actions](#actions)
  - [Flows](#flows)
  - [Shared Store](#shared-store)
- [Intermediate Patterns](#intermediate-patterns)
  - [Configuration via Closures](#configuration-via-closures)
  - [Error Handling & Retries](#error-handling--retries)
  - [Fallback on Failure](#fallback-on-failure)
  - [Conditional Branching](#conditional-branching)
- [Advanced Usage](#advanced-usage)
  - [Custom Node Types](#custom-node-types)
    - [RetryableNode Interface](#retryablenode-interface)
  - [Batch Processing](#batch-processing)
    - [Advanced Batch Configuration](#advanced-batch-configuration)
    - [Batch Error Handling](#batch-error-handling)
  - [Batch Flows](#batch-flows)
  - [Nested Flows](#nested-flows)
  - [Flow as Node](#flow-as-node)
  - [Worker Pool](#worker-pool)
  - [Utility Functions](#utility-functions)
    - [ToSlice](#toslice)
- [Best Practices](#best-practices)
- [Examples](#examples)
- [License](#license)

## Installation

```bash
go get github.com/mark3labs/flyt
```

## Getting Started

### Using the Project Template (Recommended)

The fastest way to start a new Flyt project is using the official template:

```bash
# Create a new project from the template
git clone https://github.com/mark3labs/flyt-project-template my-flyt-project
cd my-flyt-project

# Remove the template git history and start fresh
rm -rf .git
git init

# Install dependencies
go mod tidy

# Run the example
go run main.go
```

The template provides a starting point for your Flyt project with a basic structure and example code.

### Manual Setup

```go
package main

import (
    "context"
    "fmt"
    "github.com/mark3labs/flyt"
)

func main() {
    // Create a simple node using the helper
    node := flyt.NewNode(
        flyt.WithExecFunc(func(ctx context.Context, prepResult any) (any, error) {
            fmt.Println("Hello, Flyt!")
            return "done", nil
        }),
    )

    // Run it
    ctx := context.Background()
    shared := flyt.NewSharedStore()
    
    action, err := flyt.Run(ctx, node, shared)
    if err != nil {
        panic(err)
    }
    fmt.Printf("Completed with action: %s\n", action)
}
```

## Core Concepts

### Nodes

Nodes are the building blocks. Each node has three phases:

1. **Prep** - Read from shared store and prepare data
2. **Exec** - Execute main logic (can be retried)
3. **Post** - Process results and decide next action

```go
// Using the helper
node := flyt.NewNode(
    flyt.WithPrepFunc(func(ctx context.Context, shared *flyt.SharedStore) (any, error) {
        data, _ := shared.Get("input")
        return data, nil
    }),
    flyt.WithExecFunc(func(ctx context.Context, prepResult any) (any, error) {
        // Process data
        return "result", nil
    }),
    flyt.WithPostFunc(func(ctx context.Context, shared *flyt.SharedStore, prepResult, execResult any) (flyt.Action, error) {
        shared.Set("output", execResult)
        return flyt.DefaultAction, nil
    }),
)
```

### Actions

Actions are strings returned by a node's Post phase that determine what happens next:

```go
func (n *MyNode) Post(ctx context.Context, shared *flyt.SharedStore, prepResult, execResult any) (flyt.Action, error) {
    if execResult.(bool) {
        return "success", nil  // Go to node connected with "success"
    }
    return "retry", nil       // Go to node connected with "retry"
}
```

The default action is `flyt.DefaultAction` (value: "default"). If no connection exists for an action, the flow ends.

### Flows

Connect nodes to create workflows:

```go
// Create nodes
validateNode := createValidateNode()
processNode := createProcessNode()
errorNode := createErrorNode()

// Build flow with action-based routing
flow := flyt.NewFlow(validateNode)
flow.Connect(validateNode, "valid", processNode)    // If validation succeeds
flow.Connect(validateNode, "invalid", errorNode)    // If validation fails
flow.Connect(processNode, "done", nil)              // End flow after processing

// Run flow
err := flow.Run(ctx, shared)
```

### Shared Store

Thread-safe data sharing between nodes:

```go
shared := flyt.NewSharedStore()

// Set and get individual values
shared.Set("key", "value")
value, ok := shared.Get("key")

// Get all data as a map (returns a copy)
allData := shared.GetAll()

// Merge multiple values at once
shared.Merge(map[string]any{
    "user_id": 123,
    "config": map[string]any{"timeout": 30},
})
```

### Error Handling & Retries

Add retry logic to handle transient failures:

```go
node := flyt.NewNode(
    flyt.WithExecFunc(func(ctx context.Context, prepResult any) (any, error) {
        // This will be retried up to 3 times
        return callFlakeyAPI()
    }),
    flyt.WithMaxRetries(3),
    flyt.WithWait(time.Second),
)
```

### Fallback on Failure

Handle failures gracefully by implementing the `FallbackNode` interface:

```go
type CachedAPINode struct {
    *flyt.BaseNode
    cache map[string]any
}

func (n *CachedAPINode) ExecFallback(prepResult any, err error) (any, error) {
    // Return cached data when API fails
    key := prepResult.(string)
    if cached, ok := n.cache[key]; ok {
        return cached, nil
    }
    // Return default value if no cache
    return map[string]any{"status": "unavailable"}, nil
}

func (n *CachedAPINode) Exec(ctx context.Context, prepResult any) (any, error) {
    key := prepResult.(string)
    data, err := callAPI(key)
    if err == nil {
        n.cache[key] = data // Update cache on success
    }
    return data, err
}
```

The `ExecFallback` method is called automatically after all retries are exhausted, allowing you to provide degraded functionality, cached results, or default values.

### Conditional Branching

Control flow based on results:

```go
decisionNode := flyt.NewNode(
    flyt.WithExecFunc(func(ctx context.Context, prepResult any) (any, error) {
        value := prepResult.(int)
        return value > 100, nil
    }),
    flyt.WithPostFunc(func(ctx context.Context, shared *flyt.SharedStore, prepResult, execResult any) (flyt.Action, error) {
        if execResult.(bool) {
            return "high", nil
        }
        return "low", nil
    }),
)

flow := flyt.NewFlow(decisionNode)
flow.Connect(decisionNode, "high", highNode)
flow.Connect(decisionNode, "low", lowNode)
```

## Advanced Usage

### Custom Node Types

For complex nodes with state, create custom types:

```go
type RateLimitedNode struct {
    *flyt.BaseNode
    limiter *rate.Limiter
}

func NewRateLimitedNode(rps int) *RateLimitedNode {
    return &RateLimitedNode{
        BaseNode: flyt.NewBaseNode(),
        limiter:  rate.NewLimiter(rate.Limit(rps), 1),
    }
}

func (n *RateLimitedNode) Exec(ctx context.Context, prepResult any) (any, error) {
    if err := n.limiter.Wait(ctx); err != nil {
        return nil, err
    }
    // Process with rate limiting
    return process(prepResult)
}
```

#### RetryableNode Interface

For custom retry logic, implement the `RetryableNode` interface:

```go
type CustomRetryNode struct {
    *flyt.BaseNode
    attempts int
}

func (n *CustomRetryNode) GetMaxRetries() int {
    // Dynamic retry count based on state
    if n.attempts > 5 {
        return 0 // Stop retrying after 5 total attempts
    }
    return 3
}

func (n *CustomRetryNode) GetWait() time.Duration {
    // Exponential backoff
    return time.Duration(n.attempts) * time.Second
}

func (n *CustomRetryNode) Exec(ctx context.Context, prepResult any) (any, error) {
    n.attempts++
    return callAPI(prepResult)
}
```

### Batch Processing

Process multiple items concurrently:

```go
// Simple batch node for processing items
processFunc := func(ctx context.Context, item any) (any, error) {
    // Process each item
    return fmt.Sprintf("processed: %v", item), nil
}

batchNode := flyt.NewBatchNode(processFunc, true) // true for concurrent
shared.Set("items", []string{"item1", "item2", "item3"})
```

#### Advanced Batch Configuration

For more control over batch processing:

```go
config := &flyt.BatchConfig{
    BatchSize:   10,        // Process 10 items at a time
    Concurrency: 5,         // Use 5 concurrent workers
    ItemsKey:    "data",    // Custom key for input items
    ResultsKey:  "output",  // Custom key for results
    CountKey:    "total",   // Custom key for processed count
}

processFunc := func(ctx context.Context, item any) (any, error) {
    return processItem(item)
}

batchNode := flyt.NewBatchNodeWithConfig(processFunc, true, config)
```

#### Batch Error Handling

Batch operations aggregate errors:

```go
action, err := flyt.Run(ctx, batchNode, shared)
if err != nil {
    if batchErr, ok := err.(*flyt.BatchError); ok {
        // Access individual errors
        for i, e := range batchErr.Errors {
            if e != nil {
                fmt.Printf("Item %d failed: %v\n", i, e)
            }
        }
    }
}
```

### Batch Flows

Run the same flow multiple times with different parameters:

```go
// Create a flow factory - returns a new flow instance for each iteration
flowFactory := func() *flyt.Flow {
    validateNode := flyt.NewNode(
        flyt.WithPrepFunc(func(ctx context.Context, shared *flyt.SharedStore) (any, error) {
            // Each flow has its own SharedStore with merged FlowInputs
            userID, _ := shared.Get("user_id")
            email, _ := shared.Get("email")
            return map[string]any{"user_id": userID, "email": email}, nil
        }),
        flyt.WithExecFunc(func(ctx context.Context, prepResult any) (any, error) {
            data := prepResult.(map[string]any)
            // Process user data
            return processUser(data), nil
        }),
    )
    return flyt.NewFlow(validateNode)
}

// Define input parameters for each flow iteration
// Each FlowInputs map is merged into that flow's isolated SharedStore
batchFunc := func(ctx context.Context, shared *flyt.SharedStore) ([]flyt.FlowInputs, error) {
    // Could fetch from database, API, etc.
    return []flyt.FlowInputs{
        {"user_id": 1, "email": "user1@example.com"},
        {"user_id": 2, "email": "user2@example.com"},
        {"user_id": 3, "email": "user3@example.com"},
    }, nil
}

// Create and run batch flow
batchFlow := flyt.NewBatchFlow(flowFactory, batchFunc, true) // true for concurrent
err := batchFlow.Run(ctx, shared)

// Each flow runs in isolation with its own SharedStore containing the FlowInputs
```

### Nested Flows

Compose flows for complex workflows:

```go
// Sub-flow for data validation
validationFlow := createValidationFlow()

// Main flow
mainFlow := flyt.NewFlow(fetchNode)
mainFlow.Connect(fetchNode, "validate", validationFlow)
mainFlow.Connect(validationFlow, flyt.DefaultAction, processNode)
```

### Flow as Node

Flows implement the Node interface and can be used anywhere a node is expected:

```go
// Create a reusable flow
func createProcessingFlow() *flyt.Flow {
    validateNode := createValidateNode()
    transformNode := createTransformNode()
    
    flow := flyt.NewFlow(validateNode)
    flow.Connect(validateNode, "valid", transformNode)
    return flow
}

// Use the flow as a node in another flow
processingFlow := createProcessingFlow()
mainFlow := flyt.NewFlow(fetchNode)
mainFlow.Connect(fetchNode, flyt.DefaultAction, processingFlow) // Flow used as node
mainFlow.Connect(processingFlow, flyt.DefaultAction, saveNode)
```

### Worker Pool

For custom concurrent task management:

```go
// Create a worker pool with 10 workers
pool := flyt.NewWorkerPool(10)

// Submit tasks
for _, item := range items {
    item := item // Capture loop variable
    pool.Submit(func() {
        // Process item
        result := processItem(item)
        // Store result safely
        mu.Lock()
        results = append(results, result)
        mu.Unlock()
    })
}

// Wait for all tasks to complete
pool.Wait()

// Clean up
pool.Close()
```

### Utility Functions

#### ToSlice

Convert various types to slices for batch processing:

```go
// Convert different types to []any
items1 := flyt.ToSlice([]string{"a", "b", "c"})
items2 := flyt.ToSlice([]int{1, 2, 3})
items3 := flyt.ToSlice("single item") // Returns []any{"single item"}

// Useful for batch processing with mixed types
shared.Set("items", flyt.ToSlice(data))
```

## Best Practices

1. **Single Responsibility**: Each node should do one thing well
2. **Idempotency**: Nodes should be idempotent when possible
3. **Error Handling**: Always handle errors appropriately
4. **Context Awareness**: Respect context cancellation
5. **Concurrency Safety**: Don't share node instances across flows

## Examples

Check out the [cookbook](cookbook/) directory for complete, real-world examples:

- [Agent](cookbook/agent/) - AI agent with web search capabilities using LLM and search providers
- [Chat](cookbook/chat/) - Interactive chat application with conversation history
- [LLM Streaming](cookbook/llm-streaming/) - Real-time streaming of LLM responses with OpenAI SSE
- [MCP](cookbook/mcp/) - Model Context Protocol integration with OpenAI function calling
- [Summarize](cookbook/summarize/) - Text summarization with error handling and retries
- [Tracing](cookbook/tracing/) - Distributed tracing with Langfuse for observability

## License

MIT
