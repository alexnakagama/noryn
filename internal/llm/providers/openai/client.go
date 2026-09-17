package openai

import (
	"net/http"
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
