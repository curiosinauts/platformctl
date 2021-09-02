package cmd

import (
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds entity to the platform",
	Long:  `Adds entity to the platform`,
}

func init() {
	rootCmd.AddCommand(addCmd)
}
