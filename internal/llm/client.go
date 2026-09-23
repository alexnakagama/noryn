package llm

import "context"

type Client interface {
	Chat(ctx context.Context, request Request) (Response, error)
	ChatStream(ctx context.Context, request Request) (<-chan StreamChunk, error)
}
