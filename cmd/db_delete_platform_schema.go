package cmd

import (
	"github.com/curiosinauts/platformctl/internal/msg"
	"github.com/curiosinauts/platformctl/pkg/database"

	"github.com/spf13/cobra"
)

// deleteSchemaCmd represents the schema command
var dbDeletePlatformSchemaCmd = &cobra.Command{
	Use:   "platform-schema",
	Short: "Drops platform database schema",
	Long:  `Drops platform database schema`,
	Run: func(cmd *cobra.Command, args []string) {

		eh := ErrorHandler{"Dropping database schema"}
		_, err := db.Exec(database.DropSchema)
		eh.HandleError("dropping database schema", err)

		msg.Success("dropping database schema")
	},
}

func init() {
	dbDeleteCmd.AddCommand(dbDeletePlatformSchemaCmd)
}
