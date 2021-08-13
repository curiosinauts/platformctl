package cmd

import (
	"fmt"
	"github.com/curiosinauts/platformctl/internal/database"

	"github.com/spf13/cobra"
)

// createSchemaCmd represents the schema command
var createSchemaCmd = &cobra.Command{
	Use:   "schema",
	Short: "Creates database schema",
	Long:  `Creates database schema`,
	Run: func(cmd *cobra.Command, args []string) {
		_, err := db.Exec(database.CreateSchema)
		if err != nil {
			fmt.Println("creating schema failed")
			fmt.Printf("  %v", err.Error())
		}
	},
}

func init() {
	createCmd.AddCommand(createSchemaCmd)
}
