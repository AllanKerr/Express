package cmd

import (
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Args:cobra.ExactArgs(0),
	Run: handler.List,
}

func init() {
	RootCmd.AddCommand(listCmd)
}