package provider

import (
	"context"
	"errors"

	"github.com/jagostoni/copilot-cli-env/internal/config"
)

type openRouterProvider struct{}

func NewOpenRouterProvider() Provider {
	return &openRouterProvider{}
}

func (p *openRouterProvider) ID() ProviderID {
	return OpenRouter
}

func (p *openRouterProvider) Name() string {
	return "OpenRouter"
}

func (p *openRouterProvider) RequiresAPIKey() bool {
	return true
}

func (p *openRouterProvider) DefaultBaseURL() string {
	return "https://openrouter.ai/api/v1"
}

func (p *openRouterProvider) MapEnv(baseURL, model, apiKey string) (config.CopilotEnv, error) {
	if apiKey == "" {
		return config.CopilotEnv{}, errors.New("API key is required for OpenRouter")
	}
	if model == "" {
		return config.CopilotEnv{}, errors.New("model is required for OpenRouter")
	}

	finalBaseURL := baseURL
	if finalBaseURL == "" {
		finalBaseURL = p.DefaultBaseURL()
	}

	// OpenRouter acts as an OpenAI compatible endpoint
	wireAPI := ""
	if len(model) >= 5 && model[:5] == "gpt-5" {
		wireAPI = "responses"
	} else if len(model) > 7 && model[len(model)-5:] == "gpt-5" {
		// Some openrouter models might look like openai/gpt-5
		wireAPI = "responses"
	} else {
		// General substring check for gpt-5
		for i := 0; i < len(model)-4; i++ {
			if model[i:i+5] == "gpt-5" {
				wireAPI = "responses"
				break
			}
		}
	}

	return config.CopilotEnv{
		ProviderBaseURL: finalBaseURL,
		Model:           model,
		ProviderType:    "openai",
		APIKey:          apiKey,
		Offline:         false,
		WireAPI:         wireAPI,
	}, nil
}

func (p *openRouterProvider) FetchModels(ctx context.Context, baseURL, apiKey string) ([]Model, error) {
	finalBaseURL := baseURL
	if finalBaseURL == "" {
		finalBaseURL = p.DefaultBaseURL()
	}
	return fetchOpenAIModels(ctx, finalBaseURL, apiKey)
}
