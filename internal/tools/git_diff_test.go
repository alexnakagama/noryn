package tools

import (
	"os/exec"
	"strings"
	"testing"

	"github.com/alexnakagama/noryn/internal/project"
)

func TestGitDiffTool_Name(t *testing.T) {
	tool := NewGitDiffTool(nil)

	got := tool.Name()

	if got != "git_diff" {
		t.Errorf("Name() = %q, want %q", got, "git_diff")
	}
}

func TestGitDiffTool_Execute(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git is not installed")
	}

	tests := []struct {
		name            string
		arguments       string
		setup           func(t *testing.T, root string)
		want            string
		wantContains    []string
		wantNotContains []string
		wantErr         bool
		wantErrContains string
	}{
		{
			name:      "reports no diff in a clean repository",
			arguments: "{}",
			setup: func(t *testing.T, root string) {
				t.Helper()

				initGitRepo(t, root)
			},
			want: "",
		},
		{
			name:      "ignores untracked files",
			arguments: "{}",
			setup: func(t *testing.T, root string) {
				t.Helper()

				initGitRepo(t, root)
				createGitFile(t, root, "untracked.txt", "new\n")
			},
			want: "",
		},
		{
			name:      "reports an unstaged modification",
			arguments: "{}",
			setup: func(t *testing.T, root string) {
				t.Helper()

				initGitRepo(t, root)
				createGitFile(t, root, "tracked.txt", "original\n")
				commitAll(t, root, "initial")
				createGitFile(t, root, "tracked.txt", "changed\n")
			},
			wantContains: []string{"-original", "+changed"},
		},
		{
			name:      "does not report staged changes as unstaged",
			arguments: "{}",
			setup: func(t *testing.T, root string) {
				t.Helper()

				initGitRepo(t, root)
				createGitFile(t, root, "tracked.txt", "original\n")
				commitAll(t, root, "initial")
				createGitFile(t, root, "tracked.txt", "changed\n")
				runGit(t, root, "add", "tracked.txt")
			},
			want: "",
		},
		{
			name:      "filters the diff to the requested path",
			arguments: `{"path":"a.txt"}`,
			setup: func(t *testing.T, root string) {
				t.Helper()

				initGitRepo(t, root)
				createGitFile(t, root, "a.txt", "a-original\n")
				createGitFile(t, root, "b.txt", "b-original\n")
				commitAll(t, root, "initial")
				createGitFile(t, root, "a.txt", "a-changed\n")
				createGitFile(t, root, "b.txt", "b-changed\n")
			},
			wantContains:    []string{"diff --git a/a.txt", "+a-changed"},
			wantNotContains: []string{"b.txt"},
		},
		{
			name:      "reports nothing for a path with no changes",
			arguments: `{"path":"untouched.txt"}`,
			setup: func(t *testing.T, root string) {
				t.Helper()

				initGitRepo(t, root)
				createGitFile(t, root, "untouched.txt", "same\n")
				commitAll(t, root, "initial")
			},
			want: "",
		},
		{
			name:      "empty path shows the project diff",
			arguments: `{"path":""}`,
			setup: func(t *testing.T, root string) {
				t.Helper()

				initGitRepo(t, root)
				createGitFile(t, root, "tracked.txt", "original\n")
				commitAll(t, root, "initial")
				createGitFile(t, root, "tracked.txt", "changed\n")
			},
			wantContains: []string{"-original", "+changed"},
		},
		{
			name:      "rejects a path outside the project",
			arguments: `{"path":"../../outside.txt"}`,
			setup: func(t *testing.T, root string) {
				t.Helper()

				initGitRepo(t, root)
			},
			wantErr:         true,
			wantErrContains: "outside the project",
		},
		{
			name:      "rejects malformed arguments",
			arguments: `{not json`,
			setup: func(t *testing.T, root string) {
				t.Helper()

				initGitRepo(t, root)
			},
			wantErr: true,
		},
		{
			name:            "not a git repository returns an error",
			arguments:       "{}",
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

			tool := NewGitDiffTool(&project.Project{Root: root})

			got, err := tool.Execute(tt.arguments)

			if tt.wantErr {
				if err == nil {
					t.Fatal("Execute() error = nil, want an error")
				}

				if tt.wantErrContains != "" {
					combined := strings.ToLower(err.Error() + "\n" + got)
					want := strings.ToLower(tt.wantErrContains)

					if !strings.Contains(combined, want) {
						t.Errorf(
							"Execute() error/output = %q, want it to contain %q",
							combined,
							tt.wantErrContains,
						)
					}
				}

				return
			}

			if err != nil {
				t.Fatalf("Execute() error = %v", err)
			}

			for _, substr := range tt.wantContains {
				if !strings.Contains(got, substr) {
					t.Errorf(
						"Execute() output = %q, want it to contain %q",
						got,
						substr,
					)
				}
			}

			for _, substr := range tt.wantNotContains {
				if strings.Contains(got, substr) {
					t.Errorf(
						"Execute() output = %q, want it NOT to contain %q",
						got,
						substr,
					)
				}
			}

			if len(tt.wantContains) == 0 && got != tt.want {
				t.Errorf(
					"Execute() output = %q, want %q",
					got,
					tt.want,
				)
			}
		})
	}
}
