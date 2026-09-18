package main

import (
	"log"
	"os"

	"github.com/alexnakagama/noryn/internal/agent"
	"github.com/alexnakagama/noryn/internal/llm/providers/openai"
	"github.com/alexnakagama/noryn/internal/project"
	"github.com/alexnakagama/noryn/internal/tools"
)

func main() {
	project, err := project.Discover(".")
	if err != nil {
		log.Fatal(err)
	}

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENAI_API_KEY is not set")
	}

	client := openai.NewClient(apiKey)

	agent := agent.New(
		client,
		tools.NewReadFileTool(project),
		tools.NewWriteFileTool(project),
		tools.NewListDirectoryTool(project),
		tools.NewShellTool(project),
		tools.NewSearchTool(project),
		tools.NewGitStatusTool(project),
		tools.NewGitDiffTool(project),
		tools.NewGitLogTool(project),
	)

	_ = agent

	log.Println("noryn initialized")
}
