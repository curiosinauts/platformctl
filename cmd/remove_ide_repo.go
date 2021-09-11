package cmd

import (
	"github.com/curiosinauts/platformctl/internal/msg"
	"github.com/curiosinauts/platformctl/pkg/database"
	"github.com/spf13/cobra"
)

var removeIDERepoCmdRepos []string

// removeIDERepoCmd represents the remove ide repo command
var removeIDERepoCmd = &cobra.Command{
	Use:     "ide-repo {email} {ide}",
	Aliases: []string{"ide-repos"},
	Short:   "Removes ide repo from user",
	Long:    `Removes ide repo from user`,
	Args:    cobra.MinimumNArgs(2),
	Example: `platformctl remove ide-repo foo@example.com vscode`,
	Run: func(cmd *cobra.Command, args []string) {

		email := args[0]
		targetIDEName := args[1]

		eh := ErrorHandler{"removing ide repo from user"}

		userObject, dberr := database.NewUserObject(userService, email)
		eh.HandleError("initializing user object", dberr)

		hasIDE, dberr := userObject.DoesUserHaveIDE(targetIDEName)
		eh.HandleError("does user object have ide", dberr)

		if hasIDE && dberr == nil {
			ide := &database.IDE{}
			eh.HandleError("finding ide by name", dbs.FindBy(ide, "name=$1", targetIDEName))

			userIDE, _ := userObject.UserIDE(*ide)
			if len(removeIDERepoCmdRepos) > 0 {
				RemoveIDERepos(userIDE.ID, removeIDERepoCmdRepos)
			}
		}

		msg.Success("removing ide repo for user")
	},
}

func init() {
	removeCmd.AddCommand(removeIDERepoCmd)
	removeIDERepoCmd.Flags().StringArrayVarP(&removeIDERepoCmdRepos, "repo", "r", []string{}, "-r https://example-repo.com/foo")
}
