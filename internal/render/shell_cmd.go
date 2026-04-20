package render

import (
	"fmt"
	"strings"

	"github.com/jagostoni/copilot-cli-env/internal/config"
)

type shellCommandRenderer struct {
	shell string
}

func NewShellCommandRenderer(shell string) Renderer {
	return &shellCommandRenderer{
		shell: shell,
	}
}

func escapeBash(val string) string {
	// Simple escaping for demonstration; handles quotes
	return strings.ReplaceAll(val, "\"", "\\\"")
}

func escapePowerShell(val string) string {
	return strings.ReplaceAll(val, "\"", "`\"")
}

func (r *shellCommandRenderer) Render(env config.CopilotEnv) (string, error) {
	var builder strings.Builder

	if env.IsReset {
		vars := []string{
			"COPILOT_PROVIDER_BASE_URL",
			"COPILOT_MODEL",
			"COPILOT_PROVIDER_TYPE",
			"COPILOT_PROVIDER_API_KEY",
			"COPILOT_OFFLINE",
			"COPILOT_PROVIDER_MAX_PROMPT_TOKENS",
			"COPILOT_PROVIDER_MAX_OUTPUT_TOKENS",
			"COPILOT_PROVIDER_WIRE_API",
		}
		for _, k := range vars {
			switch r.shell {
			case "powershell", "pwsh":
				builder.WriteString(fmt.Sprintf("Remove-Item Env:\\%s -ErrorAction SilentlyContinue\n", k))
			case "cmd":
				builder.WriteString(fmt.Sprintf("set %s=\n", k))
			default: // bash, zsh, fish (for simplicity, using export for all unix-like)
				builder.WriteString(fmt.Sprintf("unset %s\n", k))
			}
		}
		return builder.String(), nil
	}

	vars := []struct {
		Key   string
		Value string
	}{
		{"COPILOT_PROVIDER_BASE_URL", env.ProviderBaseURL},
		{"COPILOT_MODEL", env.Model},
		{"COPILOT_PROVIDER_TYPE", env.ProviderType},
		{"COPILOT_PROVIDER_API_KEY", env.APIKey},
	}

	if env.Offline {
		vars = append(vars, struct {
			Key   string
			Value string
		}{"COPILOT_OFFLINE", "true"})
	}

	if env.WireAPI != "" {
		vars = append(vars, struct {
			Key   string
			Value string
		}{"COPILOT_PROVIDER_WIRE_API", env.WireAPI})
	}

	if env.MaxPromptTokens > 0 {
		vars = append(vars, struct {
			Key   string
			Value string
		}{"COPILOT_PROVIDER_MAX_PROMPT_TOKENS", fmt.Sprintf("%d", env.MaxPromptTokens)})
	}

	if env.MaxOutputTokens > 0 {
		vars = append(vars, struct {
			Key   string
			Value string
		}{"COPILOT_PROVIDER_MAX_OUTPUT_TOKENS", fmt.Sprintf("%d", env.MaxOutputTokens)})
	}

	for _, v := range vars {
		if v.Value == "" {
			continue
		}

		switch r.shell {
		case "powershell", "pwsh":
			builder.WriteString(fmt.Sprintf("$env:%s=\"%s\"\n", v.Key, escapePowerShell(v.Value)))
		case "cmd":
			builder.WriteString(fmt.Sprintf("set %s=%s\n", v.Key, v.Value))
		default: // bash, zsh, fish (for simplicity, using export for all unix-like)
			builder.WriteString(fmt.Sprintf("export %s=\"%s\"\n", v.Key, escapeBash(v.Value)))
		}
	}

	return builder.String(), nil
}
