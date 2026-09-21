package agent

import (
	"fmt"
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
		history[i] = userMsg(fmt.Sprintf("m%d", i))
	}

	return history
}

func TestBuildAtLimitReturnsAll(t *testing.T) {
	manager := &ContextManager{maxMessages: 5}
	history := numberedHistory(5)

	got := manager.Build(history)

	if len(got) != 5 {
		t.Fatalf("Build() returned %d messages, want 5", len(got))
	}

	if &got[0] != &history[0] {
		t.Error("Build() copied the history instead of returning the original slice")
	}
}

func TestBuildEmptyHistory(t *testing.T) {
	manager := &ContextManager{maxMessages: 5}

	if got := manager.Build(nil); got != nil {
		t.Errorf("Build(nil) = %+v, want nil", got)
	}

	if got := manager.Build([]llm.Message{}); len(got) != 0 {
		t.Errorf("Build(empty) returned %d messages, want 0", len(got))
	}
}

func TestBuildTrimsToMostRecent(t *testing.T) {
	manager := &ContextManager{maxMessages: 3}

	for _, length := range []int{4, 5, 6, 10} {
		t.Run(fmt.Sprintf("length-%d", length), func(t *testing.T) {
			history := numberedHistory(length)

			got := manager.Build(history)

			if len(got) != 3 {
				t.Fatalf("Build() returned %d messages, want 3", len(got))
			}

			wantFirst := fmt.Sprintf("m%d", length-3)
			if got[0].Content != wantFirst {
				t.Errorf("first message = %q, want %q", got[0].Content, wantFirst)
			}

			wantLast := fmt.Sprintf("m%d", length-1)
			if got[len(got)-1].Content != wantLast {
				t.Errorf("last message = %q, want %q", got[len(got)-1].Content, wantLast)
			}
		})
	}
}

func TestBuildBoundaryLengths(t *testing.T) {
	manager := NewContextManager()

	tests := []struct {
		name   string
		length int
	}{
		{name: "one below limit", length: 99},
		{name: "exactly at limit", length: 100},
		{name: "one over limit", length: 101},
		{name: "two over limit", length: 102},
		{name: "three over limit", length: 103},
		{name: "well over limit", length: 150},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			history := numberedHistory(tt.length)

			got := manager.Build(history)

			wantStart := tt.length - maxContextMessages
			if wantStart < 0 {
				wantStart = 0
			}

			wantLength := tt.length - wantStart
			if len(got) != wantLength {
				t.Fatalf("Build() returned %d messages, want %d", len(got), wantLength)
			}

			wantFirst := fmt.Sprintf("m%d", wantStart)
			if got[0].Content != wantFirst {
				t.Errorf("first message = %q, want %q", got[0].Content, wantFirst)
			}

			wantLast := fmt.Sprintf("m%d", tt.length-1)
			if got[len(got)-1].Content != wantLast {
				t.Errorf("last message = %q, want %q", got[len(got)-1].Content, wantLast)
			}
		})
	}
}

func TestBuildDoesNotSplitSingleToolResult(t *testing.T) {
	history := []llm.Message{
		userMsg("u"),
		assistantWithToolCalls("call-1"),
		toolMsg("call-1"),
	}

	tests := []struct {
		name        string
		maxMessages int
		wantLength  int
	}{
		{name: "back-off triggered", maxMessages: 1, wantLength: 2},
		{name: "start lands on assistant", maxMessages: 2, wantLength: 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := &ContextManager{maxMessages: tt.maxMessages}

			got := manager.Build(history)

			if len(got) != tt.wantLength {
				t.Fatalf("Build() returned %d messages, want %d", len(got), tt.wantLength)
			}

			if got[0].Role != "assistant" {
				t.Fatalf("first message role = %q, want %q", got[0].Role, "assistant")
			}

			if len(got[0].ToolCalls) != 1 || got[0].ToolCalls[0].ID != "call-1" {
				t.Errorf("first message tool calls = %+v, want call-1", got[0].ToolCalls)
			}

			if got[1].Role != "tool" || got[1].ToolCallID != "call-1" {
				t.Errorf("second message = %+v, want tool result call-1", got[1])
			}
		})
	}
}

func TestBuildDoesNotSplitMultipleToolResults(t *testing.T) {
	history := []llm.Message{
		userMsg("u"),
		assistantWithToolCalls("call-1", "call-2"),
		toolMsg("call-1"),
		toolMsg("call-2"),
	}
	manager := &ContextManager{maxMessages: 2}

	got := manager.Build(history)

	if len(got) != 3 {
		t.Fatalf("Build() returned %d messages, want 3", len(got))
	}

	if got[0].Role != "assistant" || len(got[0].ToolCalls) != 2 {
		t.Fatalf("first message = %+v, want assistant with 2 tool calls", got[0])
	}

	if got[0].ToolCalls[0].ID != "call-1" || got[0].ToolCalls[1].ID != "call-2" {
		t.Errorf("tool call IDs = %+v, want call-1 then call-2", got[0].ToolCalls)
	}

	if got[1].ToolCallID != "call-1" {
		t.Errorf("second message ToolCallID = %q, want %q", got[1].ToolCallID, "call-1")
	}

	if got[2].ToolCallID != "call-2" {
		t.Errorf("third message ToolCallID = %q, want %q", got[2].ToolCallID, "call-2")
	}
}

