// Package main implements an AI research agent using the Flyt workflow framework.
package main

import (
	"github.com/mark3labs/flyt"
)

func CreateAgentFlow(apiKey string) *flyt.Flow {
	decide := NewDecideActionNode()
	search := NewSearchWebNode()
	answer := NewAnswerQuestionNode()

	params := flyt.Params{"api_key": apiKey}
	decide.SetParams(params)
	answer.SetParams(params)

	flow := flyt.NewFlow(decide)

	flow.Connect(decide, "search", search)
	flow.Connect(decide, "answer", answer)
	flow.Connect(search, "decide", decide)

	return flow
}
func NewSharedState() *flyt.SharedStore {
	return flyt.NewSharedStore()
}
