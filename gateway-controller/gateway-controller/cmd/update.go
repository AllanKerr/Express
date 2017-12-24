package cmd

import (
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update <name>",
	Args:cobra.ExactArgs(1),
	Run: handler.Update,
}

func init() {

	updateCmd.Flags().String("image", "", "The new Docker image to roll out.")
	updateCmd.Flags().Int32("min", 1, "The minimum number of instances.")
	updateCmd.Flags().Int32("max", -1, "The minimum number of instances.")
	updateCmd.Flags().String("endpoint-config", "", "The configuration file for accessing the deployed API.")

	RootCmd.AddCommand(updateCmd)
}