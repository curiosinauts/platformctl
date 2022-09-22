package cmd

import (
	"github.com/spf13/cobra"
)

// rungrafanaCmdCmd represents the grafana command
var grafanaCmd = &cobra.Command{
	Use:   "grafana",
	Short: "grafana",
	Long:  `grafana`,
}

func init() {
	rootCmd.AddCommand(grafanaCmd)
}
