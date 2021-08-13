package cmd

import (
	"fmt"
	"github.com/curiosinauts/platformctl/internal/database"

	"github.com/spf13/cobra"
)

// dropSchemaCmd represents the schema command
var dropSchemaCmd = &cobra.Command{
	Use:   "schema",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		_, err := db.Exec(database.DropSchema)
		if err != nil {
			fmt.Println("dropping schema failed")
			fmt.Printf("  %v", err.Error())
		}
	},
}

func init() {
	dropCmd.AddCommand(dropSchemaCmd)
}
