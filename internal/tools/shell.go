package tools

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"time"

	"github.com/alexnakagama/noryn/internal/llm"
	"github.com/alexnakagama/noryn/internal/project"
)

const defaultShellTimeout = 30 * time.Second

type ShellTool struct {
	project *project.Project
	timeout time.Duration
}

type shellArguments struct {
	Command string `json:"command"`
}

func NewShellTool(project *project.Project) *ShellTool {
	return &ShellTool{
		project: project,
	}
}

func (t *ShellTool) Name() string {
	return "shell"
}

func (t *ShellTool) Description() string {
	return "Execute a shell command inside the project directory."
}

func (t *ShellTool) Definition() llm.ToolDefinition {
	return llm.ToolDefinition{
		Name:        t.Name(),
		Description: t.Description(),
		Parameters: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"command": map[string]any{
					"type":        "string",
					"description": "Shell command to execute.",
				},
			},
			"required": []string{"command"},
		},
	}
}

func (t *ShellTool) Execute(arguments string) (string, error) {
	var args shellArguments

	err := json.Unmarshal([]byte(arguments), &args)
	if err != nil {
		return "", err
	}

	if args.Command == "" {
		return "", fmt.Errorf("command cannot be empty")
	}

	cmd := exec.Command("sh", "-c", args.Command)
	cmd.Dir = t.project.Root

	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), err
	}

	return string(output), nil
}
