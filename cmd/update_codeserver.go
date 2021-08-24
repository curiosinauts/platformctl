package cmd

import (
	"github.com/curiosinauts/platformctl/internal/msg"
	"github.com/curiosinauts/platformctl/pkg/crypto"
	"github.com/curiosinauts/platformctl/pkg/database"
	"github.com/curiosinauts/platformctl/pkg/jenkinsutil"
	"github.com/spf13/cobra"
)

// codeserverCmd represents the codeserver command
var codeserverCmd = &cobra.Command{
	Use:     "codeserver",
	Short:   "Updates code server for given user",
	Long:    `Updates code server for given user`,
	PreRunE: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		eh := ErrorHandler{"updating code server"}

		email := args[0]
		userService := database.NewUserService(db)

		user, dberr := userService.FindUserByHashedEmail(crypto.Hashed(email))
		eh.HandleError("finding user by hashed email", dberr)

		jenkins, err := jenkinsutil.NewJenkins()
		eh.HandleError("accessing Jenkins job", err)

		option := map[string]string{
			"USERNAME": user.Username,
		}
		_, err = jenkins.BuildJob("codeserver", option)
		eh.HandleError("calling Jenkins job to build codeserver instance", err)

		msg.Success("updating code server")
	},
}

func init() {
	updateCmd.AddCommand(codeserverCmd)
}
