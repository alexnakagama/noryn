package agent

import (
	"context"

	"github.com/alexnakagama/noryn/internal/llm"
)

type Agent struct {
	client llm.Client
}

func New(client llm.Client) *Agent {
	return &Agent{
		client: client,
	}
}

func (a *Agent) Chat(ctx context.Context, request llm.Request) (llm.Response, error) {
	return a.client.Chat(ctx, request)
}
