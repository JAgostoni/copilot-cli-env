package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jagostoni/copilot-cli-env/internal/ui"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Interactive onboarding",
	Run: func(cmd *cobra.Command, args []string) {
		app := ui.InitialModel()
		p := tea.NewProgram(app)
		if _, err := p.Run(); err != nil {
			fmt.Printf("Alas, there's been an error: %v", err)
			os.Exit(1)
		}

		if app.IsCancelled() {
			os.Exit(0)
		}

		executeApply(app.ConfigData(), app.OutputMode(), app.DetectedOS(), app.DetectedShell())
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
