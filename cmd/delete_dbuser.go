package cmd

import (
	"fmt"

	"github.com/curiosinauts/platformctl/internal/msg"
	"github.com/curiosinauts/platformctl/pkg/postgresutil"

	"github.com/spf13/cobra"
)

var dropSchema bool

// deleteDBUserCmd represents the user schema command
var deleteDBUserCmd = &cobra.Command{
	Use:     "db-user <username> <host> <dbname>",
	Aliases: []string{"dbuser"},
	Short:   "Drops database user",
	Long:    `Drops database user and optinally drop user's schema`,
	Args:    cobra.MinimumNArgs(3),
	Example: "platformctl drop db-user john db.example.com devdb",
	Run: func(cmd *cobra.Command, args []string) {

		username := args[0]
		host := args[1]
		dbname := args[2]

		eh := ErrorHandler{"dropping database user"}
		psql := postgresutil.NewPSQLClientByHostAndDBName(host, dbname)

		if dropSchema {
			out, err := psql.DropUserSchema(username, debug)
			eh.HandleError("dropping database schema", err)
			fmt.Println()
			fmt.Println(out)
		}

		out, err := psql.DropUser(username, debug)
		eh.HandleError("dropping database user", err)
		fmt.Println()
		fmt.Println(out)

		msg.Success("dropping database user")
	},
}

func init() {
	deleteCmd.AddCommand(deleteDBUserCmd)
	deleteDBUserCmd.Flags().BoolVarP(&dropSchema, "schema", "s", false, "drop schema")
}
