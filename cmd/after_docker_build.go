package cmd

import (
	"fmt"
	"os"

	"github.com/curiosinauts/platformctl/internal/msg"
	"github.com/spf13/cobra"
)

// afterDockerBuildCmd represents the dockerBuild command
var afterDockerBuildCmd = &cobra.Command{
	Use:   "docker-build",
	Short: "Deletes the files that were generated during before docker-build cmd",
	Long:  `Deletes the files that were generated during before docker-build cmd`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println()

		os.Remove("./.ssh/id_rsa")

		os.Remove(".gitconfig")

		os.Remove("config.yml")

		os.Remove("gotty.sh")

		os.Remove("repositories.txt")

		os.Remove("runtime_install.sh")

		if len(args) > 0 {
			username := args[0]
			os.Remove("vscode-" + username + ".yml")
		}

		msg.Success("after docker-build")
	},
}

func init() {
	afterCmd.AddCommand(afterDockerBuildCmd)
}
