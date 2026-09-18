package openrouter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
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
		baseURL:    "https://openrouter.ai/api/v1",
	}
}

type request struct {
	Model    string    `json:"model"`
	Messages []message `json:"messages"`
	Tools    []tool    `json:"tools,omitempty"`
}

type message struct {
	Role       string     `json:"role"`
	Content    string     `json:"content"`
	ToolCalls  []toolCall `json:"tool_calls,omitempty"`
	ToolCallID string     `json:"tool_call_id,omitempty"`
}

type toolCall struct {
	ID       string   `json:"id"`
	Type     string   `json:"type"`
	Function function `json:"function"`
}

type function struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type tool struct {
	Type     string       `json:"type"`
	Function toolFunction `json:"function"`
}

type toolFunction struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Parameters  map[string]any `json:"parameters"`
}

type response struct {
	Choices []choice `json:"choices"`
}

type choice struct {
	Message message `json:"message"`
}

func convertTools(definitions []llm.ToolDefinition) []tool {
	tools := make([]tool, 0, len(definitions))

	for _, definition := range definitions {
		tools = append(tools, tool{
			Type: "function",
			Function: toolFunction{
				Name:        definition.Name,
				Description: definition.Description,
				Parameters:  definition.Parameters,
			},
		})
	}

	return tools
}

func (c *Client) Chat(ctx context.Context, req llm.Request) (llm.Response, error) {
	messages := make([]message, 0, len(req.Messages))

	for _, msg := range req.Messages {
		message := message{
			Role:       msg.Role,
			Content:    msg.Content,
			ToolCallID: msg.ToolCallID,
		}

		for _, call := range msg.ToolCalls {
			message.ToolCalls = append(message.ToolCalls, toolCall{
				ID:   call.ID,
				Type: "function",
				Function: function{
					Name:      call.Name,
					Arguments: call.Arguments,
				},
			})
		}

		messages = append(messages, message)
	}

	body := request{
		Model:    req.Model,
		Messages: messages,
		Tools:    convertTools(req.Tools),
	}

	data, err := json.Marshal(body)
	if err != nil {
		return llm.Response{}, err
	}

	httpReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.baseURL+"/chat/completions",
		nil,
	)
	if err != nil {
		return llm.Response{}, err
	}

	httpReq.Body = io.NopCloser(
		bytes.NewReader(data),
	)

	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return llm.Response{}, err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return llm.Response{}, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return llm.Response{}, fmt.Errorf(
			"openrouter returned status %d: %s",
			resp.StatusCode,
			string(responseBody),
		)
	}

	var result response

	if err := json.Unmarshal(responseBody, &result); err != nil {
		return llm.Response{}, err
	}

	if len(result.Choices) == 0 {
		return llm.Response{}, fmt.Errorf("openrouter returned no choices")
	}

	msg := result.Choices[0].Message

	var toolCalls []llm.ToolCall

	for _, call := range msg.ToolCalls {
		toolCalls = append(toolCalls, llm.ToolCall{
			ID:        call.ID,
			Name:      call.Function.Name,
			Arguments: call.Function.Arguments,
		})
	}

	return llm.Response{
		Message: llm.Message{
			Role:      "assistant",
			Content:   msg.Content,
			ToolCalls: toolCalls,
		},
		ToolCalls: toolCalls,
	}, nil
}
