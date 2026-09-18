package tools

import (
	"encoding/json"
	"os/exec"

	"github.com/alexnakagama/noryn/internal/project"
)

type GitDiffTool struct {
	project *project.Project
}

type gitDiffArguments struct {
	Path string `json:"path"`
}

func NewGitDiffTool(project *project.Project) *GitDiffTool {
	return &GitDiffTool{
		project: project,
	}
}

func (t *GitDiffTool) Name() string {
	return "git_diff"
}

func (t *GitDiffTool) Execute(arguments string) (string, error) {
	var args gitDiffArguments

	err := json.Unmarshal([]byte(arguments), &args)
	if err != nil {
		return "", err
	}

	cmdArgs := []string{"diff"}

	if args.Path != "" {
		path, err := t.project.ResolvePath(args.Path)
		if err != nil {
			return "", err
		}

		cmdArgs = append(cmdArgs, "--", path)
	}

	cmd := exec.Command("git", cmdArgs...)
	cmd.Dir = t.project.Root

	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), err
	}

	return string(output), nil
}
