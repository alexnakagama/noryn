package openai

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/alexnakagama/noryn/internal/llm"
)

const validOpenAIResponse = `{"output":[{"type":"message","status":"completed","role":"assistant","content":[{"type":"output_text","text":"ok","annotations":[]}]}]}`

func TestClient_SerializesToolResultMessage(t *testing.T) {
	var capturedBody []byte

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertRequestMetadata(t, r)

		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("read request body: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		capturedBody = body

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if _, err := w.Write([]byte(validOpenAIResponse)); err != nil {
			t.Errorf("write response: %v", err)
		}
	}))
	defer server.Close()

	client := NewClient("test-api-key")
	client.baseURL = server.URL

	if _, err := client.Chat(context.Background(), llm.Request{
		Model: "gpt-5",
		Messages: []llm.Message{
			{Role: "tool", Content: "result", ToolCallID: "call-1"},
		},
	}); err != nil {
		t.Fatalf("Chat() error = %v", err)
	}

	t.Run("tool message uses a function_call_output input", func(t *testing.T) {
		var payload struct {
			Model string `json:"model"`
			Input []struct {
				Type   string `json:"type"`
				CallID string `json:"call_id"`
				Output string `json:"output"`
			} `json:"input"`
		}

		if err := json.Unmarshal(capturedBody, &payload); err != nil {
			t.Fatalf("json.Unmarshal() error = %v, body: %s", err, capturedBody)
		}

		if len(payload.Input) != 1 {
			t.Fatalf("input length = %d, want 1", len(payload.Input))
		}

		input := payload.Input[0]

		if input.Type != "function_call_output" {
			t.Errorf("input.type = %q, want %q", input.Type, "function_call_output")
		}

		if input.CallID != "call-1" {
			t.Errorf("input.call_id = %q, want %q", input.CallID, "call-1")
		}

		if input.Output != "result" {
			t.Errorf("input.output = %q, want %q", input.Output, "result")
		}
	})

	t.Run("serialized output matches the OpenAI wire format", func(t *testing.T) {
		want := `{"type":"function_call_output","call_id":"call-1","output":"result"}`

		if !strings.Contains(string(capturedBody), want) {
			t.Errorf("request body = %s, want it to contain %s", capturedBody, want)
		}
	})
}

func assertRequestMetadata(t *testing.T, r *http.Request) {
	t.Helper()

	if r.Method != http.MethodPost {
		t.Errorf("method = %q, want %q", r.Method, http.MethodPost)
	}

	if r.URL.Path != "/responses" {
		t.Errorf("path = %q, want %q", r.URL.Path, "/responses")
	}

	if got := r.Header.Get("Authorization"); got != "Bearer test-api-key" {
		t.Errorf("Authorization = %q, want %q", got, "Bearer test-api-key")
	}

	if got := r.Header.Get("Content-Type"); got != "application/json" {
		t.Errorf("Content-Type = %q, want %q", got, "application/json")
	}
}
