package cmd

import (
	"github.com/spf13/cobra"
)

// Teardown a previously deployed application container resulting
// in all of its resources being deleted.
//
// https://github.com/AllanKerr/Express/blob/master/docs/gateway/teardown-command.md
var teardownCmd = &cobra.Command{
	Use:   "teardown <name>",
	Args:cobra.ExactArgs(1),
	Run: handler.Teardown,
}

func init() {
	RootCmd.AddCommand(teardownCmd)
}