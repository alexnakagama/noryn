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

func (t *ReadFileTool) Name() string {
	return "read_file"
}

func (t *ReadFileTool) Execute(arguments string) (string, error) {
	var args readFileArguments

	err := json.Unmarshal([]byte(arguments), &args)
	if err != nil {
		return "", err
	}

	data, err := os.ReadFile(args.Path)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
