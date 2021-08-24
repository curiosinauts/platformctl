package cmd

import (
	"fmt"

	"github.com/curiosinauts/platformctl/pkg/regutil"
	"github.com/spf13/cobra"
)

// listReposCmd represents the image command
var listReposCmd = &cobra.Command{
	Use:     "repositories",
	Aliases: []string{"repos", "repo"},
	Short:   "Lists repositories",
	Long:    `List repositories`,
	Run: func(cmd *cobra.Command, args []string) {

		eh := ErrorHandler{"getting next docker tag"}

		registryClient, err := regutil.NewRegistryClient(debug)
		eh.PrintError("getting registry client", err)

		repositories, err := registryClient.Repositories()
		eh.HandleError("checking manifest", err)

		for _, repo := range repositories {
			fmt.Println(repo)
		}
	},
}

func init() {
	listCmd.AddCommand(listReposCmd)
}
