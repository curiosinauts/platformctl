package cmd

import (
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create",
	Long:  `Create`,
}

func init() {
	rootCmd.AddCommand(createCmd)
}
