package cmd

import (
	"fmt"

	"github.com/jagostoni/copilot-cli-env/internal/envdetect"
	"github.com/spf13/cobra"
)

var detectCmd = &cobra.Command{
	Use:   "detect",
	Short: "Environment detection test command",
	Run: func(cmd *cobra.Command, args []string) {
		detector := envdetect.NewDetector()
		env, err := detector.Detect()
		if err != nil {
			fmt.Printf("Error detecting environment: %v\n", err)
			return
		}
		fmt.Printf("Detected OS: %s\n", env.OS)
		fmt.Printf("Detected Shell: %s\n", env.Shell)
		if env.Terminal != "" {
			fmt.Printf("Detected Terminal: %s\n", env.Terminal)
		}
	},
}

func init() {
	rootCmd.AddCommand(detectCmd)
}
