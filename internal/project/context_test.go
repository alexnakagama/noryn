package project

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestBuildContext(t *testing.T) {
	root := t.TempDir()

	files := []string{
		"main.go",
		"internal/agent.go",
	}

	for _, file := range files {
		path := filepath.Join(root, file)

		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			t.Fatal(err)
		}

		if err := os.WriteFile(path, []byte("test"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	p := &Project{Root: root}

	ctx, err := p.BuildContext()
	if err != nil {
		t.Fatal(err)
	}

	expected := []string{
		"internal/agent.go",
		"main.go",
	}

	if !reflect.DeepEqual(ctx.Files, expected) {
		t.Fatalf("expected %v, got %v", expected, ctx.Files)
	}
}

func TestBuildContextIgnoresDirectories(t *testing.T) {
	root := t.TempDir()

	files := []string{
		"main.go",
		".git/config",
		"node_modules/package/index.js",
		"vendor/library/library.go",
		"internal/agent/agent.go",
	}

	for _, file := range files {
		path := filepath.Join(root, file)

		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			t.Fatal(err)
		}

		if err := os.WriteFile(path, []byte("test"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	p := &Project{Root: root}

	ctx, err := p.BuildContext()
	if err != nil {
		t.Fatal(err)
	}

	expected := []string{
		"internal/agent/agent.go",
		"main.go",
	}

	if !reflect.DeepEqual(ctx.Files, expected) {
		t.Fatalf("expected %v, got %v", expected, ctx.Files)
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

	expected := "Project root: /tmp/noryn\n" +
		"Files:\n" +
		"- main.go\n" +
		"- internal/agent/agent.go\n" +
		"- go.mod\n"

	got := ctx.String()

	if got != expected {
		t.Fatalf("expected %q, got %q", expected, got)
	}
}
