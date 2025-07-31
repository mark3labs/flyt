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

// CallLLM calls OpenAI to process the prompt
func CallLLM(apiKey, prompt string) (string, error) {
	url := "https://api.openai.com/v1/chat/completions"

	requestBody := map[string]interface{}{
		"model": "gpt-4o-mini",
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": "You are a helpful assistant that analyzes questions and decides which mathematical tool to use. Always respond in the exact YAML format requested.",
			},
			{
				"role":    "user",
				"content": prompt,
			},
		},
		"temperature": 0.3,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(jsonBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var response struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}

	return response.Choices[0].Message.Content, nil
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

// NewDecideToolNode creates a node that uses LLM to analyze the question and select appropriate tool
func NewDecideToolNode(question, apiKey string) flyt.Node {
	return flyt.NewNode(
		flyt.WithPrepFunc(func(ctx context.Context, shared *flyt.SharedStore) (any, error) {
			toolInfo, _ := shared.Get("tool_info")
			return toolInfo.(string), nil
		}),
		flyt.WithExecFunc(func(ctx context.Context, prepResult any) (any, error) {
			toolInfo := prepResult.(string)

			prompt := fmt.Sprintf(`### CONTEXT
You are an assistant that can use tools via Model Context Protocol (MCP).

### ACTION SPACE
%s

### TASK
Answer this question: "%s"

## NEXT ACTION
Analyze the question, extract any numbers or parameters, and decide which tool to use.
Return your response in this format:

thinking: |
    <your step-by-step reasoning about what the question is asking and what numbers to extract>
tool: <name of the tool to use>
reason: <why you chose this tool>
parameters:
    <parameter_name>: <parameter_value>
    <parameter_name>: <parameter_value>

IMPORTANT: 
1. Extract numbers from the question properly
2. Use proper indentation (4 spaces) for multi-line fields
3. Use the | character for multi-line text fields`, toolInfo, question)

			fmt.Println("🤔 Analyzing question and deciding which tool to use...")
			response, err := CallLLM(apiKey, prompt)
			if err != nil {
				return nil, fmt.Errorf("failed to call LLM: %w", err)
			}

			// Parse the response to extract tool and parameters
			toolName, parameters := parseToolDecision(response)

			fmt.Printf("💡 Selected tool: %s\n", toolName)
			fmt.Printf("🔢 Extracted parameters: %v\n", parameters)

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

// parseToolDecision extracts tool name and parameters from LLM response
func parseToolDecision(response string) (string, map[string]interface{}) {
	var toolName string
	parameters := make(map[string]interface{})

	lines := strings.Split(response, "\n")
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		if strings.HasPrefix(line, "tool:") {
			toolName = strings.TrimSpace(strings.TrimPrefix(line, "tool:"))
		} else if strings.HasPrefix(line, "parameters:") {
			// Parse parameters
			i++
			for i < len(lines) && strings.HasPrefix(lines[i], "    ") {
				paramLine := strings.TrimSpace(lines[i])
				if strings.Contains(paramLine, ":") {
					parts := strings.SplitN(paramLine, ":", 2)
					paramName := strings.TrimSpace(parts[0])
					paramValue := strings.TrimSpace(parts[1])

					// Try to parse as number
					var value interface{}
					// For very large numbers, we need to handle them as float64
					var f float64
					if _, err := fmt.Sscanf(paramValue, "%f", &f); err == nil {
						value = f
					} else {
						value = paramValue
					}
					parameters[paramName] = value
				}
				i++
			}
			i-- // Back up one since the outer loop will increment
		}
	}

	return toolName, parameters
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
