package chat

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	userStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7AA2F7")).
			Bold(true)

	toolStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7C8594"))

	toolNameStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#D8DEE9")).
			Bold(true)

	toolRunningStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#E0AF68"))

	toolCompletedStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#9ECE6A"))

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F7768E"))
)

func (m Model) View() string {
	if len(m.messages) == 0 {
		return "Start a conversation with Noryn."
	}

	var view strings.Builder

	for _, message := range m.messages {
		switch message.Role {

		case "user":
			view.WriteString(userStyle.Render("> "))

		case "error":
			view.WriteString(errorStyle.Render("Error: "))
		}

		for _, part := range message.Parts {
			switch part := part.(type) {

			case TextPart:
				view.WriteString(part.Content)

			case ToolPart:
				view.WriteString("\n")
				view.WriteString("  ")

				view.WriteString(toolStyle.Render(part.Name))

				switch part.Status {
				case "running":
					view.WriteString(" ")
					view.WriteString(
						toolRunningStyle.Render("running"),
					)

				case "completed":
					view.WriteString(" ")
					view.WriteString(
						toolCompletedStyle.Render("✓"),
					)
				}

			case ErrorPart:
				view.WriteString(
					errorStyle.Render(part.Content),
				)
			}
		}

		view.WriteString("\n\n")
	}

	return view.String()
}
