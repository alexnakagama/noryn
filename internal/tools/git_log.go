package tools

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"

	"github.com/alexnakagama/noryn/internal/llm"
	"github.com/alexnakagama/noryn/internal/project"
)

type GitLogTool struct {
	project *project.Project
}

type gitLogArguments struct {
	Limit int `json:"limit"`
}

func NewGitLogTool(project *project.Project) *GitLogTool {
	return &GitLogTool{
		project: project,
	}
}

func (t *GitLogTool) Name() string {
	return "git_log"
}

func (t *GitLogTool) Description() string {
	return "Show the recent Git commit history of the project."
}

func (t *GitLogTool) Definition() llm.ToolDefinition {
	return llm.ToolDefinition{
		Name:        t.Name(),
		Description: t.Description(),
		Parameters: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"limit": map[string]any{
					"type":        "integer",
					"description": "Maximum number of commits to show. Must be between 1 and 50.",
				},
			},
			"required": []string{"limit"},
		},
	}
}

func (t *GitLogTool) Execute(arguments string) (string, error) {
	var args gitLogArguments

	err := json.Unmarshal([]byte(arguments), &args)
	if err != nil {
		return "", err
	}

	if args.Limit <= 0 {
		return "", fmt.Errorf("limit cannot be 0 or less than 0")
	}

	if args.Limit > 50 {
		return "", fmt.Errorf("limit cannot be greater than 50")
	}

	cmd := exec.Command(
		"git",
		"log",
		"--oneline",
		"-n",
		strconv.Itoa(args.Limit),
	)

	cmd.Dir = t.project.Root

	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), err
	}

	return string(output), nil
}
