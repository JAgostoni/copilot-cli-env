package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "copilot-cli-env",
	Short: "A tool to manage GitHub Copilot CLI environments",
	Long:  `Copilot Env Manager is a CLI tool designed to help you configure and manage environments for GitHub Copilot CLI.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Default to running the init interactive flow if no subcommands are passed
		initCmd.Run(cmd, args)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
