package cmd

import (
	"github.com/spf13/viper"

	"github.com/curiosinauts/platformctl/internal/msg"
	"github.com/curiosinauts/platformctl/pkg/crypto"
	"github.com/curiosinauts/platformctl/pkg/executil"
	"github.com/curiosinauts/platformctl/pkg/giteautil"
	"github.com/curiosinauts/platformctl/pkg/regutil"

	"github.com/curiosinauts/platformctl/pkg/database"
	"github.com/spf13/cobra"
)

// removeUserCmd represents the user command
var removeUserCmd = &cobra.Command{
	Use:     "user",
	Short:   "Removes user from the platform",
	Long:    `Removes user from the platform`,
	PreRunE: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		email := args[0]
		userService := database.NewUserService(db)

		eh := ErrorHandler{"removing user"}

		user, dberr := userService.FindUserByHashedEmail(crypto.Hashed(email))

		accessToken := viper.Get("gitea_access_token").(string)
		giteaURL := viper.Get("gitea_url").(string)
		gitClient, err := giteautil.NewGitClient(accessToken, giteaURL)
		eh.PrintError("instantiating git client", err)

		err = gitClient.DeleteUserRepo(user.Username)
		eh.PrintError("deleting user repos from gitea", err)

		err = gitClient.DeleteUserPublicKey(user, user.PublicKeyID)
		eh.PrintError("deleting user public key from gitea", err)

		err = gitClient.RemoveUser(user.Username)
		eh.PrintError("removing user from gitea", err)

		_, err = executil.Execute("kubectl delete ingress vscode-"+user.Username, debug)
		eh.PrintError("deleting ingress", err)

		_, err = executil.Execute("kubectl delete service vscode-"+user.Username, debug)
		eh.PrintError("deleting service", err)

		_, err = executil.Execute("kubectl delete deployment vscode-"+user.Username, debug)
		eh.PrintError("deleting deployment", err)

		repository := "7onetella/vscode-" + user.Username
		tags, err := regutil.ListTags(repository, false)
		eh.PrintError("listing tags", err)

		for _, tag := range tags {
			msg.Info("deleting tag " + tag)
			err = regutil.DeleteImage(repository, tag, false)
			eh.PrintError("deleting image", err)
		}

		dberr = userService.DeleteALLIDERuntimeInstallsForUser(user.ID)
		eh.PrintError("delete user ide runtime installs", dberr)

		dberr = userService.DeleteALLIDEReposForUser(user.ID)
		eh.PrintError("delete user ide repos", dberr)

		dberr = userService.DeleteALLUserIDEsForUser(user.ID)
		eh.PrintError("delete user ides", dberr)

		dberr = userService.DeleteALLUserReposForUser(user.ID)
		eh.PrintError("delete user repos", dberr)

		dberr = userService.Delete(user.ID)
		eh.PrintError("delete user", dberr)

		msg.Success("removing user")
	},
}

func init() {
	removeCmd.AddCommand(removeUserCmd)
}
