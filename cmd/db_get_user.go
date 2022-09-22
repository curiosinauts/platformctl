package cmd

import (
	"fmt"
	"strings"

	"github.com/curiosinauts/platformctl/pkg/postgresutil"

	"github.com/spf13/cobra"
)

// listDBUserCmd represents the user schema command
var listDBUserCmd = &cobra.Command{
	Use:     "user <host> <dbname>",
	Aliases: []string{"users"},
	Short:   "Lists database users",
	Long:    `Lists database users`,
	Args:    cobra.MinimumNArgs(1),
	Example: "platformctl list db-users db.example.com",
	Run: func(cmd *cobra.Command, args []string) {

		host := args[0]

		eh := ErrorHandler{"Dropping database user & user schema"}
		psql := postgresutil.NewPSQLClientByHostAndDBName(host, "non needed")
		out, err := psql.ListDBUsers(debug)
		eh.HandleError("listing database users", err)
		fmt.Println()
		fmt.Println(strings.TrimSpace(out))
	},
}

func init() {
	listCmd.AddCommand(listDBUserCmd)
}
