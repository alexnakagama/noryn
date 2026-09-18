package tools

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alexnakagama/noryn/internal/project"
)

func TestGitStatusTool_Name(t *testing.T) {
	tool := NewGitStatusTool(nil)

	got := tool.Name()

	if got != "git_status" {
		t.Errorf("Name() = %q, want %q", got, "git_status")
	}
}

func TestGitStatusTool_Execute(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git is not installed")
	}

	tests := []struct {
		name            string
		setup           func(t *testing.T, root string)
		want            string
		wantErr         bool
		wantErrContains string
	}{
		{
			name: "reports a clean repository as empty",
			setup: func(t *testing.T, root string) {
				t.Helper()

				initGitRepo(t, root)
			},
			want: "",
		},
		{
			name: "reports untracked files in path order",
			setup: func(t *testing.T, root string) {
				t.Helper()

				initGitRepo(t, root)
				createGitFile(t, root, "a.txt", "a\n")
				createGitFile(t, root, "b.txt", "b\n")
			},
			want: "?? a.txt\n?? b.txt\n",
		},
		{
			name: "reports a staged new file",
			setup: func(t *testing.T, root string) {
				t.Helper()

				initGitRepo(t, root)
				createGitFile(t, root, "added.txt", "new\n")
				runGit(t, root, "add", "added.txt")
			},
			want: "A  added.txt\n",
		},
		{
			name: "reports an unstaged modification",
			setup: func(t *testing.T, root string) {
				t.Helper()

				initGitRepo(t, root)
				createGitFile(t, root, "tracked.txt", "original\n")
				commitAll(t, root, "initial")
				createGitFile(t, root, "tracked.txt", "changed\n")
			},
			want: " M tracked.txt\n",
		},
		{
			name: "reports a staged modification",
			setup: func(t *testing.T, root string) {
				t.Helper()

				initGitRepo(t, root)
				createGitFile(t, root, "tracked.txt", "original\n")
				commitAll(t, root, "initial")
				createGitFile(t, root, "tracked.txt", "changed\n")
				runGit(t, root, "add", "tracked.txt")
			},
			want: "M  tracked.txt\n",
		},
		{
			name: "reports a deleted tracked file",
			setup: func(t *testing.T, root string) {
				t.Helper()

				initGitRepo(t, root)
				createGitFile(t, root, "gone.txt", "content\n")
				commitAll(t, root, "initial")

				if err := os.Remove(filepath.Join(root, "gone.txt")); err != nil {
					t.Fatalf("os.Remove() error = %v", err)
				}
			},
			want: " D gone.txt\n",
		},
		{
			name: "reports staged and untracked changes together",
			setup: func(t *testing.T, root string) {
				t.Helper()

				initGitRepo(t, root)
				createGitFile(t, root, "staged.txt", "staged\n")
				runGit(t, root, "add", "staged.txt")
				createGitFile(t, root, "untracked.txt", "untracked\n")
			},
			want: "A  staged.txt\n?? untracked.txt\n",
		},
		{
			name:            "not a git repository returns an error",
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

			tool := NewGitStatusTool(&project.Project{Root: root})

			got, err := tool.Execute("")

			if tt.wantErr {
				if err == nil {
					t.Fatal("Execute() error = nil, want an error")
				}

				if tt.wantErrContains != "" && !strings.Contains(got, tt.wantErrContains) {
					t.Errorf("Execute() output = %q, want it to contain %q", got, tt.wantErrContains)
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

func initGitRepo(t *testing.T, root string) {
	t.Helper()

	runGit(t, root, "init", "--quiet")
	runGit(t, root, "config", "user.email", "test@noryn.dev")
	runGit(t, root, "config", "user.name", "Noryn Test")
}

func commitAll(t *testing.T, root, message string) {
	t.Helper()

	runGit(t, root, "add", "-A")
	runGit(t, root, "commit", "--quiet", "-m", message)
}

func runGit(t *testing.T, root string, args ...string) string {
	t.Helper()

	cmd := exec.Command("git", args...)
	cmd.Dir = root

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %s error = %v\n%s", strings.Join(args, " "), err, output)
	}

	return string(output)
}

func createGitFile(t *testing.T, root, relativePath, content string) {
	t.Helper()

	path := filepath.Join(root, relativePath)

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("os.WriteFile(%q) error = %v", path, err)
	}
}
