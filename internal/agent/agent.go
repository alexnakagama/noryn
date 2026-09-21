package agent

import (
	"context"
	"fmt"

	"github.com/alexnakagama/noryn/internal/llm"
	"github.com/alexnakagama/noryn/internal/tools"
)

const maxToolResultLength = 10_000
const maxToolIterations = 20

func truncateToolResult(result string) string {
	if len(result) <= maxToolResultLength {
		return result
	}

	return result[:maxToolResultLength] + "\n[tool result truncated]"
}

type Agent struct {
	client       llm.Client
	tools        map[string]tools.Tool
	history      []llm.Message
	context      *ContextManager
	onToolCall   func(llm.ToolCall)
	onToolResult func(llm.ToolCall, string)
}

func New(client llm.Client, toolList ...tools.Tool) *Agent {
	toolMap := make(map[string]tools.Tool)

	for _, tool := range toolList {
		toolMap[tool.Name()] = tool
	}

	return &Agent{
		client:  client,
		tools:   toolMap,
		context: NewContextManager(),
	}
}

func (a *Agent) SetToolCallHandler(handler func(llm.ToolCall)) {
	a.onToolCall = handler
}

func (a *Agent) SetToolResultHandler(handler func(llm.ToolCall, string)) {
	a.onToolResult = handler
}

func (a *Agent) Chat(ctx context.Context, request llm.Request) (llm.Response, error) {
	a.history = append(a.history, request.Messages...)
	request.Messages = a.context.Build(a.history)
	request.Tools = a.toolDefinitions()

	toolIterations := 0

	for {
		response, err := a.client.Chat(ctx, request)
		if err != nil {
			return llm.Response{}, err
		}

		if len(response.ToolCalls) == 0 {
			a.history = append(a.history, response.Message)
			return response, nil
		}

		toolIterations++

		if toolIterations > maxToolIterations {
			return llm.Response{}, fmt.Errorf("maximum tool iterations exceeded")
		}

		a.history = append(a.history, response.Message)

		for _, call := range response.ToolCalls {
			if a.onToolCall != nil {
				a.onToolCall(call)
			}

			result, err := a.executeTool(call)
			if err != nil {
				result = "tool error: " + err.Error()
			}

			result = truncateToolResult(result)

			if a.onToolResult != nil {
				a.onToolResult(call, result)
			}

			a.history = append(a.history, llm.Message{
				Role:       "tool",
				Content:    result,
				ToolCallID: call.ID,
			})
		}

		request.Messages = a.context.Build(a.history)
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
