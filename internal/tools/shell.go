package tools

type ShellTool struct{}

type shellArguments struct {
	Command string `json:"command"`
}

func (t *ShellTool) Name() string {
	return "shell"
}

func (t *ShellTool) Execute(arguments string) (string, error) {}
