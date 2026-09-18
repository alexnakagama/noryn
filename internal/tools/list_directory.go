package tools

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/alexnakagama/noryn/internal/project"
)

type ListDirectoryTool struct {
	project *project.Project
}

type listDirectoryArguments struct {
	Path string `json:"path"`
}

func NewListDirectoryTool(project *project.Project) *ListDirectoryTool {
	return &ListDirectoryTool{
		project: project,
	}
}

func (t *ListDirectoryTool) Name() string {
	return "list_directory"
}

func (t *ListDirectoryTool) Execute(arguments string) (string, error) {
	var args listDirectoryArguments

	err := json.Unmarshal([]byte(arguments), &args)
	if err != nil {
		return "", err
	}

	entries, err := os.ReadDir(args.Path)
	if err != nil {
		return "", err
	}

	var result strings.Builder

	for _, entry := range entries {
		result.WriteString(entry.Name())
		result.WriteByte('\n')
	}

	return result.String(), nil
}
