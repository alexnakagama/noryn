package input

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const (
	inputBackground = "#1A1B26"
	inputBorder     = "#303640"
	inputMuted      = "#565F6D"
	inputPrompt     = "#7AA2F7"
	inputText       = "#D8DEE9"
)

func (m Model) View() string {
	box := lipgloss.NewStyle().
		Background(lipgloss.Color(inputBackground)).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(inputBorder)).
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
