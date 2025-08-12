# Tracing Example

This example demonstrates how to add distributed tracing to Flyt workflows using the CloudWeGo Eino-ext Langfuse integration.

## Features

- **Comprehensive tracing**: Traces the entire flow execution including prep, exec, and post phases
- **Hierarchical spans**: Creates nested spans for each node and phase
- **LLM generation tracking**: Specialized tracing for model calls with token usage
- **Metadata tracking**: Captures input, output, and custom metadata at each step
- **Error tracking**: Records errors and failures in traces
- **Dependency injection**: Tracer is injected through closures
- **Batched uploads**: Efficient batch processing of trace events
- **Automatic retries**: Built-in retry logic for failed uploads

## Prerequisites

- Go 1.23 or later
- Langfuse account (optional - runs in demo mode without it)
- OpenAI API key (optional - uses mock responses without it)

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

To enable real LLM calls (using GPT-4):

```bash
export OPENAI_API_KEY=your-openai-api-key
```

Without these variables, the example runs in demo mode with console logging and mock LLM responses.

## Usage

```bash
# Run with Langfuse tracing (requires environment variables)
go run .

# Run in demo mode (no Langfuse required)
go run .
```

## How It Works

### Flow Structure

The example creates a four-node pipeline:
1. **GreetingNode**: Creates a greeting message
2. **UppercaseNode**: Converts the greeting to uppercase
3. **ReverseNode**: Reverses the uppercase greeting
4. **SummarizeNode**: Uses GPT-4 to create a playful summary (includes LLM generation tracing)

### Tracing Architecture

```
BasicGreetingFlow (Trace)
├── GreetingNode (Span)
│   ├── GreetingNode.prep (Child Span)
│   ├── GreetingNode.exec (Child Span)
│   └── GreetingNode.post (Child Span)
├── UppercaseNode (Span)
│   ├── UppercaseNode.prep (Child Span)
│   ├── UppercaseNode.exec (Child Span)
│   └── UppercaseNode.post (Child Span)
├── ReverseNode (Span)
│   ├── ReverseNode.prep (Child Span)
│   ├── ReverseNode.exec (Child Span)
│   └── ReverseNode.post (Child Span)
└── SummarizeNode (Span)
    ├── SummarizeNode.prep (Child Span)
    ├── SummarizeNode.exec (Child Span)
    │   └── gpt-4-summary (Generation) - LLM-specific tracing
    └── SummarizeNode.post (Child Span)
```

Each node creates a parent span that encompasses all its phases, with prep, exec, and post phases as child spans. The LLM node uses specialized generation tracing for model calls, capturing token usage and model parameters.

### Key Components

1. **Tracer**: Wraps the Langfuse client and manages traces/spans
2. **Trace**: Represents the overall flow execution
3. **Span**: Represents individual node phases (prep, exec, post)
4. **Metadata**: Captures context at each step

## Code Structure

- `main.go`: Main application with traced flow setup
- `tracer.go`: Tracer implementation using CloudWeGo's eino-ext Langfuse client
- `README.md`: This file

## Implementation Details

This example uses the `github.com/cloudwego/eino-ext/libs/acl/langfuse` client, which provides:
- Event-based API for creating traces, spans, and events
- Built-in batching and queueing for efficient API calls
- Configurable flush intervals and batch sizes
- Automatic retry logic for failed requests
- JSON serialization for input/output data

## Tracing Pattern

Each node follows this hierarchical pattern:

```go
func CreateNode(tracer *Tracer) flyt.Node {
    var nodeSpan *Span
    
    return flyt.NewNode(
        flyt.WithPrepFunc(func(ctx context.Context, shared *flyt.SharedStore) (any, error) {
            // Create node-level parent span on first execution
            if nodeSpan == nil {
                nodeSpan = tracer.StartSpan("NodeName", nodeMetadata, input)
            }
            
            // Create prep phase as child span
            span := tracer.StartChildSpan(nodeSpan, "NodeName.prep", phaseMetadata, input)
            defer span.EndWithOutput(output, nil)
            // ... prep logic ...
        }),
        flyt.WithExecFunc(func(ctx context.Context, prepResult any) (any, error) {
            // Create exec phase as child span
            span := tracer.StartChildSpan(nodeSpan, "NodeName.exec", phaseMetadata, input)
            defer span.EndWithOutput(output, nil)
            // ... exec logic ...
        }),
        flyt.WithPostFunc(func(ctx context.Context, shared *flyt.SharedStore, prepResult, execResult any) (flyt.Action, error) {
            // Create post phase as child span
            span := tracer.StartChildSpan(nodeSpan, "NodeName.post", phaseMetadata, input)
            defer span.EndWithOutput(output, nil)
            
            // End the node-level parent span
            defer nodeSpan.EndWithOutput(finalOutput, nil)
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