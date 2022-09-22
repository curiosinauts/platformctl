package cmd

import (
	"fmt"

	"github.com/curiosinauts/platformctl/pkg/grafanautil"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

// listPanelCmd represents the panel command
var listPanelCmd = &cobra.Command{
	Use:     "panel",
	Aliases: []string{"panels"},
	Short:   "list panel",
	Long:    `list panel`,
	Run: func(cmd *cobra.Command, args []string) {

		partialPanelTitle := ""
		if len(args) > 0 {
			partialPanelTitle = args[0]
		}

		eh := ErrorHandler{"list panels"}

		panels, err := grafanautil.ListPanels("7UdvG-Mnk", partialPanelTitle)
		eh.HandleError("calling list boards", err)

		// data := [][]string{}
		for _, p := range panels {
			id := uuid.NewString()
			// data = append(data, []string{strconv.Itoa(int(p.ID)), p.Title})
			err := grafanautil.DownloadPanel(int(p.ID), 600, 300, 1, id, debug)
			eh.HandleError("downloading panel", err)
			fmt.Println(p.Title)
			fmt.Printf("https://fileserver.curiosityworks.org/%s.png\n", id)

		}

		// table := tablewriter.NewWriter(os.Stdout)
		// table.SetRowLine(true)
		// table.SetHeader([]string{"ID", "Title"})

		// for _, v := range data {
		// 	table.Append(v)
		// }
		// table.SetRowSeparator("-")
		// table.Render()
	},
}

func init() {
	listCmd.AddCommand(listPanelCmd)
}
