package instructions

import "strings"

func BuildPrompt(projectInstructions, projectContext, userPrompt string) string {
	projectContext = strings.TrimRight(projectContext, "\n")

	return projectInstructions + "\n\n" +
		projectContext + "\n\n" +
		userPrompt
}
