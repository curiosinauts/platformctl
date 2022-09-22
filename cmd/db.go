package cmd

import (
	"github.com/spf13/cobra"
)

// dbCmd represents the db command
var dbCmd = &cobra.Command{
	Use:   "backup",
	Short: "backup",
	Long:  `backup`,
}

func init() {
	rootCmd.AddCommand(dbCmd)
}
