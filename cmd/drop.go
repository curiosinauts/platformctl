package cmd

import (
	"github.com/spf13/cobra"
)

// dropCmd represents the drop command
var dropCmd = &cobra.Command{
	Use:   "drop",
	Short: "Drop",
	Long:  `Drop`,
}

func init() {
	rootCmd.AddCommand(dropCmd)
}
