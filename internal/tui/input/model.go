package input

import "github.com/charmbracelet/bubbles/textarea"

type Model struct {
	textarea textarea.Model
}

func New() Model {
	input := textarea.New()

	input.Placeholder = "Ask Noryn anything..."
	input.Prompt = "> "
	input.CharLimit = 0

	return Model{
		textarea: input,
	}
}
