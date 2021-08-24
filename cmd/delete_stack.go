package cmd

import (
	"github.com/curiosinauts/platformctl/internal/msg"
	"github.com/curiosinauts/platformctl/pkg/executil"
	"github.com/spf13/cobra"
)

// deleteStackCmd represents the stack command
var deleteStackCmd = &cobra.Command{
	Use:     "stack",
	Short:   "Deletes stack",
	Long:    `Deletes stack`,
	PreRunE: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		username := args[0]

		eh := ErrorHandler{"deleting stack"}
		output, err := executil.Execute("kubectl delete ingress vscode-"+username, debug)
		eh.HandleErrorWithOutput("deleting ingress", err, output)

		output, err = executil.Execute("kubectl delete service vscode-"+username, debug)
		eh.HandleErrorWithOutput("deleting service", err, output)

		output, err = executil.Execute("kubectl delete deployment vscode-"+username, debug)
		eh.HandleErrorWithOutput("deleting deployment", err, output)

		msg.Success("deleting stack")
	},
}

func init() {
	deleteCmd.AddCommand(deleteStackCmd)
}
