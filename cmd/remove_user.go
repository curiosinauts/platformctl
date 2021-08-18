package cmd

import (
	"fmt"

	"github.com/spf13/viper"

	"github.com/curiosinauts/platformctl/internal/msg"
	"github.com/curiosinauts/platformctl/pkg/crypto"
	"github.com/curiosinauts/platformctl/pkg/executil"
	"github.com/curiosinauts/platformctl/pkg/giteautil"

	"github.com/curiosinauts/platformctl/internal/database"
	"github.com/spf13/cobra"
)

var removeUserCmdDebug bool

// removeUserCmd represents the user command
var removeUserCmd = &cobra.Command{
	Use:   "user",
	Short: "Removes user from the platform",
	Long:  `Removes user from the platform`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println()

		email := args[0]
		userService := database.NewUserService(db)

		eh := ErrorHandler{"removing user"}

		user, dberr := userService.FindUserByHashedEmail(crypto.Hashed(email))

		dberr = userService.DeleteALLIDERuntimeInstallsForUser(user.ID)
		eh.HandleError("delete user ide runtime installs", dberr)

		dberr = userService.DeleteALLIDEReposForUser(user.ID)
		eh.HandleError("delete user ide repos", dberr)

		dberr = userService.DeleteALLUserIDEsForUser(user.ID)
		eh.HandleError("delete user ides", dberr)

		dberr = userService.DeleteALLUserReposForUser(user.ID)
		eh.HandleError("delete user repos", dberr)

		dberr = userService.Delete(user.ID)
		eh.HandleError("delete user", dberr)

		accessToken := viper.Get("gitea_access_token").(string)
		giteaURL := viper.Get("gitea_url").(string)
		gitClient, err := giteautil.NewGitClient(accessToken, giteaURL)
		eh.HandleError("instantiating git client", err)

		err = gitClient.DeleteUserRepo(user.Username)
		eh.HandleError("deleting user repos from gitea", err)

		err = gitClient.DeleteUserPublicKey(user, user.PublicKeyID)
		eh.HandleError("deleting user public key from gitea", err)

		err = gitClient.RemoveUser(user.Username)
		eh.HandleError("removing user from gitea", err)

		executil.Execute("kubectl delete ingress vscode-"+user.Username, removeUserCmdDebug)
		executil.Execute("kubectl delete service vscode-"+user.Username, removeUserCmdDebug)
		executil.Execute("kubectl delete deployment vscode-"+user.Username, removeUserCmdDebug)

		msg.Success("removing user")
	},
}

func init() {
	removeCmd.AddCommand(removeUserCmd)
	removeCmd.Flags().BoolVarP(&removeUserCmdDebug, "debug", "d", false, "Debug this command")
}
