package agent

import "github.com/alexnakagama/noryn/internal/llm"

const maxContextTokens = 32_000

type ContextManager struct {
	maxTokens    int
	tokenCounter TokenCounter
}

func NewContextManager() *ContextManager {
	return &ContextManager{
		maxTokens:    maxContextTokens,
		tokenCounter: EstimateTokenCounter{},
	}
}

func (c *ContextManager) Build(history []llm.Message) []llm.Message {
	if len(history) == 0 {
		return history
	}

	totalTokens := 0
	start := len(history)

	for i := len(history) - 1; i >= 0; i-- {
		messageTokens := c.tokenCounter.Count(history[i])

		if totalTokens+messageTokens > c.maxTokens {
			break
		}

		totalTokens += messageTokens
		start = i
	}

	// Dont start in the middle of a tool-call turn
	for start > 0 && history[start].Role == "tool" {
		start--
	}

	return history[start:]
}
