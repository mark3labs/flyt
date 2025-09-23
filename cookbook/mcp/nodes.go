package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/mark3labs/flyt"
	"github.com/mark3labs/mcp-go/mcp"
)

// NewGetToolsNode creates a node that discovers available tools from the MCP server
func NewGetToolsNode(mcpClient *MCPClient) flyt.Node {
	return flyt.NewNode(
		flyt.WithExecFuncAny(func(ctx context.Context, input any) (any, error) {
			fmt.Println("🔍 Getting available tools...")

			tools, err := mcpClient.GetTools(ctx)
			if err != nil {
				return flyt.R(nil), fmt.Errorf("failed to get tools: %w", err)
			}

			fmt.Printf("📦 Found %d tools: ", len(tools))
			toolNames := make([]string, len(tools))
			for i, tool := range tools {
				toolNames[i] = tool.Name
			}
			fmt.Println(strings.Join(toolNames, ", "))

			return flyt.R(tools), nil
		}),
		flyt.WithPostFuncAny(func(ctx context.Context, shared *flyt.SharedStore, prepResult, execResult any) (flyt.Action, error) {
			if result, ok := execResult.(flyt.Result); ok {
				shared.Set("tools", result.Value())
			} else {
				shared.Set("tools", execResult)
			}
			return flyt.Action("decide"), nil
		}),
	)
}

// NewDecideToolNode creates a node that uses LLM with function calling to select appropriate tool
func NewDecideToolNode(llm *LLM) flyt.Node {
	return flyt.NewNode(
		flyt.WithPrepFuncAny(func(ctx context.Context, shared *flyt.SharedStore) (any, error) {
			tools, _ := shared.Get("tools")
			messages, ok := shared.Get("messages")
			if !ok {
				// Initialize messages if not present
				messages = []map[string]interface{}{}
			}
			return flyt.R(map[string]interface{}{
				"tools":    tools.([]mcp.Tool),
				"messages": messages.([]map[string]interface{}),
			}), nil
		}),
		flyt.WithExecFuncAny(func(ctx context.Context, prepResult any) (any, error) {
			prepRes, ok := prepResult.(flyt.Result)
			if !ok {
				// Handle legacy case
				prepRes = flyt.R(prepResult)
			}
			data := prepRes.MustMap()
			tools := data["tools"].([]mcp.Tool)
			messages := data["messages"].([]map[string]interface{})

			fmt.Println("🤔 Calling OpenAI with function calling...")
			result, err := llm.CallWithFunctions(messages, tools)
			if err != nil {
				return flyt.R(nil), fmt.Errorf("failed to call LLM with functions: %w", err)
			}

			// Add the assistant's message to history
			message := result["message"]
			// Convert the message struct to a map
			messageMap := map[string]interface{}{
				"role":    "assistant",
				"content": "",
			}

			// Use reflection or type assertion to extract message content
			if msg, ok := message.(struct {
				Role         string `json:"role"`
				Content      string `json:"content"`
				FunctionCall *struct {
					Name      string `json:"name"`
					Arguments string `json:"arguments"`
				} `json:"function_call"`
			}); ok {
				messageMap["role"] = msg.Role
				messageMap["content"] = msg.Content
				if msg.FunctionCall != nil {
					messageMap["function_call"] = map[string]interface{}{
						"name":      msg.FunctionCall.Name,
						"arguments": msg.FunctionCall.Arguments,
					}
				}
			}

			messages = append(messages, messageMap)

			// Check if there's a function call
			if functionCall, ok := result["function_call"]; ok {
				fc := functionCall.(map[string]interface{})
				fmt.Printf("💡 OpenAI wants to call function: %s\n", fc["name"])
				fmt.Printf("🔢 Function parameters: %v\n", fc["arguments"])

				return flyt.R(map[string]interface{}{
					"function_call": fc,
					"messages":      messages,
				}), nil
			}

			// No function call, just a regular message
			fmt.Printf("💬 Assistant response: %s\n", messageMap["content"])
			return flyt.R(map[string]interface{}{
				"messages": messages,
				"done":     true,
			}), nil
		}),
		flyt.WithPostFuncAny(func(ctx context.Context, shared *flyt.SharedStore, prepResult, execResult any) (flyt.Action, error) {
			var data map[string]interface{}
			if result, ok := execResult.(flyt.Result); ok {
				data = result.MustMap()
			} else {
				data = execResult.(map[string]interface{})
			}
			shared.Set("messages", data["messages"])

			if _, ok := data["done"]; ok {
				// Conversation is complete
				return flyt.DefaultAction, nil
			}

			// There's a function call to execute
			fc := data["function_call"].(map[string]interface{})
			shared.Set("function_call", fc)
			return flyt.Action("execute"), nil
		}),
	)
}

// NewExecuteToolNode creates a node that executes the selected tool with parameters
func NewExecuteToolNode(mcpClient *MCPClient) flyt.Node {
	return flyt.NewNode(
		flyt.WithPrepFuncAny(func(ctx context.Context, shared *flyt.SharedStore) (any, error) {
			functionCall, _ := shared.Get("function_call")
			messages, _ := shared.Get("messages")
			return flyt.R(map[string]interface{}{
				"function_call": functionCall.(map[string]interface{}),
				"messages":      messages.([]map[string]interface{}),
			}), nil
		}),
		flyt.WithExecFuncAny(func(ctx context.Context, prepResult any) (any, error) {
			prepRes, ok := prepResult.(flyt.Result)
			if !ok {
				prepRes = flyt.R(prepResult)
			}
			data := prepRes.MustMap()
			fc := data["function_call"].(map[string]interface{})
			messages := data["messages"].([]map[string]interface{})

			toolName := fc["name"].(string)
			parameters := fc["arguments"].(map[string]interface{})

			fmt.Printf("🔧 Executing tool '%s' with parameters: %v\n", toolName, parameters)
			result, err := mcpClient.CallTool(ctx, toolName, parameters)
			if err != nil {
				return nil, fmt.Errorf("failed to call tool: %w", err)
			}

			fmt.Printf("📊 Tool result: %v\n", result)

			// Format the function call output message for OpenAI
			// OpenAI expects function results as assistant messages with a specific format
			functionOutput := map[string]interface{}{
				"role":    "function",
				"name":    toolName,
				"content": result,
			}

			// Add function output to messages
			messages = append(messages, functionOutput)

			return flyt.R(map[string]interface{}{
				"messages": messages,
			}), nil
		}),
		flyt.WithPostFuncAny(func(ctx context.Context, shared *flyt.SharedStore, prepResult, execResult any) (flyt.Action, error) {
			var data map[string]interface{}
			if result, ok := execResult.(flyt.Result); ok {
				data = result.MustMap()
			} else {
				data = execResult.(map[string]interface{})
			}
			shared.Set("messages", data["messages"])
			// Go back to decide node to continue the conversation
			return flyt.Action("decide"), nil
		}),
	)
}
