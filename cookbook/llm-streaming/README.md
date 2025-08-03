# LLM Streaming Example

This example demonstrates how to stream responses from an LLM using the Flyt workflow framework.

## Features

- **Real-time streaming**: Stream LLM responses token by token as they arrive
- **Two modes**:
  - Interactive chat mode (default)
  - Single prompt mode with `--single` flag
- **Dependency injection**: LLM client is injected through closures
- **Context-based streaming**: Proper handling using Go contexts

## Prerequisites

- Go 1.21 or later
- OpenAI API key

## Installation

```bash
cd cookbook/llm-streaming
go mod tidy
```

## Usage

### Interactive Chat Mode (Default)

```bash
# Set API key as environment variable
export OPENAI_API_KEY=your-api-key
go run .

# Or provide API key as argument
go run . your-api-key
```

In interactive mode:
- Type your message and press ENTER to send
- Watch the response stream in real-time
- Type "exit" to quit

### Single Prompt Mode

```bash
# Use default prompt
go run . your-api-key --single

# Use custom prompt
go run . your-api-key --single "Explain quantum computing in simple terms"
```

## How It Works

1. **LLM Struct**: Encapsulates the OpenAI API client with streaming support
2. **Stream Node**: Handles the streaming response and displays it in real-time
3. **Context Management**: Uses Go contexts for proper stream handling
4. **Channel-based Streaming**: Returns chunks through a channel for real-time processing

## Code Structure

- `main.go`: Main application with flow setup
- `llm.go`: LLM client implementation with streaming support
- `README.md`: This file

## Key Components

### LLM Client
The `LLM` struct provides:
- `Stream()`: Real streaming from OpenAI API using Server-Sent Events

### Stream Node
The streaming node:
- Streams responses chunk by chunk
- Displays each chunk as it arrives
- Collects the full response for storage

### Flow Structure
- **Interactive mode**: Input Node → Stream Node → (loop back)
- **Single mode**: Stream Node only