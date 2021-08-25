package cmd

import (
	"github.com/curiosinauts/platformctl/pkg/sshutil"
	"github.com/spf13/cobra"
)

// sshRemoteCmd represents the remote command
var sshRemoteCmd = &cobra.Command{
	Use:   "remote",
	Short: "Executes remote script over SSH",
	Long:  `Executes remote script over SSH`,
	// RunE:  cobra.MinimumNArgs(4),
	Run: func(cmd *cobra.Command, args []string) {
		host := args[0]
		port := args[1]
		user := args[2]
		script := args[3]
		sshutil.RemoteSSHExec(host, port, user, script)
	},
}

func init() {
	sshCmd.AddCommand(sshRemoteCmd)
}
