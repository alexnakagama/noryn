package instructions

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alexnakagama/noryn/internal/project"
)

func TestBuilder_Build(t *testing.T) {
	tests := []struct {
		name            string
		setup           func(t *testing.T, root string)
		want            string
		wantErr         bool
		wantErrContains string
	}{
		{
			name: "returns the AGENTS.md content",
			setup: func(t *testing.T, root string) {
				t.Helper()

				writeAGENTS(t, root, "This project is an AI coding agent.\n")
			},
			want: "This project is an AI coding agent.\n",
		},
		{
			name: "preserves multiline content",
			setup: func(t *testing.T, root string) {
				t.Helper()

				writeAGENTS(t, root, "line one\nline two\nline three\n")
			},
			want: "line one\nline two\nline three\n",
		},
		{
			name: "returns an empty string for an empty file",
			setup: func(t *testing.T, root string) {
				t.Helper()

				writeAGENTS(t, root, "")
			},
			want: "",
		},
		{
			name: "returns an error when AGENTS.md is missing",
			setup: func(t *testing.T, root string) {
				t.Helper()
			},
			wantErr:         true,
			wantErrContains: "no such file",
		},
		{
			name: "returns an error when the project does not exist",
			setup: func(t *testing.T, root string) {
				t.Helper()

				if err := os.RemoveAll(root); err != nil {
					t.Fatalf("os.RemoveAll() error = %v", err)
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := t.TempDir()

			if tt.setup != nil {
				tt.setup(t, root)
			}

			builder := NewBuilder(&project.Project{Root: root})

			got, err := builder.Build()

			if tt.wantErr {
				if err == nil {
					t.Fatal("Build() error = nil, want an error")
				}

				if tt.wantErrContains != "" && !strings.Contains(err.Error(), tt.wantErrContains) {
					t.Errorf("Build() error = %v, want it to contain %q", err, tt.wantErrContains)
				}
				return
			}

			if err != nil {
				t.Fatalf("Build() error = %v", err)
			}

			if got != tt.want {
				t.Errorf("Build() = %q, want %q", got, tt.want)
			}
		})
	}
}

func writeAGENTS(t *testing.T, root, content string) {
	t.Helper()

	path := filepath.Join(root, "AGENTS.md")

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("os.WriteFile(%q) error = %v", path, err)
	}
}
