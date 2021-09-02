package cmd

import (
	"github.com/spf13/cobra"
)

// nextCmd represents the next command
var nextCmd = &cobra.Command{
	Use:   "next",
	Short: "Next",
	Long:  `Next`,
}

func init() {
	rootCmd.AddCommand(nextCmd)
}
