# Flyt Chat Example

An interactive chat application demonstrating how to use Flyt for conversational AI with self-looping flows.

## Features

- Interactive command-line chat interface
- Maintains conversation history
- Uses OpenAI's GPT-3.5-turbo model
- Demonstrates self-looping flows in Flyt

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

## Running the Chat

```bash
go run .
```

Type your messages and press Enter. The AI will respond and maintain the conversation context.

Type `exit` to end the conversation.

## How It Works

The example demonstrates:

1. **Self-looping flow**: The chat node connects to itself with a "continue" action
2. **State management**: Conversation history is stored in SharedStore
3. **Conditional flow control**: Returns "end" action to exit the loop
4. **Using closures**: The API key is captured in the node's closure

## Code Structure

- `main.go` - Main application with the chat node implementation
- `llm.go` - OpenAI API integration
- `go.mod` - Module dependencies