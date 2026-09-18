package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Provider         string
	Model            string
	APIKey           string
	OpenRouterAPIKey string
}

func Load() (*Config, error) {
	if err := godotenv.Load(".env"); err != nil {
		return nil, err
	}

	provider := os.Getenv("NORYN_PROVIDER")
	if provider == "" {
		provider = "openai"
	}

	model := os.Getenv("NORYN_MODEL")
	if model == "" {
		model = "gpt-5"
	}

	apiKey := os.Getenv("OPENAI_API_KEY")

	return &Config{
		Provider: provider,
		Model:    model,
		APIKey:   apiKey,
	}, nil
}
