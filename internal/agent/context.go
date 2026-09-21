package agent

import "github.com/alexnakagama/noryn/internal/llm"

const (
	maxContextMessages = 100
)

type ContextManager struct {
	maxMessages int
}

func NewContextManager(maxMessages int) *ContextManager {
	return &ContextManager{
		maxMessages: maxMessages,
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
