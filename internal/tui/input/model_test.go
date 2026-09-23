package input

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestUpdateEnterSubmits(t *testing.T) {
	m := New()
	m.textarea.SetValue("hello")

	m, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Fatal("Update() returned a nil command, want a submit command")
	}

	msg := cmd()

	submit, ok := msg.(SubmitMessage)
	if !ok {
		t.Fatalf("command returned %T, want SubmitMessage", msg)
	}

	if submit.Content != "hello" {
		t.Errorf("SubmitMessage.Content = %q, want %q", submit.Content, "hello")
	}

	if m.Value() != "" {
		t.Errorf("Value() = %q, want empty after submit", m.Value())
	}
}

func TestUpdateEnterSubmitsMultiline(t *testing.T) {
	m := New()
	m.textarea.SetValue("ab\ncd")

	m, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Fatal("Update() returned a nil command, want a submit command")
	}

	msg := cmd()

	submit, ok := msg.(SubmitMessage)
	if !ok {
		t.Fatalf("command returned %T, want SubmitMessage", msg)
	}

	if submit.Content != "ab\ncd" {
		t.Errorf("SubmitMessage.Content = %q, want %q", submit.Content, "ab\ncd")
	}
}

func TestUpdateEnterEmptyDoesNotSubmit(t *testing.T) {
	m := New()

	m, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd != nil {
		t.Fatalf("Update() returned a non-nil command, want nil")
	}
}

func TestUpdateAltEnterInsertsNewline(t *testing.T) {
	m := New()
	m.textarea.SetValue("ab")

	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter, Alt: true})

	if m.Value() != "ab\n" {
		t.Errorf("Value() = %q, want %q", m.Value(), "ab\n")
	}
}

func TestUpdateCtrlJInsertsNewline(t *testing.T) {
	m := New()
	m.textarea.SetValue("ab")

	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyCtrlJ})

	if m.Value() != "ab\n" {
		t.Errorf("Value() = %q, want %q", m.Value(), "ab\n")
	}
}

func TestSetWidthStoresFullWidth(t *testing.T) {
	m := New()

	m.SetWidth(80)

	if m.width != 80 {
		t.Errorf("width = %d, want 80", m.width)
	}
}

func TestSetWidthClampsTinyWidth(t *testing.T) {
	m := New()

	m.SetWidth(2)

	if m.width != 2 {
		t.Errorf("width = %d, want 2", m.width)
	}
}

func TestSetHintRight(t *testing.T) {
	m := New()

	m.SetHintRight("gpt-4o")

	if m.hintRight != "gpt-4o" {
		t.Errorf("hintRight = %q, want %q", m.hintRight, "gpt-4o")
	}
}
