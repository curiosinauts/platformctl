package cmd

import (
	"github.com/spf13/cobra"
)

// dbCreateCmd represents the create command
var dbCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "create",
	Long:  `create`,
}

func init() {
	dbCmd.AddCommand(dbCreateCmd)
}
