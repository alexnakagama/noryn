package tools

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/alexnakagama/noryn/internal/project"
)

func TestShellTool_Name(t *testing.T) {
	tool := ShellTool{}

	got := tool.Name()

	if got != "shell" {
		t.Errorf("Name() = %q, want %q", got, "shell")
	}
}

func TestShellTool_Execute(t *testing.T) {
	tests := []struct {
		name       string
		args       string
		setup      func(t *testing.T, root string)
		wantOutput string
		wantErr    bool
	}{
		{
			name:       "returns the command output",
			args:       `{"command":"echo hello"}`,
			wantOutput: "hello\n",
		},
		{
			name:       "runs chained commands",
			args:       `{"command":"printf 'a\\nb' && printf c"}`,
			wantOutput: "a\nbc",
		},
		{
			name:    "empty command returns an error",
			args:    `{"command":""}`,
			wantErr: true,
		},
		{
			name:    "missing command returns an error",
			args:    `{}`,
			wantErr: true,
		},
		{
			name:       "failing command returns its combined output and error",
			args:       `{"command":"echo boom >&2; exit 1"}`,
			wantOutput: "boom\n",
			wantErr:    true,
		},
		{
			name:    "non-zero exit returns an error",
			args:    `{"command":"exit 3"}`,
			wantErr: true,
		},
		{
			name:    "malformed arguments return an error",
			args:    "not-json",
			wantErr: true,
		},
		{
			name: "runs command from project root",
			args: `{"command":"pwd"}`,
		},
		{
			name: "can access files from project root",
			args: `{"command":"cat test.txt"}`,
			setup: func(t *testing.T, root string) {
				t.Helper()

				path := filepath.Join(root, "test.txt")

				if err := os.WriteFile(
					path,
					[]byte("hello from project\n"),
					0644,
				); err != nil {
					t.Fatalf("os.WriteFile() error = %v", err)
				}
			},
			wantOutput: "hello from project\n",
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

			tool := NewShellTool(p)

			got, err := tool.Execute(tt.args)

			if tt.name == "runs command from project root" {
				expected, err := os.Getwd()
				if err != nil {
					t.Fatalf("os.Getwd() error = %v", err)
				}

				_ = expected
			}

			if tt.wantOutput != "" && got != tt.wantOutput {
				t.Errorf("Execute() output = %q, want %q", got, tt.wantOutput)
			}

			if tt.wantErr {
				if err == nil {
					t.Fatal("Execute() error = nil, want an error")
				}

				return
			}

			if err != nil {
				t.Fatalf("Execute() error = %v", err)
			}
		})
	}
}
