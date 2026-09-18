package providers

import (
	"github.com/alexnakagama/noryn/internal/config"
	"github.com/alexnakagama/noryn/internal/llm"
)

func New(cfg config.Config) (llm.Client, error) {}
