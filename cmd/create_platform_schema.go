package cmd

import (
	"github.com/curiosinauts/platformctl/internal/msg"
	"github.com/curiosinauts/platformctl/pkg/database"

	"github.com/spf13/cobra"
)

// createSchemaCmd represents the schema command
var createSchemaCmd = &cobra.Command{
	Use:   "platformctl-schema",
	Short: "Creates database platformctl schema",
	Long:  `Creates database paltformctl schema`,
	Run: func(cmd *cobra.Command, args []string) {

		eh := ErrorHandler{"Creating database platformctl schema"}
		_, err := db.Exec(database.CreateSchema)
		eh.HandleError("creating database platformctl schema", err)

		msg.Success("creating database platformctl schema")
	},
}

func init() {
	createCmd.AddCommand(createSchemaCmd)
}
