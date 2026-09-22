package agent

import "github.com/alexnakagama/noryn/internal/llm"

type TokenCounter interface {
	Count(message llm.Message) int
}

type EstimateTokenCounter struct{}

func (EstimateTokenCounter) Count(message llm.Message) int {
	return len(message.Content) / 4
}
