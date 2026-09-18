package llm

type Message struct {
	Role       string
	Content    string
	ToolCalls  []ToolCall
	ToolCallID string
}

type Request struct {
	Model    string
	Messages []Message
}

type Response struct {
	Message   Message
	ToolCalls []ToolCall
}

type ToolCall struct {
	ID        string
	Name      string
	Arguments string
}
