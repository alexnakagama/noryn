package tui

import (
	"github.com/alexnakagama/noryn/internal/agent"
	"github.com/alexnakagama/noryn/internal/tui/screens/chat"
	"github.com/alexnakagama/noryn/internal/tui/screens/welcome"
	"github.com/alexnakagama/noryn/internal/tui/styles"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	agent *agent.Agent

	router Router

	welcome welcome.Model
	chat    chat.Model

	model               string
	projectInstructions string
	projectContext      string

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
		agent: a,

		router: NewRouter(),

		welcome: welcome.New(),
		chat:    chat.New(),

		model:               model,
		projectInstructions: projectInstructions,
		projectContext:      projectContext,
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

	switch m.router.Current() {
	case ScreenWelcome:
		return m.updateWelcome(msg)

	case ScreenChat:
		return m.updateChat(msg)
	}

	return m, nil
}

func (m Model) updateWelcome(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "enter" {
			m.router.Navigate(ScreenChat)

			m.chat, _ = m.chat.Update(
				tea.WindowSizeMsg{
					Width:  m.width,
					Height: m.height,
				},
			)

			return m, m.chat.Init()
		}
	}

	var cmd tea.Cmd
	m.welcome, cmd = m.welcome.Update(msg)

	return m, cmd
}

func (m Model) updateChat(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	m.chat, cmd = m.chat.Update(msg)

	return m, cmd
}

func (m Model) View() string {
	var view string

	switch m.router.Current() {
	case ScreenWelcome:
		view = m.welcome.View()

	case ScreenChat:
		view = m.chat.View()
	}

	return styles.AppStyle.
		Width(m.width).
		Height(m.height).
		Render(view)
}
