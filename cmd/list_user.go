package cmd

import (
	"os"
	"strings"

	"github.com/curiosinauts/platformctl/pkg/database"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// listUserCmd represents the user command
var listUserCmd = &cobra.Command{
	Use:     "user",
	Aliases: []string{"users"},
	Short:   "List users",
	Long:    `List users`,
	Run: func(cmd *cobra.Command, args []string) {
		eh := ErrorHandler{"listing users"}

		users := &[]database.User{}
		dberr := dbs.List("platformctl.user", users)
		eh.HandleError("list users", dberr)

		var data [][]string

		for _, user := range *users {
			username := user.Username
			ides, dberr := dbs.FindUserIDEsByUserID(user.ID)
			eh.HandleError("finding ides for user "+username, dberr)
			for _, ide := range *ides {
				runtimeInstalls := []database.RuntimeInstall{}
				dberr = dbs.FindUserIDERuntimeInstallsByUsernameAndIDE(&runtimeInstalls, username, ide)
				eh.HandleError("finding runtime installs for ide "+ide, dberr)

				var runtimeInstallNames []string
				for _, runtimeInstall := range runtimeInstalls {
					runtimeInstallNames = append(runtimeInstallNames, runtimeInstall.Name)
				}
				data = append(data, []string{username, ide, strings.Join(runtimeInstallNames, ",")})
			}
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Name", "IDEs", "Runtime Installs"})

		for _, v := range data {
			table.Append(v)
		}
		table.Render() // Send output
	},
}

func init() {
	listCmd.AddCommand(listUserCmd)
}
