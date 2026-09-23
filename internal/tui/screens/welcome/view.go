package welcome

import (
	"strings"

	"github.com/alexnakagama/noryn/internal/tui/styles"
)

func (m Model) View() string {
	title := styles.TitleStyle.Render("NORYN")

	tagline := styles.SubtextStyle.Render(
		styles.Tagline("v0.1.0"),
	)

	help := styles.HelpStyle.Render(
		"Press Enter to start",
	)

	return strings.Join([]string{
		title,
		"",
		tagline,
		"",
		"",
		help,
	}, "\n")
}
