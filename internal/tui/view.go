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
		2

	if availableHeight < 0 {
		availableHeight = 0
	}

	rows := []string{header, ""}

	if availableHeight > 0 {
		rows = append(rows, lipgloss.NewStyle().
			Height(availableHeight).
			Render(chat))
	}

	rows = append(rows, input, "", status)

	content := strings.Join(rows, "\n")

	return AppStyle.
		Width(m.width).
		Height(m.height).
		Render(content)
}
