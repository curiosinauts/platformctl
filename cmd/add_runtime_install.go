package cmd

import (
	"strings"

	"github.com/curiosinauts/platformctl/pkg/crypto"
	"github.com/curiosinauts/platformctl/pkg/database"
	"github.com/spf13/cobra"
)

// addRuntimeInstallCmd represents the runtimeInstall command
var addRuntimeInstallCmd = &cobra.Command{
	Use:     "runtime-install",
	Aliases: []string{"runtime-installs"},
	Short:   "Adds runtime install to users",
	Long:    `Adds runtime install to users`,
	Example: `platformctl add runtime-installs 7onetella@gmail.com vscode tmux,poetry`,
	Run: func(cmd *cobra.Command, args []string) {

		emailOrAll := args[0]
		targetIDEName := args[1]
		runtimeInstallsStr := args[2]
		runtimeInstallNames := strings.Split(runtimeInstallsStr, ",")

		eh := ErrorHandler{"adding runtime installs for users"}

		var users []database.User
		var dberr *database.DBError

		if emailOrAll == "all" {
			users, dberr = userService.List()
			eh.HandleError("retrieving all users", dberr)
		} else {
			user, dberr := userService.FindUserByHashedEmail(crypto.Hashed(emailOrAll))
			eh.HandleError("finding user by hashed email", dberr)
			users = append(users, user)
		}

		for _, user := range users {
			userObject := database.UserObject{
				User:        user,
				UserService: userService,
			}
			eh.HandleError("initializing user object", dberr)

			hasIDE, dberr := userObject.DoesUserHaveIDE(targetIDEName)
			eh.HandleError("does user object have ide", dberr)

			if hasIDE && dberr == nil {
				ide, dberr := userService.FindIDEByName(targetIDEName)
				eh.HandleError("find ide by name", dberr)
				for _, runtimeInstallName := range runtimeInstallNames {
					hasRuntimeInstall, dberr := userObject.DoesUserHaveRuntimeInstallFor(ide, runtimeInstallName)
					eh.HandleError("has runtime installs", dberr)

					if !hasRuntimeInstall && dberr == nil {
						userIDE, _ := userObject.UserIDE(ide)
						runtimeInstall, _ := userService.FindRuntimeInstallByName(runtimeInstallName)
						ideRuntimeInstall := database.IDERuntimeInstall{
							UserIDEID:        userIDE.ID,
							RuntimeInstallID: runtimeInstall.ID,
						}
						userService.AddIDERuntimeInstall(ideRuntimeInstall)
					}
				}
			}
		}
	},
}

func init() {
	addCmd.AddCommand(addRuntimeInstallCmd)
}
