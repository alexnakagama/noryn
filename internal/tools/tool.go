package tools

import "github.com/alexnakagama/noryn/internal/llm"

type Tool interface {
	Name() string
	Execute(arguments string) (string, error)
	Definition() llm.ToolDefinition
}
