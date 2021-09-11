package cmd

import (
	"github.com/curiosinauts/platformctl/pkg/jenkinsutil"
	"github.com/spf13/cobra"
)

var runJenkinsJobCmdVars []string

// runJenkinsJobCmd represents the jenkinsJob command
var runJenkinsJobCmd = &cobra.Command{
	Use:     "jenkins-job",
	Short:   "Runs Jenkins job",
	Long:    `Runs Jenkins job`,
	Args:    cobra.MinimumNArgs(1),
	Example: "platformctl run jenkins-job upgrade-platformctl 2.3.0",
	Run: func(cmd *cobra.Command, args []string) {
		jobName := args[0]
		version := args[0]
		option := map[string]string{
			"PLATFORMCTL_VERSION": version,
		}

		eh := ErrorHandler{"running jenkins job"}

		jenkins, err := jenkinsutil.NewJenkins()
		eh.HandleError("accessing Jenkins job", err)

		_, err = jenkins.BuildJob(jobName, option)
		eh.HandleError("calling Jenkins job to release new platformctl to Jenkins environment", err)
	},
}

func init() {
	runCmd.AddCommand(runJenkinsJobCmd)
	runJenkinsJobCmd.Flags().StringArrayVarP(&addUserCmdRepos, "values", "v", []string{}, "specify ")
}
