package cmd

import (
	"github.com/spf13/cobra"
)

var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Create, update, and delete clients.",
}

func init() {
	RootCmd.AddCommand(clientCmd)
}