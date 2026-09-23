package chat

import (
	"github.com/alexnakagama/noryn/internal/tui/styles"
	"github.com/charmbracelet/lipgloss"
)

const sidebarWidth = 30

var sidebarBorder = lipgloss.Border{
	Top:         "─",
	Bottom:      "─",
	Left:        "│",
	Right:       "│",
	TopLeft:     "╭",
	TopRight:    "╮",
	BottomLeft:  "╰",
	BottomRight: "╯",
}

func (m Model) View() string {
	sidebar := lipgloss.NewStyle().
		Width(sidebarWidth).
		Height(m.height-2).
		Border(sidebarBorder).
		BorderForeground(styles.ColorBlue).
		Padding(1, 2).
		Render(
			"STATUS\n" +
				"● Ready\n\n" +
				"PROJECT\n" +
				"noryn\n\n" +
				"BRANCH\n" +
				"main\n\n" +
				"MODEL\n" +
				"gpt-5.6\n\n" +
				"USAGE\n" +
				"12.4k tokens",
		)

	mainWidth := m.width - lipgloss.Width(sidebar)

	input := lipgloss.NewStyle().
		Width(mainWidth - 4).
		Render(m.input.View())

	inputHeight := lipgloss.Height(input)

	spacerHeight := m.height - inputHeight - 5
	if spacerHeight < 0 {
		spacerHeight = 0
	}

	spacer := lipgloss.NewStyle().
		Height(spacerHeight).
		Render("")

	main := lipgloss.NewStyle().
		Width(mainWidth).
		Height(m.height).
		Padding(1, 2).
		Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				"Conversation",
				spacer,
				input,
			),
		)

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		main,
		sidebar,
	)
}
