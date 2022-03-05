package cmd

import (
	"fmt"

	"github.com/curiosinauts/platformctl/internal/msg"
	"github.com/curiosinauts/platformctl/pkg/postgresutil"
	"github.com/spf13/cobra"
)

// dbCmd represents the db command
var backUpDBSchemaCmd = &cobra.Command{
	Use:     "dbschema <username> <password> <host> <dbname> <schema>",
	Aliases: []string{"db-schema"},
	Short:   "backs up db schema with data",
	Long:    `backs up db schema with data`,
	Example: "platformctl backup dbschema john pass1234 db.example.com devdb john",
	Args:    cobra.MinimumNArgs(5),
	Run: func(cmd *cobra.Command, args []string) {
		username := args[0]
		password := args[1]
		host := args[2]
		dbname := args[3]
		schemaName := args[4]

		psql := postgresutil.NewPSQLClientByHostAndDBName(host, dbname)
		_, filepath, err := psql.BackUpSchemaWithData(password, username, host, dbname, schemaName, debug)
		if err != nil {
			fmt.Println(err)
			return
		}

		msg.Success("backing up db schema to " + filepath)
	},
}

func init() {
	backupCmd.AddCommand(backUpDBSchemaCmd)
}
