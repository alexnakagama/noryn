package input

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

func TestViewHasSinglePrompt(t *testing.T) {
	m := New()

	if strings.Contains(m.View(), "> >") {
		t.Errorf("View() renders a duplicated prompt: %q", m.View())
	}
}

func TestViewHasHintRow(t *testing.T) {
	m := New()

	if !strings.Contains(m.View(), "enter: send") {
		t.Errorf("View() missing the keybinding hint: %q", m.View())
	}

	if !strings.Contains(m.View(), "alt+enter: newline") {
		t.Errorf("View() missing the newline hint: %q", m.View())
	}
}

func TestViewRightAlignsHint(t *testing.T) {
	m := New()
	m.SetWidth(80)
	m.SetHintRight("gpt-4o")

	view := m.View()

	if !strings.Contains(view, "gpt-4o") {
		t.Errorf("View() missing the right-aligned hint: %q", view)
	}

	if strings.Index(view, "gpt-4o") <= strings.Index(view, "enter: send") {
		t.Errorf("View() right hint is not after the keybinding hint: %q", view)
	}
}

func TestViewWithoutHintRight(t *testing.T) {
	m := New()
	m.SetWidth(40)

	if strings.Contains(m.View(), "hintRight") {
		t.Errorf("View() contains unexpected text: %q", m.View())
	}
}

// countTransparentSpaces returns the number of space cells rendered while no
// background color was active.
func countTransparentSpaces(s string) int {
	transparent := 0
	activeBG := false

	i := 0
	for i < len(s) {
		if s[i] == '\x1b' {
			j := i + 1

			if j < len(s) && s[j] == '[' {
				j++
				start := j

				for j < len(s) && s[j] != 'm' {
					j++
				}

				for _, p := range strings.Split(s[start:j], ";") {
					switch {
					case p == "0" || p == "" || p == "49":
						activeBG = false
					case strings.HasPrefix(p, "48"):
						activeBG = true
					}
				}
			}

			i = j + 1
			continue
		}

		if s[i] == '\n' || s[i] == '\r' {
			activeBG = false
		} else if s[i] == ' ' && !activeBG {
			transparent++
		}

		i++
	}

	return transparent
}

func TestViewIsFullyOpaque(t *testing.T) {
	lipgloss.SetColorProfile(termenv.TrueColor)

	m := New()
	m.SetWidth(60)

	tests := map[string]func(){
		"empty placeholder": func() {},
		"short input": func() {
			m.textarea.SetValue("hello world")
		},
		"wrapped input": func() {
			m.textarea.SetValue(strings.Repeat("a very long line that wraps around plenty ", 4))
		},
		"empty after input": func() {
			m.textarea.SetValue("")
		},
	}

	for name, setup := range tests {
		setup()

		if n := countTransparentSpaces(m.View()); n != 0 {
			t.Errorf("%s: View() has %d transparent space cells", name, n)
		}
	}
}
