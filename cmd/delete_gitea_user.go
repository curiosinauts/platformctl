package cmd

import (
	"github.com/curiosinauts/platformctl/internal/msg"
	"github.com/curiosinauts/platformctl/pkg/giteautil"
	"github.com/spf13/cobra"
)

// deleteGiteaUserCmd represents the gituser command
var deleteGiteaUserCmd = &cobra.Command{
	Use:     "gitea-user",
	Short:   "Delete gitea user account",
	Long:    `Delete gitea user account`,
	PreRunE: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		username := args[0]

		eh := ErrorHandler{"deleting gitea user account"}

		gitClient, err := giteautil.NewGitClient()
		eh.PrintError("instantiating git client", err)

		err = gitClient.DeleteUserRepo(username)
		eh.PrintError("deleting user repos from gitea", err)

		err = gitClient.RemoveUser(username)
		eh.PrintError("removing user from gitea", err)

		msg.Success("deleting gitea user account")
	},
}

func init() {
	deleteCmd.AddCommand(deleteGiteaUserCmd)
}
