package cmd

import (
	"github.com/spf13/cobra"
)

// enablePiholeCmd represents the pihole command
var enablePiholeCmd = &cobra.Command{
	Use:   "pihole",
	Short: "Enables pihole blocking",
	Long:  `Enables pihole blocking`,
	Run: func(cmd *cobra.Command, args []string) {
		sshCmd.Run(sshCmd, []string{"debian@192.168.0.110:22", "/usr/local/bin/pihole enable"})
	},
}

func init() {
	enableCmd.AddCommand(enablePiholeCmd)
}
