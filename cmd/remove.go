package cmd

import (
	"github.com/spf13/cobra"
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Removes entity from the platfrom",
	Long:  `Removes entity from the platfrom`,
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
