package styles

import "github.com/charmbracelet/lipgloss"

var (
	ColorText    = lipgloss.Color("#D8DEE9")
	ColorSubtle  = lipgloss.Color("#7C8594")
	ColorPrimary = lipgloss.Color("#7AA2F7")

	WelcomeTitle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorPrimary)

	WelcomeSubtitle = lipgloss.NewStyle().
			Foreground(ColorSubtle)

	WelcomePrompt = lipgloss.NewStyle().
			Foreground(ColorText)
)
