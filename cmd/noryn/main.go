package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/alexnakagama/noryn/internal/agent"
	"github.com/alexnakagama/noryn/internal/config"
	"github.com/alexnakagama/noryn/internal/instructions"
	"github.com/alexnakagama/noryn/internal/llm"
	"github.com/alexnakagama/noryn/internal/llm/providers"
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

	projectContext, err := project.BuildContext()
	if err != nil {
		log.Fatal(err)
	}

	projectContextText := projectContext.String()

	instructionsBuilder := instructions.NewBuilder(project)

	projectInstructions, err := instructionsBuilder.Build()
	if err != nil {
		log.Fatal(err)
	}

	provider := flag.String("provider", cfg.Provider, "LLM provider")
	model := flag.String("model", cfg.Model, "LLM model")

	flag.Parse()

	cfg.Provider = *provider
	cfg.Model = *model

	client, err := providers.New(*cfg)
	if err != nil {
		log.Fatal(err)
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

	agent.SetToolCallHandler(func(call llm.ToolCall) {
		fmt.Println("→", call.Name, call.Arguments)
	})

	agent.SetToolResultHandler(func(call llm.ToolCall, result string) {
		fmt.Println("←", call.Name)
	})

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Noryn")
	fmt.Println("Type 'exit' to quit.")
	fmt.Println()

	for {
		fmt.Print("noryn> ")

		if !scanner.Scan() {
			break
		}

		prompt := strings.TrimSpace(scanner.Text())

		if prompt == "" {
			continue
		}

		if prompt == "exit" {
			break
		}

		response, err := agent.Chat(
			context.Background(),
			llm.Request{
				Model: cfg.Model,
				Messages: []llm.Message{
					{
						Role: "user",
						Content: instructions.BuildPrompt(
							projectInstructions,
							projectContextText,
							prompt,
						),
					},
				},
			},
		)

		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		fmt.Println(response.Message.Content)
		fmt.Println()
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
