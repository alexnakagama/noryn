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
