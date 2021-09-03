package cmd

import (
	"github.com/spf13/cobra"
)

// enableCmd represents the enable command
var enableCmd = &cobra.Command{
	Use:   "enable",
	Short: "Enable",
	Long:  `Enable`,
}

func init() {
	rootCmd.AddCommand(enableCmd)
}
