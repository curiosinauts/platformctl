package cmd

import (
	"fmt"
	"os"

	"github.com/curiosinauts/platformctl/pkg/database"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// listRuntimeInstallsCmd represents the runtimeInstalls command
var listRuntimeInstallsCmd = &cobra.Command{
	Use:     "runtime-installs",
	Short:   "Lists runtime installs",
	Long:    `Lists runtime installs`,
	Aliases: []string{"runtime-install"},
	Run: func(cmd *cobra.Command, args []string) {

		dbs := database.NewUserService(db)
		runtimeInstalls := &[]database.RuntimeInstall{}

		eh := ErrorHandler{"listing runtime installs"}
		dberr := dbs.List(&database.RuntimeInstall{}, runtimeInstalls)
		eh.HandleError("listing entities", dberr)

		data := [][]string{}
		for _, r := range *runtimeInstalls {
			data = append(data, []string{r.Name, r.ScriptBody})
			fmt.Println(r.Name)
			fmt.Println()
			fmt.Println(r.ScriptBody)
			fmt.Println("-------------------------------------------------------------------------------")
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetRowLine(true)
		table.SetHeader([]string{"Name", "Runtime Install"})

		for _, v := range data {
			table.Append(v)
		}
		table.SetRowSeparator("-")
		// table.Render() // Send output
	},
}

func init() {
	listCmd.AddCommand(listRuntimeInstallsCmd)
}
