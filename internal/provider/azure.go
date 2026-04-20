package provider

import (
	"context"
	"errors"
	"net/url"
	"strings"

	"github.com/jagostoni/copilot-cli-env/internal/config"
)

type azureProvider struct{}

func NewAzureProvider() Provider {
	return &azureProvider{}
}

func (p *azureProvider) ID() ProviderID {
	return Azure
}

func (p *azureProvider) Name() string {
	return "Azure OpenAI"
}

func (p *azureProvider) RequiresAPIKey() bool {
	return true
}

func (p *azureProvider) DefaultBaseURL() string {
	return ""
}

func (p *azureProvider) MapEnv(baseURL, model, apiKey string) (config.CopilotEnv, error) {
	if baseURL == "" {
		return config.CopilotEnv{}, errors.New("base URL is required for Azure")
	}
	if apiKey == "" {
		return config.CopilotEnv{}, errors.New("API key is required for Azure")
	}
	if model == "" {
		return config.CopilotEnv{}, errors.New("model is required for Azure")
	}

	parsedURL, err := url.Parse(baseURL)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		return config.CopilotEnv{}, errors.New("invalid base URL for Azure")
	}
	
	if !strings.HasSuffix(parsedURL.Host, ".openai.azure.com") {
		return config.CopilotEnv{}, errors.New("base URL must be a valid Azure OpenAI endpoint (e.g., https://your-resource.openai.azure.com/)")
	}

	wireAPI := ""
	for i := 0; i < len(model)-4; i++ {
		if model[i:i+5] == "gpt-5" {
			wireAPI = "responses"
			break
		}
	}

	return config.CopilotEnv{
		ProviderBaseURL: baseURL,
		Model:           model,
		ProviderType:    "azure",
		APIKey:          apiKey,
		Offline:         false,
		WireAPI:         wireAPI,
	}, nil
}

func (p *azureProvider) FetchModels(ctx context.Context, baseURL, apiKey string) ([]Model, error) {
	return nil, ErrDiscoveryUnsupported
}
