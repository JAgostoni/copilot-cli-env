package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/jagostoni/copilot-cli-env/internal/provider"
	"github.com/spf13/cobra"
)

var (
	modelsProvider string
	modelsBaseUrl  string
	modelsApiKey   string
)

var modelsCmd = &cobra.Command{
	Use:   "models",
	Short: "Model discovery command",
	Run: func(cmd *cobra.Command, args []string) {
		var p provider.Provider

		switch modelsProvider {
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
		default:
			fmt.Printf("Unknown provider: %s\n", modelsProvider)
			os.Exit(1)
		}

		ctx := context.Background()
		models, err := p.FetchModels(ctx, modelsBaseUrl, modelsApiKey)
		if err != nil {
			if err == provider.ErrDiscoveryUnsupported {
				fmt.Printf("Model discovery is not supported for provider %s\n", p.Name())
				return
			}
			fmt.Printf("Failed to fetch models: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Available models for %s:\n", p.Name())
		for _, m := range models {
			if m.Description != "" {
				fmt.Printf(" - %s (%s)\n", m.ID, m.Description)
			} else {
				fmt.Printf(" - %s\n", m.ID)
			}
		}
	},
}

func init() {
	modelsCmd.Flags().StringVar(&modelsProvider, "provider", "", "Provider type (e.g., openai, openrouter, anthropic)")
	modelsCmd.Flags().StringVar(&modelsBaseUrl, "base-url", "", "Provider base URL")
	modelsCmd.Flags().StringVar(&modelsApiKey, "api-key", "", "Provider API Key")
	modelsCmd.MarkFlagRequired("provider")
	
	rootCmd.AddCommand(modelsCmd)
}
