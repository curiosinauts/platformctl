package cmd

import (
	"fmt"
	"strconv"

	"github.com/curiosinauts/platformctl/pkg/grafanautil"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

// showGraphCmd represents the graph command
var showGraphCmd = &cobra.Command{
	Use:     "graph {panel partial name} {hours ago}",
	Short:   "Shows grafana graphs",
	Long:    `Shows grafana graphs`,
	Aliases: []string{"graphs"},
	Example: "platformctl show graph foo 1",
	Run: func(cmd *cobra.Command, args []string) {
		// panelPartialName := args[0]
		hoursAgo, _ := strconv.Atoi(args[1])

		id := uuid.NewString()
		eh := ErrorHandler{"showing graphs"}

		err := grafanautil.DownloadPanel(15, 600, 300, hoursAgo, id, debug)
		eh.HandleError("downloading panel", err)

		fmt.Printf("https://fileserver.curiosityworks.org/%s.png\n", id)
	},
}

func init() {
	showCmd.AddCommand(showGraphCmd)
}
