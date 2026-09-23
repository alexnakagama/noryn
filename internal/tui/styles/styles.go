package styles

import "github.com/charmbracelet/lipgloss"

// Tokyo Night color palette.
var (
	ColorBase    = lipgloss.Color("#0D0E12")
	ColorSurface = lipgloss.Color("#16161E")
	ColorOverlay = lipgloss.Color("#24283B")

	ColorText    = lipgloss.Color("#C0CAF5")
	ColorSubtext = lipgloss.Color("#565F89")

	ColorBlue   = lipgloss.Color("#7AA2F7")
	ColorCyan   = lipgloss.Color("#7DCFFF")
	ColorPurple = lipgloss.Color("#BB9AF7")
	ColorGreen  = lipgloss.Color("#9ECE6A")
	ColorYellow = lipgloss.Color("#E0AF68")
	ColorRed    = lipgloss.Color("#F7768E")
	ColorOrange = lipgloss.Color("#FF9E64")
)

const Cursor = "▸ "

func Tagline(version string) string {
	return "Noryn " + version + " — AI coding assistant"
}

var (
	AppStyle = lipgloss.NewStyle().
			Background(ColorBase).
			Foreground(ColorText)

	TitleStyle = lipgloss.NewStyle().
			Foreground(ColorBlue).
			Bold(true)

	HeadingStyle = lipgloss.NewStyle().
			Foreground(ColorPurple).
			Bold(true)

	HelpStyle = lipgloss.NewStyle().
			Foreground(ColorSubtext)

	SubtextStyle = lipgloss.NewStyle().
			Foreground(ColorSubtext)

	SelectedStyle = lipgloss.NewStyle().
			Foreground(ColorBlue).
			Bold(true)

	UnselectedStyle = lipgloss.NewStyle().
			Foreground(ColorText)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(ColorGreen)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(ColorRed)

	WarningStyle = lipgloss.NewStyle().
			Foreground(ColorYellow)
)
