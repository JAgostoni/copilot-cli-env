package render

import (
	"fmt"
	"strings"

	"github.com/jagostoni/copilot-cli-env/internal/config"
)

type envFileRenderer struct{}

func NewEnvFileRenderer() Renderer {
	return &envFileRenderer{}
}

func (r *envFileRenderer) Render(env config.CopilotEnv) (string, error) {
	if env.IsReset {
		return "", nil
	}

	var builder strings.Builder

	if env.ProviderBaseURL != "" {
		builder.WriteString(fmt.Sprintf("COPILOT_PROVIDER_BASE_URL=%s\n", env.ProviderBaseURL))
	}
	if env.Model != "" {
		builder.WriteString(fmt.Sprintf("COPILOT_MODEL=%s\n", env.Model))
	}
	if env.ProviderType != "" {
		builder.WriteString(fmt.Sprintf("COPILOT_PROVIDER_TYPE=%s\n", env.ProviderType))
	}
	if env.APIKey != "" {
		builder.WriteString(fmt.Sprintf("COPILOT_PROVIDER_API_KEY=%s\n", env.APIKey))
	}
	if env.Offline {
		builder.WriteString("COPILOT_OFFLINE=true\n")
	}
	if env.WireAPI != "" {
		builder.WriteString(fmt.Sprintf("COPILOT_PROVIDER_WIRE_API=%s\n", env.WireAPI))
	}
	if env.MaxPromptTokens > 0 {
		builder.WriteString(fmt.Sprintf("COPILOT_PROVIDER_MAX_PROMPT_TOKENS=%d\n", env.MaxPromptTokens))
	}
	if env.MaxOutputTokens > 0 {
		builder.WriteString(fmt.Sprintf("COPILOT_PROVIDER_MAX_OUTPUT_TOKENS=%d\n", env.MaxOutputTokens))
	}

	return builder.String(), nil
}
