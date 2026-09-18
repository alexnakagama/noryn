package tools

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/alexnakagama/noryn/internal/project"
)

func TestReadFileTool_Name(t *testing.T) {
	tool := ReadFileTool{}

	got := tool.Name()

	if got != "read_file" {
		t.Errorf("Name() = %q, want %q", got, "read_file")
	}
}

func TestReadFileTool_Execute(t *testing.T) {
	root := t.TempDir()

	p := &project.Project{
		Root: root,
	}

	existing := filepath.Join(root, "note.txt")
	createSourceFile(t, existing, "first line\nsecond line\n")

	empty := filepath.Join(root, "empty.txt")
	createSourceFile(t, empty, "")

	tests := []struct {
		name    string
		args    string
		want    string
		wantErr bool
	}{
		{
			name: "returns the file contents",
			args: `{"path":"note.txt"}`,
			want: "first line\nsecond line\n",
		},
		{
			name: "returns an empty string for an empty file",
			args: `{"path":"empty.txt"}`,
			want: "",
		},
		{
			name:    "missing file returns an error",
			args:    `{"path":"missing.txt"}`,
			wantErr: true,
		},
		{
			name:    "directory path returns an error",
			args:    `{"path":"."}`,
			wantErr: true,
		},
		{
			name:    "empty path returns an error",
			args:    `{"path":""}`,
			wantErr: true,
		},
		{
			name:    "missing path argument returns an error",
			args:    `{"other":"value"}`,
			wantErr: true,
		},
		{
			name:    "malformed arguments return an error",
			args:    "not-json",
			wantErr: true,
		},
		{
			name:    "path outside project returns an error",
			args:    `{"path":"../outside.txt"}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tool := NewReadFileTool(p)

			got, err := tool.Execute(tt.args)

			if tt.wantErr {
				if err == nil {
					t.Fatal("Execute() error = nil, want an error")
				}
				return
			}

			if err != nil {
				t.Fatalf("Execute() error = %v", err)
			}

			if got != tt.want {
				t.Errorf("Execute() = %q, want %q", got, tt.want)
			}
		})
	}
}

func createSourceFile(t *testing.T, path, content string) {
	t.Helper()

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("os.WriteFile(%q) error = %v", path, err)
	}
}
