package cmd

import (
	"github.com/curiosinauts/platformctl/pkg/database"
	"github.com/spf13/cobra"
)

var addIDERepoCmdRepos []string

// addIDERepoCmd represents the add ide repo command
var addIDERepoCmd = &cobra.Command{
	Use:     "ide-repo",
	Aliases: []string{"ide-repos"},
	Short:   "Adds ide repo to user",
	Long:    `Adds ide repo to user`,
	Run: func(cmd *cobra.Command, args []string) {

		email := args[0]
		targetIDEName := args[1]

		eh := ErrorHandler{"adding runtime installs for users"}

		userObject, dberr := database.NewUserObject(userService, email)
		eh.HandleError("initializing user object", dberr)

		hasIDE, dberr := userObject.DoesUserHaveIDE(targetIDEName)
		eh.HandleError("does user object have ide", dberr)

		if hasIDE && dberr == nil {
			ide, dberr := userService.FindIDEByName(targetIDEName)
			eh.HandleError("finding ide by name", dberr)

			userIDE, _ := userObject.UserIDE(ide)
			if len(addIDERepoCmdRepos) > 0 {
				AddUserRepos(userObject.ID, addIDERepoCmdRepos)
				AddIDERepos(userIDE.ID, addIDERepoCmdRepos)
			}
		}
	},
}

func init() {
	addCmd.AddCommand(addIDERepoCmd)
	addIDERepoCmd.Flags().StringArrayVarP(&addIDERepoCmdRepos, "repo", "r", []string{}, "-r https://example-repo.com/foo")
}
