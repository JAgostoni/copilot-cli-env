package provider

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/jagostoni/copilot-cli-env/internal/config"
)

type anthropicProvider struct{}

func NewAnthropicProvider() Provider {
	return &anthropicProvider{}
}

func (p *anthropicProvider) ID() ProviderID {
	return Anthropic
}

func (p *anthropicProvider) Name() string {
	return "Anthropic"
}

func (p *anthropicProvider) RequiresAPIKey() bool {
	return true
}

func (p *anthropicProvider) DefaultBaseURL() string {
	return "https://api.anthropic.com/v1"
}

func (p *anthropicProvider) MapEnv(baseURL, model, apiKey string) (config.CopilotEnv, error) {
	if apiKey == "" {
		return config.CopilotEnv{}, errors.New("API key is required for Anthropic")
	}
	if model == "" {
		return config.CopilotEnv{}, errors.New("model is required for Anthropic")
	}

	finalBaseURL := baseURL
	if finalBaseURL == "" {
		finalBaseURL = p.DefaultBaseURL()
	}

	return config.CopilotEnv{
		ProviderBaseURL: finalBaseURL,
		Model:           model,
		ProviderType:    "anthropic",
		APIKey:          apiKey,
		Offline:         false,
	}, nil
}

func (p *anthropicProvider) FetchModels(ctx context.Context, baseURL, apiKey string) ([]Model, error) {
	finalBaseURL := baseURL
	if finalBaseURL == "" {
		finalBaseURL = p.DefaultBaseURL()
	}

	if apiKey == "" {
		return nil, fmt.Errorf("API key is required")
	}

	endpoint := strings.TrimRight(finalBaseURL, "/") + "/models"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status: %s", resp.Status)
	}

	var result struct {
		Data []struct {
			ID          string `json:"id"`
			DisplayName string `json:"display_name"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	var models []Model
	for _, m := range result.Data {
		desc := m.DisplayName
		models = append(models, Model{ID: m.ID, Description: desc})
	}
	return models, nil
}
