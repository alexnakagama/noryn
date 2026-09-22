package agent

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/alexnakagama/noryn/internal/llm"
)

func userMsg(content string) llm.Message {
	return llm.Message{Role: "user", Content: content}
}

func toolMsg(id string) llm.Message {
	return llm.Message{Role: "tool", Content: "result-" + id, ToolCallID: id}
}

func assistantWithToolCalls(ids ...string) llm.Message {
	calls := make([]llm.ToolCall, 0, len(ids))
	for _, id := range ids {
		calls = append(calls, llm.ToolCall{ID: id, Name: "fake"})
	}

	return llm.Message{Role: "assistant", Content: "calling tools", ToolCalls: calls}
}

func numberedHistory(length int) []llm.Message {
	history := make([]llm.Message, length)
	for i := range history {
		history[i] = userMsg(fmt.Sprintf("m%03d", i))
	}

	return history
}

func newTestManager(maxTokens int) *ContextManager {
	return &ContextManager{maxTokens: maxTokens, tokenCounter: EstimateTokenCounter{}}
}

func TestBuildEmptyHistory(t *testing.T) {
	manager := NewContextManager()

	if got := manager.Build(nil); got != nil {
		t.Errorf("Build(nil) = %+v, want nil", got)
	}

	if got := manager.Build([]llm.Message{}); len(got) != 0 {
		t.Errorf("Build(empty) returned %d messages, want 0", len(got))
	}
}

func TestBuildReturnsSubslice(t *testing.T) {
	history := numberedHistory(3)
	manager := newTestManager(5)

	got := manager.Build(history)

	if len(got) != 3 {
		t.Fatalf("Build() returned %d messages, want 3", len(got))
	}

	if &got[0] != &history[0] {
		t.Error("Build() copied the history instead of returning the original slice")
	}
}

func TestBuildDoesNotMutateInput(t *testing.T) {
	history := []llm.Message{
		userMsg("u"),
		assistantWithToolCalls("call-1", "call-2"),
		toolMsg("call-1"),
		toolMsg("call-2"),
	}
	before := append([]llm.Message(nil), history...)

	manager := newTestManager(6)
	manager.Build(history)

	if !reflect.DeepEqual(history, before) {
		t.Errorf("Build() mutated the input history: got %+v, want %+v", history, before)
	}
}

func TestBuild(t *testing.T) {
	zeroTokenHistory := make([]llm.Message, 10)
	for i := range zeroTokenHistory {
		zeroTokenHistory[i] = userMsg("a")
	}

	tests := []struct {
		name               string
		maxTokens          int
		history            []llm.Message
		wantLength         int
		wantFirstRole      string
		wantFirstToolCalls int
		wantFirstContent   string
		wantLastToolCallID string
	}{
		{
			name:             "below limit returns all",
			maxTokens:        5,
			history:          numberedHistory(3),
			wantLength:       3,
			wantFirstRole:    "user",
			wantFirstContent: "m000",
		},
		{
			name:             "exactly at limit",
			maxTokens:        3,
			history:          numberedHistory(3),
			wantLength:       3,
			wantFirstRole:    "user",
			wantFirstContent: "m000",
		},
		{
			name:             "one over limit",
			maxTokens:        3,
			history:          numberedHistory(4),
			wantLength:       3,
			wantFirstRole:    "user",
			wantFirstContent: "m001",
		},
		{
			name:             "well over limit",
			maxTokens:        3,
			history:          numberedHistory(10),
			wantLength:       3,
			wantFirstRole:    "user",
			wantFirstContent: "m007",
		},
		{
			name:          "zero-token messages all fit",
			maxTokens:     100,
			history:       zeroTokenHistory,
			wantLength:    10,
			wantFirstRole: "user",
		},
		{
			name:               "backoff starts on assistant",
			maxTokens:          3,
			history:            []llm.Message{userMsg("u"), assistantWithToolCalls("call-1"), toolMsg("call-1")},
			wantLength:         2,
			wantFirstRole:      "assistant",
			wantFirstToolCalls: 1,
			wantLastToolCallID: "call-1",
		},
		{
			name:               "multiple tool results kept together",
			maxTokens:          6,
			history:            []llm.Message{userMsg("u"), assistantWithToolCalls("call-1", "call-2"), toolMsg("call-1"), toolMsg("call-2")},
			wantLength:         3,
			wantFirstRole:      "assistant",
			wantFirstToolCalls: 2,
			wantLastToolCallID: "call-2",
		},
		{
			name:               "backoff reaches start of history",
			maxTokens:          9,
			history:            []llm.Message{assistantWithToolCalls("call-1"), toolMsg("call-1"), toolMsg("call-1"), toolMsg("call-1")},
			wantLength:         4,
			wantFirstRole:      "assistant",
			wantFirstToolCalls: 1,
			wantLastToolCallID: "call-1",
		},
		{
			name:               "all tool messages",
			maxTokens:          3,
			history:            []llm.Message{toolMsg("call-1"), toolMsg("call-1"), toolMsg("call-1"), toolMsg("call-1")},
			wantLength:         4,
			wantFirstRole:      "tool",
			wantLastToolCallID: "call-1",
		},
		{
			name:               "backoff stops at non-tool message",
			maxTokens:          3,
			history:            []llm.Message{userMsg("u"), llm.Message{Role: "assistant", Content: "plain"}, toolMsg("call-1"), toolMsg("call-1")},
			wantLength:         3,
			wantFirstRole:      "assistant",
			wantLastToolCallID: "call-1",
		},
		{
			name:          "exactly at token budget",
			maxTokens:     maxContextTokens,
			history:       []llm.Message{userMsg(strings.Repeat("a", maxContextTokens*4))},
			wantLength:    1,
			wantFirstRole: "user",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := newTestManager(tt.maxTokens)

			got := manager.Build(tt.history)

			if len(got) != tt.wantLength {
				t.Fatalf("Build() returned %d messages, want %d", len(got), tt.wantLength)
			}

			if got[0].Role != tt.wantFirstRole {
				t.Errorf("first message role = %q, want %q", got[0].Role, tt.wantFirstRole)
			}

			if len(got[0].ToolCalls) != tt.wantFirstToolCalls {
				t.Errorf("first message tool calls = %d, want %d", len(got[0].ToolCalls), tt.wantFirstToolCalls)
			}

			if tt.wantFirstContent != "" && got[0].Content != tt.wantFirstContent {
				t.Errorf("first message content = %q, want %q", got[0].Content, tt.wantFirstContent)
			}

			if tt.wantLastToolCallID != "" && got[len(got)-1].ToolCallID != tt.wantLastToolCallID {
				t.Errorf("last message ToolCallID = %q, want %q", got[len(got)-1].ToolCallID, tt.wantLastToolCallID)
			}
		})
	}
}
