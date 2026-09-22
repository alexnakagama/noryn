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
	if len(history) <= c.maxMessages {
		return history
	}

	start := len(history) - c.maxMessages

	// Dont start in the middle of a tool-call turn
	for start > 0 && history[start].Role == "tool" {
		start--
	}

	return history[start:]
}
