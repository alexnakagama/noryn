package chat

import "strings"

func (m Model) View() string {
	if len(m.messages) == 0 {
		return "Start a conversation with Noryn."
	}

	var view strings.Builder

	for _, message := range m.messages {
		switch message.Role {
		case "user":
			view.WriteString("> ")
			view.WriteString(message.Content)
			view.WriteString("\n\n")

		case "assistant":
			view.WriteString(message.Content)
			view.WriteString("\n\n")
		}
	}

	return view.String()
}
