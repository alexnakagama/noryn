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
