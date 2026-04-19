package provider

import (
	"context"
	"errors"

	"github.com/jagostoni/copilot-cli-env/internal/config"
)

var ErrDiscoveryUnsupported = errors.New("model discovery is not supported for this provider")

type ProviderID string

const (
	OpenAI     ProviderID = "openai"
	Anthropic  ProviderID = "anthropic"
	Azure      ProviderID = "azure"
	OpenRouter ProviderID = "openrouter"
	Ollama     ProviderID = "ollama"
	Custom     ProviderID = "custom"
	Reset      ProviderID = "reset"
)

type Model struct {
	ID                  string
	Description         string // Optional friendly name or context window info
	ContextLength       int    // Optional max prompt tokens
	MaxCompletionTokens int    // Optional max output tokens
}

type Provider interface {
	ID() ProviderID
	Name() string
	RequiresAPIKey() bool
	DefaultBaseURL() string
	// MapEnv takes user inputs and outputs the exact Copilot Env shape
	MapEnv(baseURL, model, apiKey string) (config.CopilotEnv, error)
	// FetchModels retrieves the list of available models from the provider
	FetchModels(ctx context.Context, baseURL, apiKey string) ([]Model, error)
}
