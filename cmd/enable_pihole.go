package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// enablePiholeCmd represents the pihole command
var enablePiholeCmd = &cobra.Command{
	Use:   "pihole",
	Short: "Enables pihole blocking",
	Long:  `Enables pihole blocking`,
	Run: func(cmd *cobra.Command, args []string) {
		piholeHost := viper.Get("pihole_host").(string)
		sshCmd.Run(sshCmd, []string{"debian@" + piholeHost, "/usr/local/bin/pihole enable"})
	},
}

func init() {
	enableCmd.AddCommand(enablePiholeCmd)
}
