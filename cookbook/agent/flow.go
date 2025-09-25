// Package main implements an AI research agent using the Flyt workflow framework.
package main

import (
	"github.com/mark3labs/flyt"
)

func CreateAgentFlow(llm *LLM, searcher Searcher) *flyt.Flow {
	decide := NewDecideActionNode(llm)
	search := NewSearchWebNode(searcher)
	answer := NewAnswerQuestionNode(llm)

	return flyt.NewFlow(decide).
		Connect(decide, "search", search).
		Connect(decide, "answer", answer).
		Connect(search, "decide", decide)
}
func NewSharedState() *flyt.SharedStore {
	return flyt.NewSharedStore()
}
