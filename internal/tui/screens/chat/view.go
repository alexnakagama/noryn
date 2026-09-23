package chat

import (
	"github.com/alexnakagama/noryn/internal/tui/styles"
	"github.com/charmbracelet/lipgloss"
)

const sidebarWidth = 28

func (m Model) View() string {
	mainWidth := m.width - sidebarWidth

	main := lipgloss.NewStyle().
		Width(mainWidth).
		Height(m.height).
		Padding(1, 2).
		Render(
			"Conversation\n\n> _",
		)

	sidebar := lipgloss.NewStyle().
		Width(sidebarWidth).
		Height(m.height).
		Padding(1, 2).
		Background(styles.ColorSurface).
		Render(
			"STATUS\n\n" +
				"● Ready\n\n" +
				"PROJECT\n\n" +
				"noryn\n\n" +
				"BRANCH\n\n" +
				"main\n\n" +
				"MODEL\n\n" +
				"gpt-5.6\n\n" +
				"USAGE\n\n" +
				"12.4k tokens",
		)

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		main,
		sidebar,
	)
}
