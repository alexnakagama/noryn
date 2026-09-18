package tools

import (
	"encoding/json"
	"os"

	"github.com/alexnakagama/noryn/internal/project"
)

type ReadFileTool struct {
	project *project.Project
}

type readFileArguments struct {
	Path string `json:"path"`
}

func NewReadFileTool(project *project.Project) *ReadFileTool {
	return &ReadFileTool{
		project: project,
	}
}

func (t *ReadFileTool) Name() string {
	return "read_file"
}

func (t *ReadFileTool) Execute(arguments string) (string, error) {
	var args readFileArguments

	err := json.Unmarshal([]byte(arguments), &args)
	if err != nil {
		return "", err
	}

	path, err := t.project.ResolvePath(args.Path)
	if err != nil {
		return "", err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
