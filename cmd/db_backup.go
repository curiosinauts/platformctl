package cmd

import (
	"github.com/spf13/cobra"
)

// dbBackupCmd represents the backup command
var dbBackupCmd = &cobra.Command{
	Use:   "backup",
	Short: "backup",
	Long:  `backup`,
}

func init() {
	dbCmd.AddCommand(dbBackupCmd)
}
