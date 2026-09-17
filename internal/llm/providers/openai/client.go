package openai

import "net/http"

type Client struct {
	apiKey     string
	httpClient *http.Client
	baseURL    string
}
