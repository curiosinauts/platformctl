package cmd

import (
	"github.com/spf13/cobra"
)

// disableCmd represents the disable command
var disableCmd = &cobra.Command{
	Use:   "disable",
	Short: "Disable",
	Long:  `Disable`,
}

func init() {
	rootCmd.AddCommand(disableCmd)
}
