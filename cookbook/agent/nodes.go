// Package main implements an AI research agent using the Flyt workflow framework.
package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/mark3labs/flyt"
)

// NewDecideActionNode creates a node that decides whether to search for more information or provide an answer.
// It uses an LLM to analyze the question and context to make this decision.
func NewDecideActionNode(apiKey string) flyt.Node {
	return flyt.NewNode(
		flyt.WithPrepFunc(func(ctx context.Context, shared *flyt.SharedStore) (any, error) {
			question, _ := shared.Get("question")
			context, _ := shared.Get("context")
			if context == nil {
				context = "No previous search"
			}

			searchCount, _ := shared.Get("search_count")
			if searchCount == nil {
				searchCount = 0
			}

			return map[string]any{
				"question":     question,
				"context":      context,
				"search_count": searchCount,
			}, nil
		}),
		flyt.WithExecFunc(func(ctx context.Context, prepResult any) (any, error) {
			data := prepResult.(map[string]any)
			question := data["question"].(string)
			contextStr := data["context"].(string)
			searchCount := data["search_count"].(int)

			fmt.Println("🤔 Agent deciding what to do next...")

			forceAnswer := ""
			if searchCount >= 3 {
				forceAnswer = "\nIMPORTANT: You have already searched 3 times. You MUST choose 'answer' now with the information you have."
			}

			prompt := fmt.Sprintf(`You are a research assistant that can search the web.
Question: %s
Number of searches performed: %d

Previous Research:
%s

Analyze the search results above. If you can see specific names, dates, or facts that answer the question, choose "answer".
Only search again if you truly need different information.

Decide whether to:
1. "search" - Look up more information on the web (only if current results don't contain the answer)
2. "answer" - Answer the question with current knowledge

Respond with ONLY one of these two words followed by a colon and your reasoning.
If you choose "search", also provide the search query after another colon.
Format: ACTION:REASON:QUERY (if search)

Example responses:
- "search:Need current information about Nobel Prize winners:Nobel Prize Physics 2024 winners"
- "answer:I have enough information to answer the question"%s`, question, searchCount, contextStr, forceAnswer)

			// Call LLM with the captured API key
			response, err := CallLLM(apiKey, prompt)
			if err != nil {
				return nil, fmt.Errorf("LLM call failed: %w", err)
			}

			return map[string]any{
				"response":    response,
				"question":    question,
				"contextStr":  contextStr,
				"searchCount": searchCount,
			}, nil
		}),
		flyt.WithPostFunc(func(ctx context.Context, shared *flyt.SharedStore, prepResult, execResult any) (flyt.Action, error) {
			data := execResult.(map[string]any)
			response := data["response"].(string)

			// Parse the response
			parts := strings.Split(response, ":")
			if len(parts) < 2 {
				return "", fmt.Errorf("invalid LLM response format: %s", response)
			}

			action := strings.TrimSpace(strings.ToLower(parts[0]))
			// reason := strings.TrimSpace(parts[1]) // Currently unused but available if needed

			if action == "search" && len(parts) >= 3 {
				searchQuery := strings.TrimSpace(parts[2])
				shared.Set("search_query", searchQuery)
				fmt.Printf("🔍 Agent decided to search for: %s\n", searchQuery)

				searchCount := data["searchCount"].(int)
				shared.Set("search_count", searchCount+1)
			} else if action == "answer" {
				fmt.Println("💡 Agent decided to answer the question")
			}

			return flyt.Action(action), nil
		}),
	)
}

// NewSearchWebNode creates a node that searches the web for information
func NewSearchWebNode() flyt.Node {
	return flyt.NewNode(
		flyt.WithPrepFunc(func(ctx context.Context, shared *flyt.SharedStore) (any, error) {
			query, _ := shared.Get("search_query")
			braveKey, _ := shared.Get("brave_api_key")

			return map[string]any{
				"query":     query,
				"brave_key": braveKey,
			}, nil
		}),
		flyt.WithExecFunc(func(ctx context.Context, prepResult any) (any, error) {
			data := prepResult.(map[string]any)
			query := data["query"].(string)

			fmt.Printf("🌐 Searching the web for: %s\n", query)

			// Use real web search APIs
			var results string
			var err error

			// Try Brave first if API key is available
			if braveKey, ok := data["brave_key"].(string); ok && braveKey != "" {
				results, err = SearchWebBrave(query, braveKey)
				if err != nil {
					fmt.Printf("⚠️  Brave search failed: %v, falling back to DuckDuckGo\n", err)
					// Fall back to DuckDuckGo
					results, err = SearchWeb(query)
					if err != nil {
						return "", fmt.Errorf("all search methods failed: %w", err)
					}
				}
			} else {
				// Use default search (DuckDuckGo)
				results, err = SearchWeb(query)
				if err != nil {
					return "", fmt.Errorf("search failed: %w", err)
				}
			}

			return map[string]any{
				"results": results,
				"query":   query,
			}, nil
		}),
		flyt.WithPostFunc(func(ctx context.Context, shared *flyt.SharedStore, prepResult, execResult any) (flyt.Action, error) {
			data := execResult.(map[string]any)
			results := data["results"].(string)
			query := data["query"].(string)

			previousContext, _ := shared.Get("context")
			previousStr := ""
			if previousContext != nil {
				previousStr = previousContext.(string)
			}

			newContext := fmt.Sprintf("%s\n\nSEARCH: %s\nRESULTS: %s", previousStr, query, results)
			shared.Set("context", newContext)

			fmt.Println("📚 Found information, analyzing results...")

			return "decide", nil
		}),
	)
}

// NewAnswerQuestionNode creates a node that generates the final answer
func NewAnswerQuestionNode(apiKey string) flyt.Node {
	return flyt.NewNode(
		flyt.WithPrepFunc(func(ctx context.Context, shared *flyt.SharedStore) (any, error) {
			question, _ := shared.Get("question")
			context, _ := shared.Get("context")

			return map[string]any{
				"question": question,
				"context":  context,
			}, nil
		}),
		flyt.WithExecFunc(func(ctx context.Context, prepResult any) (any, error) {
			data := prepResult.(map[string]any)
			question := data["question"].(string)
			contextStr := ""
			if data["context"] != nil {
				contextStr = data["context"].(string)
			}

			fmt.Println("✍️ Crafting final answer...")

			prompt := fmt.Sprintf(`Based on the following information, answer the question.
Question: %s
Research: %s

Provide a comprehensive answer using the research results.`, question, contextStr)

			// Call LLM with the captured API key
			answer, err := CallLLM(apiKey, prompt)
			if err != nil {
				return nil, fmt.Errorf("failed to generate answer: %w", err)
			}

			return answer, nil
		}),
		flyt.WithPostFunc(func(ctx context.Context, shared *flyt.SharedStore, prepResult, execResult any) (flyt.Action, error) {
			answer := execResult.(string)

			shared.Set("answer", answer)
			fmt.Println("✅ Answer generated successfully")

			// End the flow
			return "done", nil
		}),
	)
}
