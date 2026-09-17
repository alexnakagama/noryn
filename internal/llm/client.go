package llm

import "context"

type Client interface {
	Chat(ctx context.Context, request Request) (Response, error)
}
