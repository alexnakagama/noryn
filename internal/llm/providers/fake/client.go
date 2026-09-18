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
				Content: "I need to see the files.",
			},
			ToolCalls: []llm.ToolCall{
				{
					ID:        "call-1",
					Name:      "list_directory",
					Arguments: `{"path":"."}`,
				},
			},
		}, nil
	}

	if c.callCount == 2 {
		return llm.Response{
			Message: llm.Message{
				Role:    "assistant",
				Content: "I need to read the main file.",
			},
			ToolCalls: []llm.ToolCall{
				{
					ID:        "call-2",
					Name:      "read_file",
					Arguments: `{"path":"cmd/noryn/main.go"}`,
				},
			},
		}, nil
	}

	return llm.Response{
		Message: llm.Message{
			Role:    "assistant",
			Content: "I inspected the project successfully.",
		},
	}, nil
}
