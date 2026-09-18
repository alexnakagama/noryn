package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/alexnakagama/noryn/internal/agent"
	"github.com/alexnakagama/noryn/internal/llm"
	"github.com/alexnakagama/noryn/internal/llm/providers/fake"
	"github.com/alexnakagama/noryn/internal/llm/providers/openai"
	"github.com/alexnakagama/noryn/internal/project"
	"github.com/alexnakagama/noryn/internal/tools"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("error loading .env file")
	}

	project, err := project.Discover(".")
	if err != nil {
		log.Fatal(err)
	}

	fakeMode := flag.Bool("fake", false, "use the fake LLM provider")
	model := flag.String("model", "gpt-5", "model to use")

	flag.Parse()

	if flag.NArg() < 1 {
		log.Fatal("usage: noryn [--fake] [--model MODEL] <prompt>")
	}

	prompt := flag.Arg(0)

	var client llm.Client

	if *fakeMode {
		client = fake.NewClient()
	} else {
		apiKey := os.Getenv("OPENAI_API_KEY")

		if apiKey == "" {
			log.Fatal("OPENAI_API_KEY is not set")
		}

		client = openai.NewClient(apiKey)
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
			Model: *model,
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
