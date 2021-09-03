package cmd

import (
	"github.com/spf13/cobra"
)

// disablePiholeCmd represents the pihole command
var disablePiholeCmd = &cobra.Command{
	Use:   "pihole",
	Short: "Disables pihole blocking",
	Long:  `Disables pihole blocking`,
	Run: func(cmd *cobra.Command, args []string) {
		sshCmd.Run(sshCmd, []string{"debian@192.168.0.110:22", "/usr/local/bin/pihole disable"})
	},
}

func init() {
	disableCmd.AddCommand(disablePiholeCmd)
}
