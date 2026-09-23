package chat

import (
	"github.com/alexnakagama/noryn/internal/tui/styles"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	width  int
	height int

	input textinput.Model
}

func New() Model {
	input := textinput.New()

	input.Placeholder = "Ask Noryn anything..."
	input.Prompt = "> "
	input.CharLimit = 4000
	input.Focus()

	input.PromptStyle = lipgloss.NewStyle().
		Foreground(styles.ColorBlue)

	input.TextStyle = lipgloss.NewStyle().
		Foreground(styles.ColorText)

	input.PlaceholderStyle = lipgloss.NewStyle().
		Foreground(styles.ColorSubtext)

	return Model{
		input: input,
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		m.input.Width = m.width
	}

	var cmd tea.Cmd

	m.input, cmd = m.input.Update(msg)

	return m, cmd
}
