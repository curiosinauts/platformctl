package cmd

import (
	"os"
	"strconv"

	"github.com/curiosinauts/platformctl/pkg/grafanautil"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// listBoardCmd represents the board command
var listBoardCmd = &cobra.Command{
	Use:     "board",
	Aliases: []string{"boards"},
	Short:   "list board",
	Long:    `list board`,
	Run: func(cmd *cobra.Command, args []string) {

		eh := ErrorHandler{"list boards"}

		foundboards, err := grafanautil.ListDashboards("")
		eh.HandleError("calling list boards", err)

		data := [][]string{}
		for _, fb := range foundboards {
			data = append(data, []string{strconv.Itoa(int(fb.ID)), fb.UID, fb.Title})
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetRowLine(true)
		table.SetHeader([]string{"ID", "UID", "Title"})

		for _, v := range data {
			table.Append(v)
		}
		table.SetRowSeparator("-")
		table.Render()
	},
}

func init() {
	listCmd.AddCommand(listBoardCmd)
}
