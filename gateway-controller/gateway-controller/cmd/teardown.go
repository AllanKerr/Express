package cmd

import (
	"github.com/spf13/cobra"
)

var teardownCmd = &cobra.Command{
	Use:   "teardown <name>",
	Args:cobra.ExactArgs(1),
	Run: handler.Teardown,
}

func init() {
	RootCmd.AddCommand(teardownCmd)
}