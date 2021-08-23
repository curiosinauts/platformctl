package cmd

import (
	"fmt"

	"github.com/curiosinauts/platformctl/internal/msg"
	"github.com/curiosinauts/platformctl/pkg/regutil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// listReposCmd represents the image command
var listReposCmd = &cobra.Command{
	Use:     "repositories",
	Aliases: []string{"repos", "repo"},
	Short:   "Lists repositories",
	Long:    `List repositories`,
	Run: func(cmd *cobra.Command, args []string) {

		url, ok := viper.Get("docker_registry_url").(string)
		if !ok {
			msg.Failure("getting tag list: PLATFORMCTL_DOCKER_REGISTRY_URL env is required")
		}
		eh := ErrorHandler{"getting next docker tag"}
		hub := regutil.NewRegistryClient(url, nextTagCmdDebug)

		repositories, err := hub.Repositories()
		eh.HandleError("checking manifest", err)

		for _, repo := range repositories {
			fmt.Println(repo)
		}
	},
}

func init() {
	listCmd.AddCommand(listReposCmd)
}
