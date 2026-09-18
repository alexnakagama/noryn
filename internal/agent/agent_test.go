package agent

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/alexnakagama/noryn/internal/llm"
)

// stubClient returns queued responses, one per Chat call, and records every
// request it receives so tests can assert on the full conversation.
type stubClient struct {
	responses []llm.Response
	requests  []llm.Request
	callCount int
}

func (c *stubClient) Chat(_ context.Context, request llm.Request) (llm.Response, error) {
	c.requests = append(c.requests, request)

	if c.callCount >= len(c.responses) {
		return llm.Response{}, nil
	}

	response := c.responses[c.callCount]
	c.callCount++
	return response, nil
}

// failingClient always returns the configured error.
type failingClient struct {
	err error
}

func (c failingClient) Chat(context.Context, llm.Request) (llm.Response, error) {
	return llm.Response{}, c.err
}

// fakeTool is a deterministic tool used to observe tool execution,
// the tool results and the tool errors inside the agent loop.
type fakeTool struct{}

func (fakeTool) Name() string {
	return "fake"
}

func (fakeTool) Description() string {
	return "A fake tool for testing."
}

func (fakeTool) Definition() llm.ToolDefinition {
	return llm.ToolDefinition{
		Name:        "fake",
		Description: "A fake tool for testing.",
		Parameters: map[string]any{
			"type":       "object",
			"properties": map[string]any{},
		},
	}
}

func (fakeTool) Execute(arguments string) (string, error) {
	if arguments == "boom" {
		return "", errors.New("fake tool exploded")
	}
	return "result(" + arguments + ")", nil
}

// finalResponse creates a plain response that contains no tool calls.
func finalResponse(content string) llm.Response {
	return llm.Response{
		Message: llm.Message{Role: "assistant", Content: content},
	}
}

// toolCallResponse creates a response that requests the given tool calls.
func toolCallResponse(calls ...llm.ToolCall) llm.Response {
	return llm.Response{
		Message:   llm.Message{Role: "assistant", Content: "calling tools"},
		ToolCalls: calls,
	}
}

func TestNewRegistersToolsByName(t *testing.T) {
	agent := New(&stubClient{}, fakeTool{})

	got, err := agent.executeTool(llm.ToolCall{Name: "fake", Arguments: "hello"})
	if err != nil {
		t.Fatalf("executeTool() error = %v", err)
	}
	if got != "result(hello)" {
		t.Errorf("executeTool() = %q, want %q", got, "result(hello)")
	}
}

func TestNewUnknownToolCallReturnsAnError(t *testing.T) {
	agent := New(&stubClient{})

	_, err := agent.executeTool(llm.ToolCall{Name: "unknown", Arguments: "x"})
	if err == nil {
		t.Fatal("executeTool() error = nil, want an error")
	}
}

func TestChatReturnsResponseWithoutToolCalls(t *testing.T) {
	final := finalResponse("done")
	client := &stubClient{responses: []llm.Response{final}}
	agent := New(client)
	request := llm.Request{
		Model:    "test-model",
		Messages: []llm.Message{{Role: "user", Content: "hi"}},
	}

	got, err := agent.Chat(context.Background(), request)
	if err != nil {
		t.Fatalf("Chat() error = %v", err)
	}

	if !reflect.DeepEqual(got, final) {
		t.Errorf("Chat() = %+v, want %+v", got, final)
	}

	if len(client.requests) != 1 {
		t.Fatalf("Chat() called the client %d times, want 1", len(client.requests))
	}

	received := client.requests[0]
	if received.Model != "test-model" {
		t.Errorf("Chat() model = %q, want %q", received.Model, "test-model")
	}
	if len(received.Messages) != 1 {
		t.Errorf("Chat() sent %d messages, want 1", len(received.Messages))
	}
}

func TestChatExecutesToolAndFeedsResultBack(t *testing.T) {
	client := &stubClient{responses: []llm.Response{
		toolCallResponse(llm.ToolCall{ID: "call-1", Name: "fake", Arguments: "hello"}),
		finalResponse("done"),
	}}
	agent := New(client, fakeTool{})
	request := llm.Request{
		Messages: []llm.Message{{Role: "user", Content: "run it"}},
	}

	got, err := agent.Chat(context.Background(), request)
	if err != nil {
		t.Fatalf("Chat() error = %v", err)
	}

	if got.Message.Content != "done" {
		t.Errorf("Chat() content = %q, want %q", got.Message.Content, "done")
	}

	if len(client.requests) != 2 {
		t.Fatalf("Chat() called the client %d times, want 2", len(client.requests))
	}

	wantMessages := []llm.Message{
		{Role: "user", Content: "run it"},
		{Role: "assistant", Content: "calling tools"},
		{Role: "tool", Content: "result(hello)", ToolCallID: "call-1"},
	}
	gotMessages := client.requests[1].Messages

	if len(gotMessages) != len(wantMessages) {
		t.Fatalf("second turn sent %d messages: %+v, want %d", len(gotMessages), gotMessages, len(wantMessages))
	}
	for i, want := range wantMessages {
		if !reflect.DeepEqual(gotMessages[i], want) {
			t.Errorf("message[%d] = %+v, want %+v", i, gotMessages[i], want)
		}
	}
}

