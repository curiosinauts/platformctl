package cmd

import (
	"fmt"
	"strings"

	"github.com/curiosinauts/platformctl/pkg/postgresutil"

	"github.com/spf13/cobra"
)

// listDBUserCmd represents the user schema command
var listDBUserCmd = &cobra.Command{
	Use:     "db-users <host> <dbname>",
	Aliases: []string{"dbusers"},
	Short:   "Lists database users",
	Long:    `Lists database users`,
	Args:    cobra.MinimumNArgs(2),
	Example: "platformctl list db-users db.example.com devdb",
	Run: func(cmd *cobra.Command, args []string) {

		host := args[0]
		dbname := args[1]

		eh := ErrorHandler{"Dropping database user & user schema"}
		psql := postgresutil.NewPSQLClientByHostAndDBName(host, dbname)
		out, err := psql.ListDBUsers(debug)
		eh.HandleError("listing database users", err)
		fmt.Println()
		fmt.Println(strings.TrimSpace(out))
	},
}

func init() {
	listCmd.AddCommand(listDBUserCmd)
}
