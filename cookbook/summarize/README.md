# Flyt Summarize Example

A text summarization example demonstrating error handling and retry logic with Flyt.

## Features

- Text summarization using OpenAI's GPT-3.5-turbo
- Automatic retry on failures with configurable attempts
- Demonstrates the NewNode helper syntax
- Shows error handling patterns

## Prerequisites

- Go 1.21 or later
- OpenAI API key

## Setup

1. Set your OpenAI API key:
```bash
export OPENAI_API_KEY="your-api-key-here"
```

2. Install dependencies:
```bash
go mod tidy
```

## Running the Example

Basic usage:
```bash
go run .
```

To see retry behavior in action:
```bash
# This simulates failures to demonstrate retry logic
SIMULATE_FAILURE=true USE_FALLBACK_DEMO=true go run .
```

## How It Works

The example demonstrates:

1. **Simple node creation**: Using `flyt.NewNode()` with custom functions
2. **Error handling**: Automatic retries with `WithMaxRetries()` and `WithWait()`
3. **State management**: Storing input and output in SharedStore
4. **Fallback behavior**: Providing a default message when all retries fail

## Code Structure

- `main.go` - Main application with two node implementations
- `llm.go` - OpenAI API integration
- `go.mod` - Module dependencies

## Key Concepts

### Basic Summarization Node
```go
flyt.NewNode(
    flyt.WithPrepFunc(...),    // Read text from SharedStore
    flyt.WithExecFunc(...),    // Call LLM for summarization
    flyt.WithPostFunc(...),    // Store result
    flyt.WithMaxRetries(3),    // Retry up to 3 times
    flyt.WithWait(time.Second), // Wait between retries
)
```

### Error Handling
The framework automatically retries failed operations based on the node configuration. If all retries fail, you can provide a fallback in your main flow logic.