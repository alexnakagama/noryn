package instructions

import "strings"

func BuildPrompt(systemPrompt, projectInstructions, projectContext, userPrompt string) string {
	systemPrompt = strings.TrimRight(systemPrompt, "\n")
	projectInstructions = strings.TrimRight(projectInstructions, "\n")
	projectContext = strings.TrimRight(projectContext, "\n")

	return systemPrompt + "\n\n" +
		projectInstructions + "\n\n" +
		projectContext + "\n\n" +
		userPrompt
}
