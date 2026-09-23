package chat

type Message struct {
	Role    string
	Content string
}

type Model struct {
	messages  []Message
	streaming bool
}

func New() Model {
	return Model{
		messages: []Message{},
	}
}

func (m Model) AddMessage(message Message) Model {
	m.messages = append(m.messages, message)

	return m
}

func (m Model) Messages() []Message {
	return m.messages
}

func (m Model) StartStreaming() Model {
	m.streaming = true

	m.messages = append(m.messages, Message{
		Role: "assistant",
	})

	return m
}

func (m Model) AppendAssistantContent(content string) Model {
	if len(m.messages) == 0 {
		return m
	}

	last := &m.messages[len(m.messages)-1]

	if last.Role != "assistant" {
		return m
	}

	last.Content += content

	return m
}

func (m Model) FinishStreaming() Model {
	m.streaming = false

	return m
}

func (m Model) IsStreaming() bool {
	return m.streaming
}

func (m Model) AddToolMessage(content string) Model {
	m.messages = append(m.messages, Message{
		Role:    "tool",
		Content: content,
	})

	return m
}
