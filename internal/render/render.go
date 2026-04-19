package render

import (
	"strings"

	"github.com/jagostoni/copilot-cli-env/internal/config"
)

type OutputMode string

const (
	ModeEnvFile OutputMode = "env"
	ModeConsole OutputMode = "console" // direct dump
	ModeCommand OutputMode = "command" // shell-specific commands
)

type Renderer interface {
	Render(env config.CopilotEnv) (string, error)
}

// MaskSecret obfuscates a secret string, revealing only a prefix and appending asterisks.
func MaskSecret(secret string) string {
	if len(secret) <= 4 {
		return strings.Repeat("*", len(secret))
	}

	prefixLen := 4
	if strings.HasPrefix(secret, "sk-") {
		if len(secret) >= 6 {
			prefixLen = 6
		} else {
			prefixLen = len(secret)
		}
	}

	return secret[:prefixLen] + "..." + strings.Repeat("*", 4)
}
