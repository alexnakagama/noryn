package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
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
	Tools []tool  `json:"tools,omitempty"`
}

type tool struct {
	Type        string         `json:"type"`
	Name        string         `json:"name"`
	Description string         `json:"description,omitempty"`
	Parameters  map[string]any `json:"parameters"`
}

type input struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type response struct {
	Output []output `json:"output"`
}

type output struct {
	Type      string `json:"type"`
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
	CallID    string `json:"call_id"`

	Content []content `json:"content"`
}

type content struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

func convertTools(definitions []llm.ToolDefinition) []tool {
	tools := make([]tool, 0, len(definitions))

	for _, definition := range definitions {
		tools = append(tools, tool{
			Type:        "function",
			Name:        definition.Name,
			Description: definition.Description,
			Parameters:  definition.Parameters,
		})
	}

	return tools
}

func (c *Client) Chat(ctx context.Context, req llm.Request) (llm.Response, error) {
	inputs := make([]input, 0, len(req.Messages))

	for _, message := range req.Messages {
		inputs = append(inputs, input{
			Role:    message.Role,
			Content: message.Content,
		})
	}

	body := request{
		Model: req.Model,
		Input: inputs,
		Tools: convertTools(req.Tools),
	}

	data, err := json.Marshal(body)
	if err != nil {
		return llm.Response{}, fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.baseURL+"/responses",
		bytes.NewReader(data),
	)
	if err != nil {
		return llm.Response{}, fmt.Errorf("create req: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return llm.Response{}, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return llm.Response{}, fmt.Errorf("openai returned status %d", resp.StatusCode)
	}

	var result response

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return llm.Response{}, fmt.Errorf("decode response: %w", err)
	}

	for _, output := range result.Output {
		for _, content := range output.Content {
			if content.Type == "output_text" {
				return llm.Response{
					Message: llm.Message{
						Role:    "assistant",
						Content: content.Text,
					},
				}, nil
			}
		}
	}

	return llm.Response{}, fmt.Errorf("openai response contains no text")
}
