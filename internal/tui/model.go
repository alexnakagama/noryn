package tui

import (
	"github.com/alexnakagama/noryn/internal/agent"
	"github.com/alexnakagama/noryn/internal/tui/screens/welcome"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	agent *agent.Agent

	model               string
	projectInstructions string
	projectContext      string

	welcome welcome.Model
}

func New(a *agent.Agent, model string, projectInstructions string, projectContext string) Model {
	return Model{
		agent: a,

		model:               model,
		projectInstructions: projectInstructions,
		projectContext:      projectContext,

		welcome: welcome.New(),
	}
}

func (m Model) Init() tea.Cmd {
	return m.welcome.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	m.welcome, cmd = m.welcome.Update(msg)

	return m, cmd
}

func (m Model) View() string {
	return m.welcome.View()
}
