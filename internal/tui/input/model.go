package input

import (
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	textarea textarea.Model
}

func New() Model {
	input := textarea.New()

	input.Placeholder = "Ask Noryn anything..."
	input.Prompt = "> "
	input.CharLimit = 0

	return Model{
		textarea: input,
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd

	m.textarea, cmd = m.textarea.Update(msg)

	return m, cmd
}

func (m Model) View() string {
	return m.textarea.View()
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
