package cmd

import (
	"fmt"

	"github.com/curiosinauts/platformctl/internal/msg"
	"github.com/curiosinauts/platformctl/pkg/database"

	"github.com/spf13/cobra"
)

// createSchemaCmd represents the schema command
var createSchemaCmd = &cobra.Command{
	Use:   "schema",
	Short: "Creates database schema",
	Long:  `Creates database schema`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println()
		eh := ErrorHandler{"Creating database schema"}
		_, err := db.Exec(database.CreateSchema)
		eh.HandleError("creating database schema", err)

		msg.Success("creating database schema")
	},
}

func init() {
	createCmd.AddCommand(createSchemaCmd)
}
