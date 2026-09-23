package tui

import (
	"github.com/alexnakagama/noryn/internal/agent"
	"github.com/alexnakagama/noryn/internal/tui/screens/welcome"
	"github.com/alexnakagama/noryn/internal/tui/styles"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	agent *agent.Agent

	model               string
	projectInstructions string
	projectContext      string

	welcome welcome.Model

	width  int
	height int
}

func New(
	a *agent.Agent,
	model string,
	projectInstructions string,
	projectContext string,
) Model {
	return Model{
		agent:               a,
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
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	var cmd tea.Cmd

	m.welcome, cmd = m.welcome.Update(msg)

	return m, cmd
}

func (m Model) View() string {
	return styles.AppStyle.
		Width(m.width).
		Height(m.height).
		Render(m.welcome.View())
}
