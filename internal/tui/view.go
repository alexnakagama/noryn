package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	var sections []string

	sections = append(sections, LogoStyle.Render("NORYN"))
	sections = append(sections, "")

	sections = append(sections, m.chat.View())

	sections = append(sections, "")
	sections = append(sections, m.input.View())

	sections = append(sections, "")
	sections = append(sections, m.status.View())

	content := strings.Join(sections, "\n")

	return lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Render(content)
}
