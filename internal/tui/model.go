package tui

import (
	"github.com/alexnakagama/noryn/internal/tui/chat"
	"github.com/alexnakagama/noryn/internal/tui/input"
	"github.com/alexnakagama/noryn/internal/tui/status"
)

type Model struct {
	width  int
	height int

	chat   chat.Model
	input  input.Model
	status status.Model
}