func TestChatHandlesMultipleToolCallsInOneTurn(t *testing.T) {
	client := &stubClient{responses: []llm.Response{
		toolCallResponse(
			llm.ToolCall{ID: "call-1", Name: "fake", Arguments: "a"},
			llm.ToolCall{ID: "call-2", Name: "fake", Arguments: "b"},
		),
		finalResponse("done"),
	}}
	agent := New(client, fakeTool{})
	request := llm.Request{Messages: []llm.Message{{Role: "user", Content: "run it"}}}

	_, err := agent.Chat(context.Background(), request)
	if err != nil {
		t.Fatalf("Chat() error = %v", err)
	}

	gotMessages := client.requests[1].Messages
	wantMessages := []llm.Message{
		{Role: "user", Content: "run it"},
		{Role: "assistant", Content: "calling tools"},
		{Role: "tool", Content: "result(a)", ToolCallID: "call-1"},
		{Role: "tool", Content: "result(b)", ToolCallID: "call-2"},
	}

	if len(gotMessages) != len(wantMessages) {
		t.Fatalf("second turn sent %d messages: %+v, want %d", len(gotMessages), gotMessages, len(wantMessages))
	}
	for i, want := range wantMessages {
		if !reflect.DeepEqual(gotMessages[i], want) {
			t.Errorf("message[%d] = %+v, want %+v", i, gotMessages[i], want)
		}
	}
}

func TestChatReportsToolFailuresBackToTheModel(t *testing.T) {
	client := &stubClient{responses: []llm.Response{
		toolCallResponse(llm.ToolCall{ID: "call-1", Name: "fake", Arguments: "boom"}),
		finalResponse("done"),
	}}
	agent := New(client, fakeTool{})
	request := llm.Request{Messages: []llm.Message{{Role: "user", Content: "run it"}}}

	_, err := agent.Chat(context.Background(), request)
	if err != nil {
		t.Fatalf("Chat() error = %v", err)
	}

	gotMessages := client.requests[1].Messages
	wantMessages := []llm.Message{
		{Role: "user", Content: "run it"},
		{Role: "assistant", Content: "calling tools"},
		{Role: "tool", Content: "tool error: fake tool exploded", ToolCallID: "call-1"},
	}

	if len(gotMessages) != len(wantMessages) {
		t.Fatalf("second turn sent %d messages: %+v, want %d", len(gotMessages), gotMessages, len(wantMessages))
	}
	for i, want := range wantMessages {
		if !reflect.DeepEqual(gotMessages[i], want) {
			t.Errorf("message[%d] = %+v, want %+v", i, gotMessages[i], want)
		}
	}
}

func TestChatReportsUnknownToolCallsBackToTheModel(t *testing.T) {
	client := &stubClient{responses: []llm.Response{
		toolCallResponse(llm.ToolCall{ID: "call-1", Name: "ghost", Arguments: "x"}),
		finalResponse("done"),
	}}
	agent := New(client, fakeTool{})
	request := llm.Request{Messages: []llm.Message{{Role: "user", Content: "run it"}}}

	_, err := agent.Chat(context.Background(), request)
	if err != nil {
		t.Fatalf("Chat() error = %v", err)
	}

	gotMessages := client.requests[1].Messages
	wantMessages := []llm.Message{
		{Role: "user", Content: "run it"},
		{Role: "assistant", Content: "calling tools"},
		{Role: "tool", Content: "tool error: tool not found: ghost", ToolCallID: "call-1"},
	}

	if len(gotMessages) != len(wantMessages) {
		t.Fatalf("second turn sent %d messages: %+v, want %d", len(gotMessages), gotMessages, len(wantMessages))
	}
	for i, want := range wantMessages {
		if !reflect.DeepEqual(gotMessages[i], want) {
			t.Errorf("message[%d] = %+v, want %+v", i, gotMessages[i], want)
		}
	}
}

func TestChatReturnsClientErrors(t *testing.T) {
	wantErr := errors.New("model unavailable")
	agent := New(&failingClient{err: wantErr})

	_, err := agent.Chat(context.Background(), llm.Request{})
	if !errors.Is(err, wantErr) {
		t.Fatalf("Chat() error = %v, want %v", err, wantErr)
	}
}

func TestChatSendsToolDefinitions(t *testing.T) {
	client := &stubClient{
		responses: []llm.Response{
			finalResponse("done"),
		},
	}

	agent := New(client, fakeTool{})

	_, err := agent.Chat(context.Background(), llm.Request{})
	if err != nil {
		t.Fatalf("Chat() error = %v", err)
	}

	if len(client.requests) != 1 {
		t.Fatalf("client received %d requests, want 1", len(client.requests))
	}

	if len(client.requests[0].Tools) != 1 {
		t.Fatalf("request contains %d tools, want 1", len(client.requests[0].Tools))
	}

	got := client.requests[0].Tools[0]

	if got.Name != "fake" {
		t.Errorf("tool name = %q, want %q", got.Name, "fake")
	}

	if got.Description != "A fake tool for testing." {
		t.Errorf("tool description = %q, want %q", got.Description, "A fake tool for testing.")
	}
}

