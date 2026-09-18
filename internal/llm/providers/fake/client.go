package fake

import (
	"context"

	"github.com/alexnakagama/noryn/internal/llm"
)

type Client struct{}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) Chat(ctx context.Context, request llm.Request) (llm.Response, error) {
	return llm.Response{
		Message: llm.Message{
			Role:    "assistant",
			Content: "I need to read a file.",
		},
		ToolCalls: []llm.ToolCall{
			{
				ID:        "call-1",
				Name:      "read_file",
				Arguments: `{"path":"cmd/noryn/main.go"}`,
			},
		},
	}, nil
}
