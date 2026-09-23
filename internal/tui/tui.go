package tui

import (
	"github.com/alexnakagama/noryn/internal/agent"
	tea "github.com/charmbracelet/bubbletea"
)

func Run(a *agent.Agent) error {
	model := New(a)

	program := tea.NewProgram(model)

	_, err := program.Run()

	return err
}
