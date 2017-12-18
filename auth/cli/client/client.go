package client

import (
	"github.com/spf13/cobra"
	"cli/client/grant"
)

var ClientCmd = &cobra.Command{
	Use:   "client",
	Short: "Create, update, and delete clients.",
}

func init() {
	ClientCmd.AddCommand(grant.ClientGrantCmd)
}