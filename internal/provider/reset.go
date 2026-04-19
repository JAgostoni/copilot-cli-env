package provider

import (
	"context"

	"github.com/jagostoni/copilot-cli-env/internal/config"
)

type resetProvider struct{}

func NewResetProvider() Provider {
	return &resetProvider{}
}

func (p *resetProvider) ID() ProviderID {
	return Reset
}

func (p *resetProvider) Name() string {
	return "Reset / Clear Configuration"
}

func (p *resetProvider) RequiresAPIKey() bool {
	return false
}

func (p *resetProvider) DefaultBaseURL() string {
	return ""
}

func (p *resetProvider) MapEnv(baseURL, model, apiKey string) (config.CopilotEnv, error) {
	return config.CopilotEnv{IsReset: true}, nil
}

func (p *resetProvider) FetchModels(ctx context.Context, baseURL, apiKey string) ([]Model, error) {
	return nil, ErrDiscoveryUnsupported
}
