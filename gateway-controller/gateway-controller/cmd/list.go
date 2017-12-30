package cmd

import (
	"github.com/spf13/cobra"
)

// The list command lists the deployed application containers
// that have been deployed using the deploy command
//
// https://github.com/AllanKerr/Express/blob/master/docs/gateway/list-command.md
var listCmd = &cobra.Command{
	Use:   "list",
	Args:cobra.ExactArgs(0),
	Run: handler.List,
}

func init() {
	RootCmd.AddCommand(listCmd)
}