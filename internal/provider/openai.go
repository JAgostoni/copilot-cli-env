package provider

import (
	"context"
	"errors"

	"github.com/jagostoni/copilot-cli-env/internal/config"
)

type openaiProvider struct{}

func NewOpenAIProvider() Provider {
	return &openaiProvider{}
}

func (p *openaiProvider) ID() ProviderID {
	return OpenAI
}

func (p *openaiProvider) Name() string {
	return "OpenAI"
}

func (p *openaiProvider) RequiresAPIKey() bool {
	return true
}

func (p *openaiProvider) DefaultBaseURL() string {
	return "https://api.openai.com/v1"
}

func (p *openaiProvider) MapEnv(baseURL, model, apiKey string) (config.CopilotEnv, error) {
	if apiKey == "" {
		return config.CopilotEnv{}, errors.New("API key is required for OpenAI")
	}
	if model == "" {
		return config.CopilotEnv{}, errors.New("model is required for OpenAI")
	}

	finalBaseURL := baseURL
	if finalBaseURL == "" {
		finalBaseURL = p.DefaultBaseURL()
	}

	return config.CopilotEnv{
		ProviderBaseURL: finalBaseURL,
		Model:           model,
		ProviderType:    "openai",
		APIKey:          apiKey,
		Offline:         false,
	}, nil
}

func (p *openaiProvider) FetchModels(ctx context.Context, baseURL, apiKey string) ([]Model, error) {
	finalBaseURL := baseURL
	if finalBaseURL == "" {
		finalBaseURL = p.DefaultBaseURL()
	}
	return fetchOpenAIModels(ctx, finalBaseURL, apiKey)
}
