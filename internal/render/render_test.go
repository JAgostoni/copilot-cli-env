package render

import (
	"strings"
	"testing"

	"github.com/jagostoni/copilot-cli-env/internal/config"
)

func TestMaskSecret(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"sk-1234567890", "sk-123...****"},
		{"abcdef", "abcd...****"},
		{"abc", "***"},
		{"sk-12", "sk-12...****"},
		{"1234", "****"},
	}

	for _, tt := range tests {
		if got := MaskSecret(tt.input); got != tt.expected {
			t.Errorf("MaskSecret(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

func TestEnvFileRenderer(t *testing.T) {
	env := config.CopilotEnv{
		ProviderBaseURL: "https://api.example.com",
		Model:           "gpt-4",
		ProviderType:    "openai",
		APIKey:          "secret-key",
		Offline:         true,
	}

	renderer := NewEnvFileRenderer()
	out, err := renderer.Render(env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedLines := []string{
		"COPILOT_PROVIDER_BASE_URL=https://api.example.com",
		"COPILOT_MODEL=gpt-4",
		"COPILOT_PROVIDER_TYPE=openai",
		"COPILOT_PROVIDER_API_KEY=secret-key",
		"COPILOT_OFFLINE=true",
	}

	for _, line := range expectedLines {
		if !strings.Contains(out, line) {
			t.Errorf("expected output to contain %q, got:\n%s", line, out)
		}
	}
}

func TestShellCommandRenderer(t *testing.T) {
	env := config.CopilotEnv{
		ProviderBaseURL: "https://api.example.com",
		Model:           "gpt-4\"quote\"",
		ProviderType:    "openai",
		APIKey:          "secret-key",
		Offline:         true,
	}

	tests := []struct {
		shell         string
		expectedLines []string
	}{
		{
			shell: "bash",
			expectedLines: []string{
				"export COPILOT_PROVIDER_BASE_URL=\"https://api.example.com\"",
				"export COPILOT_MODEL=\"gpt-4\\\"quote\\\"\"",
				"export COPILOT_PROVIDER_TYPE=\"openai\"",
				"export COPILOT_PROVIDER_API_KEY=\"secret-key\"",
				"export COPILOT_OFFLINE=\"true\"",
			},
		},
		{
			shell: "powershell",
			expectedLines: []string{
				"$env:COPILOT_PROVIDER_BASE_URL=\"https://api.example.com\"",
				"$env:COPILOT_MODEL=\"gpt-4`\"quote`\"\"",
				"$env:COPILOT_PROVIDER_TYPE=\"openai\"",
				"$env:COPILOT_PROVIDER_API_KEY=\"secret-key\"",
				"$env:COPILOT_OFFLINE=\"true\"",
			},
		},
		{
			shell: "cmd",
			expectedLines: []string{
				"set COPILOT_PROVIDER_BASE_URL=https://api.example.com",
				"set COPILOT_MODEL=gpt-4\"quote\"",
				"set COPILOT_PROVIDER_TYPE=openai",
				"set COPILOT_PROVIDER_API_KEY=secret-key",
				"set COPILOT_OFFLINE=true",
			},
		},
	}

	for _, tt := range tests {
		renderer := NewShellCommandRenderer(tt.shell)
		out, err := renderer.Render(env)
		if err != nil {
			t.Fatalf("shell %s: unexpected error: %v", tt.shell, err)
		}

		for _, line := range tt.expectedLines {
			if !strings.Contains(out, line) {
				t.Errorf("shell %s: expected output to contain %q, got:\n%s", tt.shell, line, out)
			}
		}
	}
}
