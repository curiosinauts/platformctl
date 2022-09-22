package cmd

import (
	"github.com/spf13/cobra"
)

// registryCmd represents the registry command
var registryCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete",
	Long:  `Delete`,
}

func init() {
	rootCmd.AddCommand(registryCmd)
}
