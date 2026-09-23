package input

import (
	"strings"
	"testing"
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
