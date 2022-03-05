package cmd

import (
	"fmt"

	"github.com/curiosinauts/platformctl/internal/msg"
	"github.com/curiosinauts/platformctl/pkg/postgresutil"
	"github.com/spf13/cobra"
)

// dbCmd represents the db command
var restoreDBSchemaCmd = &cobra.Command{
	Use:     "dbschema <username> <password> <host> <dbname> <schema>",
	Short:   "restores db schema",
	Long:    `restores db schema`,
	Example: "platformctl restore dbschema john pass1234 db.example.com devdb john",
	Args:    cobra.MinimumNArgs(5),
	Run: func(cmd *cobra.Command, args []string) {
		username := args[0]
		password := args[1]
		host := args[2]
		dbname := args[3]
		schemaName := args[4]

		psql := postgresutil.NewPSQLClientByHostAndDBName(host, dbname)
		_, filepath, err := psql.RestoreSchemaWithData(password, username, host, dbname, schemaName, debug)
		if err != nil {
			fmt.Println(err)
			return
		}

		msg.Success("restoring backup from " + filepath)
	},
}

func init() {
	restoreCmd.AddCommand(restoreDBSchemaCmd)
}
