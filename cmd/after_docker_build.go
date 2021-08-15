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
	Short: "Deletes files that were generated during before docker-build cmd",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println()

		os.RemoveAll("./.ssh")

		os.Remove(".gitconfig")

		os.Remove("config.yml")

		os.Remove("gotty.sh")

		os.Remove("repositories.txt")

		os.Remove("runtime_install.sh")

		msg.Success("after docker-build")
	},
}

func init() {
	afterCmd.AddCommand(afterDockerBuildCmd)
}
