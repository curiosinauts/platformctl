package cmd

import (
	"fmt"
	"strings"

	"github.com/curiosinauts/platformctl/pkg/crypto"
	"github.com/curiosinauts/platformctl/pkg/database"
	"github.com/spf13/cobra"
)

// describeUserCmd represents the user command
var describeUserCmd = &cobra.Command{
	Use:     "user",
	Short:   "Describes the given user",
	Long:    `Describes the given user`,
	PreRunE: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		email := args[0]

		eh := ErrorHandler{"describing user"}

		user := database.User{}
		dberr := dbs.FindBy(&user, "hashed_email=$1", crypto.Hashed(email))
		eh.HandleError("retrieving user by email", dberr)

		fmt.Println()
		fmt.Println("ID          : ", user.ID)
		fmt.Println("Username    : ", user.Username)
		fmt.Println("GoogleID    : ", user.GoogleID)
		fmt.Println("HashedEmail : ", user.HashedEmail)

		username := user.Username
		ides, dberr := dbs.FindUserIDEsByUserID(user.ID)
		eh.HandleError("finding ides for user "+username, dberr)
		runtimeInstallNames := &[]string{}
		for _, ide := range *ides {
			runtimeInstallNames, dberr = dbs.FindUserIDERuntimeInstallNamesByUserAndIDE(username, ide)
			eh.HandleError("finding runtime installs for ide "+ide, dberr)
			fmt.Println("IDE         : ", ide)
			fmt.Println("Runtime Ins : ", strings.Join(*runtimeInstallNames, ","))

			ideS := database.IDE{}
			dberr = dbs.FindBy(&ideS, "name=$1", ide)
			eh.HandleError("finding ide by name", dberr)

			userIDE := database.UserIDE{}
			dbs.FindBy(&userIDE, "user_id=$1 AND ide_id=$2", user.ID, ideS.ID)
			eh.HandleError("finding user ide by user and ide", dberr)

			ideRepos := []database.IDERepo{}
			dberr = dbs.ListBy("ide_repo", &ideRepos, "user_ide_id=$1", userIDE.ID)
			eh.HandleError("listing ide repos", dberr)

			fmt.Println("IDE Repos   : ")
			for _, ideRepo := range ideRepos {
				fmt.Println("               " + ideRepo.URI)
			}
		}
	},
}

func init() {
	describeCmd.AddCommand(describeUserCmd)
}
