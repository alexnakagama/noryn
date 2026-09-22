package agent

import (
	"strings"
	"testing"

	"github.com/alexnakagama/noryn/internal/llm"
)

var _ TokenCounter = EstimateTokenCounter{}
var _ TokenCounter = (*EstimateTokenCounter)(nil)

func TestEstimateTokenCounterCount(t *testing.T) {
	counter := EstimateTokenCounter{}

	tests := []struct {
		name    string
		message llm.Message
		want    int
	}{
		{name: "zero value message", message: llm.Message{}, want: 0},
		{name: "empty content", message: llm.Message{Content: ""}, want: 0},
		{name: "one byte", message: llm.Message{Content: "a"}, want: 0},
		{name: "three bytes", message: llm.Message{Content: "abc"}, want: 0},
		{name: "four bytes", message: llm.Message{Content: "abcd"}, want: 1},
		{name: "eight bytes", message: llm.Message{Content: "abcdefgh"}, want: 2},
		{name: "five bytes", message: llm.Message{Content: "abcde"}, want: 1},
		{name: "seven bytes", message: llm.Message{Content: "abcdefg"}, want: 1},
		{name: "twelve bytes", message: llm.Message{Content: "abcdefghijkl"}, want: 3},
		{name: "four spaces", message: llm.Message{Content: "    "}, want: 1},
		{name: "four newlines", message: llm.Message{Content: "\n\n\n\n"}, want: 1},
		{name: "four tabs", message: llm.Message{Content: "\t\t\t\t"}, want: 1},
		{name: "two spaces", message: llm.Message{Content: "  "}, want: 0},
		{name: "accented char", message: llm.Message{Content: "é"}, want: 0},
		{name: "japanese", message: llm.Message{Content: "日本語"}, want: 2},
		{name: "emoji", message: llm.Message{Content: "😀"}, want: 1},
		{name: "mixed ascii and unicode", message: llm.Message{Content: "héllo wörld"}, want: 3},
		{name: "role ignored", message: llm.Message{Role: "assistant", Content: "abcd"}, want: 1},
		{name: "tool call id ignored", message: llm.Message{Role: "tool", Content: "abcd", ToolCallID: "call-1"}, want: 1},
		{
			name: "tool calls with empty content",
			message: llm.Message{
				Role:      "assistant",
				ToolCalls: []llm.ToolCall{{ID: "call-1", Name: "fake", Arguments: "very long arguments"}},
			},
			want: 0,
		},
		{
			name:    "tool calls with content",
			message: llm.Message{Role: "assistant", Content: "abcd", ToolCalls: []llm.ToolCall{{ID: "call-1"}}},
			want:    1,
		},
		{name: "one million bytes", message: llm.Message{Content: strings.Repeat("a", 1_000_000)}, want: 250_000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := counter.Count(tt.message); got != tt.want {
				t.Errorf("Count(%+v) = %d, want %d", tt.message, got, tt.want)
			}
		})
	}
}

func TestEstimateTokenCounterCountBoundary(t *testing.T) {
	counter := EstimateTokenCounter{}

	tests := []struct {
		name    string
		content string
		want    int
	}{
		{name: "exactly maxContextTokens", content: strings.Repeat("a", maxContextTokens*4), want: maxContextTokens},
		{name: "one byte under next token", content: strings.Repeat("a", maxContextTokens*4+3), want: maxContextTokens},
		{name: "one token over", content: strings.Repeat("a", maxContextTokens*4+4), want: maxContextTokens + 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := counter.Count(llm.Message{Content: tt.content}); got != tt.want {
				t.Errorf("Count() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestEstimateTokenCounterIsStateless(t *testing.T) {
	counter := EstimateTokenCounter{}
	message := llm.Message{Content: "abcdefgh"}
	first := counter.Count(message)

	for i := 0; i < 3; i++ {
		if got := counter.Count(message); got != first {
			t.Fatalf("Count() not deterministic: got %d, want %d", got, first)
		}
	}
}
