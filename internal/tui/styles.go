package tui

import "github.com/charmbracelet/lipgloss"

var (
	// ─────────────────────────────────────────────
	// Palette
	// ─────────────────────────────────────────────

	ColorText       = lipgloss.Color("#D8DEE9")
	ColorSubtle     = lipgloss.Color("#7C8594")
	ColorMuted      = lipgloss.Color("#565F6D")
	ColorBorder     = lipgloss.Color("#303640")
	ColorBackground = lipgloss.Color("#1A1B26")

	ColorPrimary   = lipgloss.Color("#7AA2F7")
	ColorSecondary = lipgloss.Color("#BB9AF7")
	ColorSuccess   = lipgloss.Color("#9ECE6A")
	ColorWarning   = lipgloss.Color("#E0AF68")
	ColorError     = lipgloss.Color("#F7768E")
	ColorInfo      = lipgloss.Color("#7DCFFF")

	// ─────────────────────────────────────────────
	// Application
	// ─────────────────────────────────────────────

	AppStyle = lipgloss.NewStyle().
			Background(ColorBackground).
			Foreground(ColorText)

	HeaderStyle = lipgloss.NewStyle().
			Foreground(ColorText).
			Bold(true)

	LogoStyle = lipgloss.NewStyle().
			Foreground(ColorPrimary).
			Bold(true)

	SubtitleStyle = lipgloss.NewStyle().
			Foreground(ColorMuted)

	// ─────────────────────────────────────────────
	// Conversation
	// ─────────────────────────────────────────────

	UserLabelStyle = lipgloss.NewStyle().
			Foreground(ColorPrimary).
			Bold(true)

	AssistantLabelStyle = lipgloss.NewStyle().
				Foreground(ColorSecondary).
				Bold(true)

	UserMessageStyle = lipgloss.NewStyle().
				Foreground(ColorText)

	AssistantMessageStyle = lipgloss.NewStyle().
				Foreground(ColorText)

	// ─────────────────────────────────────────────
	// Input
	// ─────────────────────────────────────────────

	InputStyle = lipgloss.NewStyle().
			Foreground(ColorText)

	InputPromptStyle = lipgloss.NewStyle().
				Foreground(ColorPrimary).
				Bold(true)

	InputBorderStyle = lipgloss.NewStyle().
				Foreground(ColorBorder)

	// ─────────────────────────────────────────────
	// Tools
	// ─────────────────────────────────────────────

	ToolStyle = lipgloss.NewStyle().
			Foreground(ColorInfo)

	ToolRunningStyle = lipgloss.NewStyle().
				Foreground(ColorWarning)

	ToolSuccessStyle = lipgloss.NewStyle().
				Foreground(ColorSuccess)

	ToolErrorStyle = lipgloss.NewStyle().
			Foreground(ColorError)

	ToolNameStyle = lipgloss.NewStyle().
			Foreground(ColorText).
			Bold(true)

	ToolPathStyle = lipgloss.NewStyle().
			Foreground(ColorMuted)

	// ─────────────────────────────────────────────
	// Status
	// ─────────────────────────────────────────────

	StatusBarStyle = lipgloss.NewStyle().
			Foreground(ColorMuted)

	StatusReadyStyle = lipgloss.NewStyle().
				Foreground(ColorSuccess)

	StatusBusyStyle = lipgloss.NewStyle().
			Foreground(ColorWarning)

	StatusErrorStyle = lipgloss.NewStyle().
				Foreground(ColorError)

	StatusModelStyle = lipgloss.NewStyle().
				Foreground(ColorSubtle)

	StatusPathStyle = lipgloss.NewStyle().
			Foreground(ColorMuted)

	// ─────────────────────────────────────────────
	// Markdown / Code
	// ─────────────────────────────────────────────

	CodeStyle = lipgloss.NewStyle().
			Foreground(ColorText)

	CodeBlockStyle = lipgloss.NewStyle().
			Foreground(ColorText).
			Background(lipgloss.Color("#16171F")).
			Padding(1, 2)

	InlineCodeStyle = lipgloss.NewStyle().
			Foreground(ColorInfo)

	// ─────────────────────────────────────────────
	// Errors / Feedback
	// ─────────────────────────────────────────────

	ErrorStyle = lipgloss.NewStyle().
			Foreground(ColorError)

	WarningStyle = lipgloss.NewStyle().
			Foreground(ColorWarning)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(ColorSuccess)

	InfoStyle = lipgloss.NewStyle().
			Foreground(ColorInfo)

	// ─────────────────────────────────────────────
	// Separators
	// ─────────────────────────────────────────────

	SeparatorStyle = lipgloss.NewStyle().
			Foreground(ColorBorder)

	DividerStyle = lipgloss.NewStyle().
			Foreground(ColorMuted)
)
