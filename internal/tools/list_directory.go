package tools

type ListDirectoryTool struct{}

type listDirectoryArguments struct {
	Path string `json:"path"`
}

func (t *ListDirectoryTool) Name() string {
	return "list_directory"
}
