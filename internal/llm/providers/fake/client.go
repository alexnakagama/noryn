package fake

import (
	"context"

	"github.com/alexnakagama/noryn/internal/llm"
)

type Client struct {
	callCount int
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) Chat(ctx context.Context, request llm.Request) (llm.Response, error) {
	c.callCount++

	if c.callCount == 1 {
		return llm.Response{
			Message: llm.Message{
				Role:    "assistant",
				Content: "I need to run a command.",
			},
			ToolCalls: []llm.ToolCall{
				{
					ID:        "call-1",
					Name:      "shell",
					Arguments: `{"command":"echo \"Hello from Noryn\""}`,
				},
			},
		}, nil
	}

	return llm.Response{
		Message: llm.Message{
			Role:    "assistant",
			Content: "The command executed successfully.",
		},
	}, nil
}
