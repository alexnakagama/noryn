package tools

import (
	"encoding/json"
	"os"

	"github.com/alexnakagama/noryn/internal/llm"
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

func (t *ReadFileTool) Description() string {
	return "Read the contents of a file inside the project."
}

func (t *ReadFileTool) Definition() llm.ToolDefinition {
	return llm.ToolDefinition{
		Name:        t.Name(),
		Description: t.Description(),
		Parameters: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"path": map[string]any{
					"type":        "string",
					"description": "Path of the file to read.",
				},
			},
			"required": []string{"path"},
		},
	}
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
