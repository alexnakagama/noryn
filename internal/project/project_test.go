package project

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestResolvePathRelativeInsideProject(t *testing.T) {
	root := t.TempDir()

	p := &Project{Root: root}

	got, err := p.ResolvePath("main.go")
	if err != nil {
		t.Fatalf("ResolvePath() error = %v", err)
	}

	want := filepath.Join(root, "main.go")

	if got != want {
		t.Errorf("ResolvePath() = %q, want %q", got, want)
	}
}

func TestResolvePathReturnsAbsolutePath(t *testing.T) {
	root := t.TempDir()

	p := &Project{Root: root}

	got, err := p.ResolvePath("main.go")
	if err != nil {
		t.Fatalf("ResolvePath() error = %v", err)
	}

	if !filepath.IsAbs(got) {
		t.Errorf("ResolvePath() returned relative path %q, want absolute", got)
	}
}

func TestResolvePathNestedSubdirectory(t *testing.T) {
	root := t.TempDir()

	nested := filepath.Join(root, "internal", "agent")

	if err := os.MkdirAll(nested, 0755); err != nil {
		t.Fatal(err)
	}

	p := &Project{Root: root}

	got, err := p.ResolvePath("internal/agent/agent.go")
	if err != nil {
		t.Fatalf("ResolvePath() error = %v", err)
	}

	want := filepath.Join(root, "internal", "agent", "agent.go")

	if got != want {
		t.Errorf("ResolvePath() = %q, want %q", got, want)
	}
}

func TestResolvePathEmptyPath(t *testing.T) {
	p := &Project{Root: t.TempDir()}

	got, err := p.ResolvePath("")
	if err == nil {
		t.Fatalf("ResolvePath() error = nil, got %q, want an error", got)
	}
}

func TestResolvePathEscapingProject(t *testing.T) {
	root := t.TempDir()

	p := &Project{Root: root}

	paths := []string{
		"..",
		"../secret.txt",
		"../../etc/passwd",
		"subdir/../../secret.txt",
	}

	for _, path := range paths {
		got, err := p.ResolvePath(path)
		if err == nil {
			t.Errorf("ResolvePath(%q) error = nil, got %q, want an error", path, got)
		}
	}
}

func TestResolvePathDotReturnsProjectRoot(t *testing.T) {
	root := t.TempDir()

	p := &Project{Root: root}

	got, err := p.ResolvePath(".")
	if err != nil {
		t.Fatalf("ResolvePath() error = %v", err)
	}

	if got != root {
		t.Errorf("ResolvePath() = %q, want %q", got, root)
	}
}

func TestResolvePathDotDotInsideProjectStaysInside(t *testing.T) {
	root := t.TempDir()

	p := &Project{Root: root}

	got, err := p.ResolvePath("subdir/..")
	if err != nil {
		t.Fatalf("ResolvePath() error = %v", err)
	}

	if got != root {
		t.Errorf("ResolvePath() = %q, want %q", got, root)
	}
}

func TestResolvePathInternalDotDotCleaned(t *testing.T) {
	root := t.TempDir()

	p := &Project{Root: root}

	got, err := p.ResolvePath("a/../b")
	if err != nil {
		t.Fatalf("ResolvePath() error = %v", err)
	}

	want := filepath.Join(root, "b")

	if got != want {
		t.Errorf("ResolvePath() = %q, want %q", got, want)
	}
}

func TestResolvePathLeadingDotSlash(t *testing.T) {
	root := t.TempDir()

	p := &Project{Root: root}

	got, err := p.ResolvePath("./main.go")
	if err != nil {
		t.Fatalf("ResolvePath() error = %v", err)
	}

	want := filepath.Join(root, "main.go")

	if got != want {
		t.Errorf("ResolvePath() = %q, want %q", got, want)
	}
}

func TestResolvePathRedundantSeparators(t *testing.T) {
	root := t.TempDir()

	p := &Project{Root: root}

	got, err := p.ResolvePath("a//b")
	if err != nil {
		t.Fatalf("ResolvePath() error = %v", err)
	}

	want := filepath.Join(root, "a", "b")

	if got != want {
		t.Errorf("ResolvePath() = %q, want %q", got, want)
	}
}

func TestResolvePathTrailingSlash(t *testing.T) {
	root := t.TempDir()

	p := &Project{Root: root}

	got, err := p.ResolvePath("main.go/")
	if err != nil {
		t.Fatalf("ResolvePath() error = %v", err)
	}

	want := filepath.Join(root, "main.go")

	if got != want {
		t.Errorf("ResolvePath() = %q, want %q", got, want)
	}
}

func TestResolvePathNonexistentPathSucceeds(t *testing.T) {
	root := t.TempDir()

	p := &Project{Root: root}

	got, err := p.ResolvePath("missing/file.go")
	if err != nil {
		t.Fatalf("ResolvePath() error = %v", err)
	}

	want := filepath.Join(root, "missing", "file.go")

	if got != want {
		t.Errorf("ResolvePath() = %q, want %q", got, want)
	}
}

