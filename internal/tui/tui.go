package tui

import (
	"github.com/alexnakagama/noryn/internal/agent"
	tea "github.com/charmbracelet/bubbletea"
)

func Run(a *agent.Agent, model string, projectInstructions string, projectContext string) error {
	m := New(
		a,
		model,
		projectInstructions,
		projectContext,
	)

	program := tea.NewProgram(m)

	_, err := program.Run()

	return err
}
