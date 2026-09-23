package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	if m.width == 0 || m.height == 0 {
		return ""
	}

	header := LogoStyle.Render("NORYN")

	chat := m.chat.View()

	input := m.input.View()

	status := StatusBarStyle.Render(m.status.View())

	availableHeight := m.height -
		lipgloss.Height(header) -
		lipgloss.Height(input) -
		lipgloss.Height(status) -
		6

	if availableHeight < 1 {
		availableHeight = 1
	}

	chatArea := lipgloss.NewStyle().
		Height(availableHeight).
		Render(chat)

	content := strings.Join([]string{
		header,
		"",
		chatArea,
		input,
		"",
		status,
	}, "\n")

	return AppStyle.
		Width(m.width).
		Height(m.height).
		Render(content)
}
