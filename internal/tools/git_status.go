package tools

import "github.com/alexnakagama/noryn/internal/project"

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
