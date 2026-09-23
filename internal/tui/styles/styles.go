package styles

import "github.com/charmbracelet/lipgloss"

// Tokyo Night color palette.
var (
	ColorBase    = lipgloss.Color("#1A1B26")
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

// Cursor is the prefix used for the currently focused item.
const Cursor = "▸ "

// Tagline returns the welcome screen tagline with the given version.
func Tagline(version string) string {
	return "Noryn " + version + " — AI coding assistant"
}

// Pre-built reusable styles.
var (
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

	FrameStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorOverlay).
			Padding(1, 2)

	PanelStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(ColorOverlay).
			Padding(0, 1)

	ProgressFilled = lipgloss.NewStyle().
			Foreground(ColorGreen)

	ProgressEmpty = lipgloss.NewStyle().
			Foreground(ColorOverlay)

	PercentStyle = lipgloss.NewStyle().
			Foreground(ColorOrange).
			Bold(true)
)
