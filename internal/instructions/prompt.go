package instructions

func BuildPrompt(projectInstructions, userPrompt string) string {
	return projectInstructions + "\n\n" + userPrompt
}
