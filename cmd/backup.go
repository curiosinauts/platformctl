package cmd

import (
	"github.com/spf13/cobra"
)

// backupCmd represents the backup command
var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "backup",
	Long:  `backup`,
}

func init() {
	rootCmd.AddCommand(backupCmd)
}
