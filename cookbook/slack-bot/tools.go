package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// ToolExecutor handles tool execution
type ToolExecutor struct {
	httpClient *http.Client
}

func NewToolExecutor() *ToolExecutor {
	return &ToolExecutor{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// ExecuteTool executes a tool based on its name and arguments
func (te *ToolExecutor) ExecuteTool(toolName string, arguments string) (string, error) {
	switch toolName {
	case "calculator":
		return te.executeCalculator(arguments)
	case "chuck_norris_fact":
		return te.executeChuckNorrisFact(arguments)
	default:
		return "", fmt.Errorf("unknown tool: %s", toolName)
	}
}

// Calculator tool implementation
func (te *ToolExecutor) executeCalculator(arguments string) (string, error) {
	var args struct {
		Expression string `json:"expression"`
	}

	if err := json.Unmarshal([]byte(arguments), &args); err != nil {
		return "", fmt.Errorf("failed to parse calculator arguments: %w", err)
	}

	result, err := evaluateExpression(args.Expression)
	if err != nil {
		return fmt.Sprintf("Error: %v", err), nil
	}

	return fmt.Sprintf("Result: %v", result), nil
}

// Simple expression evaluator for basic math operations
func evaluateExpression(expr string) (float64, error) {
	expr = strings.TrimSpace(expr)

	// Handle basic functions
	if strings.HasPrefix(expr, "sqrt(") && strings.HasSuffix(expr, ")") {
		inner := expr[5 : len(expr)-1]
		val, err := evaluateExpression(inner)
		if err != nil {
			return 0, err
		}
		return math.Sqrt(val), nil
	}

	if strings.HasPrefix(expr, "pow(") && strings.HasSuffix(expr, ")") {
		inner := expr[4 : len(expr)-1]
		parts := strings.Split(inner, ",")
		if len(parts) != 2 {
			return 0, fmt.Errorf("pow requires two arguments")
		}
		base, err := evaluateExpression(strings.TrimSpace(parts[0]))
		if err != nil {
			return 0, err
		}
		exp, err := evaluateExpression(strings.TrimSpace(parts[1]))
		if err != nil {
			return 0, err
		}
		return math.Pow(base, exp), nil
	}

	// Handle basic arithmetic operations
	// Try multiplication and division first
	for _, op := range []string{"*", "/", "+", "-"} {
		if idx := strings.LastIndex(expr, op); idx > 0 && idx < len(expr)-1 {
			left, err := evaluateExpression(expr[:idx])
			if err != nil {
				continue // Try next operator
			}
			right, err := evaluateExpression(expr[idx+1:])
			if err != nil {
				continue // Try next operator
			}

			switch op {
			case "+":
				return left + right, nil
			case "-":
				return left - right, nil
			case "*":
				return left * right, nil
			case "/":
				if right == 0 {
					return 0, fmt.Errorf("division by zero")
				}
				return left / right, nil
			}
		}
	}

	// Handle parentheses
	if strings.HasPrefix(expr, "(") && strings.HasSuffix(expr, ")") {
		return evaluateExpression(expr[1 : len(expr)-1])
	}

	// Try to parse as number
	val, err := strconv.ParseFloat(expr, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid expression: %s", expr)
	}
	return val, nil
}

// Chuck Norris fact tool implementation
func (te *ToolExecutor) executeChuckNorrisFact(arguments string) (string, error) {
	var args struct {
		Category string `json:"category"`
	}

	// Parse arguments if provided
	if arguments != "" && arguments != "{}" {
		if err := json.Unmarshal([]byte(arguments), &args); err != nil {
			return "", fmt.Errorf("failed to parse chuck norris arguments: %w", err)
		}
	}

	// Use Chuck Norris API
	url := "https://api.chucknorris.io/jokes/random"
	if args.Category != "" && args.Category != "random" {
		url = fmt.Sprintf("https://api.chucknorris.io/jokes/random?category=%s", args.Category)
	}

	resp, err := te.httpClient.Get(url)
	if err != nil {
		// Fallback to hardcoded facts if API fails
		return te.getHardcodedChuckNorrisFact(), nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return te.getHardcodedChuckNorrisFact(), nil
	}

	var joke struct {
		Value string `json:"value"`
	}

	if err := json.Unmarshal(body, &joke); err != nil {
		return te.getHardcodedChuckNorrisFact(), nil
	}

	if joke.Value == "" {
		return te.getHardcodedChuckNorrisFact(), nil
	}

	return joke.Value, nil
}

// Hardcoded Chuck Norris facts as fallback
func (te *ToolExecutor) getHardcodedChuckNorrisFact() string {
	facts := []string{
		"Chuck Norris doesn't read books. He stares them down until he gets the information he wants.",
		"Chuck Norris counted to infinity. Twice.",
		"Chuck Norris can divide by zero.",
		"Chuck Norris doesn't use web standards. The web conforms to Chuck Norris.",
		"Chuck Norris can unit test an entire application with a single assert statement.",
		"Chuck Norris's keyboard doesn't have a Ctrl key because nothing controls Chuck Norris.",
		"Chuck Norris can access private methods.",
		"Chuck Norris's code doesn't have bugs. It has features that weren't requested yet.",
		"Chuck Norris doesn't need garbage collection because he doesn't create garbage.",
		"Chuck Norris's programs never have memory leaks. Memory is too scared to leak.",
		"Chuck Norris doesn't debug. The bugs fix themselves out of fear.",
		"Chuck Norris can solve NP-complete problems in O(1) time.",
		"Chuck Norris writes code that optimizes itself.",
		"Chuck Norris's database queries are so fast, they return results before you execute them.",
		"Chuck Norris doesn't use version control. His code is always perfect the first time.",
	}

	// Return a random fact (using current time as seed)
	index := int(time.Now().UnixNano()) % len(facts)
	return facts[index]
}

// ExecuteToolCalls executes multiple tool calls and returns their results
func ExecuteToolCalls(toolCalls []ToolCall) (map[string]string, error) {
	executor := NewToolExecutor()
	results := make(map[string]string)

	for _, toolCall := range toolCalls {
		result, err := executor.ExecuteTool(toolCall.Function.Name, toolCall.Function.Arguments)
		if err != nil {
			results[toolCall.ID] = fmt.Sprintf("Error executing tool: %v", err)
		} else {
			results[toolCall.ID] = result
		}
	}

	return results, nil
}
