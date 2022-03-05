package cmd

import (
	"fmt"

	"github.com/curiosinauts/platformctl/internal/msg"
	"github.com/curiosinauts/platformctl/pkg/postgresutil"

	"github.com/spf13/cobra"
)

var createSchema bool

// createUserSchemaCmd represents the user schema command
var createUserSchemaCmd = &cobra.Command{
	Use:     "db-user <username> <password> <host> <dbname>",
	Aliases: []string{"dbuser"},
	Short:   "Creates database user",
	Long:    `Creates database user without any association to schema`,
	Args:    cobra.MinimumNArgs(4),
	Example: "platformctl create dbuser john pass1234 db.example.com devdb",
	Run: func(cmd *cobra.Command, args []string) {

		username := args[0]
		password := args[1]

		psql := postgresutil.NewPSQLClientForSharedDB()
		if len(args) == 4 {
			hostname := args[2]
			dbname := args[3]
			psql = postgresutil.NewPSQLClientByHostAndDBName(hostname, dbname)
		}

		eh := ErrorHandler{"Creating database user"}
		if createSchema {
			out, err := psql.CreateUser(username, password, debug)
			eh.HandleError("creating database user", err)
			fmt.Println()
			fmt.Println(out)
		}

		out, err := psql.CreateUserSchema(username, debug)
		eh.HandleError("creating database user schema", err)

		fmt.Println(out)

		msg.Success("creating database user")
	},
}

func init() {
	createCmd.AddCommand(createUserSchemaCmd)
	createUserSchemaCmd.Flags().BoolVarP(&createSchema, "schema", "s", false, "create schema")
}
