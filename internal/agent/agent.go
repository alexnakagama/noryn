package agent

import (
	"context"
	"fmt"

	"github.com/alexnakagama/noryn/internal/llm"
	"github.com/alexnakagama/noryn/internal/tools"
)

type Agent struct {
	client llm.Client
	tools  map[string]tools.Tool
}

func New(client llm.Client, toolList ...tools.Tool) *Agent {
	toolMap := make(map[string]tools.Tool)

	for _, tool := range toolList {
		toolMap[tool.Name()] = tool
	}

	return &Agent{
		client: client,
		tools:  toolMap,
	}
}

func (a *Agent) Chat(ctx context.Context, request llm.Request) (llm.Response, error) {
	return a.client.Chat(ctx, request)
}

func (a *Agent) executeTool(call llm.ToolCall) (string, error) {
	tool, ok := a.tools[call.Name]
	if !ok {
		return "", fmt.Errorf("tool not found: %s", call.Name)
	}

	return tool.Execute(call.Arguments)
}
