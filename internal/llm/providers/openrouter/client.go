package openrouter

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

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
	Stream   bool      `json:"stream,omitempty"`
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

type streamResponse struct {
	Choices []streamChoice `json:"choices"`
}

type streamChoice struct {
	Delta streamDelta `json:"delta"`
}

type streamDelta struct {
	Content   string           `json:"content"`
	ToolCalls []streamToolCall `json:"tool_calls"`
}

type streamToolCall struct {
	Index    int                `json:"index"`
	ID       string             `json:"id"`
	Type     string             `json:"type"`
	Function streamToolFunction `json:"function"`
}

type streamToolFunction struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type accumulatedToolCall struct {
	ID        string
	Name      string
	Arguments strings.Builder
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

func convertMessages(messages []llm.Message) []message {
	result := make([]message, 0, len(messages))

	for _, msg := range messages {
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

		result = append(result, message)
	}

	return result
}

func (c *Client) Chat(ctx context.Context, req llm.Request) (llm.Response, error) {
	body := request{
		Model:    req.Model,
		Messages: convertMessages(req.Messages),
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
		bytes.NewReader(data),
	)
	if err != nil {
		return llm.Response{}, err
	}

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

func (c *Client) ChatStream(
	ctx context.Context,
	req llm.Request,
) (<-chan llm.StreamChunk, error) {
	body := request{
		Model:    req.Model,
		Messages: convertMessages(req.Messages),
		Tools:    convertTools(req.Tools),
		Stream:   true,
	}

	data, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.baseURL+"/chat/completions",
		bytes.NewReader(data),
	)
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "text/event-stream")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		defer resp.Body.Close()

		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		return nil, fmt.Errorf(
			"openrouter returned status %d: %s",
			resp.StatusCode,
			string(responseBody),
		)
	}

	chunks := make(chan llm.StreamChunk)

	go func() {
		defer close(chunks)
		defer resp.Body.Close()

		c.streamResponse(ctx, resp.Body, chunks)
	}()

	return chunks, nil
}

func (c *Client) streamResponse(
	ctx context.Context,
	body io.Reader,
	chunks chan<- llm.StreamChunk,
) {
	scanner := bufio.NewScanner(body)

	// Increase the scanner limit because tool arguments can be large.
	scanner.Buffer(make([]byte, 64*1024), 1024*1024)

	toolCalls := make(map[int]*accumulatedToolCall)

	for scanner.Scan() {
		line := scanner.Text()

		if !strings.HasPrefix(line, "data:") {
			continue
		}

		data := strings.TrimSpace(strings.TrimPrefix(line, "data:"))

		if data == "[DONE]" {
			c.sendCompletedToolCalls(ctx, chunks, toolCalls)

			select {
			case chunks <- llm.StreamChunk{Done: true}:
			case <-ctx.Done():
			}

			return
		}

		var event streamResponse

		if err := json.Unmarshal([]byte(data), &event); err != nil {
			c.sendStreamError(ctx, chunks, err)
			return
		}

		for _, choice := range event.Choices {
			if choice.Delta.Content != "" {
				select {
				case chunks <- llm.StreamChunk{
					Content: choice.Delta.Content,
				}:
				case <-ctx.Done():
					return
				}
			}

			for _, call := range choice.Delta.ToolCalls {
				accumulated, exists := toolCalls[call.Index]

				if !exists {
					accumulated = &accumulatedToolCall{}
					toolCalls[call.Index] = accumulated
				}

				if call.ID != "" {
					accumulated.ID = call.ID
				}

				if call.Function.Name != "" {
					accumulated.Name = call.Function.Name
				}

				accumulated.Arguments.WriteString(call.Function.Arguments)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		c.sendStreamError(ctx, chunks, err)
	}
}

func (c *Client) sendCompletedToolCalls(
	ctx context.Context,
	chunks chan<- llm.StreamChunk,
	toolCalls map[int]*accumulatedToolCall,
) {
	for i := 0; ; i++ {
		call, ok := toolCalls[i]
		if !ok {
			break
		}

		toolCall := llm.ToolCall{
			ID:        call.ID,
			Name:      call.Name,
			Arguments: call.Arguments.String(),
		}

		select {
		case chunks <- llm.StreamChunk{
			ToolCall: &toolCall,
		}:
		case <-ctx.Done():
			return
		}
	}
}

func (c *Client) sendStreamError(
	ctx context.Context,
	chunks chan<- llm.StreamChunk,
	err error,
) {
	select {
	case chunks <- llm.StreamChunk{
		Err: err,
	}:
	case <-ctx.Done():
	}
}
