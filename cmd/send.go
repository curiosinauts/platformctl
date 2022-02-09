package cmd

import (
	"github.com/spf13/cobra"
)

// sendCmd represents the send command
var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Sends",
	Long:  `Sends`,
}

func init() {
	rootCmd.AddCommand(sendCmd)
}
