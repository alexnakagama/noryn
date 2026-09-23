package input

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#303640")).
		Padding(0, 1).
		Render(m.textarea.View())

	hintStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#565F6D"))

	hint := hintStyle.Render("enter: send • alt+enter: newline")

	if m.hintRight != "" {
		right := hintStyle.Render(m.hintRight)

		pad := m.width - lipgloss.Width(hint) - lipgloss.Width(right)

		if pad < 1 {
			pad = 1
		}

		hint += strings.Repeat(" ", pad) + right
	}

	return box + "\n" + hint
}
