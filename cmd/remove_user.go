package cmd

import (
	"fmt"
	"github.com/curiosinauts/platformctl/internal/msg"
	"github.com/curiosinauts/platformctl/pkg/crypto"

	"github.com/curiosinauts/platformctl/internal/database"
	"github.com/spf13/cobra"
)

// removeUserCmd represents the user command
var removeUserCmd = &cobra.Command{
	Use:   "user",
	Short: "Remove user from the platforms",
	Long:  `Remove user from the platforms`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println()

		email := args[0]
		userService := database.NewUserService(db)

		eh := ErrorHandler{"removing user"}

		user, dberr := userService.FindUserByHashedEmail(crypto.Hashed(email))

		dberr = userService.DeleteALLIDERuntimeInstallsForUser(user.ID)
		eh.HandleError("delete user ide runtime installs", dberr)

		dberr = userService.DeleteALLIDEReposForUser(user.ID)
		eh.HandleError("delete user ide repos", dberr)

		dberr = userService.DeleteALLUserIDEsForUser(user.ID)
		eh.HandleError("delete user ides", dberr)

		dberr = userService.DeleteALLUserReposForUser(user.ID)
		eh.HandleError("delete user repos", dberr)

		dberr = userService.Delete(user.ID)
		eh.HandleError("delete user", dberr)

		msg.Success("removing user")
	},
}

func init() {
	removeCmd.AddCommand(removeUserCmd)
}
