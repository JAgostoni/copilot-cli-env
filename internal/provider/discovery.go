package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type openAIModelResponse struct {
	Data []struct {
		ID            string `json:"id"`
		ContextLength int    `json:"context_length,omitempty"`
		TopProvider   struct {
			MaxCompletionTokens int `json:"max_completion_tokens,omitempty"`
		} `json:"top_provider,omitempty"`
	} `json:"data"`
}

// fetchOpenAIModels is a shared helper to fetch models from any OpenAI-compatible /v1/models endpoint
func fetchOpenAIModels(ctx context.Context, baseURL, apiKey string) ([]Model, error) {
	if baseURL == "" {
		return nil, fmt.Errorf("base URL is required")
	}
	if apiKey == "" {
		return nil, fmt.Errorf("API key is required")
	}

	endpoint := strings.TrimRight(baseURL, "/") + "/models"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status: %s", resp.Status)
	}

	var result openAIModelResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	var models []Model
	for _, m := range result.Data {
		models = append(models, Model{
			ID:                  m.ID,
			ContextLength:       m.ContextLength,
			MaxCompletionTokens: m.TopProvider.MaxCompletionTokens,
		})
	}
	return models, nil
}
