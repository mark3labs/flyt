package main

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// createMathServer creates an MCP server with mathematical operations
func createMathServer() *server.MCPServer {
	s := server.NewMCPServer(
		"Math Operations Server",
		"1.0.0",
		server.WithToolCapabilities(false),
	)

	// Add tool
	addTool := mcp.NewTool("add",
		mcp.WithDescription("Add two numbers together"),
		mcp.WithNumber("a",
			mcp.Required(),
			mcp.Description("First number"),
		),
		mcp.WithNumber("b",
			mcp.Required(),
			mcp.Description("Second number"),
		),
	)

	s.AddTool(addTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		a, err := request.RequireFloat("a")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		b, err := request.RequireFloat("b")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		result := a + b
		return mcp.NewToolResultText(fmt.Sprintf("%.2f", result)), nil
	})

	// Subtract tool
	subtractTool := mcp.NewTool("subtract",
		mcp.WithDescription("Subtract b from a"),
		mcp.WithNumber("a",
			mcp.Required(),
			mcp.Description("First number"),
		),
		mcp.WithNumber("b",
			mcp.Required(),
			mcp.Description("Second number"),
		),
	)

	s.AddTool(subtractTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		a, err := request.RequireFloat("a")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		b, err := request.RequireFloat("b")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		result := a - b
		return mcp.NewToolResultText(fmt.Sprintf("%.2f", result)), nil
	})

	// Multiply tool
	multiplyTool := mcp.NewTool("multiply",
		mcp.WithDescription("Multiply two numbers together"),
		mcp.WithNumber("a",
			mcp.Required(),
			mcp.Description("First number"),
		),
		mcp.WithNumber("b",
			mcp.Required(),
			mcp.Description("Second number"),
		),
	)

	s.AddTool(multiplyTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		a, err := request.RequireFloat("a")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		b, err := request.RequireFloat("b")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		result := a * b
		return mcp.NewToolResultText(fmt.Sprintf("%.2f", result)), nil
	})

	// Divide tool
	divideTool := mcp.NewTool("divide",
		mcp.WithDescription("Divide a by b"),
		mcp.WithNumber("a",
			mcp.Required(),
			mcp.Description("First number"),
		),
		mcp.WithNumber("b",
			mcp.Required(),
			mcp.Description("Second number"),
		),
	)

	s.AddTool(divideTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		a, err := request.RequireFloat("a")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		b, err := request.RequireFloat("b")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		if b == 0 {
			return mcp.NewToolResultError("Division by zero is not allowed"), nil
		}

		result := a / b
		return mcp.NewToolResultText(fmt.Sprintf("%.2f", result)), nil
	})

	return s
}
