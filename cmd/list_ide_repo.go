package cmd

import (
	"github.com/curiosinauts/platformctl/internal/msg"
	"github.com/curiosinauts/platformctl/pkg/database"
	"github.com/spf13/cobra"
)

// listIDERepoCmd represents the add ide repo command
var listIDERepoCmd = &cobra.Command{
	Use:     "ide-repo {email} {ide}",
	Aliases: []string{"ide-repos"},
	Short:   "Lists ide repos of a user",
	Long:    `Lists ide repos of a user`,
	Args:    cobra.MinimumNArgs(2),
	Example: `platformctl list ide-repos foo@example.com vscode`,
	Run: func(cmd *cobra.Command, args []string) {

		email := args[0]
		targetIDEName := args[1]

		eh := ErrorHandler{"listing ide repos for user"}

		userObject, dberr := database.NewUserObject(userService, email)
		eh.HandleError("initializing user object", dberr)

		hasIDE, dberr := userObject.DoesUserHaveIDE(targetIDEName)
		eh.HandleError("does user object have ide", dberr)

		if hasIDE && dberr == nil {
			ide := &database.IDE{}
			eh.HandleError("finding ide by name", dbs.FindBy(ide, "name=$1", targetIDEName))

			userIDE, _ := userObject.UserIDE(*ide)

			ideRepos := []database.IDERepo{}
			dberr = dbs.ListBy("ide_repo", &ideRepos, "user_ide_id=$1", userIDE.ID)
			eh.HandleError("listing ide repos", dberr)

			for _, ideRepo := range ideRepos {
				msg.Info(ideRepo.URI)
			}
		}
	},
}

func init() {
	listCmd.AddCommand(listIDERepoCmd)
}
