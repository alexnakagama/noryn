package agent

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
