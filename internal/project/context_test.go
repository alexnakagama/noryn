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

func TestBuildContextIgnoresLargeFiles(t *testing.T) {
	root := t.TempDir()

	smallFile := filepath.Join(root, "main.go")
	largeFile := filepath.Join(root, "large.txt")

	if err := os.WriteFile(smallFile, []byte("package main"), 0644); err != nil {
		t.Fatal(err)
	}

	largeContent := make([]byte, 1_000_001)

	if err := os.WriteFile(largeFile, largeContent, 0644); err != nil {
		t.Fatal(err)
	}

	p := &Project{Root: root}

	ctx, err := p.BuildContext()
	if err != nil {
		t.Fatal(err)
	}

	expected := []string{
		"main.go",
	}

	if !reflect.DeepEqual(ctx.Files, expected) {
		t.Fatalf("expected %v, got %v", expected, ctx.Files)
	}
}

func TestIsBinaryFile(t *testing.T) {
	root := t.TempDir()

	textPath := filepath.Join(root, "main.go")
	binaryPath := filepath.Join(root, "image.png")

	if err := os.WriteFile(textPath, []byte("package main\n"), 0644); err != nil {
		t.Fatal(err)
	}

	binaryContent := []byte{
		0x89, 0x50, 0x4E, 0x47,
		0x0D, 0x0A, 0x1A, 0x0A,
	}

	if err := os.WriteFile(binaryPath, binaryContent, 0644); err != nil {
		t.Fatal(err)
	}

	isBinary, err := isBinaryFile(textPath)
	if err != nil {
		t.Fatal(err)
	}

	if isBinary {
		t.Fatal("expected text file to not be binary")
	}

	isBinary, err = isBinaryFile(binaryPath)
	if err != nil {
		t.Fatal(err)
	}

	if !isBinary {
		t.Fatal("expected binary file to be binary")
	}
}
