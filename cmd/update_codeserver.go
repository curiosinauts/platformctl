package cmd

import (
	"fmt"

	"github.com/curiosinauts/platformctl/internal/msg"
	"github.com/curiosinauts/platformctl/pkg/crypto"
	"github.com/curiosinauts/platformctl/pkg/database"
	"github.com/curiosinauts/platformctl/pkg/jenkinsutil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// codeserverCmd represents the codeserver command
var codeserverCmd = &cobra.Command{
	Use:   "codeserver",
	Short: "Updates code server for given user",
	Long:  `Updates code server for given user`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println()

		eh := ErrorHandler{"updating code server"}

		email := args[0]
		userService := database.NewUserService(db)

		user, dberr := userService.FindUserByHashedEmail(crypto.Hashed(email))
		eh.HandleError("finding user by hashed email", dberr)

		option := map[string]string{
			"USERNAME": user.Username,
		}
		jenkinsAPIKey := viper.Get("jenkins_api_key").(string)
		jenkinsURL := viper.Get("jenkins_url").(string)
		jenkins, err := jenkinsutil.NewJenkins(jenkinsURL, "admin", jenkinsAPIKey)
		eh.HandleError("accessing Jenkins job", err)

		_, err = jenkins.BuildJob("codeserver", option)
		eh.HandleError("calling Jenkins job to build codeserver instance", err)

		msg.Success("updating code server")
	},
}

func init() {
	updateCmd.AddCommand(codeserverCmd)
}
