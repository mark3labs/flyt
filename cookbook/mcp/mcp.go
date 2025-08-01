package main

import (
	"context"
	"fmt"
	"time"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
)

// MCPClient wraps the MCP client for easier use
type MCPClient struct {
	client *client.Client
}

// NewMCPClient creates a new MCP client with an in-process server
func NewMCPClient(ctx context.Context) (*MCPClient, error) {
	srv := createMathServer()
	c, err := client.NewInProcessClient(srv)
	if err != nil {
		return nil, fmt.Errorf("failed to create in-process client: %w", err)
	}

	// Start the client
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := c.Start(ctxWithTimeout); err != nil {
		return nil, fmt.Errorf("failed to start client: %w", err)
	}

	// Initialize the client
	initRequest := mcp.InitializeRequest{}
	initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = mcp.Implementation{
		Name:    "Flyt MCP Client",
		Version: "1.0.0",
	}
	initRequest.Params.Capabilities = mcp.ClientCapabilities{}

	_, err = c.Initialize(ctx, initRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize: %w", err)
	}

	return &MCPClient{client: c}, nil
}

// GetTools returns available tools from the MCP server
func (m *MCPClient) GetTools(ctx context.Context) ([]mcp.Tool, error) {
	toolsRequest := mcp.ListToolsRequest{}
	toolsResp, err := m.client.ListTools(ctx, toolsRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to list tools: %w", err)
	}
	return toolsResp.Tools, nil
}

// CallTool calls a tool on the MCP server
func (m *MCPClient) CallTool(ctx context.Context, toolName string, arguments map[string]interface{}) (string, error) {
	toolRequest := mcp.CallToolRequest{}
	toolRequest.Params.Name = toolName
	toolRequest.Params.Arguments = arguments

	result, err := m.client.CallTool(ctx, toolRequest)
	if err != nil {
		return "", fmt.Errorf("failed to call tool: %w", err)
	}

	// Extract the result
	if len(result.Content) > 0 {
		if textContent, ok := result.Content[0].(mcp.TextContent); ok {
			return textContent.Text, nil
		}
	}

	return "", fmt.Errorf("unexpected result format")
}
