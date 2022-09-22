package cmd

import (
	"github.com/spf13/cobra"
)

// dbRestoreCmd represents the restore command
var dbRestoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "restore",
	Long:  `restore`,
}

func init() {
	dbCmd.AddCommand(dbRestoreCmd)
}
