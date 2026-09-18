package tools

import (
	"os/exec"

	"github.com/alexnakagama/noryn/internal/project"
)

type GitStatusTool struct {
	project *project.Project
}

func NewGitStatusTool(project *project.Project) *GitStatusTool {
	return &GitStatusTool{
		project: project,
	}
}

func (t *GitStatusTool) Name() string {
	return "git_status"
}

func (t *GitStatusTool) Execute(arguements string) (string, error) {
	cmd := exec.Command("git", "status", "--short")
	cmd.Dir = t.project.Root

	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), err
	}

	return string(output), nil
}
