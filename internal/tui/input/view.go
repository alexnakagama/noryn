package input

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const (
	inputBackground     = "#16181F"
	inputSideBackground = "#1A1B26"
	inputBorder         = "#303640"
	inputBorderFocused  = "#7AA2F7"
	inputMuted          = "#7C8594"
	inputPrompt         = "#7AA2F7"
	inputText           = "#D8DEE9"
	inputSideMargin     = 8
	inputHeight         = 3
)

func (m Model) View() string {
	boxWidth := m.width - 2*inputSideMargin

	if boxWidth < 10 {
		boxWidth = m.width
	}

	innerWidth := boxWidth - 4

	if innerWidth < 1 {
		innerWidth = 1
	}

	leftPad := lipgloss.NewStyle().
		Background(lipgloss.Color(inputSideBackground)).
		Render(strings.Repeat(" ", (m.width-boxWidth)/2))

	borderColor := inputBorder
	if m.textarea.Focused() {
		borderColor = inputBorderFocused
	}

	panel := lipgloss.NewStyle().
		Background(lipgloss.Color(inputBackground))

	border := lipgloss.NewStyle().
		Foreground(lipgloss.Color(borderColor)).
		Background(lipgloss.Color(inputBackground))

	lines := strings.Split(recolorTransparentSpaces(m.textarea.View(), panel), "\n")
	if len(lines) > inputHeight {
		lines = lines[:inputHeight]
	}
	for len(lines) < inputHeight {
		lines = append(lines, "")
	}

	var b strings.Builder

	b.WriteString(leftPad)
	b.WriteString(border.Render("╭" + strings.Repeat("─", innerWidth+2) + "╮"))
	b.WriteString("\n")

	for _, line := range lines {
		visible := lipgloss.Width(line)

		if visible < innerWidth {
			line += panel.Render(strings.Repeat(" ", innerWidth-visible))
		}

		b.WriteString(leftPad)
		b.WriteString(border.Render("│"))
		b.WriteString(panel.Render(" "))
		b.WriteString(line)
		b.WriteString(panel.Render(" "))
		b.WriteString(border.Render("│"))
		b.WriteString("\n")
	}

	b.WriteString(leftPad)
	b.WriteString(border.Render("╰" + strings.Repeat("─", innerWidth+2) + "╯"))

	hintStyle := lipgloss.NewStyle().
		Background(lipgloss.Color(inputBackground)).
		Foreground(lipgloss.Color(inputMuted))

	hintText := "enter: send • alt+enter: newline"

	if m.hintRight != "" {
		pad := boxWidth - lipgloss.Width(hintText) - lipgloss.Width(m.hintRight)

		if pad < 1 {
			pad = 1
		}

		hintText += strings.Repeat(" ", pad) + m.hintRight
	}

	return b.String() + "\n" + leftPad + hintStyle.Render(hintText)
}

// recolorTransparentSpaces re-emits any space cell that was drawn without an
// active background so the input area is fully opaque even though the
// underlying textarea may pad its lines with unstyled spaces.
func recolorTransparentSpaces(s string, panel lipgloss.Style) string {
	var out strings.Builder
	out.Grow(len(s))

	activeBG := false
	i := 0

	for i < len(s) {
		if s[i] == '\x1b' {
			if i+1 < len(s) && s[i+1] == '[' {
				j := i + 2

				for j < len(s) && s[j] != 'm' {
					j++
				}

				if j < len(s) {
					out.WriteString(s[i : j+1])

					for _, p := range strings.Split(s[i+2:j], ";") {
						switch {
						case p == "0" || p == "" || p == "49":
							activeBG = false
						case strings.HasPrefix(p, "48"):
							activeBG = true
						}
					}

					i = j + 1
					continue
				}
			}

			out.WriteByte(s[i])
			i++
			continue
		}

		switch {
		case s[i] == '\n':
			out.WriteByte(s[i])
			activeBG = false
		case s[i] == ' ' && !activeBG:
			out.WriteString(panel.Render(" "))
		default:
			out.WriteByte(s[i])
		}

		i++
	}

	return out.String()
}
