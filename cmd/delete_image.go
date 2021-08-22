package cmd

import (
	"github.com/curiosinauts/platformctl/internal/msg"
	"github.com/curiosinauts/platformctl/pkg/regutil"
	"github.com/spf13/cobra"
)

// deleteImageCmd represents the image command
var deleteImageCmd = &cobra.Command{
	Use:   "image",
	Short: "Deletes docker image from private registry",
	Long:  `Deletes docker image from private registry`,
	Run: func(cmd *cobra.Command, args []string) {
		repository := args[0]

		eh := ErrorHandler{"deleting docker image from private repository"}
		tags, err := regutil.ListTags(repository, false)
		eh.PrintError("listing tags", err)

		for _, tag := range tags {
			msg.Info("deleting tag " + tag)
			err = regutil.DeleteImage(repository, tag, false)
			eh.PrintError("deleting image", err)
		}
	},
}

func init() {
	deleteCmd.AddCommand(deleteImageCmd)
}
