# Tracing Example

This example demonstrates how to add distributed tracing to Flyt workflows using Langfuse.

## Features

- **Comprehensive tracing**: Traces the entire flow execution including prep, exec, and post phases
- **Hierarchical spans**: Creates nested spans for each node and phase
- **Metadata tracking**: Captures input, output, and custom metadata at each step
- **Error tracking**: Records errors and failures in traces
- **Dependency injection**: Tracer is injected through closures

## Prerequisites

- Go 1.21 or later
- Langfuse account (optional - runs in demo mode without it)

## Installation

```bash
cd cookbook/tracing
go mod tidy
```

## Configuration

Set the following environment variables to enable Langfuse tracing:

```bash
export LANGFUSE_HOST=https://cloud.langfuse.com  # or your self-hosted instance
export LANGFUSE_PUBLIC_KEY=your-public-key
export LANGFUSE_SECRET_KEY=your-secret-key
```

Without these variables, the example runs in demo mode and logs traces to the console.

## Usage

```bash
# Run with Langfuse tracing (requires environment variables)
go run .

# Run in demo mode (no Langfuse required)
go run .
```

## How It Works

### Flow Structure

The example creates a three-node pipeline:
1. **GreetingNode**: Creates a greeting message
2. **UppercaseNode**: Converts the greeting to uppercase
3. **ReverseNode**: Reverses the uppercase greeting

### Tracing Architecture

```
BasicGreetingFlow (Trace)
├── GreetingNode.prep (Span)
├── GreetingNode.exec (Span)
├── GreetingNode.post (Span)
├── UppercaseNode.prep (Span)
├── UppercaseNode.exec (Span)
├── UppercaseNode.post (Span)
├── ReverseNode.prep (Span)
├── ReverseNode.exec (Span)
└── ReverseNode.post (Span)
```

### Key Components

1. **Tracer**: Wraps the Langfuse client and manages traces/spans
2. **Trace**: Represents the overall flow execution
3. **Span**: Represents individual node phases (prep, exec, post)
4. **Metadata**: Captures context at each step

## Code Structure

- `main.go`: Main application with traced flow setup
- `tracer.go`: Tracer implementation using langfuse-go
- `README.md`: This file

## Tracing Pattern

Each node follows this pattern:

```go
func CreateNode(tracer *Tracer) flyt.Node {
    return flyt.NewNode(
        flyt.WithPrepFunc(func(ctx context.Context, shared *flyt.SharedStore) (any, error) {
            span := tracer.StartSpan("NodeName.prep", metadata)
            defer span.End(nil)
            // ... prep logic ...
        }),
        flyt.WithExecFunc(func(ctx context.Context, prepResult any) (any, error) {
            span := tracer.StartSpan("NodeName.exec", metadata)
            defer span.End(nil)
            // ... exec logic ...
        }),
        flyt.WithPostFunc(func(ctx context.Context, shared *flyt.SharedStore, prepResult, execResult any) (flyt.Action, error) {
            span := tracer.StartSpan("NodeName.post", metadata)
            defer span.End(nil)
            // ... post logic ...
        }),
    )
}
```

## Demo Mode Output

When running without Langfuse configuration:

```
🚀 Starting Flyt Tracing Example
⚠️  Warning: Langfuse environment variables not set
🔍 [TRACE START] BasicGreetingFlow
  📍 [SPAN START] GreetingNode.prep
  📍 [SPAN END] GreetingNode.prep - Duration: 100ms - Status: SUCCESS
  📍 [SPAN START] GreetingNode.exec
  📍 [SPAN END] GreetingNode.exec - Duration: 100ms - Status: SUCCESS
  ...
🔍 [TRACE END] BasicGreetingFlow - Duration: 500ms - Status: SUCCESS
✅ Flow completed successfully!
```

## Langfuse Dashboard

When properly configured, traces appear in your Langfuse dashboard with:
- Flow execution timeline
- Node-level performance metrics
- Input/output data at each step
- Error tracking and debugging information

## Benefits

1. **Observability**: Full visibility into flow execution
2. **Performance Analysis**: Identify bottlenecks and slow operations
3. **Debugging**: Track data flow and transformations
4. **Error Tracking**: Capture and analyze failures
5. **Zero Code Changes**: Tracing via dependency injection