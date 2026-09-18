package tools

type WriteFileTool struct{}

type writeFileArguments struct {
	Path    string `json:"path"`
	Content string `json:"content"`
}
