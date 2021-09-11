package cmd

import (
	"fmt"

	"github.com/curiosinauts/platformctl/internal/msg"
	"github.com/curiosinauts/platformctl/pkg/postgresutil"

	"github.com/spf13/cobra"
)

// createUserSchemaCmd represents the user schema command
var createUserSchemaCmd = &cobra.Command{
	Use:     "user-schema",
	Short:   "Creates database user and user schema",
	Long:    `Creates database user and user schema`,
	Args:    cobra.MinimumNArgs(2),
	Example: "platformctl create user-schema foo-1234 pass1234",
	Run: func(cmd *cobra.Command, args []string) {

		username := args[0]
		password := args[1]

		eh := ErrorHandler{"Creating database user & user schema"}
		psql := postgresutil.NewPSQLClient()
		out, err := psql.CreateUser(username, password)
		eh.HandleError("creating database user", err)

		fmt.Println()
		fmt.Println(out)

		out, err = psql.CreateUserSchema(username)
		eh.HandleError("creating database user schema", err)

		fmt.Println(out)

		msg.Success("creating database usre & user schema")
	},
}

func init() {
	createCmd.AddCommand(createUserSchemaCmd)
}
