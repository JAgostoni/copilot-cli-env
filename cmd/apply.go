package cmd

import (
	"fmt"
	"os"

	"github.com/jagostoni/copilot-cli-env/internal/config"
	"github.com/jagostoni/copilot-cli-env/internal/profile"
	"github.com/jagostoni/copilot-cli-env/internal/render"
	"github.com/spf13/cobra"
)

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Persistent profile application",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Use 'configure' or 'init' to apply settings.")
	},
}

func init() {
	rootCmd.AddCommand(applyCmd)
}

func executeApply(env config.CopilotEnv, mode string, detectedOS, detectedShell string) {
	if mode == "console" || mode == "command" || mode == "bash" || mode == "powershell" || mode == "cmd" {
		shellMode := detectedShell
		if mode == "bash" || mode == "powershell" || mode == "cmd" {
			shellMode = mode
		}
		renderer := render.NewShellCommandRenderer(shellMode)
		out, err := renderer.Render(env)
		if err != nil {
			fmt.Printf("Error rendering: %v\n", err)
			os.Exit(1)
		}
		fmt.Print(out)
		fmt.Println("\n# Please copy the commands above and paste them into your shell terminal to set the variables.")
		return
	}

	if mode == "env" {
		renderer := render.NewEnvFileRenderer()
		out, err := renderer.Render(env)
		if err != nil {
			fmt.Printf("Error rendering .env: %v\n", err)
			os.Exit(1)
		}
		err = profile.UpdateEnvFile(".env", out)
		if err != nil {
			fmt.Printf("Error writing .env: %v\n", err)
			os.Exit(1)
		}
		if env.IsReset {
			fmt.Println("Successfully removed Copilot configurations from .env")
		} else {
			fmt.Println("Successfully updated .env file")
			fmt.Println("Ensure your environment loads this file (e.g. 'export $(cat .env | xargs)').")
		}
		return
	}

	if mode == "profile" {
		profilePath, err := profile.ResolveProfilePath(detectedOS, detectedShell)
		if err != nil {
			fmt.Printf("Failed to resolve profile path: %v\n", err)
			os.Exit(1)
		}

		renderer := render.NewShellCommandRenderer(detectedShell)
		out, err := renderer.Render(env)
		if err != nil {
			fmt.Printf("Error rendering profile block: %v\n", err)
			os.Exit(1)
		}

		err = profile.UpdateProfile(profilePath, out)
		if err != nil {
			fmt.Printf("Failed to update profile %s: %v\n", profilePath, err)
			os.Exit(1)
		}
		
		if env.IsReset {
			fmt.Printf("Successfully removed configurations from %s\n", profilePath)
		} else {
			fmt.Printf("Successfully updated profile at %s\n", profilePath)
		}
		
		if detectedShell == "powershell" || detectedShell == "pwsh" {
			fmt.Printf("\nPlease restart your shell or run: . %s\n", profilePath)
		} else if detectedShell == "cmd" {
			fmt.Printf("\nPlease restart your command prompt.\n")
		} else {
			fmt.Printf("\nPlease restart your shell or run: source %s\n", profilePath)
		}
		return
	}
	
	fmt.Printf("Unknown output mode: %s\n", mode)
	os.Exit(1)
}
