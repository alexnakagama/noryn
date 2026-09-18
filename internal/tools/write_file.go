package tools

import (
	"encoding/json"
	"os"
)

type WriteFileTool struct{}

type writeFileArguments struct {
	Path    string `json:"path"`
	Content string `json:"content"`
}

func (t *WriteFileTool) Name() string {
	return "write_file"
}

func (t *WriteFileTool) Execute(arguments string) (string, error) {
	var args writeFileArguments

	err := json.Unmarshal([]byte(arguments), &args)
	if err != nil {
		return "", err
	}

	err = os.WriteFile(args.Path, []byte(args.Content), 0644)
	if err != nil {
		return "", err
	}

	return "file written successfully", nil
}
