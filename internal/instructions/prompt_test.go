package instructions

import "testing"

func TestBuildPrompt(t *testing.T) {
	got := BuildPrompt(
		"You are noryn",
		"Be careful with files.",
		"Project context here.",
		"Read main.go",
	)

	expected := "You are noryn\n\n" +
		"Be careful with files.\n\n" +
		"Project context here.\n\n" +
		"Read main.go"

	if got != expected {
		t.Errorf("BuildPrompt() = %q, want %q", got, expected)
	}
}

func TestBuildPromptWithProjectContext(t *testing.T) {
	systemPrompt := "You are noryn"
	projectInstructions := "Be careful with files."
	projectContext := "Project root: /tmp/noryn\nFiles:\n- main.go\n"
	userPrompt := "Read main.go"

	got := BuildPrompt(
		systemPrompt,
		projectInstructions,
		projectContext,
		userPrompt,
	)

	expected := "You are noryn\n\n" +
		"Be careful with files.\n\n" +
		"Project root: /tmp/noryn\n" +
		"Files:\n" +
		"- main.go\n\n" +
		"Read main.go"

	if got != expected {
		t.Errorf(
			"BuildPrompt() = %q, want %q",
			got,
			expected,
		)
	}
}
