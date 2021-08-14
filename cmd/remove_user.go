package cmd

import (
	"fmt"
	"github.com/curiosinauts/platformctl/internal/database"
	"github.com/spf13/cobra"
)

// removeUserCmd represents the user command
var removeUserCmd = &cobra.Command{
	Use:   "user",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println()

		userService := database.NewUserService(db)

		eh := ErrorHandler{"removing user"}

		dberr := userService.RemoveIDERuntimeInstall(1)
		eh.HandleError("delete ide runtime install", dberr)
	},
}

func init() {
	removeCmd.AddCommand(removeUserCmd)
}
