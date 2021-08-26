package cmd

import (
	"fmt"
	"github.com/curiosinauts/platformctl/internal/msg"
	"github.com/curiosinauts/platformctl/pkg/sshutil"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

// sshCmd represents the remote command
var sshCmd = &cobra.Command{
	Use:     "ssh",
	Short:   "Executes remote script over SSH",
	Long:    `Executes remote script over SSH`,
	Example: `platformctl ssh debian@192.168.0.107:22 "sudo rm -rf /tmp"`,
	//RunE:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		s := args[0]
		script := args[1]

		if !strings.Contains(s, "@") {
			msg.Failure("missing user")
			os.Exit(1)
		}

		ts := strings.Split(s, "@")
		if len(ts) < 2 {
			msg.Failure("invalid format")
			os.Exit(1)
		}

		user := ts[0]

		port := "22"
		host := ts[1]

		if strings.Contains(ts[1], ":") {
			terms := strings.Split(ts[1], ":")
			host = terms[0]
			port = terms[1]
		}

		out, err := sshutil.RemoteSSHExec(host, port, user, script)
		if err != nil {
			msg.Failure(err.Error())
		}

		fmt.Println(out)
	},
}

func init() {
	rootCmd.AddCommand(sshCmd)
	// platformctl ssh debian@192.168.0.107:22 "sudo rm -rf /foo"
}
