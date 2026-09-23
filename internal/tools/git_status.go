package tools

import (
	"os/exec"

	"github.com/alexnakagama/noryn/internal/llm"
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

func (t *GitStatusTool) Description() string {
	return "Show the current Git status of the project."
}

func (t *GitStatusTool) Definition() llm.ToolDefinition {
	return llm.ToolDefinition{
		Name:        t.Name(),
		Description: t.Description(),
		Parameters: map[string]any{
			"type":       "object",
			"properties": map[string]any{},
		},
	}
}

func (t *GitStatusTool) Execute(arguments string) (string, error) {
	cmd := exec.Command("git", "status", "--short")
	cmd.Dir = t.project.Root

	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), err
	}

	return string(output), nil
}
