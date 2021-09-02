package cmd

import (
	"github.com/spf13/cobra"
)

// afterCmd represents the after command
var afterCmd = &cobra.Command{
	Use:   "after",
	Short: "After",
	Long:  `After`,
}

func init() {
	rootCmd.AddCommand(afterCmd)
}
