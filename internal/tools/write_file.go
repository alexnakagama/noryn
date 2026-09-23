package tools

import (
	"encoding/json"
	"os"

	"github.com/alexnakagama/noryn/internal/llm"
	"github.com/alexnakagama/noryn/internal/project"
)

type WriteFileTool struct {
	project *project.Project
}

type writeFileArguments struct {
	Path    string `json:"path"`
	Content string `json:"content"`
}

func NewWriteFileTool(project *project.Project) *WriteFileTool {
	return &WriteFileTool{
		project: project,
	}
}

func (t *WriteFileTool) Name() string {
	return "write_file"
}

func (t *WriteFileTool) Description() string {
	return "Write content to a file inside the project."
}

func (t *WriteFileTool) Definition() llm.ToolDefinition {
	return llm.ToolDefinition{
		Name:        t.Name(),
		Description: t.Description(),
		Parameters: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"path": map[string]any{
					"type":        "string",
					"description": "Path of the file to write.",
				},
				"content": map[string]any{
					"type":        "string",
					"description": "Content to write to the file.",
				},
			},
			"required": []string{"path", "content"},
		},
	}
}

func (t *WriteFileTool) Execute(arguments string) (string, error) {
	var args writeFileArguments

	err := json.Unmarshal([]byte(arguments), &args)
	if err != nil {
		return "", err
	}

	path, err := t.project.ResolvePath(args.Path)
	if err != nil {
		return "", err
	}

	err = os.WriteFile(path, []byte(args.Content), 0644)
	if err != nil {
		return "", err
	}

	return "file written successfully", nil
}
