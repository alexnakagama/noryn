package tools

import (
	"fmt"
	"os/exec"
	"strings"
	"testing"

	"github.com/alexnakagama/noryn/internal/project"
)

func TestGitLogTool_Name(t *testing.T) {
	tool := NewGitLogTool(nil)

	got := tool.Name()

	if got != "git_log" {
		t.Errorf("Name() = %q, want %q", got, "git_log")
	}
}

func TestGitLogTool_Execute(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git is not installed")
	}

	tests := []struct {
		name            string
		arguments       string
		setup           func(t *testing.T, root string)
		wantContains    []string
		wantNotContains []string
		wantOrder       []string
		wantErr         bool
		wantErrContains string
	}{
		{
			name:      "rejects empty arguments",
			arguments: "",
			setup: func(t *testing.T, root string) {
				t.Helper()

				initGitRepo(t, root)
			},
			wantErr: true,
		},
		{
			name:            "rejects a zero limit",
			arguments:       `{"limit":0}`,
			wantErr:         true,
			wantErrContains: "limit cannot be 0 or less than 0",
		},
		{
			name:            "rejects a negative limit",
			arguments:       `{"limit":-1}`,
			wantErr:         true,
			wantErrContains: "limit cannot be 0 or less than 0",
		},
		{
			name:            "rejects a limit greater than fifty",
			arguments:       `{"limit":51}`,
			wantErr:         true,
			wantErrContains: "limit cannot be greater than 50",
		},
		{
			name:      "lists commits in reverse chronological order",
			arguments: `{"limit":10}`,
			setup: func(t *testing.T, root string) {
				t.Helper()

				initGitRepo(t, root)
				commitChange(t, root, "first commit")
				commitChange(t, root, "second commit")
				commitChange(t, root, "third commit")
			},
			wantOrder: []string{"third commit", "second commit", "first commit"},
		},
		{
			name:      "respects the limit",
			arguments: `{"limit":2}`,
			setup: func(t *testing.T, root string) {
				t.Helper()

				initGitRepo(t, root)
				commitChange(t, root, "first commit")
				commitChange(t, root, "second commit")
				commitChange(t, root, "third commit")
			},
			wantContains:    []string{"third commit", "second commit"},
			wantNotContains: []string{"first commit"},
		},
		{
			name:      "returns the most recent commit for a limit of one",
			arguments: `{"limit":1}`,
			setup: func(t *testing.T, root string) {
				t.Helper()

				initGitRepo(t, root)
				commitChange(t, root, "first commit")
				commitChange(t, root, "second commit")
			},
			wantContains:    []string{"second commit"},
			wantNotContains: []string{"first commit"},
		},
		{
			name:      "accepts the maximum allowed limit",
			arguments: `{"limit":50}`,
			setup: func(t *testing.T, root string) {
				t.Helper()

				initGitRepo(t, root)
				commitChange(t, root, "first commit")
				commitChange(t, root, "second commit")
			},
			wantContains: []string{"second commit", "first commit"},
		},
		{
			name:      "reports an empty repository",
			arguments: `{"limit":10}`,
			setup: func(t *testing.T, root string) {
				t.Helper()

				initGitRepo(t, root)
			},
			wantErr:         true,
			wantErrContains: "does not have any commits",
		},
		{
			name:            "not a git repository returns an error",
			arguments:       `{"limit":10}`,
			wantErr:         true,
			wantErrContains: "not a git repository",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := t.TempDir()

			if tt.setup != nil {
				tt.setup(t, root)
			}

			tool := NewGitLogTool(&project.Project{Root: root})

			got, err := tool.Execute(tt.arguments)

			if tt.wantErr {
				if err == nil {
					t.Fatal("Execute() error = nil, want an error")
				}

				if tt.wantErrContains != "" {
					combined := strings.ToLower(err.Error() + "\n" + got)
					want := strings.ToLower(tt.wantErrContains)
					if !strings.Contains(combined, want) {
						t.Errorf("Execute() error/output = %q, want it to contain %q", combined, tt.wantErrContains)
					}
				}
				return
			}

			if err != nil {
				t.Fatalf("Execute() error = %v", err)
			}

			for _, substr := range tt.wantContains {
				if !strings.Contains(got, substr) {
					t.Errorf("Execute() output = %q, want it to contain %q", got, substr)
				}
			}

			for _, substr := range tt.wantNotContains {
				if strings.Contains(got, substr) {
					t.Errorf("Execute() output = %q, want it NOT to contain %q", got, substr)
				}
			}

			if err := assertOrder(t, got, tt.wantOrder); err != nil {
				t.Error(err)
			}
		})
	}
}

// 'commitChange' writes a file with a unique, human-readable content and
// commits it, so each commit gets a distinct subject suitable for assertions.
func commitChange(t *testing.T, root, subject string) {
	t.Helper()

	createGitFile(t, root, "file.txt", subject+"\n")
	commitAll(t, root, subject)
}

// assertOrder verifies that every subject in want appears in output and that
// they appear in the given order. It returns nil when no order is required.
func assertOrder(t *testing.T, output string, want []string) error {
	t.Helper()

	if len(want) == 0 {
		return nil
	}

	previous := -1

	for _, subject := range want {
		index := strings.Index(output, subject)
		if index == -1 {
			return fmt.Errorf("Execute() output = %q, missing %q", output, subject)
		}
		if index < previous {
			return fmt.Errorf("Execute() output = %q, %q is out of order", output, subject)
		}
		previous = index
	}

	return nil
}
