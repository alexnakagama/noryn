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
