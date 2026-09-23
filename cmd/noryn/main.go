package main

import (
	"flag"
	"log"

	"github.com/alexnakagama/noryn/internal/agent"
	"github.com/alexnakagama/noryn/internal/config"
	"github.com/alexnakagama/noryn/internal/instructions"
	"github.com/alexnakagama/noryn/internal/llm/providers"
	"github.com/alexnakagama/noryn/internal/project"
	"github.com/alexnakagama/noryn/internal/tools"
	"github.com/alexnakagama/noryn/internal/tui"
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

	registry := tools.NewRegistry(
		tools.NewReadFileTool(project),
		tools.NewWriteFileTool(project),
		tools.NewListDirectoryTool(project),
		tools.NewShellTool(project),
		tools.NewSearchTool(project),
		tools.NewGitStatusTool(project),
		tools.NewGitDiffTool(project),
		tools.NewGitLogTool(project),
	)

	agent := agent.New(
		client,
		registry,
	)

	if err := tui.Run(
		agent,
		cfg.Model,
		projectInstructions,
		projectContextText,
	); err != nil {
		log.Fatal(err)
	}
}
