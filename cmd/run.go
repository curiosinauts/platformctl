package cmd

import (
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Runs",
	Long:  `Runs`,
}

func init() {
	rootCmd.AddCommand(runCmd)
}
