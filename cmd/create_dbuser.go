package cmd

import (
	"fmt"

	"github.com/curiosinauts/platformctl/internal/msg"
	"github.com/curiosinauts/platformctl/pkg/postgresutil"

	"github.com/spf13/cobra"
)

var schemaOnly bool

// createUserSchemaCmd represents the user schema command
var createUserSchemaCmd = &cobra.Command{
	Use:     "dbuser <username> <password> [<host>] [<dbname>]",
	Short:   "Creates database user and user schema",
	Long:    `Creates database user and user schema`,
	Args:    cobra.MinimumNArgs(2),
	Example: "platformctl create dbuser foo-1234 pass1234",
	Run: func(cmd *cobra.Command, args []string) {

		username := args[0]
		password := args[1]

		psql := postgresutil.NewPSQLClient()
		if len(args) == 4 {
			hostname := args[2]
			dbname := args[3]
			psql = postgresutil.NewPSQLClientByHostAndDBName(hostname, dbname)
		}

		eh := ErrorHandler{"Creating database user & user schema"}
		if !schemaOnly {
			out, err := psql.CreateUser(username, password)
			eh.HandleError("creating database user", err)
			fmt.Println()
			fmt.Println(out)
		}

		out, err := psql.CreateUserSchema(username)
		eh.HandleError("creating database user schema", err)

		fmt.Println(out)

		msg.Success("creating database usre & user schema")
	},
}

func init() {
	createCmd.AddCommand(createUserSchemaCmd)
	createUserSchemaCmd.Flags().BoolVarP(&schemaOnly, "schema", "s", false, "create schema only")
}
