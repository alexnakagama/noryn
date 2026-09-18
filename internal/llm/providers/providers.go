package providers

import (
	"fmt"

	"github.com/alexnakagama/noryn/internal/config"
	"github.com/alexnakagama/noryn/internal/llm"
	"github.com/alexnakagama/noryn/internal/llm/providers/fake"
	"github.com/alexnakagama/noryn/internal/llm/providers/openai"
	"github.com/alexnakagama/noryn/internal/llm/providers/openrouter"
)

func New(cfg config.Config) (llm.Client, error) {
	switch cfg.Provider {
	case "fake":
		return fake.NewClient(), nil

	case "openai":
		if cfg.APIKey == "" {
			return nil, fmt.Errorf("OPENAI_API_KEY is not set")
		}

		return openai.NewClient(cfg.APIKey), nil

	case "openrouter":
		if cfg.OpenRouterAPIKey == "" {
			return nil, fmt.Errorf("OPENROUTER_API_KEY is not set")
		}

		return openrouter.NewClient(cfg.OpenRouterAPIKey), nil

	default:
		return nil, fmt.Errorf("unsupported provider: %s", cfg.Provider)
	}
}
