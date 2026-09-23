package openrouter

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alexnakagama/noryn/internal/llm"
)

func TestClient_ChatStream(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")

		fmt.Fprintln(w, `data: {"choices":[{"delta":{"content":"Hello"}}]}`)
		fmt.Fprintln(w, `data: {"choices":[{"delta":{"content":" world"}}]}`)
		fmt.Fprintln(w, `data: [DONE]`)
	}))
	defer server.Close()

	client := NewClient("test-key")
	client.baseURL = server.URL

	stream, err := client.ChatStream(context.Background(), llm.Request{
		Model: "test-model",
		Messages: []llm.Message{
			{
				Role:    "user",
				Content: "Say hello",
			},
		},
	})

	if err != nil {
		t.Fatalf("ChatStream() error = %v", err)
	}

	var content string
	var done bool

	for chunk := range stream {
		if chunk.Err != nil {
			t.Fatalf("stream error = %v", chunk.Err)
		}

		content += chunk.Content

		if chunk.Done {
			done = true
		}
	}

	if content != "Hello world" {
		t.Errorf("content = %q, want %q", content, "Hello world")
	}

	if !done {
		t.Error("stream never returned Done=true")
	}
}

func TestClient_ChatStream_ToolCall(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")

		fmt.Fprintln(w, `data: {"choices":[{"delta":{"tool_calls":[{"index":0,"id":"call_1","type":"function","function":{"name":"read_file"}}]}}]}`)
		fmt.Fprintln(w, `data: {"choices":[{"delta":{"tool_calls":[{"index":0,"function":{"arguments":"{\"path\":"}}]}}]}`)
		fmt.Fprintln(w, `data: {"choices":[{"delta":{"tool_calls":[{"index":0,"function":{"arguments":"\"main.go\"}"}}]}}]}`)
		fmt.Fprintln(w, `data: [DONE]`)
	}))
	defer server.Close()

	client := NewClient("test-key")
	client.baseURL = server.URL

	stream, err := client.ChatStream(context.Background(), llm.Request{
		Model: "test-model",
		Messages: []llm.Message{
			{
				Role:    "user",
				Content: "Read main.go",
			},
		},
	})

	if err != nil {
		t.Fatalf("ChatStream() error = %v", err)
	}

	var toolCalls []llm.ToolCall

	for chunk := range stream {
		if chunk.Err != nil {
			t.Fatalf("stream error = %v", chunk.Err)
		}

		if chunk.ToolCall != nil {
			toolCalls = append(toolCalls, *chunk.ToolCall)
		}
	}

	if len(toolCalls) != 1 {
		t.Fatalf("got %d tool calls, want 1", len(toolCalls))
	}

	call := toolCalls[0]

	if call.ID != "call_1" {
		t.Errorf("ToolCall.ID = %q, want %q", call.ID, "call_1")
	}

	if call.Name != "read_file" {
		t.Errorf("ToolCall.Name = %q, want %q", call.Name, "read_file")
	}

	if call.Arguments != `{"path":"main.go"}` {
		t.Errorf(
			"ToolCall.Arguments = %q, want %q",
			call.Arguments,
			`{"path":"main.go"}`,
		)
	}
}
