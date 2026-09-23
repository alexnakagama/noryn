package tui

import (
	"github.com/alexnakagama/noryn/internal/agent"
	tea "github.com/charmbracelet/bubbletea"
)

func Run(
	a *agent.Agent,
	model string,
	projectInstructions string,
	projectContext string,
) error {
	app := New(
		a,
		model,
		projectInstructions,
		projectContext,
	)

	_, err := tea.NewProgram(
		app,
		tea.WithAltScreen(),
	).Run()

	return err
}
