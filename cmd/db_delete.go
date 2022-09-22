package cmd

import (
	"github.com/spf13/cobra"
)

// backupCmd represents the backup command
var dbDeleteCmd = &cobra.Command{
	Use:   "backup",
	Short: "backup",
	Long:  `backup`,
}

func init() {
	dbCmd.AddCommand(dbDeleteCmd)
}
