package cmd

import (
	"github.com/spf13/cobra"
)

var deployCmd = &cobra.Command{
	Use:   "deploy <name> <image>",
	Args:cobra.ExactArgs(2),
	Run: handler.Deploy,
}

func init() {

	deployCmd.Flags().Int32("min", 1, "The minimum number of instances.")
	deployCmd.Flags().Int32("max", -1, "The minimum number of instances.")
	deployCmd.Flags().Int("cpu-percent", 80, "The CPU usage threshold at which to scale up")
	deployCmd.Flags().Int32("port", 80, "The port exposed by the Docker image.")
	deployCmd.Flags().String("endpoint-config", "", "The configuration file for protecting the deployed API.")

	RootCmd.AddCommand(deployCmd)
}