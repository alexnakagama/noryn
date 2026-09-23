package tui

import (
	"github.com/alexnakagama/noryn/internal/agent"
	"github.com/alexnakagama/noryn/internal/tui/chat"
	"github.com/alexnakagama/noryn/internal/tui/input"
	"github.com/alexnakagama/noryn/internal/tui/status"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	width  int
	height int

	agent  *agent.Agent
	chat   chat.Model
	input  input.Model
	status status.Model
}

func New(a *agent.Agent) Model {
	return Model{
		agent:  a,
		chat:   chat.New(),
		input:  input.New(),
		status: status.New(),
	}
}

func (m Model) Init() tea.Cmd {
	return m.input.Focus()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	var cmd tea.Cmd

	m.input, cmd = m.input.Update(msg)

	return m, cmd
}
