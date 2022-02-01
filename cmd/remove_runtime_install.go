package cmd

import (
	"strings"

	"github.com/curiosinauts/platformctl/internal/msg"
	"github.com/curiosinauts/platformctl/pkg/crypto"
	"github.com/curiosinauts/platformctl/pkg/database"
	"github.com/spf13/cobra"
)

var removeRuntimeInstallCmdUpdateNow bool

// removeRuntimeInstallCmd represents the runtimeInstall command
var removeRuntimeInstallCmd = &cobra.Command{
	Use:     "runtime-install {email | username | all} {ide} {runtime install}...",
	Aliases: []string{"runtime-installs"},
	Short:   "Removes runtime installs from users",
	Long:    `Removes runtime installs from users`,
	Example: `platformctl remove runtime-installs admin@curiosityworks.org vscode tmux,poetry`,
	Args:    cobra.MinimumNArgs(3),
	Run: func(cmd *cobra.Command, args []string) {

		emailOrAll := args[0]
		targetIDEName := args[1]
		runtimeInstallsStr := args[2]
		runtimeInstallNames := strings.Split(runtimeInstallsStr, ",")

		eh := ErrorHandler{"removing runtime installs for users"}

		var users []database.User
		var dberr *database.DBError

		if emailOrAll == "all" {
			dberr = dbs.List("curiosity.user", &users)
			eh.HandleError("get all users", dberr)
		} else {
			user := database.User{}
			dberr := dbs.FindBy(&user, "hashed_email=$1", crypto.Hashed(emailOrAll))
			eh.HandleError("finding user by hashed email", dberr)
			users = append(users, user)
		}

		for _, user := range users {
			userObject := database.UserObject{
				User:        user,
				UserService: dbs,
			}
			eh.HandleError("new user object", dberr)

			hasIDE, dberr := userObject.DoesUserHaveIDE(targetIDEName)
			eh.HandleError("does user have ide", dberr)

			if hasIDE && dberr == nil {
				ide := database.IDE{}
				eh.HandleError("find ide by name", dbs.FindBy(&ide, "name=$1", targetIDEName))
				for _, runtimeInstallName := range runtimeInstallNames {
					hasRuntimeInstall, dberr := userObject.DoesUserHaveRuntimeInstallFor(ide, runtimeInstallName)
					eh.HandleError("does user has runtime install", dberr)

					if hasRuntimeInstall && dberr == nil {
						userIDE, _ := userObject.UserIDE(ide)
						runtimeInstall := database.RuntimeInstall{}
						dberr = dbs.FindBy(&runtimeInstall, "name=$1", runtimeInstallName)
						eh.HandleError("finding runtime install by name", dberr)

						ideRuntimeInstall := database.IDERuntimeInstall{}
						dberr = dbs.FindBy(&ideRuntimeInstall, "user_ide_id=$1 AND runtime_install_id=$2", userIDE.ID, runtimeInstall.ID)
						eh.HandleError("finding runtime install by name", dberr)

						dberr = dbs.Del(&ideRuntimeInstall)
						eh.HandleError("finding runtime install by name", dberr)
					}
				}
			}
		}

		if removeRuntimeInstallCmdUpdateNow {
			msg.Info("updating users' ides")
			updateCodeserverCmd.Run(updateCodeserverCmd, []string{emailOrAll})
		}

		msg.Success("removing runtime installs for users")
	},
}

func init() {
	removeCmd.AddCommand(removeRuntimeInstallCmd)
	removeRuntimeInstallCmd.Flags().BoolVarP(&removeRuntimeInstallCmdUpdateNow, "now", "n", false, "update users' ide or not")
}
