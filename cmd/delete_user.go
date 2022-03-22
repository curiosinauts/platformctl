package cmd

import (
	"strings"

	"github.com/curiosinauts/platformctl/internal/msg"
	"github.com/curiosinauts/platformctl/pkg/crypto"
	"github.com/curiosinauts/platformctl/pkg/database"
	"github.com/curiosinauts/platformctl/pkg/executil"
	"github.com/curiosinauts/platformctl/pkg/giteautil"
	"github.com/curiosinauts/platformctl/pkg/postgresutil"
	"github.com/curiosinauts/platformctl/pkg/regutil"
	"github.com/curiosinauts/platformctl/pkg/sshutil"
	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

var deleteUserCmdUsername string

// deleteUserCmd represents the user command
var deleteUserCmd = &cobra.Command{
	Use:     "user <email>",
	Short:   "Removes user from the platform",
	Long:    `Removes user from the platform`,
	PreRunE: cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		var email string
		if len(args) > 0 {
			email = args[0]
		}

		eh := ErrorHandler{"removing user"}

		user := database.User{}
		hashedEmail := crypto.Hashed(email)
		if debug {
			msg.Info("hashed email = " + hashedEmail)
		}

		var dberr *database.DBError

		if len(deleteUserCmdUsername) > 0 {
			dberr = dbs.FindBy(&user, "username=$1", deleteUserCmdUsername)
			eh.PrintError("finding by username", dberr)
		} else {
			dberr := dbs.FindBy(&user, "hashed_email=$1", hashedEmail)
			eh.PrintError("finding by hashed email", dberr)
		}

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

		repository := "curiosinauts/vscode-" + user.Username
		tags, err := registryClient.ListTags(repository, debug)
		eh.PrintError("listing tags", err)

		for _, tag := range tags {
			err = registryClient.DeleteImage(repository, tag, debug)
			eh.PrintError("deleting image", err)
		}

		dockerRegistryHost := viper.Get("docker_registry_host").(string)
		output, err = sshutil.RemoteSSHExec(dockerRegistryHost, "22", "debian",
			"sudo rm -rf /var/lib/registry/docker/registry/v2/repositories/curiosinauts/vscode-"+user.Username)
		eh.PrintErrorWithOutput("deleting docker repo folder", err, output)

		postgresUsername := strings.Replace(user.Username, "-", "", -1)
		psql := postgresutil.NewPSQLClientForSharedDB()
		_, err = psql.DropUserSchema(postgresUsername, debug)
		eh.PrintError("dropping database user schema", err)

		_, err = psql.DropUser(postgresUsername, debug)
		eh.PrintError("dropping database user", err)

		dberr = dbs.Del(&user)
		eh.PrintError("delete user", dberr)

		msg.Success("removing user")
	},
}

func init() {
	deleteCmd.AddCommand(deleteUserCmd)
	deleteUserCmd.Flags().StringVarP(&deleteUserCmdUsername, "username", "u", "", "userename")
}
