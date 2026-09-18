package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/alexnakagama/noryn/internal/agent"
	"github.com/alexnakagama/noryn/internal/config"
	"github.com/alexnakagama/noryn/internal/llm"
	"github.com/alexnakagama/noryn/internal/llm/providers/fake"
	"github.com/alexnakagama/noryn/internal/llm/providers/openai"
	"github.com/alexnakagama/noryn/internal/project"
	"github.com/alexnakagama/noryn/internal/tools"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	project, err := project.Discover(".")
	if err != nil {
		log.Fatal(err)
	}

	fakeMode := flag.Bool("fake", false, "use the fake LLM provider")

	flag.Parse()

	if flag.NArg() < 1 {
		log.Fatal("usage: noryn [--fake] <prompt>")
	}

	prompt := flag.Arg(0)

	var client llm.Client

	if *fakeMode {
		client = fake.NewClient()
	} else {
		if cfg.APIKey == "" {
			log.Fatal("OPENAI_API_KEY is not set")
		}

		client = openai.NewClient(cfg.APIKey)
	}

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

	response, err := agent.Chat(
		context.Background(),
		llm.Request{
			Model: cfg.Model,
			Messages: []llm.Message{
				{
					Role:    "user",
					Content: prompt,
				},
			},
		},
	)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(response.Message.Content)
}
