package input

import (
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	textarea textarea.Model
}

type SubmitMessage struct {
	Content string
}

func New() Model {
	input := textarea.New()

	input.Placeholder = "Ask Noryn anything..."
	input.Prompt = "> "
	input.CharLimit = 0
	input.ShowLineNumbers = false
	input.SetHeight(5)

	input.FocusedStyle.CursorLine = lipgloss.NewStyle()

	input.FocusedStyle.Text = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#D8DEE9"))

	input.FocusedStyle.Prompt = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#7AA2F7")).
		Bold(true)

	input.BlurredStyle.CursorLine = lipgloss.NewStyle()

	input.BlurredStyle.Text = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#D8DEE9"))

	input.BlurredStyle.Prompt = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#7AA2F7")).
		Bold(true)

	input.Focus()

	return Model{
		textarea: input,
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok {
		if key.Type == tea.KeyEnter {
			content := m.Value()

			if content == "" {
				return m, nil
			}

			m.Reset()

			return m, func() tea.Msg {
				return SubmitMessage{
					Content: content,
				}
			}
		}
	}

	var cmd tea.Cmd

	m.textarea, cmd = m.textarea.Update(msg)

	return m, cmd
}

func (m Model) Value() string {
	return m.textarea.Value()
}

func (m *Model) Reset() {
	m.textarea.Reset()
}

func (m *Model) Focus() tea.Cmd {
	return m.textarea.Focus()
}

func (m *Model) SetWidth(width int) {
	m.textarea.SetWidth(width)
}
