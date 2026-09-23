package welcome

import (
	"strings"

	"github.com/alexnakagama/noryn/internal/tui/styles"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	title := styles.TitleStyle.Render("NORYN")

	tagline := styles.SubtextStyle.Render(
		styles.Tagline("v0.1.0"),
	)

	help := styles.HelpStyle.Render(
		"Press Enter to start",
	)

	content := strings.Join([]string{
		title,
		"",
		tagline,
		"",
		"",
		help,
	}, "\n")

	return lipgloss.NewStyle().
		Align(lipgloss.Center).
		AlignVertical(lipgloss.Center).
		Render(content)
}
