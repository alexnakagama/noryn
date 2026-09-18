package agent

import (
	"context"
	"errors"
	"reflect"
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

func (fakeTool) Name() string { return "fake" }

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
