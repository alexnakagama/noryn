package tools

import "github.com/alexnakagama/noryn/internal/project"

type GitLogTool struct {
	project *project.Project
}

type gitLogArguments struct {
	Limit int `json:"limit"`
}

func NewGitLogTool(project *project.Project) *GitLogTool {
	return &GitLogTool{
		project: project,
	}
}

func (t *GitLogTool) Name() string {
	return "git_log"
}

func (t *GitLogTool) Execute(arguments string) (string, error) {}