func TestResolvePathDirectoryPathSucceeds(t *testing.T) {
	root := t.TempDir()

	if err := os.MkdirAll(filepath.Join(root, "internal"), 0755); err != nil {
		t.Fatal(err)
	}

	p := &Project{Root: root}

	got, err := p.ResolvePath("internal")
	if err != nil {
		t.Fatalf("ResolvePath() error = %v", err)
	}

	want := filepath.Join(root, "internal")

	if got != want {
		t.Errorf("ResolvePath() = %q, want %q", got, want)
	}
}

func TestResolvePathPathWithSpaces(t *testing.T) {
	root := t.TempDir()

	p := &Project{Root: root}

	got, err := p.ResolvePath("my file.txt")
	if err != nil {
		t.Fatalf("ResolvePath() error = %v", err)
	}

	want := filepath.Join(root, "my file.txt")

	if got != want {
		t.Errorf("ResolvePath() = %q, want %q", got, want)
	}
}

func TestResolvePathRootWithTrailingSeparator(t *testing.T) {
	root := t.TempDir()

	p := &Project{Root: root + string(filepath.Separator)}

	got, err := p.ResolvePath("main.go")
	if err != nil {
		t.Fatalf("ResolvePath() error = %v", err)
	}

	want := filepath.Join(root, "main.go")

	if got != want {
		t.Errorf("ResolvePath() = %q, want %q", got, want)
	}
}

func TestResolvePathAbsoluteInputTreatedAsRelative(t *testing.T) {
	root := t.TempDir()

	p := &Project{Root: root}

	got, err := p.ResolvePath("/etc/passwd")
	if err != nil {
		t.Fatalf("ResolvePath() error = %v", err)
	}

	want := filepath.Join(root, "etc", "passwd")

	if got != want {
		t.Errorf("ResolvePath() = %q, want %q", got, want)
	}
}

func TestResolvePathRootSlashInput(t *testing.T) {
	root := t.TempDir()

	p := &Project{Root: root}

	got, err := p.ResolvePath("/")
	if err != nil {
		t.Fatalf("ResolvePath() error = %v", err)
	}

	if got != root {
		t.Errorf("ResolvePath() = %q, want %q", got, root)
	}
}

func TestResolvePathAbsoluteInputEscapingStillRejected(t *testing.T) {
	root := t.TempDir()

	p := &Project{Root: root}

	sep := string(filepath.Separator)

	paths := []string{
		sep + ".." + sep + "x",
		sep + ".." + sep + ".." + sep + "etc" + sep + "passwd",
	}

	for _, path := range paths {
		got, err := p.ResolvePath(path)
		if err == nil {
			t.Errorf("ResolvePath(%q) error = nil, got %q, want an error", path, got)
		}
	}
}

func TestResolvePathRelativeRootReturnsError(t *testing.T) {
	p := &Project{Root: "relative/dir"}

	got, err := p.ResolvePath("main.go")
	if err == nil {
		t.Fatalf("ResolvePath() error = nil, got %q, want an error", got)
	}
}

func TestResolvePathEmptyRootReturnsError(t *testing.T) {
	p := &Project{Root: ""}

	got, err := p.ResolvePath("main.go")
	if err == nil {
		t.Fatalf("ResolvePath() error = nil, got %q, want an error", got)
	}
}

func TestResolvePathSymlinkInsideProjectSucceeds(t *testing.T) {
	root := t.TempDir()

	if err := os.Symlink("/etc", filepath.Join(root, "link")); err != nil {
		t.Skipf("symlink not supported: %v", err)
	}

	p := &Project{Root: root}

	got, err := p.ResolvePath("link")
	if err != nil {
		t.Fatalf("ResolvePath() error = %v", err)
	}

	want := filepath.Join(root, "link")

	if got != want {
		t.Errorf("ResolvePath() = %q, want %q", got, want)
	}
}

func TestResolvePathErrorMessages(t *testing.T) {
	root := t.TempDir()

	p := &Project{Root: root}

	_, err := p.ResolvePath("")
	if err == nil {
		t.Fatal("ResolvePath() error = nil, want an error")
	}

	if !strings.Contains(err.Error(), "path cannot be empty") {
		t.Errorf("ResolvePath() empty path error = %q, want mention of %q", err, "path cannot be empty")
	}

	_, err = p.ResolvePath("..")
	if err == nil {
		t.Fatal("ResolvePath() error = nil, want an error")
	}

	if !strings.Contains(err.Error(), "path is outside the project: ..") {
		t.Errorf("ResolvePath() outside error = %q, want mention of %q", err, "path is outside the project: ..")
	}
}
