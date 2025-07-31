package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/mark3labs/flyt"
	"github.com/mark3labs/mcp-go/client"
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

// CallLLMWithFunctions calls OpenAI with function calling
func CallLLMWithFunctions(apiKey, question string, tools []mcp.Tool) (map[string]interface{}, error) {
	url := "https://api.openai.com/v1/chat/completions"

	// Convert MCP tools to OpenAI function format
	functions := make([]map[string]interface{}, 0, len(tools))
	for _, tool := range tools {
		function := map[string]interface{}{
			"name":        tool.Name,
			"description": tool.Description,
			"parameters":  tool.InputSchema,
		}
		functions = append(functions, function)
	}

	requestBody := map[string]interface{}{
		"model": "gpt-4o-mini",
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": "You are a helpful assistant that performs mathematical calculations using the available functions.",
			},
			{
				"role":    "user",
				"content": question,
			},
		},
		"functions":     functions,
		"function_call": "auto",
		"temperature":   0.3,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var response struct {
		Choices []struct {
			Message struct {
				Content      string `json:"content"`
				FunctionCall *struct {
					Name      string `json:"name"`
					Arguments string `json:"arguments"`
				} `json:"function_call"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(response.Choices) == 0 {
		return nil, fmt.Errorf("no choices in response")
	}

	choice := response.Choices[0]
	if choice.Message.FunctionCall == nil {
		return nil, fmt.Errorf("no function call in response")
	}

	// Parse the function arguments
	var args map[string]interface{}
	if err := json.Unmarshal([]byte(choice.Message.FunctionCall.Arguments), &args); err != nil {
		return nil, fmt.Errorf("failed to parse function arguments: %w", err)
	}

	return map[string]interface{}{
		"tool_name":  choice.Message.FunctionCall.Name,
		"parameters": args,
	}, nil
}

// NewGetToolsNode creates a node that discovers available tools from the MCP server
func NewGetToolsNode(mcpClient *MCPClient) flyt.Node {
	return flyt.NewNode(
		flyt.WithExecFunc(func(ctx context.Context, input any) (any, error) {
			fmt.Println("🔍 Getting available tools...")

			tools, err := mcpClient.GetTools(ctx)
			if err != nil {
				return nil, fmt.Errorf("failed to get tools: %w", err)
			}

			// Format tool information for display
			var toolInfo []string
			for i, tool := range tools {
				properties := tool.InputSchema.Properties
				required := tool.InputSchema.Required

				var params []string
				for paramName, paramInfo := range properties {
					paramType := "unknown"
					if pInfo, ok := paramInfo.(map[string]interface{}); ok {
						if t, ok := pInfo["type"].(string); ok {
							paramType = t
						}
					}
					reqStatus := "(Optional)"
					for _, req := range required {
						if req == paramName {
							reqStatus = "(Required)"
							break
						}
					}
					params = append(params, fmt.Sprintf("    - %s (%s): %s", paramName, paramType, reqStatus))
				}

				toolInfo = append(toolInfo, fmt.Sprintf("[%d] %s\n  Description: %s\n  Parameters:\n%s",
					i+1, tool.Name, tool.Description, strings.Join(params, "\n")))
			}

			return map[string]interface{}{
				"tools":     tools,
				"tool_info": strings.Join(toolInfo, "\n\n"),
			}, nil
		}),
		flyt.WithPostFunc(func(ctx context.Context, shared *flyt.SharedStore, prepResult, execResult any) (flyt.Action, error) {
			data := execResult.(map[string]interface{})
			shared.Set("tools", data["tools"])
			shared.Set("tool_info", data["tool_info"])
			return flyt.Action("decide"), nil
		}),
	)
}

// NewDecideToolNode creates a node that uses LLM with function calling to select appropriate tool
func NewDecideToolNode(question, apiKey string) flyt.Node {
	return flyt.NewNode(
		flyt.WithPrepFunc(func(ctx context.Context, shared *flyt.SharedStore) (any, error) {
			tools, _ := shared.Get("tools")
			return tools.([]mcp.Tool), nil
		}),
		flyt.WithExecFunc(func(ctx context.Context, prepResult any) (any, error) {
			tools := prepResult.([]mcp.Tool)

			fmt.Println("🤔 Using OpenAI function calling to analyze question...")
			result, err := CallLLMWithFunctions(apiKey, question, tools)
			if err != nil {
				return nil, fmt.Errorf("failed to call LLM with functions: %w", err)
			}

			toolName := result["tool_name"].(string)
			parameters := result["parameters"].(map[string]interface{})

			fmt.Printf("💡 OpenAI selected tool: %s\n", toolName)
			fmt.Printf("🔢 Function parameters: %v\n", parameters)

			return map[string]interface{}{
				"tool_name":  toolName,
				"parameters": parameters,
			}, nil
		}),
		flyt.WithPostFunc(func(ctx context.Context, shared *flyt.SharedStore, prepResult, execResult any) (flyt.Action, error) {
			data := execResult.(map[string]interface{})
			shared.Set("tool_name", data["tool_name"])
			shared.Set("parameters", data["parameters"])
			return flyt.Action("execute"), nil
		}),
	)
}

// NewExecuteToolNode creates a node that executes the selected tool with parameters
func NewExecuteToolNode(mcpClient *MCPClient) flyt.Node {
	return flyt.NewNode(
		flyt.WithPrepFunc(func(ctx context.Context, shared *flyt.SharedStore) (any, error) {
			toolName, _ := shared.Get("tool_name")
			parameters, _ := shared.Get("parameters")
			return map[string]interface{}{
				"tool_name":  toolName.(string),
				"parameters": parameters.(map[string]interface{}),
			}, nil
		}),
		flyt.WithExecFunc(func(ctx context.Context, prepResult any) (any, error) {
			data := prepResult.(map[string]interface{})
			toolName := data["tool_name"].(string)
			parameters := data["parameters"].(map[string]interface{})

			fmt.Printf("🔧 Executing tool '%s' with parameters: %v\n", toolName, parameters)
			result, err := mcpClient.CallTool(ctx, toolName, parameters)
			if err != nil {
				return nil, fmt.Errorf("failed to call tool: %w", err)
			}

			fmt.Printf("\n✅ Final Answer: %v\n", result)
			return result, nil
		}),
	)
}

func main() {
	// Parse command line flags
	var question string
	var apiKey string
	flag.StringVar(&question, "q", "What is 10 plus 20?", "Question to ask")
	flag.StringVar(&apiKey, "key", "", "OpenAI API key (or set OPENAI_API_KEY env var)")
	flag.Parse()

	// Get API key from flag or environment
	if apiKey == "" {
		apiKey = os.Getenv("OPENAI_API_KEY")
	}
	if apiKey == "" {
		log.Fatal("Please provide OpenAI API key via -key flag or OPENAI_API_KEY environment variable")
	}

	fmt.Printf("🤔 Processing question: %s\n", question)

	ctx := context.Background()

	// Create MCP client with in-process server
	mcpClient, err := NewMCPClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create MCP client: %v", err)
	}

	// Create nodes
	getToolsNode := NewGetToolsNode(mcpClient)
	decideNode := NewDecideToolNode(question, apiKey)
	executeNode := NewExecuteToolNode(mcpClient)

	// Create flow
	flow := flyt.NewFlow(getToolsNode)

	// Connect nodes
	flow.Connect(getToolsNode, "decide", decideNode)
	flow.Connect(decideNode, "execute", executeNode)

	// Create shared store
	shared := flyt.NewSharedStore()

	// Run the flow
	if err := flow.Run(ctx, shared); err != nil {
		log.Fatalf("Flow failed: %v", err)
	}
}
