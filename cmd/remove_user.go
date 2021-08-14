package cmd

import (
	"fmt"

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

		user, dberr := userService.FindByEmail(email)
		userIDEIDs, dberr := userService.FindUserIDEsByUserID(user.ID)
		// dberr := userService.RemoveIDERuntimeInstall(1)
		eh.HandleError("delete ide runtime install", dberr)

		fmt.Println(userIDEIDs)
	},
}

func init() {
	removeCmd.AddCommand(removeUserCmd)
}
