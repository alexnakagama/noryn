package main

import (
	"context"
	"fmt"
	"log"

	"github.com/alexnakagama/noryn/internal/agent"
	"github.com/alexnakagama/noryn/internal/llm"
	"github.com/alexnakagama/noryn/internal/llm/providers/fake"
	"github.com/alexnakagama/noryn/internal/tools"
)

func main() {
	client := fake.NewClient()

	agent := agent.New(
		client,
		&tools.ReadFileTool{},
	)

	request := llm.Request{
		Model: "fake",
		Messages: []llm.Message{
			{
				Role:    "user",
				Content: "Read cmd/noryn/main.go",
			},
		},
	}

	response, err := agent.Chat(context.Background(), request)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(response.Message.Content)
}
