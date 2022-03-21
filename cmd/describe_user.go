package cmd

import (
	"fmt"

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
		fmt.Println("Git Repo URI: ", user.GitRepoURI)
		fmt.Println("Runtime Instals : ", user.RuntimeInstalls)
	},
}

func init() {
	describeCmd.AddCommand(describeUserCmd)
}
