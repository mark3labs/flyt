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
func NewDecideActionNode(llm *LLM) flyt.Node {
	return flyt.NewNode(
		flyt.WithPrepFunc(func(ctx context.Context, shared *flyt.SharedStore) (any, error) {
			question := shared.GetString("question")
			context := shared.GetStringOr("context", "No previous search")
			searchCount := shared.GetInt("search_count")

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

			// Call LLM
			response, err := llm.Call(prompt)
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
func NewSearchWebNode(searcher Searcher) flyt.Node {
	return flyt.NewNode(
		flyt.WithPrepFunc(func(ctx context.Context, shared *flyt.SharedStore) (any, error) {
			query := shared.GetString("search_query")

			return map[string]any{
				"query": query,
			}, nil
		}),
		flyt.WithExecFunc(func(ctx context.Context, prepResult any) (any, error) {
			data := prepResult.(map[string]any)
			query := data["query"].(string)

			fmt.Printf("🌐 Searching the web for: %s\n", query)

			// Use the injected searcher
			results, err := searcher.Search(query)
			if err != nil {
				return "", fmt.Errorf("search failed: %w", err)
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

			previousContext := shared.GetString("context")

			newContext := fmt.Sprintf("%s\n\nSEARCH: %s\nRESULTS: %s", previousContext, query, results)
			shared.Set("context", newContext)

			fmt.Println("📚 Found information, analyzing results...")

			return "decide", nil
		}),
	)
}

// NewAnswerQuestionNode creates a node that generates the final answer
func NewAnswerQuestionNode(llm *LLM) flyt.Node {
	return flyt.NewNode(
		flyt.WithPrepFunc(func(ctx context.Context, shared *flyt.SharedStore) (any, error) {
			question := shared.GetString("question")
			context := shared.GetString("context")

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

Provide a short concise answer using the research results.`, question, contextStr)

			// Call LLM
			answer, err := llm.Call(prompt)
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
