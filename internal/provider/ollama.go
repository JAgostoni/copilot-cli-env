package provider

import (
	"context"
	"errors"

	"github.com/jagostoni/copilot-cli-env/internal/config"
)

type ollamaProvider struct{}

func NewOllamaProvider() Provider {
	return &ollamaProvider{}
}

func (p *ollamaProvider) ID() ProviderID {
	return Ollama
}

func (p *ollamaProvider) Name() string {
	return "Ollama"
}

func (p *ollamaProvider) RequiresAPIKey() bool {
	return false
}

func (p *ollamaProvider) DefaultBaseURL() string {
	return "http://localhost:11434/v1"
}

func (p *ollamaProvider) MapEnv(baseURL, model, apiKey string) (config.CopilotEnv, error) {
	if model == "" {
		return config.CopilotEnv{}, errors.New("model is required for Ollama")
	}

	finalBaseURL := baseURL
	if finalBaseURL == "" {
		finalBaseURL = p.DefaultBaseURL()
	}
	
	finalAPIKey := apiKey
	if finalAPIKey == "" {
		finalAPIKey = "dummy" // Often required to be non-empty for Copilot
	}

	return config.CopilotEnv{
		ProviderBaseURL: finalBaseURL,
		Model:           model,
		// Local endpoints act as an OpenAI compatible endpoint
		ProviderType:    "openai",
		APIKey:          finalAPIKey,
		Offline:         true,
	}, nil
}

func (p *ollamaProvider) FetchModels(ctx context.Context, baseURL, apiKey string) ([]Model, error) {
	return nil, ErrDiscoveryUnsupported
}
