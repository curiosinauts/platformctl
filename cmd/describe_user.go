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
		fmt.Println()

		email := args[0]
		userService := database.NewUserService(db)

		eh := ErrorHandler{"removing user"}

		user, dberr := userService.FindUserByHashedEmail(crypto.Hashed(email))
		eh.HandleError("retrieving user by email", dberr)

		fmt.Println("ID          : ", user.ID)
		fmt.Println("Username    : ", user.Username)
		fmt.Println("GoogleID    : ", user.GoogleID)
		fmt.Println("HashedEmail : ", user.HashedEmail)
	},
}

func init() {
	describeCmd.AddCommand(describeUserCmd)
}
