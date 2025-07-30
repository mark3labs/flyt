# Flyt

Flyt (Norwegian for "flow") is a minimalist workflow framework for Go, inspired by [Pocket Flow](https://github.com/The-Pocket/PocketFlow). It provides a simple graph-based abstraction for orchestrating tasks with just the essentials - no bloat, no dependencies, no vendor lock-in.

## Features

- **Ultra-Lightweight**: Core framework in ~250 lines, with optional helpers
- **Graph-based**: Model workflows as directed graphs of nodes
- **Flexible**: Support for sequential, branching, and looping workflows
- **Go-Idiomatic**: Use goroutines, channels, and other Go concurrency primitives naturally
- **Batch Helpers**: Convenient functions for common batch processing patterns
- **Retry support**: Built-in retry mechanism with configurable attempts and delays
- **Zero Dependencies**: Only uses Go standard library

## Installation

```bash
go get github.com/mark3labs/flyt
```

## Quick Start

### Basic Example

```go
package main

import (
    "fmt"
    "log"
    "github.com/mark3labs/flyt"
)

// Define a simple node
type HelloNode struct {
    *flyt.BaseNode
}

func (n *HelloNode) Exec(prepResult any) (any, error) {
    return "Hello, World!", nil
}

func (n *HelloNode) Post(shared flyt.Shared, prepResult, execResult any) (flyt.Action, error) {
    shared["message"] = execResult
    return flyt.DefaultAction, nil
}

func main() {
    // Create node and flow
    hello := &HelloNode{BaseNode: flyt.NewBaseNode()}
    flow := flyt.NewFlow(hello)
    
    // Run the flow
    shared := flyt.Shared{}
    if err := flow.Run(shared); err != nil {
        log.Fatal(err)
    }
    
    fmt.Println(shared["message"]) // Output: Hello, World!
}
```

### Workflow with Branching

```go
// Create nodes for a simple approval workflow
submit := &SubmitNode{BaseNode: flyt.NewBaseNode()}
review := &ReviewNode{BaseNode: flyt.NewBaseNode()}
approve := &ApproveNode{BaseNode: flyt.NewBaseNode()}
reject := &RejectNode{BaseNode: flyt.NewBaseNode()}

// Create workflow with branching
flow := flyt.NewFlow(submit)
flow.Connect(submit, "review", review)
flow.Connect(review, "approve", approve)
flow.Connect(review, "reject", reject)

// Run the workflow
shared := flyt.Shared{"amount": 500.0}
if err := flow.Run(shared); err != nil {
    log.Fatal(err)
}
```

### Concurrent Processing with Goroutines

```go
type ConcurrentFetchNode struct {
    *flyt.BaseNode
}

func (n *ConcurrentFetchNode) Prep(shared flyt.Shared) (any, error) {
    return shared["urls"], nil
}

func (n *ConcurrentFetchNode) Exec(prepResult any) (any, error) {
    urls := prepResult.([]string)
    
    // Use goroutines for concurrent fetching
    type result struct {
        url  string
        data string
        err  error
    }
    
    results := make(chan result, len(urls))
    
    // Launch goroutines
    for _, url := range urls {
        go func(u string) {
            // Fetch data (simulated)
            data, err := fetchData(u)
            results <- result{url: u, data: data, err: err}
        }(url)
    }
    
    // Collect results
    var allResults []map[string]string
    for i := 0; i < len(urls); i++ {
        r := <-results
        if r.err != nil {
            return nil, r.err
        }
        allResults = append(allResults, map[string]string{
            "url":  r.url,
            "data": r.data,
        })
    }
    
    return allResults, nil
}

func (n *ConcurrentFetchNode) Post(shared flyt.Shared, prepResult, execResult any) (flyt.Action, error) {
    shared["results"] = execResult
    return flyt.DefaultAction, nil
}
```

### Batch Processing

```go
type BatchProcessNode struct {
    *flyt.BaseNode
}

func (n *BatchProcessNode) Prep(shared flyt.Shared) (any, error) {
    return shared["items"], nil
}

func (n *BatchProcessNode) Exec(prepResult any) (any, error) {
    items := prepResult.([]string)
    results := make([]string, len(items))
    
    // Process items in parallel using goroutines
    var wg sync.WaitGroup
    for i, item := range items {
        wg.Add(1)
        go func(idx int, itm string) {
            defer wg.Done()
            results[idx] = processItem(itm)
        }(i, item)
    }
    wg.Wait()
    
    return results, nil
}

func (n *BatchProcessNode) Post(shared flyt.Shared, prepResult, execResult any) (flyt.Action, error) {
    shared["processed"] = execResult
    return flyt.DefaultAction, nil
}
```

## Core Concepts

### Node

A Node is the basic building block that implements a three-phase lifecycle:

1. **Prep**: Read and preprocess data from shared store
2. **Exec**: Execute the main logic (with optional retries)
3. **Post**: Process results and decide the next action

```go
type Node interface {
    Prep(shared Shared) (any, error)
    Exec(prepResult any) (any, error)
    Post(shared Shared, prepResult, execResult any) (Action, error)
    SetParams(params Params)
    GetParams() Params
}
```

### Flow

A Flow orchestrates nodes by:
- Connecting nodes with action-based transitions
- Managing the execution order
- Propagating parameters to child nodes

### Shared Store

A map (`flyt.Shared`) for sharing data between nodes throughout the workflow execution.

### Actions

Nodes return actions that determine the next step in the workflow. Use `flyt.DefaultAction` for simple sequential flows or custom actions for branching.

## Advanced Features

### Batch Helpers

Flyt provides convenient helper functions for common batch processing patterns:

#### NewBatchNode

Creates a node that automatically processes multiple items:

```go
// Sequential batch processing
processItems := flyt.NewBatchNode(func(item any) (any, error) {
    // Process each item
    return processItem(item), nil
}, false) // false = sequential

// Concurrent batch processing  
fetchURLs := flyt.NewBatchNode(func(item any) (any, error) {
    url := item.(string)
    return fetchData(url), nil
}, true) // true = concurrent

// Use in a flow
flow := flyt.NewFlow(loadDataNode)
flow.Connect(loadDataNode, flyt.DefaultAction, processItems)

// The batch node automatically looks for data in shared["items"], shared["data"], or shared["batch"]
shared := flyt.Shared{
    "items": []string{"item1", "item2", "item3"},
}
flow.Run(shared)
// Results are stored in shared["results"]
```

#### NewBatchFlow

Creates a flow that runs another flow multiple times with different parameters:

```go
// Create a flow for processing a single file
processFileFlow := flyt.NewFlow(readNode)
processFileFlow.Connect(readNode, flyt.DefaultAction, transformNode)

// Create batch flow that runs for each file
batchFlow := flyt.NewBatchFlow(processFileFlow, func(shared flyt.Shared) ([]flyt.Params, error) {
    files := shared["files"].([]string)
    params := make([]flyt.Params, len(files))
    for i, file := range files {
        params[i] = flyt.Params{"filename": file}
    }
    return params, nil
}, true) // true = process files concurrently

// Run it
shared := flyt.Shared{"files": []string{"file1.txt", "file2.txt", "file3.txt"}}
batchFlow.Run(shared)
// shared["batch_count"] contains the number of iterations
```

**Thread Safety Note**: When using concurrent batch processing:
- `NewBatchNode` with `concurrent=true` is safe as each item is processed independently
- `NewBatchFlow` with `concurrent=true` requires thread-safe nodes since the same flow instance is reused
- If your nodes maintain state, consider using sequential processing or creating separate flow instances for each concurrent execution

### Retry Mechanism

```go
node := &MyNode{
    BaseNode: flyt.NewBaseNode(
        flyt.WithMaxRetries(3),
        flyt.WithWait(time.Second),
    ),
}
```

### Context and Cancellation

Since Flyt embraces Go's idioms, you can use context for cancellation:

```go
type ContextAwareNode struct {
    *flyt.BaseNode
}

func (n *ContextAwareNode) Exec(prepResult any) (any, error) {
    // Get context from params or shared
    ctx := n.GetParams()["ctx"].(context.Context)
    
    select {
    case result := <-doWork():
        return result, nil
    case <-ctx.Done():
        return nil, ctx.Err()
    }
}

// Usage
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

node := &ContextAwareNode{BaseNode: flyt.NewBaseNode()}
node.SetParams(flyt.Params{"ctx": ctx})
```

### Worker Pool Pattern

```go
func (n *WorkerPoolNode) Exec(prepResult any) (any, error) {
    items := prepResult.([]Item)
    
    // Create channels for job distribution
    jobs := make(chan Item, len(items))
    results := make(chan Result, len(items))
    
    // Start worker pool
    numWorkers := runtime.NumCPU()
    var wg sync.WaitGroup
    
    for w := 0; w < numWorkers; w++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for item := range jobs {
                results <- processItem(item)
            }
        }()
    }
    
    // Send jobs
    for _, item := range items {
        jobs <- item
    }
    close(jobs)
    
    // Wait and collect results
    go func() {
        wg.Wait()
        close(results)
    }()
    
    var allResults []Result
    for result := range results {
        allResults = append(allResults, result)
    }
    
    return allResults, nil
}
```

## Design Philosophy

Flyt follows Go's philosophy of simplicity:

1. **No Magic**: Just interfaces and structs - what you see is what you get
2. **Use Go's Strengths**: Leverage goroutines, channels, and other concurrency primitives directly
3. **Minimal API Surface**: Only the essential abstractions needed for workflow orchestration
4. **Composable**: Nodes and flows can be easily combined and reused

## Why Flyt?

- **You already know how to use it**: If you know Go, you know Flyt
- **No abstraction overhead**: Use goroutines and channels directly, not through layers of abstraction
- **Tiny footprint**: ~250 lines of code you can read and understand in minutes
- **Zero dependencies**: No external packages, just the Go standard library

## License

MIT License - see LICENSE file for details.