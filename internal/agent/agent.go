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
	for {
		response, err := a.client.Chat(ctx, request)
		if err != nil {
			return llm.Response{}, err
		}

		if len(response.ToolCalls) == 0 {
			return response, nil
		}

		request.Messages = append(request.Messages, response.Message)

		for _, call := range response.ToolCalls {
			result, err := a.executeTool(call)
			if err != nil {
				result = "tool error: " + err.Error()
			}

			request.Messages = append(request.Messages, llm.Message{
				Role:       "tool",
				Content:    result,
				ToolCallID: call.ID,
			})
		}
	}
}

func (a *Agent) executeTool(call llm.ToolCall) (string, error) {
	tool, ok := a.tools[call.Name]
	if !ok {
		return "", fmt.Errorf("tool not found: %s", call.Name)
	}

	return tool.Execute(call.Arguments)
}

func (a *Agent) toolDefinitions() []llm.ToolDefinition {
	definitions := make([]llm.ToolDefinition, 0, len(a.tools))

	for _, tool := range a.tools {
		definitions = append(definitions, tool.Definition())
	}

	return definitions
}
