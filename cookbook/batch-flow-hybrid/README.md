# Batch Flow Hybrid Example

This example demonstrates how to combine `BatchNode` with `FlowFactory` to create a "batch flow" pattern - where each item in a batch is processed through its own flow instance.

## Key Concepts

### BatchFlowNode
A custom node that combines:
- **BatchNode**: For managing batch processing with Prep/Exec/Post phases
- **FlowFactory**: For creating isolated flow instances for each item
- **Lambda functions**: For defining processing logic inline

### Use Cases

1. **Process items through complex pipelines**: When each item needs to go through a multi-step flow
2. **Isolation**: Each item gets its own flow instance with isolated state
3. **Conditional flows**: Different items can use different flows based on their properties
4. **Concurrent processing**: Items are processed in parallel with configurable concurrency

## Examples

### Example 1: Basic BatchFlowNode
Creates a reusable `BatchFlowNode` that processes items through a flow factory:
```go
batchFlowNode := NewBatchFlowNode("BatchProcessor", flowFactory, 3)
```

### Example 2: Inline Lambda BatchNode
Shows how to create a BatchNode with inline lambda functions that use a flow factory:
```go
simpleBatchFlow := flyt.NewBatchNode(
    "SimpleBatchFlow",
    func(ctx context.Context, shared flyt.Shared) []flyt.Result { /* prep */ },
    func(ctx context.Context, shared flyt.Shared, item flyt.Result) flyt.Result { 
        // Create and run flow for this item
        flow := flowFactory()
        // ...
    },
    func(ctx context.Context, shared flyt.Shared, results []flyt.Result) flyt.Params { /* post */ },
)
```

### Example 3: Conditional Flow Selection
Demonstrates selecting different flows based on item properties:
```go
// In Exec function
if itemType == "priority" {
    flow = createPriorityFlow()  // Use priority processing
} else {
    flow = flowFactory()         // Use standard processing
}
```

## Running the Example

```bash
go run main.go
```

## Key Advantages

1. **Reusability**: Flow factories allow reusing complex flow definitions
2. **Isolation**: Each item gets its own flow instance and shared store
3. **Flexibility**: Can use different flows for different item types
4. **Concurrency**: Built-in concurrent processing with configurable worker pool
5. **Error handling**: Automatic error collection and reporting

## When to Use This Pattern

Use this pattern when you need to:
- Process batch items through complex, multi-step workflows
- Maintain isolation between items during processing
- Apply different processing logic based on item properties
- Leverage existing flow definitions for batch processing

This effectively recreates the "batch flow" concept using the simplified BatchNode API combined with FlowFactory for maximum flexibility.