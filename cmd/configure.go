package cmd

import (
	"fmt"
	"os"

	"github.com/jagostoni/copilot-cli-env/internal/envdetect"
	"github.com/jagostoni/copilot-cli-env/internal/provider"
	"github.com/spf13/cobra"
)

var (
	configureProvider string
	model             string
	baseUrl           string
	apiKey            string
	outputFlag        string
	maxPromptTokens   int
	maxOutputTokens   int
)

var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Non-interactive configuration",
	Run: func(cmd *cobra.Command, args []string) {
		var p provider.Provider

		switch configureProvider {
		case "openai":
			p = provider.NewOpenAIProvider()
		case "anthropic":
			p = provider.NewAnthropicProvider()
		case "azure":
			p = provider.NewAzureProvider()
		case "openrouter":
			p = provider.NewOpenRouterProvider()
		case "ollama":
			p = provider.NewOllamaProvider()
		case "reset":
			p = provider.NewResetProvider()
		default:
			fmt.Printf("Unknown provider: %s\n", configureProvider)
			os.Exit(1)
		}

		if configureProvider != "reset" && model == "" {
			fmt.Println("Error: required flag(s) \"model\" not set")
			os.Exit(1)
		}

		env, err := p.MapEnv(baseUrl, model, apiKey)
		if err != nil {
			fmt.Printf("Configuration error: %v\n", err)
			os.Exit(1)
		}
		
		env.MaxPromptTokens = maxPromptTokens
		env.MaxOutputTokens = maxOutputTokens

		mode := outputFlag
		if mode == "" {
			mode = "console"
		}

		detector := envdetect.NewDetector()
		detectedEnv, err := detector.Detect()
		if err != nil {
			fmt.Printf("Warning: failed to detect environment: %v\n", err)
		}

		executeApply(env, mode, detectedEnv.OS, detectedEnv.Shell)
	},
}

func init() {
	configureCmd.Flags().StringVar(&configureProvider, "provider", "", "Provider type (e.g., openai)")
	configureCmd.Flags().StringVar(&model, "model", "", "Model name (e.g., gpt-4o)")
	configureCmd.Flags().StringVar(&baseUrl, "base-url", "", "Provider base URL")
	configureCmd.Flags().StringVar(&apiKey, "api-key", "", "Provider API Key")
	configureCmd.Flags().StringVar(&outputFlag, "output", "", "Output format (e.g., env, profile, console, bash, powershell, cmd)")
	configureCmd.Flags().IntVar(&maxPromptTokens, "max-prompt-tokens", 0, "Maximum prompt tokens")
	configureCmd.Flags().IntVar(&maxOutputTokens, "max-output-tokens", 0, "Maximum output tokens")
	
	configureCmd.MarkFlagRequired("provider")

	rootCmd.AddCommand(configureCmd)
}
