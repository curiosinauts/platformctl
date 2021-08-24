package cmd

import (
	"github.com/curiosinauts/platformctl/internal/msg"
	"github.com/curiosinauts/platformctl/pkg/giteautil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// eleteGiteaUserCmd represents the gituser command
var deleteGiteaUserCmd = &cobra.Command{
	Use:     "gitea-user",
	Short:   "Delete gitea user account",
	Long:    `Delete gitea user account`,
	PreRunE: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		username := args[0]

		eh := ErrorHandler{"deleting gitea user account"}

		accessToken := viper.Get("gitea_access_token").(string)
		giteaURL := viper.Get("gitea_url").(string)
		gitClient, err := giteautil.NewGitClient(accessToken, giteaURL)
		eh.PrintError("instantiating git client", err)

		err = gitClient.DeleteUserRepo(username)
		eh.PrintError("deleting user repos from gitea", err)

		// err = gitClient.DeleteUserPublicKey(user, user.PublicKeyID)
		// eh.PrintError("deleting user public key from gitea", err)

		err = gitClient.RemoveUser(username)
		eh.PrintError("removing user from gitea", err)

		msg.Success("deleting gitea user account")
	},
}

func init() {
	deleteCmd.AddCommand(deleteGiteaUserCmd)
}
