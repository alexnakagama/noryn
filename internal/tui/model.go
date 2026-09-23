package tui

import (
	"context"

	"github.com/alexnakagama/noryn/internal/agent"
	"github.com/alexnakagama/noryn/internal/llm"
	"github.com/alexnakagama/noryn/internal/tui/chat"
	"github.com/alexnakagama/noryn/internal/tui/input"
	"github.com/alexnakagama/noryn/internal/tui/status"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	width  int
	height int

	agent               *agent.Agent
	model               string
	projectInstructions string
	projectContext      string

	chat   chat.Model
	input  input.Model
	status status.Model

	agentEvents <-chan agent.Event
}

type agentStartedMsg struct {
	events <-chan agent.Event
}

type agentEventMsg struct {
	event agent.Event
}

func New(
	a *agent.Agent,
	model string,
	projectInstructions string,
	projectContext string,
) Model {
	m := Model{
		agent:               a,
		model:               model,
		projectInstructions: projectInstructions,
		projectContext:      projectContext,
		chat:                chat.New(),
		input:               input.New(),
		status:              status.New(),
	}

	m.input.SetHintRight(model)

	return m
}

func (m Model) Init() tea.Cmd {
	return m.input.Focus()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		m.input.SetWidth(msg.Width)

	case input.SubmitMessage:
		m.chat = m.chat.AddMessage(chat.Message{
			Role:    "user",
			Content: msg.Content,
		})

		return m, m.startAgent(msg.Content)

	case agentStartedMsg:
		m.agentEvents = msg.events

		m.chat = m.chat.StartStreaming()
		m.status = m.status.SetStatus("working")

		return m, waitForAgentEvent(m.agentEvents)

	case agentEventMsg:
		return m.handleAgentEvent(msg.event)
	}

	var cmd tea.Cmd

	m.input, cmd = m.input.Update(msg)

	return m, cmd
}

func (m Model) startAgent(content string) tea.Cmd {
	return func() tea.Msg {
		request := llm.Request{
			Model: m.model,
			Messages: []llm.Message{
				{
					Role:    "system",
					Content: m.projectInstructions + "\n\n" + m.projectContext,
				},
				{
					Role:    "user",
					Content: content,
				},
			},
		}

		events, err := m.agent.ChatStream(
			context.Background(),
			request,
		)

		if err != nil {
			return agentEventMsg{
				event: agent.Event{
					Type:    agent.EventError,
					Content: err.Error(),
				},
			}
		}

		return agentStartedMsg{
			events: events,
		}
	}
}

func waitForAgentEvent(events <-chan agent.Event) tea.Cmd {
	return func() tea.Msg {
		event, ok := <-events

		if !ok {
			return agentEventMsg{
				event: agent.Event{
					Type: agent.EventDone,
				},
			}
		}

		return agentEventMsg{
			event: event,
		}
	}
}

func (m Model) handleAgentEvent(event agent.Event) (tea.Model, tea.Cmd) {
	switch event.Type {

	case agent.EventText:
		m.chat = m.chat.AppendAssistantContent(event.Content)

	case agent.EventToolCall:
		if event.ToolCall != nil {
			m.chat = m.chat.AddToolMessage(
				"Running " + event.ToolCall.Name + "...",
			)
		}

	case agent.EventToolResult:
		m.chat = m.chat.AddToolMessage(event.Content)

	case agent.EventDone:
		m.chat = m.chat.FinishStreaming()
		m.status = m.status.SetStatus("ready")

		return m, nil

	case agent.EventError:
		m.chat = m.chat.FinishStreaming()
		m.status = m.status.SetStatus("error")

		m.chat = m.chat.AddMessage(chat.Message{
			Role:    "error",
			Content: event.Content,
		})

		return m, nil
	}

	return m, waitForAgentEvent(m.agentEvents)
}
