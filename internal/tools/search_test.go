package tools

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/alexnakagama/noryn/internal/project"
)

func TestSearchTool_Name(t *testing.T) {
	tool := NewSearchTool(nil)

	got := tool.Name()

	if got != "search" {
		t.Errorf("Name() = %q, want %q", got, "search")
	}
}

func TestSearchTool_Execute(t *testing.T) {
	tests := []struct {
		name    string
		args    string
		setup   func(t *testing.T, root string)
		want    string
		wantErr bool
	}{
		{
			name: "finds a query across files in walk order",
			args: `{"query":"TODO","path":"."}`,
			want: searchMatch("README.md", 3, "TODO: document the tool") +
				searchMatch("main.go", 4, "\t// TODO: greet") +
				searchMatch(filepath.Join("sub", "helper.go"), 3, "// TODO: implement helper"),
		},
		{
			name: "matches every line that contains the query",
			args: `{"query":"package","path":"."}`,
			want: searchMatch("main.go", 1, "package main") +
				searchMatch(filepath.Join("sub", "helper.go"), 1, "package sub"),
		},
		{
			name: "scopes the search to a subdirectory",
			args: `{"query":"package","path":"sub"}`,
			want: searchMatch(filepath.Join("sub", "helper.go"), 1, "package sub"),
		},
		{
			name: "does not report directories as matches",
			args: `{"query":"needle","path":"."}`,
			setup: func(t *testing.T, root string) {
				t.Helper()

				createSearchFile(t, root, "needle-dir/inner.txt", "a needle in the haystack\n")
			},
			want: searchMatch(filepath.Join("needle-dir", "inner.txt"), 1, "a needle in the haystack"),
		},
		{
			name: "query is case sensitive",
			args: `{"query":"todo","path":"."}`,
			want: "",
		},
		{
			name: "query with no matches returns an empty result",
			args: `{"query":"absent","path":"."}`,
			want: "",
		},
		{
			name:    "empty query returns an error",
			args:    `{"query":"","path":"."}`,
			wantErr: true,
		},
		{
			name:    "empty path returns an error",
			args:    `{"query":"TODO","path":""}`,
			wantErr: true,
		},
		{
			name:    "missing query returns an error",
			args:    `{"path":"."}`,
			wantErr: true,
		},
		{
			name:    "path outside project returns an error",
			args:    `{"query":"TODO","path":"../outside"}`,
			wantErr: true,
		},
		{
			name:    "malformed arguments return an error",
			args:    "not-json",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := t.TempDir()

			createSearchFixtures(t, root)

			if tt.setup != nil {
				tt.setup(t, root)
			}

			tool := NewSearchTool(&project.Project{Root: root})

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

func createSearchFixtures(t *testing.T, root string) {
	t.Helper()

	createSearchFile(t, root, "README.md", "# Noryn\n\nTODO: document the tool\n")
	createSearchFile(t, root, "main.go", "package main\n\nfunc main() {\n\t// TODO: greet\n}\n")
	createSearchFile(t, root, "sub/helper.go", "package sub\n\n// TODO: implement helper\n")
}

func createSearchFile(t *testing.T, root, relativePath, content string) {
	t.Helper()

	path := filepath.Join(root, relativePath)

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatalf("os.MkdirAll(%q) error = %v", filepath.Dir(path), err)
	}

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("os.WriteFile(%q) error = %v", path, err)
	}
}

func searchMatch(path string, lineNumber int, line string) string {
	return fmt.Sprintf("%s:%d:%s\n", path, lineNumber, line)
}
