package tools

import (
	"encoding/json"
	"os/exec"

	"github.com/alexnakagama/noryn/internal/project"
)

type ShellTool struct {
	project *project.Project
}

type shellArguments struct {
	Command string `json:"command"`
}

func (t *ShellTool) Name() string {
	return "shell"
}

func (t *ShellTool) Execute(arguments string) (string, error) {
	var args shellArguments

	err := json.Unmarshal([]byte(arguments), &args)
	if err != nil {
		return "", err
	}

	cmd := exec.Command("sh", "-c", args.Command)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), err
	}

	return string(output), nil
}
