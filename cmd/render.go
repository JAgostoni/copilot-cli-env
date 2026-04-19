package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var renderCmd = &cobra.Command{
	Use:   "render",
	Short: "Output generation (dry run)",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("render command called")
	},
}

func init() {
	rootCmd.AddCommand(renderCmd)
}
