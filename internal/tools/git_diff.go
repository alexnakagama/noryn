package tools

import "github.com/alexnakagama/noryn/internal/project"

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

func (t *GitDiffTool) Execute(arguments string) (string, error) {}
