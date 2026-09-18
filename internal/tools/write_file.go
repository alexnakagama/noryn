package tools

type WriteFileTool struct{}

type writeFileArguments struct {
	Path    string `json:"path"`
	Content string `json:"content"`
}

func (t *WriteFileTool) Name() string {
	return "write_file"
}

func (t *WriteFileTool) Execute(arguments string) (string, error) {}
