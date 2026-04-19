package provider

import (
	"testing"
)

func TestOpenAIProvider(t *testing.T) {
	p := NewOpenAIProvider()
	
	if p.ID() != OpenAI {
		t.Errorf("Expected ID %s, got %s", OpenAI, p.ID())
	}

	_, err := p.MapEnv("", "gpt-4o", "")
	if err == nil {
		t.Error("Expected error for missing API key")
	}

	_, err = p.MapEnv("", "", "test-key")
	if err == nil {
		t.Error("Expected error for missing model")
	}

	env, err := p.MapEnv("", "gpt-4o", "test-key")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if env.ProviderType != "openai" {
		t.Errorf("Expected ProviderType openai, got %s", env.ProviderType)
	}
	if env.ProviderBaseURL != p.DefaultBaseURL() {
		t.Errorf("Expected BaseURL %s, got %s", p.DefaultBaseURL(), env.ProviderBaseURL)
	}
	if env.Model != "gpt-4o" {
		t.Errorf("Expected Model gpt-4o, got %s", env.Model)
	}
	if env.APIKey != "test-key" {
		t.Errorf("Expected APIKey test-key, got %s", env.APIKey)
	}
	if env.Offline != false {
		t.Errorf("Expected Offline false, got true")
	}
}

func TestAnthropicProvider(t *testing.T) {
	p := NewAnthropicProvider()
	
	if p.ID() != Anthropic {
		t.Errorf("Expected ID %s, got %s", Anthropic, p.ID())
	}

	env, err := p.MapEnv("", "claude-3-opus", "test-key")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if env.ProviderType != "anthropic" {
		t.Errorf("Expected ProviderType anthropic, got %s", env.ProviderType)
	}
}

func TestAzureProvider(t *testing.T) {
	p := NewAzureProvider()
	
	if p.ID() != Azure {
		t.Errorf("Expected ID %s, got %s", Azure, p.ID())
	}

	_, err := p.MapEnv("", "model", "key")
	if err == nil {
		t.Error("Expected error for missing base URL")
	}

	_, err = p.MapEnv("invalid-url", "model", "key")
	if err == nil {
		t.Error("Expected error for invalid base URL")
	}

	_, err = p.MapEnv("https://example.com", "model", "key")
	if err == nil {
		t.Error("Expected error for non-Azure domain")
	}

	env, err := p.MapEnv("https://my-resource.openai.azure.com/", "gpt-4", "test-key")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if env.ProviderType != "azure" {
		t.Errorf("Expected ProviderType azure, got %s", env.ProviderType)
	}
}

func TestOpenRouterProvider(t *testing.T) {
	p := NewOpenRouterProvider()
	
	env, err := p.MapEnv("", "meta-llama/llama-3-8b-instruct", "test-key")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if env.ProviderType != "openai" {
		t.Errorf("Expected ProviderType openai, got %s", env.ProviderType)
	}
	if env.ProviderBaseURL != "https://openrouter.ai/api/v1" {
		t.Errorf("Expected BaseURL https://openrouter.ai/api/v1, got %s", env.ProviderBaseURL)
	}
}

func TestOllamaProvider(t *testing.T) {
	p := NewOllamaProvider()
	
	env, err := p.MapEnv("", "llama3", "")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if env.ProviderType != "openai" {
		t.Errorf("Expected ProviderType openai, got %s", env.ProviderType)
	}
	if env.Offline != true {
		t.Errorf("Expected Offline true, got %v", env.Offline)
	}
	if env.APIKey != "dummy" {
		t.Errorf("Expected dummy APIKey, got %s", env.APIKey)
	}
	if env.ProviderBaseURL != "http://localhost:11434/v1" {
		t.Errorf("Expected BaseURL http://localhost:11434/v1, got %s", env.ProviderBaseURL)
	}
}
