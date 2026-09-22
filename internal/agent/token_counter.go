package agent

import "github.com/alexnakagama/noryn/internal/llm"

type TokenCounter interface {
	Count(message llm.Message) int
}
