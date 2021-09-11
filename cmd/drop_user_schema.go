package cmd

import (
	"fmt"

	"github.com/curiosinauts/platformctl/internal/msg"
	"github.com/curiosinauts/platformctl/pkg/postgresutil"

	"github.com/spf13/cobra"
)

// dropUserSchemaCmd represents the user schema command
var dropUserSchemaCmd = &cobra.Command{
	Use:     "user-schema",
	Short:   "Drops database user and user schema",
	Long:    `Drops database user and user schema`,
	Args:    cobra.MinimumNArgs(1),
	Example: "platformctl drop user-schema foo-1234 pass1234",
	Run: func(cmd *cobra.Command, args []string) {

		username := args[0]

		eh := ErrorHandler{"Dropping database user & user schema"}
		psql := postgresutil.NewPSQLClient()
		out, err := psql.DropUserSchema(username)
		eh.HandleError("dropping database user schema", err)

		fmt.Println()
		fmt.Println(out)

		out, err = psql.DropUser(username)
		eh.HandleError("dropping database user", err)

		fmt.Println(out)

		msg.Success("dropping database usre & user schema")
	},
}

func init() {
	dropCmd.AddCommand(dropUserSchemaCmd)
}
