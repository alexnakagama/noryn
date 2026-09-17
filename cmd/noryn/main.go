package main

import (
	"context"
	"fmt"
	"log"

	"github.com/alexnakagama/noryn/internal/agent"
	"github.com/alexnakagama/noryn/internal/llm"
	"github.com/alexnakagama/noryn/internal/llm/providers/fake"
)

func main() {
	client := fake.NewClient()
	agent := agent.New(client)

	request := llm.Request{
		Model: "fake",
		Messages: []llm.Message{
			{
				Role:    "user",
				Content: "Hello",
			},
		},
	}

	response, err := agent.Chat(context.Background(), request)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(response.Message.Content)

	for _, toolCall := range response.ToolCall {
		fmt.Println("Tool:", toolCall.Name)
		fmt.Println("Arguments:", toolCall.Arguments)
	}
}
