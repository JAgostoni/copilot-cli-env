# Session 3: Provider Adapters & Env Mapping

## Goal
Create the abstraction layer for supported AI providers and map them to Copilot CLI variables.

## Tooling & Library Versions
- **Go Version:** `1.26.x`
- **Libraries:** Standard library only (`errors`, `fmt`, `net/url`). We avoid heavy provider SDKs here to keep the binary small and focused.

## Setup Instructions
Create the following structure:
```text
internal/
  provider/
    provider.go
    openai.go
    anthropic.go
    azure.go
    openrouter.go
    ollama.go
    provider_test.go
```

## Tasks

### 1. Define the Provider Interface
In `internal/provider/provider.go`, define how the system interacts with a provider:

```go
package provider

import "github.com/yourusername/copilot-cli-env/internal/config"

type ProviderID string

const (
	OpenAI     ProviderID = "openai"
	Anthropic  ProviderID = "anthropic"
	Azure      ProviderID = "azure"
	OpenRouter ProviderID = "openrouter"
	Ollama     ProviderID = "ollama"
	Custom     ProviderID = "custom"
)

type Provider interface {
	ID() ProviderID
	Name() string
	RequiresAPIKey() bool
	DefaultBaseURL() string
	// MapEnv takes user inputs and outputs the exact Copilot Env shape
	MapEnv(baseURL, model, apiKey string) (config.CopilotEnv, error)
}
```

### 2. Implement OpenAI & Anthropic Adapters
- `openai.go`: Maps `ProviderType` to `"openai"`.
- `anthropic.go`: Maps `ProviderType` to `"anthropic"`.

### 3. Implement Azure Adapter
Azure requires specific URL construction. Ensure `MapEnv` validates that the `baseURL` provided actually looks like a valid Azure endpoint before returning.

### 4. Implement OpenRouter Adapter
- `openrouter.go`: Crucially, this must return a `config.CopilotEnv` with `ProviderType: "openai"`, because GitHub Copilot CLI treats OpenAI-compatible endpoints as the `openai` provider type. The `ProviderBaseURL` should default to `https://openrouter.ai/api/v1`.

### 5. Implement Ollama / Local Adapter
- `ollama.go`: Maps to `ProviderType: "openai"`. Sets `Offline: true`. Sets `APIKey: "dummy"` (or whatever local auth requires, often it just needs to be non-empty for Copilot).

## Unit Testing
Write table-driven unit tests in `provider_test.go` ensuring near 100% test coverage for this logic. Ensure that:
1. `OpenRouter.MapEnv` returns `ProviderType: "openai"`.
2. `Anthropic.MapEnv` returns `ProviderType: "anthropic"`.
3. `Ollama.MapEnv` sets `Offline: true`.