func TestChatPreservesHistoryBetweenCalls(t *testing.T) {
	client := &stubClient{
		responses: []llm.Response{
			{
				Message: llm.Message{
					Role:    "assistant",
					Content: "first response",
				},
			},
			{
				Message: llm.Message{
					Role:    "assistant",
					Content: "second response",
				},
			},
		},
	}

	agent := New(client)

	_, err := agent.Chat(context.Background(), llm.Request{
		Messages: []llm.Message{
			{
				Role:    "user",
				Content: "first message",
			},
		},
	})
	if err != nil {
		t.Fatalf("first Chat() error = %v", err)
	}

	_, err = agent.Chat(context.Background(), llm.Request{
		Messages: []llm.Message{
			{
				Role:    "user",
				Content: "second message",
			},
		},
	})
	if err != nil {
		t.Fatalf("second Chat() error = %v", err)
	}

	if len(client.requests) != 2 {
		t.Fatalf("got %d requests, want 2", len(client.requests))
	}

	history := client.requests[1].Messages

	if len(history) != 3 {
		t.Fatalf("got %d messages in history, want 3", len(history))
	}

	if history[0].Content != "first message" {
		t.Errorf("history[0] = %q, want %q", history[0].Content, "first message")
	}

	if history[1].Content != "first response" {
		t.Errorf("history[1] = %q, want %q", history[1].Content, "first response")
	}

	if history[2].Content != "second message" {
		t.Errorf("history[2] = %q, want %q", history[2].Content, "second message")
	}
}

func TestTruncateToolResult(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "result shorter than limit",
			input: "hello",
			want:  "hello",
		},
		{
			name:  "empty result",
			input: "",
			want:  "",
		},
		{
			name:  "result exactly at limit",
			input: strings.Repeat("a", maxToolResultLength),
			want:  strings.Repeat("a", maxToolResultLength),
		},
		{
			name:  "result longer than limit",
			input: strings.Repeat("a", maxToolResultLength+100),
			want:  strings.Repeat("a", maxToolResultLength) + "\n[tool result truncated]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := truncateToolResult(tt.input)

			if got != tt.want {
				t.Errorf("truncateToolResult() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestRecentHistory(t *testing.T) {
	tests := []struct {
		name          string
		historyLength int
		wantLength    int
	}{
		{
			name:          "history shorter than limit",
			historyLength: 10,
			wantLength:    10,
		},
		{
			name:          "history exactly at limit",
			historyLength: maxHistoryMessages,
			wantLength:    maxHistoryMessages,
		},
		{
			name:          "history longer than limit",
			historyLength: maxHistoryMessages + 10,
			wantLength:    maxHistoryMessages,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			history := make([]llm.Message, tt.historyLength)

			got := recentHistory(history)

			if len(got) != tt.wantLength {
				t.Errorf("recentHistory() returned %d messages, want %d", len(got), tt.wantLength)
			}
		})
	}
}

func TestRecentHistoryKeepsRecentMessages(t *testing.T) {
	history := make([]llm.Message, 120)

	for i := range history {
		history[i].Content = fmt.Sprintf("message-%d", i)
	}

	got := recentHistory(history)

	if len(got) != maxHistoryMessages {
		t.Fatalf("recentHistory() returned %d messages, want %d", len(got), maxHistoryMessages)
	}

	if got[0].Content != "message-20" {
		t.Errorf("first message = %q, want %q", got[0].Content, "message-20")
	}

	if got[len(got)-1].Content != "message-119" {
		t.Errorf(
			"last message = %q, want %q",
			got[len(got)-1].Content,
			"message-119",
		)
	}
}

func TestRecentHistoryPreservesToolCallTurn(t *testing.T) {
	history := make([]llm.Message, 0, 102)

	for i := 0; i < 100; i++ {
		history = append(history, llm.Message{
			Role:    "user",
			Content: fmt.Sprintf("message-%d", i),
		})
	}

	history = append(history,
		llm.Message{
			Role:    "assistant",
			Content: "I need to run a command.",
			ToolCalls: []llm.ToolCall{
				{
					ID:   "call-1",
					Name: "shell",
				},
			},
		},
		llm.Message{
			Role:       "tool",
			Content:    "command output",
			ToolCallID: "call-1",
		},
	)

	// We expect recentHistory to preserve the complete tool turn.
	got := recentHistory(history)

	if len(got) != 100 {
		t.Fatalf("recentHistory() returned %d messages, want %d", len(got), 100)
	}

	if len(got[97].ToolCalls) != 1 {
		t.Fatal("tool call was separated from its result")
	}

	if got[98].ToolCallID != "call-1" {
		t.Errorf(
			"tool result ToolCallID = %q, want %q",
			got[98].ToolCallID,
			"call-1",
		)
	}
}
