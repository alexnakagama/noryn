package tools

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alexnakagama/noryn/internal/project"
)

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

func (t *SearchTool) Execute(arguments string) (string, error) {
	var args searchArguments

	err := json.Unmarshal([]byte(arguments), &args)
	if err != nil {
		return "", err
	}

	if args.Query == "" {
		return "", fmt.Errorf("query cannot be empty")
	}

	if args.Path == "" {
		return "", fmt.Errorf("path cannot be empty")
	}

	root, err := t.project.ResolvePath(args.Path)
	if err != nil {
		return "", err
	}

	var result strings.Builder

	err = filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if entry.IsDir() {
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		content := string(data)

		for lineNumber, line := range strings.Split(content, "\n") {
			if strings.Contains(line, args.Query) {
				relativePath, err := filepath.Rel(t.project.Root, path)
				if err != nil {
					return err
				}

				fmt.Fprintf(
					&result,
					"%s:%d:%s\n",
					relativePath,
					lineNumber+1,
					line,
				)
			}
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	return result.String(), nil
}
