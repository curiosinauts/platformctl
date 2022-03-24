package cmd

import (
	"fmt"
	"os"

	"github.com/curiosinauts/platformctl/internal/msg"
	"github.com/curiosinauts/platformctl/pkg/crypto"
	"github.com/curiosinauts/platformctl/pkg/database"
	"github.com/spf13/cobra"
)

// updateCodeserverCmd represents the codeserver command
var updateCodeserverCmd = &cobra.Command{
	Use:     "codeserver <email|all>",
	Aliases: []string{"code-server", "pod"},
	Short:   "Updates code server for given user",
	Long:    `Updates code server for given user`,
	PreRunE: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		eh := ErrorHandler{"updating code server"}

		emailOrAll := args[0]

		var users []database.User
		var dberr *database.DBError

		if emailOrAll == "all" {
			dberr = dbs.List("users", &users)
			eh.HandleError("retrieving all users", dberr)
		} else {
			user := database.User{}
			dberr := dbs.FindBy(&user, "hashed_email=$1", crypto.Hashed(emailOrAll))
			eh.HandleError("finding user by hashed email", dberr)
			users = append(users, user)
		}

		for _, user := range users {
			CreateDeploymentServiceIngressYamlFile(user)

			CreateUserSecretsFile(user)

			ApplySecrets(user, eh)

			ApplyDeployment(user, eh)

			os.Remove("vscode-" + user.Username + "-secrets.yml")

			os.Remove("vscode-" + user.Username + ".yml")

			msg.Success(fmt.Sprintf("updating code server for %s", user.Username))
		}
	},
}

func init() {
	updateCmd.AddCommand(updateCodeserverCmd)
}
