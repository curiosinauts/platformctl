package cmd

import (
	"strings"

	"github.com/curiosinauts/platformctl/internal/msg"
	"github.com/curiosinauts/platformctl/pkg/jenkinsutil"
	"github.com/spf13/cobra"
)

var runJenkinsJobCmdParams []string

// runJenkinsJobCmd represents the jenkinsJob command
var runJenkinsJobCmd = &cobra.Command{
	Use:     "jenkins-job",
	Short:   "Runs Jenkins job",
	Long:    `Runs Jenkins job`,
	Args:    cobra.MinimumNArgs(1),
	Example: "platformctl run jenkins-job upgrade-platformctl 2.3.0",
	Run: func(cmd *cobra.Command, args []string) {
		jobName := args[0]

		params := map[string]string{}
		for _, p := range runJenkinsJobCmdParams {
			terms := strings.Split(p, "=")
			key := strings.TrimSpace(terms[0])
			value := strings.TrimSpace(terms[1])
			params[key] = value
		}

		eh := ErrorHandler{"running jenkins job"}

		jenkins, err := jenkinsutil.NewJenkins()
		eh.HandleError("accessing Jenkins job", err)

		_, err = jenkins.BuildJob(jobName, params)
		eh.HandleError("calling Jenkins job to release new platformctl to Jenkins environment", err)

		msg.Success("running Jenkins job")
	},
}

func init() {
	runCmd.AddCommand(runJenkinsJobCmd)
	runJenkinsJobCmd.Flags().StringArrayVarP(&runJenkinsJobCmdParams, "parameters", "p", []string{}, "specify jenkins job parameters")
}
