package cmd

import (
	"github.com/spf13/cobra"
)

// beforeCmd represents the before command
var beforeCmd = &cobra.Command{
	Use:   "before",
	Short: "Before",
	Long:  `Before`,
}

func init() {
	rootCmd.AddCommand(beforeCmd)
}
