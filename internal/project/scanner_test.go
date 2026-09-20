package project

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindRootWhenGoModInStartDirectory(t *testing.T) {
	root := t.TempDir()

	if err := os.WriteFile(filepath.Join(root, "go.mod"), []byte("module example\n"), 0644); err != nil {
		t.Fatal(err)
	}

	got, err := FindRoot(root)
	if err != nil {
		t.Fatalf("FindRoot() error = %v", err)
	}

	if got != root {
		t.Errorf("FindRoot() = %q, want %q", got, root)
	}
}

func TestFindRootFromNestedDirectory(t *testing.T) {
	root := t.TempDir()

	if err := os.WriteFile(filepath.Join(root, "go.mod"), []byte("module example\n"), 0644); err != nil {
		t.Fatal(err)
	}

	nested := filepath.Join(root, "internal", "agent")

	if err := os.MkdirAll(nested, 0755); err != nil {
		t.Fatal(err)
	}

	got, err := FindRoot(nested)
	if err != nil {
		t.Fatalf("FindRoot() error = %v", err)
	}

	if got != root {
		t.Errorf("FindRoot() = %q, want %q", got, root)
	}
}

func TestFindRootReturnsAbsolutePath(t *testing.T) {
	root := t.TempDir()

	if err := os.WriteFile(filepath.Join(root, "go.mod"), []byte("module example\n"), 0644); err != nil {
		t.Fatal(err)
	}

	current, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	rel, err := filepath.Rel(current, root)
	if err != nil {
		t.Fatal(err)
	}

	got, err := FindRoot(rel)
	if err != nil {
		t.Fatalf("FindRoot() error = %v", err)
	}

	if filepath.IsAbs(got) != true {
		t.Errorf("FindRoot() returned relative path %q, want absolute", got)
	}
}

func TestFindRootErrorWhenGoModNotFound(t *testing.T) {
	dir := t.TempDir()

	got, err := FindRoot(dir)
	if err == nil {
		t.Fatalf("FindRoot() error = nil, got root %q, want an error", got)
	}
}

func TestDiscoverReturnsProject(t *testing.T) {
	root := t.TempDir()

	if err := os.WriteFile(filepath.Join(root, "go.mod"), []byte("module example\n"), 0644); err != nil {
		t.Fatal(err)
	}

	nested := filepath.Join(root, "cmd", "app")

	if err := os.MkdirAll(nested, 0755); err != nil {
		t.Fatal(err)
	}

	p, err := Discover(nested)
	if err != nil {
		t.Fatalf("Discover() error = %v", err)
	}

	if p == nil {
		t.Fatal("Discover() returned nil project")
	}

	if p.Root != root {
		t.Errorf("Discover() root = %q, want %q", p.Root, root)
	}
}

func TestDiscoverErrorWhenGoModNotFound(t *testing.T) {
	dir := t.TempDir()

	p, err := Discover(dir)
	if err == nil {
		t.Fatalf("Discover() error = nil, got project %+v, want an error", p)
	}

	if p != nil {
		t.Errorf("Discover() returned non-nil project %+v with error, want nil", p)
	}
}
