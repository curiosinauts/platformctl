package cmd

import (
	"fmt"

	"github.com/curiosinauts/platformctl/pkg/postgresutil"
	"github.com/spf13/cobra"
)

// dbCmd represents the db command
var restoreDbCmd = &cobra.Command{
	Use:   "db <username> <password> <host> <dbname> <schema>",
	Short: "restores db",
	Long:  `restores db`,
	Args:  cobra.MinimumNArgs(5),
	Run: func(cmd *cobra.Command, args []string) {
		username := args[0]
		password := args[1]
		host := args[2]
		dbname := args[3]
		schemaName := args[4]

		psql := postgresutil.NewPSQLClientByHostAndDBName(host, dbname)
		out, err := psql.RestoreSchemaWithData(password, username, host, dbname, schemaName, debug)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(out)
	},
}

func init() {
	restoreCmd.AddCommand(restoreDbCmd)
}
