package input

import (
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	textarea  textarea.Model
	width     int
	hintRight string
}

type SubmitMessage struct {
	Content string
}

func New() Model {
	input := textarea.New()

	input.Placeholder = "Describe the task you want accomplished..."
	input.Prompt = "> "
	input.CharLimit = 0
	input.ShowLineNumbers = false
	input.SetHeight(3)

	input.SetPromptFunc(2, func(line int) string {
		if line == 0 {
			return "> "
		}
		return ""
	})

	panel := lipgloss.NewStyle().
		Background(lipgloss.Color(inputBackground))

	input.FocusedStyle.Base = panel
	input.BlurredStyle.Base = panel

	input.FocusedStyle.CursorLine = lipgloss.NewStyle().
		Background(lipgloss.Color(cursorLineBackground))

	input.FocusedStyle.Text = lipgloss.NewStyle().
		Foreground(lipgloss.Color(inputText))

	input.FocusedStyle.Prompt = lipgloss.NewStyle().
		Foreground(lipgloss.Color(inputPrompt)).
		Bold(true)

	input.FocusedStyle.Placeholder = lipgloss.NewStyle().
		Foreground(lipgloss.Color(inputMuted))

	input.BlurredStyle.CursorLine = lipgloss.NewStyle().
		Background(lipgloss.Color(cursorLineBackground))

	input.BlurredStyle.Text = lipgloss.NewStyle().
		Foreground(lipgloss.Color(inputText))

	input.BlurredStyle.Prompt = lipgloss.NewStyle().
		Foreground(lipgloss.Color(inputPrompt)).
		Bold(true)

	input.BlurredStyle.Placeholder = lipgloss.NewStyle().
		Foreground(lipgloss.Color(inputMuted))

	input.Focus()

	return Model{
		textarea: input,
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok {
		switch {
		case key.Type == tea.KeyEnter && key.Alt:
			var cmd tea.Cmd
			m.textarea, cmd = m.textarea.Update(tea.KeyMsg{Type: tea.KeyEnter})
			return m, cmd

		case key.Type == tea.KeyCtrlJ:
			var cmd tea.Cmd
			m.textarea, cmd = m.textarea.Update(tea.KeyMsg{Type: tea.KeyEnter})
			return m, cmd

		case key.Type == tea.KeyEnter:
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
	m.width = width

	inner := width - 4

	if inner < 1 {
		inner = 1
	}

	m.textarea.SetWidth(inner)
}

func (m *Model) SetHintRight(hint string) {
	m.hintRight = hint
}
