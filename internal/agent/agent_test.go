package agent

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/alexnakagama/noryn/internal/llm"
	"github.com/alexnakagama/noryn/internal/tools"
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

// toolLoopClient always requests the same tool, which lets tests exercise the
// maximum tool iteration guard.
type toolLoopClient struct {
	requests  []llm.Request
	callCount int
}

func (c *toolLoopClient) Chat(_ context.Context, request llm.Request) (llm.Response, error) {
	c.requests = append(c.requests, request)
	c.callCount++

	return toolCallResponse(llm.ToolCall{ID: "call-1", Name: "fake", Arguments: "x"}), nil
}

// clientStep is a single scripted response or error for scriptedClient.
type clientStep struct {
	response llm.Response
	err      error
}

// scriptedClient returns a scripted sequence of responses and errors.
type scriptedClient struct {
	steps []clientStep
	calls int
}

func (c *scriptedClient) Chat(_ context.Context, _ llm.Request) (llm.Response, error) {
	if c.calls >= len(c.steps) {
		return llm.Response{}, nil
	}

	step := c.steps[c.calls]
	c.calls++
	return step.response, step.err
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

// longResultTool returns a result larger than maxToolResultLength so tests can
// observe truncation inside the agent loop.
type longResultTool struct{}

func (longResultTool) Name() string {
	return "long"
}

func (longResultTool) Description() string {
	return "A tool that returns a very long result."
}

func (longResultTool) Definition() llm.ToolDefinition {
	return llm.ToolDefinition{
		Name:        "long",
		Description: "A tool that returns a very long result.",
	}
}

func (longResultTool) Execute(string) (string, error) {
	return strings.Repeat("x", maxToolResultLength+100), nil
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

// assertMessages compares two message slices field by field.
func assertMessages(t *testing.T, got, want []llm.Message) {
	t.Helper()

	if len(got) != len(want) {
		t.Fatalf("got %d messages: %+v, want %d", len(got), got, len(want))
	}

	for i := range want {
		if !reflect.DeepEqual(got[i], want[i]) {
			t.Errorf("message[%d] = %+v, want %+v", i, got[i], want[i])
		}
	}
}

func TestChatReturnsResponseWithoutToolCalls(t *testing.T) {
	final := finalResponse("done")
	client := &stubClient{responses: []llm.Response{final}}
	agent := New(client, tools.NewRegistry())
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
	agent := New(client, tools.NewRegistry(fakeTool{}))
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
	assertMessages(t, client.requests[1].Messages, wantMessages)
}

func TestChatHandlesMultipleToolCallsInOneTurn(t *testing.T) {
	client := &stubClient{responses: []llm.Response{
		toolCallResponse(
			llm.ToolCall{ID: "call-1", Name: "fake", Arguments: "a"},
			llm.ToolCall{ID: "call-2", Name: "fake", Arguments: "b"},
		),
		finalResponse("done"),
	}}
	agent := New(client, tools.NewRegistry(fakeTool{}))
	request := llm.Request{Messages: []llm.Message{{Role: "user", Content: "run it"}}}

	_, err := agent.Chat(context.Background(), request)
	if err != nil {
		t.Fatalf("Chat() error = %v", err)
	}

	wantMessages := []llm.Message{
		{Role: "user", Content: "run it"},
		{Role: "assistant", Content: "calling tools"},
		{Role: "tool", Content: "result(a)", ToolCallID: "call-1"},
		{Role: "tool", Content: "result(b)", ToolCallID: "call-2"},
	}
	assertMessages(t, client.requests[1].Messages, wantMessages)
}

func TestChatReportsToolFailuresBackToTheModel(t *testing.T) {
	client := &stubClient{responses: []llm.Response{
		toolCallResponse(llm.ToolCall{ID: "call-1", Name: "fake", Arguments: "boom"}),
		finalResponse("done"),
	}}
	agent := New(client, tools.NewRegistry(fakeTool{}))
	request := llm.Request{Messages: []llm.Message{{Role: "user", Content: "run it"}}}

	_, err := agent.Chat(context.Background(), request)
	if err != nil {
		t.Fatalf("Chat() error = %v", err)
	}

	wantMessages := []llm.Message{
		{Role: "user", Content: "run it"},
		{Role: "assistant", Content: "calling tools"},
		{Role: "tool", Content: "tool error: fake tool exploded", ToolCallID: "call-1"},
	}
	assertMessages(t, client.requests[1].Messages, wantMessages)
}

func TestChatReportsUnknownToolCallsBackToTheModel(t *testing.T) {
	client := &stubClient{responses: []llm.Response{
		toolCallResponse(llm.ToolCall{ID: "call-1", Name: "ghost", Arguments: "x"}),
		finalResponse("done"),
	}}
	agent := New(client, tools.NewRegistry(fakeTool{}))
	request := llm.Request{Messages: []llm.Message{{Role: "user", Content: "run it"}}}

	_, err := agent.Chat(context.Background(), request)
	if err != nil {
		t.Fatalf("Chat() error = %v", err)
	}

	wantMessages := []llm.Message{
		{Role: "user", Content: "run it"},
		{Role: "assistant", Content: "calling tools"},
		{Role: "tool", Content: "tool error: tool not found: ghost", ToolCallID: "call-1"},
	}
	assertMessages(t, client.requests[1].Messages, wantMessages)
}

func TestChatReturnsClientErrors(t *testing.T) {
	wantErr := errors.New("model unavailable")
	agent := New(&failingClient{err: wantErr}, tools.NewRegistry())

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

	agent := New(client, tools.NewRegistry(fakeTool{}))

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

func TestChatSendsToolDefinitionsOnEveryIteration(t *testing.T) {
	client := &stubClient{responses: []llm.Response{
		toolCallResponse(llm.ToolCall{ID: "call-1", Name: "fake", Arguments: "x"}),
		finalResponse("done"),
	}}
	agent := New(client, tools.NewRegistry(fakeTool{}))

	_, err := agent.Chat(context.Background(), llm.Request{
		Messages: []llm.Message{{Role: "user", Content: "run it"}},
	})
	if err != nil {
		t.Fatalf("Chat() error = %v", err)
	}

	if len(client.requests) != 2 {
		t.Fatalf("client received %d requests, want 2", len(client.requests))
	}

	for i, request := range client.requests {
		if len(request.Tools) != 1 {
			t.Errorf("request[%d] contains %d tools, want 1", i, len(request.Tools))
		}
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

	agent := New(client, tools.NewRegistry())

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

func TestChatAccumulatesHistoryAcrossToolCalls(t *testing.T) {
	client := &stubClient{responses: []llm.Response{
		toolCallResponse(llm.ToolCall{ID: "call-1", Name: "fake", Arguments: "hello"}),
		finalResponse("first done"),
		finalResponse("second done"),
	}}
	agent := New(client, tools.NewRegistry(fakeTool{}))

	_, err := agent.Chat(context.Background(), llm.Request{
		Messages: []llm.Message{{Role: "user", Content: "first message"}},
	})
	if err != nil {
		t.Fatalf("first Chat() error = %v", err)
	}

	_, err = agent.Chat(context.Background(), llm.Request{
		Messages: []llm.Message{{Role: "user", Content: "second message"}},
	})
	if err != nil {
		t.Fatalf("second Chat() error = %v", err)
	}

	if len(client.requests) != 3 {
		t.Fatalf("client received %d requests, want 3", len(client.requests))
	}

	assertMessages(t, client.requests[0].Messages, []llm.Message{
		{Role: "user", Content: "first message"},
	})

	assertMessages(t, client.requests[1].Messages, []llm.Message{
		{Role: "user", Content: "first message"},
		{Role: "assistant", Content: "calling tools"},
		{Role: "tool", Content: "result(hello)", ToolCallID: "call-1"},
	})

	assertMessages(t, client.requests[2].Messages, []llm.Message{
		{Role: "user", Content: "first message"},
		{Role: "assistant", Content: "calling tools"},
		{Role: "tool", Content: "result(hello)", ToolCallID: "call-1"},
		{Role: "assistant", Content: "first done"},
		{Role: "user", Content: "second message"},
	})
}

func TestChatReturnsClientErrorMidLoop(t *testing.T) {
	wantErr := errors.New("model unavailable")
	client := &scriptedClient{steps: []clientStep{
		{response: toolCallResponse(llm.ToolCall{ID: "call-1", Name: "fake", Arguments: "x"})},
		{err: wantErr},
	}}
	agent := New(client, tools.NewRegistry(fakeTool{}))

	_, err := agent.Chat(context.Background(), llm.Request{
		Messages: []llm.Message{{Role: "user", Content: "run it"}},
	})
	if !errors.Is(err, wantErr) {
		t.Fatalf("Chat() error = %v, want %v", err, wantErr)
	}

	if client.calls != 2 {
		t.Fatalf("client called %d times, want 2", client.calls)
	}

	assertMessages(t, agent.history, []llm.Message{
		{Role: "user", Content: "run it"},
		{Role: "assistant", Content: "calling tools"},
		{Role: "tool", Content: "result(x)", ToolCallID: "call-1"},
	})
}

func TestChatEnforcesMaxToolIterations(t *testing.T) {
	client := &toolLoopClient{}
	agent := New(client, tools.NewRegistry(fakeTool{}))

	var callCount, resultCount int
	agent.SetToolCallHandler(func(llm.ToolCall) { callCount++ })
	agent.SetToolResultHandler(func(llm.ToolCall, string) { resultCount++ })

	_, err := agent.Chat(context.Background(), llm.Request{
		Messages: []llm.Message{{Role: "user", Content: "loop"}},
	})
	if err == nil {
		t.Fatal("Chat() error = nil, want maximum tool iterations error")
	}

	if err.Error() != "maximum tool iterations exceeded" {
		t.Errorf("Chat() error = %q, want %q", err.Error(), "maximum tool iterations exceeded")
	}

	if client.callCount != maxToolIterations+1 {
		t.Errorf("client called %d times, want %d", client.callCount, maxToolIterations+1)
	}

	if callCount != maxToolIterations {
		t.Errorf("tool call handler invoked %d times, want %d", callCount, maxToolIterations)
	}

	if resultCount != maxToolIterations {
		t.Errorf("tool result handler invoked %d times, want %d", resultCount, maxToolIterations)
	}

	if len(agent.history) != 1+2*maxToolIterations {
		t.Errorf("history has %d messages, want %d", len(agent.history), 1+2*maxToolIterations)
	}
}

func TestChatToolCallHandlerReceivesEachCall(t *testing.T) {
	client := &stubClient{responses: []llm.Response{
		toolCallResponse(
			llm.ToolCall{ID: "call-1", Name: "fake", Arguments: "a"},
			llm.ToolCall{ID: "call-2", Name: "fake", Arguments: "b"},
		),
		finalResponse("done"),
	}}
	agent := New(client, tools.NewRegistry(fakeTool{}))

	var calls []llm.ToolCall
	agent.SetToolCallHandler(func(call llm.ToolCall) {
		calls = append(calls, call)
	})

	_, err := agent.Chat(context.Background(), llm.Request{
		Messages: []llm.Message{{Role: "user", Content: "run it"}},
	})
	if err != nil {
		t.Fatalf("Chat() error = %v", err)
	}

	want := []llm.ToolCall{
		{ID: "call-1", Name: "fake", Arguments: "a"},
		{ID: "call-2", Name: "fake", Arguments: "b"},
	}

	if !reflect.DeepEqual(calls, want) {
		t.Errorf("tool call handler received %+v, want %+v", calls, want)
	}
}

func TestChatToolResultHandlerReceivesResults(t *testing.T) {
	type result struct {
		call   llm.ToolCall
		result string
	}

	client := &stubClient{responses: []llm.Response{
		toolCallResponse(
			llm.ToolCall{ID: "call-1", Name: "fake", Arguments: "a"},
			llm.ToolCall{ID: "call-2", Name: "fake", Arguments: "b"},
		),
		finalResponse("done"),
	}}
	agent := New(client, tools.NewRegistry(fakeTool{}))

	var results []result
	agent.SetToolResultHandler(func(call llm.ToolCall, value string) {
		results = append(results, result{call: call, result: value})
	})

	_, err := agent.Chat(context.Background(), llm.Request{
		Messages: []llm.Message{{Role: "user", Content: "run it"}},
	})
	if err != nil {
		t.Fatalf("Chat() error = %v", err)
	}

	want := []result{
		{call: llm.ToolCall{ID: "call-1", Name: "fake", Arguments: "a"}, result: "result(a)"},
		{call: llm.ToolCall{ID: "call-2", Name: "fake", Arguments: "b"}, result: "result(b)"},
	}

	if !reflect.DeepEqual(results, want) {
		t.Errorf("tool result handler received %+v, want %+v", results, want)
	}
}

func TestChatToolResultHandlerReceivesErrorResults(t *testing.T) {
	type result struct {
		call   llm.ToolCall
		result string
	}

	client := &stubClient{responses: []llm.Response{
		toolCallResponse(llm.ToolCall{ID: "call-1", Name: "fake", Arguments: "boom"}),
		finalResponse("done"),
	}}
	agent := New(client, tools.NewRegistry(fakeTool{}))

	var results []result
	agent.SetToolResultHandler(func(call llm.ToolCall, value string) {
		results = append(results, result{call: call, result: value})
	})

	_, err := agent.Chat(context.Background(), llm.Request{
		Messages: []llm.Message{{Role: "user", Content: "run it"}},
	})
	if err != nil {
		t.Fatalf("Chat() error = %v", err)
	}

	want := []result{
		{
			call:   llm.ToolCall{ID: "call-1", Name: "fake", Arguments: "boom"},
			result: "tool error: fake tool exploded",
		},
	}

	if !reflect.DeepEqual(results, want) {
		t.Errorf("tool result handler received %+v, want %+v", results, want)
	}
}

func TestChatToolResultHandlerReceivesUnknownToolError(t *testing.T) {
	type result struct {
		call   llm.ToolCall
		result string
	}

	client := &stubClient{responses: []llm.Response{
		toolCallResponse(llm.ToolCall{ID: "call-1", Name: "ghost", Arguments: "x"}),
		finalResponse("done"),
	}}
	agent := New(client, tools.NewRegistry(fakeTool{}))

	var results []result
	agent.SetToolResultHandler(func(call llm.ToolCall, value string) {
		results = append(results, result{call: call, result: value})
	})

	_, err := agent.Chat(context.Background(), llm.Request{
		Messages: []llm.Message{{Role: "user", Content: "run it"}},
	})
	if err != nil {
		t.Fatalf("Chat() error = %v", err)
	}

	want := []result{
		{
			call:   llm.ToolCall{ID: "call-1", Name: "ghost", Arguments: "x"},
			result: "tool error: tool not found: ghost",
		},
	}

	if !reflect.DeepEqual(results, want) {
		t.Errorf("tool result handler received %+v, want %+v", results, want)
	}
}

func TestChatHandlersNotInvokedWithoutToolCalls(t *testing.T) {
	client := &stubClient{responses: []llm.Response{finalResponse("done")}}
	agent := New(client, tools.NewRegistry(fakeTool{}))

	var callCount, resultCount int
	agent.SetToolCallHandler(func(llm.ToolCall) { callCount++ })
	agent.SetToolResultHandler(func(llm.ToolCall, string) { resultCount++ })

	_, err := agent.Chat(context.Background(), llm.Request{
		Messages: []llm.Message{{Role: "user", Content: "hi"}},
	})
	if err != nil {
		t.Fatalf("Chat() error = %v", err)
	}

	if callCount != 0 {
		t.Errorf("tool call handler invoked %d times, want 0", callCount)
	}

	if resultCount != 0 {
		t.Errorf("tool result handler invoked %d times, want 0", resultCount)
	}
}

func TestChatTruncatesLongToolResultsInHistory(t *testing.T) {
	client := &stubClient{responses: []llm.Response{
		toolCallResponse(llm.ToolCall{ID: "call-1", Name: "long"}),
		finalResponse("done"),
	}}
	agent := New(client, tools.NewRegistry(longResultTool{}))

	var handledResult string
	agent.SetToolResultHandler(func(_ llm.ToolCall, result string) {
		handledResult = result
	})

	_, err := agent.Chat(context.Background(), llm.Request{
		Messages: []llm.Message{{Role: "user", Content: "run it"}},
	})
	if err != nil {
		t.Fatalf("Chat() error = %v", err)
	}

	wantContent := strings.Repeat("x", maxToolResultLength) + "\n[tool result truncated]"

	if handledResult != wantContent {
		t.Errorf("tool result handler received a result of length %d, want %d", len(handledResult), len(wantContent))
	}

	messages := client.requests[1].Messages
	last := messages[len(messages)-1]

	if last.Role != "tool" {
		t.Fatalf("last message role = %q, want %q", last.Role, "tool")
	}

	if last.Content != wantContent {
		t.Errorf("tool message content length = %d, want %d", len(last.Content), len(wantContent))
	}
}

func TestChatPassesBuiltHistoryToClient(t *testing.T) {
	client := &stubClient{responses: []llm.Response{finalResponse("done")}}
	agent := New(client, tools.NewRegistry())
	agent.context = newTestManager(3)

	history := numberedHistory(5)

	_, err := agent.Chat(context.Background(), llm.Request{Messages: history})
	if err != nil {
		t.Fatalf("Chat() error = %v", err)
	}

	assertMessages(t, client.requests[0].Messages, []llm.Message{
		{Role: "user", Content: "m002"},
		{Role: "user", Content: "m003"},
		{Role: "user", Content: "m004"},
	})
}

func TestChatRebuildsHistoryAfterToolTurn(t *testing.T) {
	client := &stubClient{responses: []llm.Response{
		toolCallResponse(llm.ToolCall{ID: "call-1", Name: "fake", Arguments: "x"}),
		finalResponse("done"),
	}}
	agent := New(client, tools.NewRegistry(fakeTool{}))
	agent.context = newTestManager(6)

	_, err := agent.Chat(context.Background(), llm.Request{
		Messages: []llm.Message{
			{Role: "user", Content: "u000"},
			{Role: "user", Content: "u001"},
			{Role: "user", Content: "u002"},
		},
	})
	if err != nil {
		t.Fatalf("Chat() error = %v", err)
	}

	if len(client.requests) != 2 {
		t.Fatalf("client received %d requests, want 2", len(client.requests))
	}

	assertMessages(t, client.requests[0].Messages, []llm.Message{
		{Role: "user", Content: "u000"},
		{Role: "user", Content: "u001"},
		{Role: "user", Content: "u002"},
	})

	assertMessages(t, client.requests[1].Messages, []llm.Message{
		{Role: "user", Content: "u002"},
		{Role: "assistant", Content: "calling tools"},
		{Role: "tool", Content: "result(x)", ToolCallID: "call-1"},
	})
}

func TestChatTrimsHistoryAtRealLimit(t *testing.T) {
	client := &stubClient{responses: []llm.Response{finalResponse("done")}}
	agent := New(client, tools.NewRegistry())

	history := make([]llm.Message, 0, 105)
	for i := 0; i < 105; i++ {
		history = append(history, llm.Message{
			Role:    "user",
			Content: strings.Repeat("a", maxContextTokens*4/100),
		})
	}

	_, err := agent.Chat(context.Background(), llm.Request{Messages: history})
	if err != nil {
		t.Fatalf("Chat() error = %v", err)
	}

	got := client.requests[0].Messages
	if len(got) != 100 {
		t.Fatalf("client received %d messages, want %d", len(got), 100)
	}

	if got[0].Content != history[5].Content {
		t.Errorf("first message = %q, want %q", got[0].Content, history[5].Content)
	}

	if got[len(got)-1].Content != history[104].Content {
		t.Errorf("last message = %q, want %q", got[len(got)-1].Content, history[104].Content)
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
