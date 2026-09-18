package main

import (
	"fmt"
	"log"

	"github.com/alexnakagama/noryn/internal/project"
)

func main() {
	project, err := project.Discover(".")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("project root: ", project.Root)
}
