package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// disablePiholeCmd represents the pihole command
var disablePiholeCmd = &cobra.Command{
	Use:   "pihole",
	Short: "Disables pihole blocking",
	Long:  `Disables pihole blocking`,
	Run: func(cmd *cobra.Command, args []string) {
		piholeHost := viper.Get("pihole_host").(string)
		sshCmd.Run(sshCmd, []string{"debian@" + piholeHost, "/usr/local/bin/pihole disable"})
	},
}

func init() {
	disableCmd.AddCommand(disablePiholeCmd)
}
