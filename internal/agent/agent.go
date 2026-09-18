package agent

import (
	"context"
	"fmt"

	"github.com/alexnakagama/noryn/internal/llm"
	"github.com/alexnakagama/noryn/internal/tools"
)

const maxToolResultLength = 10_000
const maxHistoryMessages = 100

func recentHistory(history []llm.Message) []llm.Message {
	if len(history) <= maxHistoryMessages {
		return history
	}

	start := len(history) - maxHistoryMessages

	for start < len(history) && history[start].Role == "tool" {
		start--
	}

	return history[start:]
}

func truncateToolResult(result string) string {
	if len(result) <= maxToolResultLength {
		return result
	}

	return result[:maxToolResultLength] + "\n[tool result truncated]"
}

type Agent struct {
	client  llm.Client
	tools   map[string]tools.Tool
	history []llm.Message
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
	a.history = append(a.history, request.Messages...)
	request.Messages = recentHistory(a.history)
	request.Tools = a.toolDefinitions()

	for {
		response, err := a.client.Chat(ctx, request)
		if err != nil {
			return llm.Response{}, err
		}

		if len(response.ToolCalls) == 0 {
			a.history = append(a.history, response.Message)
			return response, nil
		}

		a.history = append(a.history, response.Message)
		request.Messages = recentHistory(a.history)

		for _, call := range response.ToolCalls {
			result, err := a.executeTool(call)
			if err != nil {
				result = "tool error: " + err.Error()
			}

			result = truncateToolResult(result)

			a.history = append(a.history, llm.Message{
				Role:       "tool",
				Content:    result,
				ToolCallID: call.ID,
			})

			request.Messages = recentHistory(a.history)
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
