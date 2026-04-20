package config

// CopilotEnv represents the final set of environment variables needed.
type CopilotEnv struct {
	ProviderBaseURL string // COPILOT_PROVIDER_BASE_URL (Required)
	Model           string // COPILOT_MODEL (Required)
	ProviderType    string // COPILOT_PROVIDER_TYPE (Optional: 'openai', 'azure', 'anthropic')
	APIKey          string // COPILOT_PROVIDER_API_KEY (Optional for some endpoints)
	Offline         bool   // COPILOT_OFFLINE (Optional: usually 'true' for BYOK)
	IsReset         bool   // Indicates if this is a configuration reset action
	MaxPromptTokens int    // COPILOT_PROVIDER_MAX_PROMPT_TOKENS
	MaxOutputTokens int    // COPILOT_PROVIDER_MAX_OUTPUT_TOKENS
	WireAPI         string // COPILOT_PROVIDER_WIRE_API (Optional: 'completions' or 'responses')
}
