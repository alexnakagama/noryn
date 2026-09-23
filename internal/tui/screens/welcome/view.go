package welcome

import (
	"strings"

	"github.com/alexnakagama/noryn/internal/tui/styles"
	"github.com/charmbracelet/lipgloss"
)

const logo = `
███╗   ██╗ ██████╗ ██████╗ ██╗   ██╗███╗   ██╗
████╗  ██║██╔═══██╗██╔══██╗╚██╗ ██╔╝████╗  ██║
██╔██╗ ██║██║   ██║██████╔╝ ╚████╔╝ ██╔██╗ ██║
██║╚██╗██║██║   ██║██╔══██╗  ╚██╔╝  ██║╚██╗██║
██║ ╚████║╚██████╔╝██║  ██║   ██║   ██║ ╚████║
╚═╝  ╚═══╝ ╚═════╝ ╚═╝  ╚═╝   ╚═╝   ╚═╝  ╚═══╝
`

func (m Model) View() string {
	logo := styles.TitleStyle.Render(strings.TrimSpace(logo))

	tagline := styles.SubtextStyle.Render(
		styles.Tagline("v0.1.0"),
	)

	help := styles.HelpStyle.Render(
		"Press Enter to start",
	)

	content := strings.Join([]string{
		logo,
		"",
		tagline,
		"",
		"",
		help,
	}, "\n")

	return lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Align(lipgloss.Center).
		AlignVertical(lipgloss.Center).
		Render(content)
}
