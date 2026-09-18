package tools

import "github.com/alexnakagama/noryn/internal/project"

type SearchTool struct {
	project *project.Project
}

type searchArguments struct {
	Query string `json:"query"`
	Path  string `json:"path"`
}

func NewSearchTool(project *project.Project) *SearchTool {
	return &SearchTool{
		project: project,
	}
}

func (t *SearchTool) Name() string {
	return "search"
}
