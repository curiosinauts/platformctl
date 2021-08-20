package cmd

import (
	"github.com/curiosinauts/platformctl/internal/msg"
	"github.com/curiosinauts/platformctl/pkg/database"

	"github.com/spf13/cobra"
)

// dropSchemaCmd represents the schema command
var dropSchemaCmd = &cobra.Command{
	Use:   "schema",
	Short: "Drops database schema",
	Long:  `Drops database schema`,
	Run: func(cmd *cobra.Command, args []string) {

		eh := ErrorHandler{"Dropping database schema"}
		_, err := db.Exec(database.DropSchema)
		eh.HandleError("dropping database schema", err)

		msg.Success("dropping database schema")
	},
}

func init() {
	dropCmd.AddCommand(dropSchemaCmd)
}
