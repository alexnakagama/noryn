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

func (c *ContextManager) Build(history []llm.Message) []llm.Message {}
