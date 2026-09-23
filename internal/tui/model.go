package tui

type Model struct {
	width  int
	height int

	chat   chat.Model
	input  input.Model
	status status.Model
}
