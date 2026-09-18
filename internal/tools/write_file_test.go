package tools

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/alexnakagama/noryn/internal/project"
)

func TestWriteFileTool_Name(t *testing.T) {
	tool := WriteFileTool{}

	got := tool.Name()

	if got != "write_file" {
		t.Errorf("Name() = %q, want %q", got, "write_file")
	}
}

func TestWriteFileTool_Execute(t *testing.T) {
	tests := []struct {
		name        string
		args        string
		seed        func(t *testing.T, path string)
		wantContent string
		wantErr     bool
	}{
		{
			name:        "creates a new file with the given content",
			args:        `{"path":"output.txt","content":"hello\nworld"}`,
			wantContent: "hello\nworld",
		},
		{
			name:        "creates an empty file",
			args:        `{"path":"output.txt","content":""}`,
			wantContent: "",
		},
		{
			name: "overwrites an existing file",
			args: `{"path":"output.txt","content":"new content"}`,
			seed: func(t *testing.T, path string) {
				t.Helper()
				createDestinationFile(t, path, "old content")
			},
			wantContent: "new content",
		},
		{
			name: "directory target returns an error",
			args: `{"path":"output.txt","content":"x"}`,
			seed: func(t *testing.T, path string) {
				t.Helper()

				if err := os.Mkdir(path, 0755); err != nil {
					t.Fatalf("os.Mkdir(%q) error = %v", path, err)
				}
			},
			wantErr: true,
		},
		{
			name:    "empty path returns an error",
			args:    `{"path":"","content":"x"}`,
			wantErr: true,
		},
		{
			name:    "malformed arguments return an error",
			args:    "not-json",
			wantErr: true,
		},
		{
			name:    "path outside project returns an error",
			args:    `{"path":"../outside.txt","content":"should fail"}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := t.TempDir()

			p := &project.Project{
				Root: root,
			}

			path := filepath.Join(root, "output.txt")

			if tt.seed != nil {
				tt.seed(t, path)
			}

			tool := NewWriteFileTool(p)

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

			if got != "file written successfully" {
				t.Errorf(
					"Execute() message = %q, want %q",
					got,
					"file written successfully",
				)
			}

			content, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("os.ReadFile(%q) error = %v", path, err)
			}

			if string(content) != tt.wantContent {
				t.Errorf(
					"file content = %q, want %q",
					string(content),
					tt.wantContent,
				)
			}
		})
	}
}

func createDestinationFile(t *testing.T, path, content string) {
	t.Helper()

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("os.WriteFile(%q) error = %v", path, err)
	}
}
