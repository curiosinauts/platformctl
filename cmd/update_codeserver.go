package cmd

import (
	"fmt"

	"github.com/curiosinauts/platformctl/internal/msg"
	"github.com/curiosinauts/platformctl/pkg/crypto"
	"github.com/curiosinauts/platformctl/pkg/database"
	"github.com/curiosinauts/platformctl/pkg/jenkinsutil"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

// updateCodeserverCmd represents the codeserver command
var updateCodeserverCmd = &cobra.Command{
	Use:     "codeserver",
	Aliases: []string{"code-server"},
	Short:   "Updates code server for given user",
	Long:    `Updates code server for given user`,
	PreRunE: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		eh := ErrorHandler{"updating code server"}

		emailOrAll := args[0]

		var users []database.User
		var dberr *database.DBError

		if emailOrAll == "all" {
			dberr = dbs.List("platformctl.user", &users)
			eh.HandleError("retrieving all users", dberr)
		} else {
			user := database.User{}
			dberr := dbs.FindBy(&user, "hashed_email=$1", crypto.Hashed(emailOrAll))
			eh.HandleError("finding user by hashed email", dberr)
			users = append(users, user)
		}

		for _, user := range users {
			jenkins, err := jenkinsutil.NewJenkins()
			eh.HandleError("accessing Jenkins job", err)

			option := map[string]string{
				"USERNAME": user.Username,
				"VERSION":  uuid.NewString(),
			}
			_, err = jenkins.BuildJob("codeserver", option)
			eh.HandleError("calling Jenkins job to build codeserver instance", err)

			msg.Success(fmt.Sprintf("updating code server for %s", user.Username))
		}
	},
}

func init() {
	updateCmd.AddCommand(updateCodeserverCmd)
}
