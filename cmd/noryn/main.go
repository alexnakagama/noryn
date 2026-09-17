package main

import (
	"fmt"
	"log"

	"github.com/alexnakagama/noryn/internal/tools"
)

func main() {
	readFile := tools.ReadFileTool{}

	content, err := readFile.Execute(`{"path":"cmd/noryn/main.go"}`)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(content)
}
