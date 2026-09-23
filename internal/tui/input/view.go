package input

import "github.com/charmbracelet/lipgloss"

func (m Model) View() string {
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#303640")).
		Padding(0, 1).
		Render(m.textarea.View())
}
