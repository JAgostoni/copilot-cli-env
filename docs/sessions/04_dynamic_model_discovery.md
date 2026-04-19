# Session 4: Dynamic Model Discovery

## Goal
Implement the API clients to fetch live model lists from the providers, enabling a discoverable TUI experience.

## Tooling & Library Versions
- **Go Version:** `1.26.x`
- **Libraries:** Standard library `net/http`, `encoding/json`, `context`. We will use standard `net/http` to avoid bringing in large third-party SDKs just for a single `/v1/models` endpoint.

## Setup Instructions
Extend the provider module:
```text
internal/
  provider/
    discovery.go
    openai_models.go
    ...
```

## Tasks

### 1. Extend the Provider Interface
Add a discovery method to the `Provider` interface in `internal/provider/provider.go`:

```go
import "context"

type Model struct {
	ID          string
	Description string // Optional friendly name or context window info
}

type Provider interface {
	// ... previous methods ...
	FetchModels(ctx context.Context, baseURL, apiKey string) ([]Model, error)
}
```

### 2. Implement OpenAI / OpenRouter Discovery
Both use the standard OpenAI-compatible `/v1/models` endpoint.
- Create an HTTP GET request to `baseURL + "/v1/models"`.
- Set the `Authorization: Bearer <apiKey>` header.
- Parse the JSON response:
  ```json
  {
    "data": [
      {"id": "gpt-4o"},
      {"id": "gpt-3.5-turbo"}
    ]
  }
  ```
- Map to `[]Model`.

### 3. Implement Anthropic Discovery
Anthropic recently added a `/v1/models` endpoint.
- Header: `x-api-key: <apiKey>`, `anthropic-version: 2023-06-01`.
- Parse response and map to `[]Model`.

### 4. Handle Offline/Custom Gracefully
For `Ollama` or `Custom` providers where `/v1/models` might not be standardized or available, return an empty list or a specific `ErrDiscoveryUnsupported` error, so the UI knows to fallback directly to a manual text input field.

### 5. CLI Hook
In `cmd/models.go`, hook this up so `go run main.go models --provider openrouter --api-key sk-123` prints the live list of models to the console.

## Manual Testing
By the end of this session, perform the following manual tests:
- Run `go run main.go models ...` with a real API key (temporarily) to verify live network fetching works.
- Verify that if the API key is incorrect or absent, a sensible error message is printed.
