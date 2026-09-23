package input

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const (
	inputBackground      = "#16181F"
	cursorLineBackground = "#1A1D2E"
	inputBorder          = "#303640"
	inputBorderFocused   = "#7AA2F7"
	inputMuted           = "#7C8594"
	inputPrompt          = "#7AA2F7"
	inputText            = "#D8DEE9"
)

func (m Model) View() string {
	border := inputBorder
	if m.textarea.Focused() {
		border = inputBorderFocused
	}

	box := lipgloss.NewStyle().
		Background(lipgloss.Color(inputBackground)).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(border)).
		Padding(0, 1).
		Render(m.textarea.View())

	hintStyle := lipgloss.NewStyle().
		Background(lipgloss.Color(inputBackground)).
		Foreground(lipgloss.Color(inputMuted))

	hintText := "enter: send • alt+enter: newline"

	if m.hintRight != "" {
		pad := m.width - lipgloss.Width(hintText) - lipgloss.Width(m.hintRight)

		if pad < 1 {
			pad = 1
		}

		hintText += strings.Repeat(" ", pad) + m.hintRight
	}

	hint := hintStyle.Render(hintText)

	return box + "\n" + hint
}
