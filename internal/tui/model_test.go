package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestInitReturnsFocusCmd(t *testing.T) {
	m := New(nil, "", "", "")

	if m.Init() == nil {
		t.Error("Init() returned a nil command, want the input focus command")
	}
}

func TestUpdateWindowSizeMsgSetsDimensions(t *testing.T) {
	m := New(nil, "test-model", "", "")

	updated, _ := m.Update(tea.WindowSizeMsg{Width: 100, Height: 40})

	model, ok := updated.(Model)
	if !ok {
		t.Fatalf("Update() returned %T, want Model", updated)
	}

	if model.width != 100 {
		t.Errorf("width = %d, want 100", model.width)
	}

	if model.height != 40 {
		t.Errorf("height = %d, want 40", model.height)
	}
}
