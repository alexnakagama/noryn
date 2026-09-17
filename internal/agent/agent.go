package agent

import "github.com/alexnakagama/noryn/internal/llm"

type Agent struct {
	client llm.Client
}

func New(client llm.Client) *Agent {
	return &Agent{
		client: client,
	}
}
