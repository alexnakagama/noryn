package openai

import (
	"context"
	"net/http"

	"github.com/alexnakagama/noryn/internal/llm"
)

type Client struct {
	apiKey     string
	httpClient *http.Client
	baseURL    string
}

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey:     apiKey,
		httpClient: &http.Client{},
		baseURL:    "https://api.openai.com/v1",
	}
}

type request struct {
	Model string  `json:"model"`
	Input []input `json:"input"`
}

type input struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func (c *Client) Chat(ctx context.Context, req llm.Request) (llm.Response, error) {}
