package chat

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

	last.Parts = append(last.Parts, TextPart{
		Content: content,
	})

	return m
}

func (m Model) FinishStreaming() Model {
	m.streaming = false

	return m
}

func (m Model) IsStreaming() bool {
	return m.streaming
}

func (m Model) AddToolCall(name string, arguments string) Model {
	if len(m.messages) == 0 {
		return m
	}

	last := &m.messages[len(m.messages)-1]

	if last.Role != "assistant" {
		return m
	}

	last.Parts = append(last.Parts, ToolPart{
		Name:      name,
		Arguments: arguments,
		Status:    "running",
	})

	return m
}

func (m Model) CompleteTool(name string, result string) Model {
	if len(m.messages) == 0 {
		return m
	}

	last := &m.messages[len(m.messages)-1]

	for i := len(last.Parts) - 1; i >= 0; i-- {
		part, ok := last.Parts[i].(ToolPart)
		if !ok {
			continue
		}

		if part.Name != name || part.Status != "running" {
			continue
		}

		part.Status = "completed"
		part.Result = result

		last.Parts[i] = part

		return m
	}

	return m
}
