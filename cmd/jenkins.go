package cmd

import (
	"fmt"

	"github.com/curiosinauts/platformctl/pkg/jenkinsutil"
	"github.com/spf13/cobra"
)

// jekinsCmd represents the jekins command
var jenkinsCmd = &cobra.Command{
	Use:   "jenkins",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		username := args[0]

		option := map[string]string{
			"USERNAME": username,
		}

		jenkins, err := jenkinsutil.NewJenkins("https://jenkins.int.curiosityworks.org/", "admin", "116350724f59868b6efc6e8cef07f18bbd")
		if err != nil {
			fmt.Println(err)
		}

		jenkins.BuildJob("codeserver", option)
	},
}

func init() {
	rootCmd.AddCommand(jenkinsCmd)
}
