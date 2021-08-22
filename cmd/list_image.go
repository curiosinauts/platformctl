package cmd

import (
	"github.com/spf13/cobra"
)

// imageCmd represents the image command
var imageCmd = &cobra.Command{
	Use:   "image",
	Short: "Lists images",
	Long:  `List images`,
	Run: func(cmd *cobra.Command, args []string) {

		// url, ok := viper.Get("docker_registry_url").(string)
		// if !ok {
		// 	msg.Failure("getting tag list: PLATFORM_DOCKER_REGISTRY_URL env is required")
		// }
		// eh := ErrorHandler{"getting next docker tag"}
		// hub := NewRegistryClient(url, nextTagCmdDebug)

		// digest, err := hub.ManifestDigest("7onetella/base", "1.0.0")
		// eh.HandleError("checking manifest", err)

		// msg.Info(digest.String())
	},
}

func init() {
	listCmd.AddCommand(imageCmd)
}
