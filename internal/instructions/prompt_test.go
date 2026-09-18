package instructions

import "testing"

func TestBuildPrompt(t *testing.T) {
	projectInstructions := "Follow the project instructions."
	userPrompt := "Run the tests."

	want := "Follow the project instructions.\n\nRun the tests."

	got := BuildPrompt(projectInstructions, userPrompt)

	if got != want {
		t.Fatalf("BuildPrompt() = %q, want %q", got, want)
	}
}
