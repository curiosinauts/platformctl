package cmd

import (
	"strings"

	"github.com/curiosinauts/platformctl/pkg/database"
	"github.com/spf13/cobra"
)

// addRuntimeInstallCmd represents the runtimeInstall command
var addRuntimeInstallCmd = &cobra.Command{
	Use:     "runtime-install",
	Aliases: []string{"runtime-installs"},
	Short:   "Adds runtime install to users",
	Long:    `Adds runtime install to users`,
	Run: func(cmd *cobra.Command, args []string) {
		// platformctl add runtime-installs scott.seo@gmail.com vscode tmux,nodejs
		email := args[0]
		targetIDEName := args[1]
		runtimeInstallsStr := args[2]
		runtimeInstallNames := strings.Split(runtimeInstallsStr, ",")

		eh := ErrorHandler{"adding runtime installs for users"}

		userService := database.NewUserService(db)

		userObject, dberr := database.NewUserObject(userService, email)
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
	},
}

func init() {
	addCmd.AddCommand(addRuntimeInstallCmd)
}
