package cmd

import (
	"github.com/curiosinauts/platformctl/internal/msg"
	"github.com/curiosinauts/platformctl/pkg/crypto"
	"github.com/curiosinauts/platformctl/pkg/database"
	"github.com/curiosinauts/platformctl/pkg/executil"
	"github.com/curiosinauts/platformctl/pkg/giteautil"
	"github.com/curiosinauts/platformctl/pkg/regutil"
	"github.com/curiosinauts/platformctl/pkg/sshutil"

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

		eh := ErrorHandler{"removing user"}

		user := database.User{}
		dberr := dbs.FindBy(&user, "hashed_email=$1", crypto.Hashed(email))

		gitClient, err := giteautil.NewGitClient()
		eh.PrintError("instantiating git client", err)

		err = gitClient.DeleteUserRepo(user.Username)
		eh.PrintError("deleting user repos from gitea", err)

		err = gitClient.DeleteUserPublicKey(user, user.PublicKeyID)
		eh.PrintError("deleting user public key from gitea", err)

		err = gitClient.RemoveUser(user.Username)
		eh.PrintError("removing user from gitea", err)

		output, err := executil.Execute("kubectl delete ingress vscode-"+user.Username, debug)
		eh.PrintErrorWithOutput("deleting ingress", err, output)

		output, err = executil.Execute("kubectl delete service vscode-"+user.Username, debug)
		eh.PrintErrorWithOutput("deleting service", err, output)

		output, err = executil.Execute("kubectl delete deployment vscode-"+user.Username, debug)
		eh.PrintErrorWithOutput("deleting deployment", err, output)

		registryClient, err := regutil.NewRegistryClient(debug)
		eh.PrintError("getting registry client", err)

		repository := "7onetella/vscode-" + user.Username
		tags, err := registryClient.ListTags(repository, debug)
		eh.PrintError("listing tags", err)

		for _, tag := range tags {
			err = registryClient.DeleteImage(repository, tag, debug)
			eh.PrintError("deleting image", err)
		}

		output, err = sshutil.RemoteSSHExec("vm-docker-registry.curiosityworks.org", "22", "debian",
			"sudo rm -rf /var/lib/registry/docker/registry/v2/repositories/7onetella/vscode-"+user.Username)
		eh.PrintErrorWithOutput("deleting docker repo folder", err, output)

		dberr = userService.DeleteALLIDERuntimeInstallsForUser(user.ID)
		eh.PrintError("delete user ide runtime installs", dberr)

		dberr = userService.DeleteALLIDEReposForUser(user.ID)
		eh.PrintError("delete user ide repos", dberr)

		dberr = userService.DeleteALLUserIDEsForUser(user.ID)
		eh.PrintError("delete user ides", dberr)

		dberr = userService.DeleteALLUserReposForUser(user.ID)
		eh.PrintError("delete user repos", dberr)

		dberr = userService.Del(user)
		eh.PrintError("delete user", dberr)

		msg.Success("removing user")
	},
}

func init() {
	removeCmd.AddCommand(removeUserCmd)
}
