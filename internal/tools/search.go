package tools

import "github.com/alexnakagama/noryn/internal/project"

type SearchTool struct {
	project *project.Project
}

type searchArguments struct {
	Query string `json:"query"`
	Path  string `json:"path"`
}
