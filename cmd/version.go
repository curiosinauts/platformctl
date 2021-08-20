package cmd

import (
	"github.com/curiosinauts/platformctl/internal/msg"
	"github.com/spf13/cobra"
)

var version = "0.1.6"

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints version",
	Long:  `Prints version`,
	Run: func(cmd *cobra.Command, args []string) {
		msg.Info(version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
