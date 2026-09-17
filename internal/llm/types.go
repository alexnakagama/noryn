package llm

type Message struct {
	Role    string
	Content string
}

type Request struct {
	Model    string
	Messages []Message
}

type Response struct {
	Message  Message
	ToolCall []ToolCall
}

type ToolCall struct {
	ID        string
	Name      string
	Arguments string
}
