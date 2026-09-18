package project

import (
	"os"
	"path/filepath"
	"testing"
)

func TestBuildContext(t *testing.T) {
	root := t.TempDir()

	err := os.WriteFile(
		filepath.Join(root, "main.go"),
		[]byte("package main"),
		0644,
	)
	if err != nil {
		t.Fatal(err)
	}

	err = os.Mkdir(filepath.Join(root, "internal"), 0755)
	if err != nil {
		t.Fatal(err)
	}

	err = os.WriteFile(
		filepath.Join(root, "internal", "agent.go"),
		[]byte("package agent"),
		0644,
	)
	if err != nil {
		t.Fatal(err)
	}

	project := &Project{
		Root: root,
	}

	ctx, err := project.BuildContext()
	if err != nil {
		t.Fatalf("BuildContext() returned error: %v", err)
	}

	if ctx.Root != root {
		t.Errorf(
			"Context.Root = %q, want %q",
			ctx.Root,
			root,
		)
	}

	if len(ctx.Files) != 2 {
		t.Fatalf(
			"Context.Files has %d files, want %d",
			len(ctx.Files),
			2,
		)
	}

	expectedFiles := map[string]bool{
		"main.go":           true,
		"internal/agent.go": true,
	}

	for _, file := range ctx.Files {
		if !expectedFiles[file] {
			t.Errorf("unexpected file in context: %q", file)
		}
	}
}

func TestContextString(t *testing.T) {
	ctx := Context{
		Root: "/tmp/noryn",
		Files: []string{
			"main.go",
			"internal/agent/agent.go",
			"go.mod",
		},
	}

	got := ctx.String()

	expected := "Project root: /tmp/noryn\n" +
		"Files:\n" +
		"- main.go\n" +
		"- internal/agent/agent.go\n" +
		"- go.mod\n"

	if got != expected {
		t.Errorf(
			"Context.String() = %q, want %q",
			got,
			expected,
		)
	}
}