func TestBuildBackoffReachesStartOfHistory(t *testing.T) {
	history := []llm.Message{
		assistantWithToolCalls("call-1"),
		toolMsg("call-1"),
		toolMsg("call-1"),
		toolMsg("call-1"),
	}
	manager := &ContextManager{maxMessages: 2}

	got := manager.Build(history)

	if len(got) != len(history) {
		t.Fatalf("Build() returned %d messages, want %d", len(got), len(history))
	}

	if got[0].Role != "assistant" || len(got[0].ToolCalls) != 1 {
		t.Errorf("first message = %+v, want assistant with 1 tool call", got[0])
	}
}

func TestBuildAllToolMessages(t *testing.T) {
	history := []llm.Message{
		toolMsg("call-1"),
		toolMsg("call-1"),
		toolMsg("call-1"),
		toolMsg("call-1"),
	}
	manager := &ContextManager{maxMessages: 3}

	got := manager.Build(history)

	if len(got) != len(history) {
		t.Fatalf("Build() returned %d messages, want %d", len(got), len(history))
	}

	for i, message := range got {
		if message.Role != "tool" {
			t.Errorf("message[%d] role = %q, want %q", i, message.Role, "tool")
		}
	}
}

func TestBuildStopsBackoffAtNonToolMessage(t *testing.T) {
	history := []llm.Message{
		userMsg("u"),
		{Role: "assistant", Content: "plain"},
		toolMsg("call-1"),
		toolMsg("call-1"),
	}
	manager := &ContextManager{maxMessages: 2}

	got := manager.Build(history)

	if len(got) != 3 {
		t.Fatalf("Build() returned %d messages, want 3", len(got))
	}

	if got[0].Role != "assistant" || len(got[0].ToolCalls) != 0 {
		t.Errorf("first message = %+v, want assistant without tool calls", got[0])
	}
}

func TestBuildTableDrivenSmallN(t *testing.T) {
	tests := []struct {
		name               string
		maxMessages        int
		history            []llm.Message
		wantLength         int
		wantFirstRole      string
		wantFirstToolCalls int
		wantLastToolCallID string
	}{
		{
			name:        "start lands on assistant",
			maxMessages: 2,
			history: []llm.Message{
				userMsg("u"),
				assistantWithToolCalls("call-1"),
				toolMsg("call-1"),
			},
			wantLength:         2,
			wantFirstRole:      "assistant",
			wantFirstToolCalls: 1,
			wantLastToolCallID: "call-1",
		},
		{
			name:        "start lands mid tool run",
			maxMessages: 1,
			history: []llm.Message{
				userMsg("u"),
				assistantWithToolCalls("call-1"),
				toolMsg("call-1"),
			},
			wantLength:         2,
			wantFirstRole:      "assistant",
			wantFirstToolCalls: 1,
			wantLastToolCallID: "call-1",
		},
		{
			name:        "multiple tool results kept together",
			maxMessages: 2,
			history: []llm.Message{
				userMsg("u"),
				assistantWithToolCalls("call-1", "call-2"),
				toolMsg("call-1"),
				toolMsg("call-2"),
			},
			wantLength:         3,
			wantFirstRole:      "assistant",
			wantFirstToolCalls: 2,
			wantLastToolCallID: "call-2",
		},
		{
			name:        "back-off stops at history start",
			maxMessages: 2,
			history: []llm.Message{
				assistantWithToolCalls("call-1"),
				toolMsg("call-1"),
				toolMsg("call-1"),
				toolMsg("call-1"),
			},
			wantLength:         4,
			wantFirstRole:      "assistant",
			wantFirstToolCalls: 1,
			wantLastToolCallID: "call-1",
		},
		{
			name:        "plain history trims without back-off",
			maxMessages: 3,
			history: []llm.Message{
				userMsg("m0"),
				userMsg("m1"),
				userMsg("m2"),
				userMsg("m3"),
				userMsg("m4"),
			},
			wantLength:    3,
			wantFirstRole: "user",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := &ContextManager{maxMessages: tt.maxMessages}

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

			if tt.wantLastToolCallID != "" {
				last := got[len(got)-1]
				if last.ToolCallID != tt.wantLastToolCallID {
					t.Errorf("last message ToolCallID = %q, want %q", last.ToolCallID, tt.wantLastToolCallID)
				}
			}
		})
	}
}
