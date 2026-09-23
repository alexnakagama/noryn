package agent

import "github.com/alexnakagama/noryn/internal/llm"

type EventType int

const (
	EventText EventType = iota
	EventToolCall
	EventToolResult
	EventDone
	EventError
)

type Event struct {
	Type     EventType
	Content  string
	ToolCall *llm.ToolCall
}
