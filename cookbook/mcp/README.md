# Flyt MCP (Model Context Protocol) Example

This example demonstrates how to integrate Flyt with the Model Context Protocol (MCP) using the `github.com/mark3labs/mcp-go` library. It shows how to build an agent that can discover and use tools exposed by an MCP server, all running in a single process using the InProcess transport.

## Overview

The example creates a math operations MCP server and a Flyt-based client that:

1. **Discovers Tools** - Automatically discovers available math operations from the MCP server
2. **Analyzes Questions** - Uses an LLM to analyze questions and select appropriate tools
3. **Executes Tools** - Calls the selected tool with extracted parameters via MCP

## Features

- **In-Process Architecture**: Server and client run in the same process for simplicity
- **Tool Discovery**: Automatically discovers available tools from the MCP server
- **Dynamic Tool Selection**: Uses an LLM to analyze questions and select appropriate tools
- **Math Operations**: Supports addition, subtraction, multiplication, and division
- **Type-Safe Parameters**: MCP handles parameter validation and type conversion

## Architecture

```mermaid
flowchart LR
    start[Get Tools] -->|discover| decide[Decide Tool]
    decide -->|analyze| execute[Execute Tool]
    execute -->|result| done[Done]
```

The flow consists of three nodes:

1. **Get Tools Node**: Connects to the MCP server and retrieves available tools
2. **Decide Node**: Uses an LLM to analyze the question and select the appropriate tool
3. **Execute Node**: Calls the selected tool with extracted parameters

## How to Run

First, set your OpenAI API key:

```bash
export OPENAI_API_KEY="your-api-key-here"
```

Then run the example:

```bash
cd cookbook/mcp
go run main.go -q "What is 10 plus 20?"
```

Or use the default question:

```bash
go run main.go
```

You can also pass the API key directly:

```bash
go run main.go -key "your-api-key" -q "What is 100 divided by 4?"
```

## Code Structure

The example is contained in a single file (`main.go`) that includes:

- **MCP Server Creation** (`createMathServer`): Defines math operation tools
- **MCP Client** (`MCPClient`): Wraps the MCP client for easier use
- **Flyt Nodes**: Three nodes that implement the workflow
- **LLM Integration**: Simulated LLM for demo purposes

## Key Concepts

### MCP Integration

The example shows how to:
- Create an MCP server with tool definitions
- Use InProcess transport for single-process operation
- Initialize the MCP client with proper protocol version
- List available tools from the server
- Call tools with typed parameters
- Handle tool results

### Flyt Flow

The example demonstrates:
- Using `NewNode` helper for cleaner node creation
- Passing data between nodes via `SharedStore`
- Error handling at each stage
- Dynamic flow based on tool discovery

### InProcess Transport

The InProcess transport allows the MCP server and client to run in the same process:
- No need for separate server process
- Direct function calls instead of IPC
- Simplified deployment and testing
- Same API as other transports

## Extending the Example

To add more tools:

1. Add new tool definitions in `createMathServer()`
2. Implement the tool handler function
3. The client will automatically discover the new tools

Example:
```go
// Power tool
powerTool := mcp.NewTool("power",
    mcp.WithDescription("Raise a to the power of b"),
    mcp.WithNumber("a", mcp.Required(), mcp.Description("Base number")),
    mcp.WithNumber("b", mcp.Required(), mcp.Description("Exponent")),
)

s.AddTool(powerTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    a, _ := request.RequireFloat("a")
    b, _ := request.RequireFloat("b")
    result := math.Pow(a, b)
    return mcp.NewToolResultText(fmt.Sprintf("%.2f", result)), nil
})
```

## Dependencies

- `github.com/mark3labs/flyt` - The Flyt workflow engine
- `github.com/mark3labs/mcp-go` - Go implementation of MCP

## Notes

- The example uses OpenAI's GPT-4o-mini model for tool selection
- The InProcess transport is ideal for single-application scenarios
- For distributed systems, use stdio or HTTP transports
- Error handling is simplified for clarity. Production code should be more robust
- The MCP server and client run in the same process, simplifying deployment