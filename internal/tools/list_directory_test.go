package tools

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/alexnakagama/noryn/internal/project"
)

func TestListDirectoryTool_Name(t *testing.T) {
	tool := ListDirectoryTool{}

	got := tool.Name()

	if got != "list_directory" {
		t.Errorf("Name() = %q, want %q", got, "list_directory")
	}
}

func TestListDirectoryTool_Execute(t *testing.T) {
	tests := []struct {
		name    string
		args    string
		setup   func(t *testing.T, root string)
		want    string
		wantErr bool
	}{
		{
			name: "lists directory entries sorted by name",
			args: `{"path":"."}`,
			setup: func(t *testing.T, root string) {
				createDirEntry(t, root, "README.md", false)
				createDirEntry(t, root, "load.go", false)
				createDirEntry(t, root, "models", true)
				createDirEntry(t, root, "z_last.txt", false)
			},
			want: "README.md\nload.go\nmodels\nz_last.txt\n",
		},
		{
			name: "reports an empty directory as empty",
			args: `{"path":"empty"}`,
			setup: func(t *testing.T, root string) {
				if err := os.Mkdir(filepath.Join(root, "empty"), 0755); err != nil {
					t.Fatalf("os.Mkdir() error = %v", err)
				}
			},
			want: "",
		},
		{
			name:    "missing directory returns an error",
			args:    `{"path":"missing"}`,
			wantErr: true,
		},
		{
			name: "file path returns an error",
			args: `{"path":"load.go"}`,
			setup: func(t *testing.T, root string) {
				createDirEntry(t, root, "load.go", false)
			},
			wantErr: true,
		},
		{
			name:    "empty path returns an error",
			args:    `{"path":""}`,
			wantErr: true,
		},
		{
			name:    "malformed arguments return an error",
			args:    "not-json",
			wantErr: true,
		},
		{
			name:    "path outside project returns an error",
			args:    `{"path":"../outside"}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := t.TempDir()

			p := &project.Project{
				Root: root,
			}

			if tt.setup != nil {
				tt.setup(t, root)
			}

			tool := NewListDirectoryTool(p)

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

func createDirEntry(t *testing.T, dir, name string, isDir bool) {
	t.Helper()

	path := filepath.Join(dir, name)

	if isDir {
		if err := os.MkdirAll(path, 0755); err != nil {
			t.Fatalf("os.MkdirAll(%q) error = %v", path, err)
		}
		return
	}

	if err := os.WriteFile(path, nil, 0644); err != nil {
		t.Fatalf("os.WriteFile(%q) error = %v", path, err)
	}
}
