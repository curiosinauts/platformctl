package cmd

import (
	"github.com/spf13/cobra"
)

// describeCmd represents the describe command
var describeCmd = &cobra.Command{
	Use:   "describe",
	Short: "Describe",
	Long:  `Describe`,
}

func init() {
	rootCmd.AddCommand(describeCmd)
}
