package cmd

import (
	"fmt"
	"github.com/curiosinauts/platformctl/pkg/database"
	"github.com/spf13/cobra"
)

// listUserCmd represents the user command
var listUserCmd = &cobra.Command{
	Use:     "user",
	Aliases: []string{"users"},
	Short:   "List users",
	Long:    `List users`,
	Run: func(cmd *cobra.Command, args []string) {
		userService := database.NewUserService(db)
		users, _ := userService.List()
		for _, user := range users {
			fmt.Println(user)
		}
	},
}

func init() {
	listCmd.AddCommand(listUserCmd)
}
