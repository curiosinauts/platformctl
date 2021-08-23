package cmd

import (
	"github.com/curiosinauts/platformctl/internal/msg"
	"github.com/curiosinauts/platformctl/pkg/regutil"
	"github.com/spf13/cobra"
)

// deleteImageCmd represents the image command
var deleteImagesCmd = &cobra.Command{
	Use:        "images",
	Short:      "Deletes docker image from private registry",
	Long:       `Deletes docker image from private registry`,
	ArgAliases: []string{"repository"},
	PreRunE:    cobra.MinimumNArgs(1),
	ValidArgs:  []string{"repository"},
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
	deleteCmd.AddCommand(deleteImagesCmd)
}
