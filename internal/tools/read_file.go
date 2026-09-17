package tools

type ReadFileTool struct{}

type readFileArguments struct {
	Path string `json:"path"`
}

func (t *ReadFileTool) Name() string {
	return "read_file"
}
