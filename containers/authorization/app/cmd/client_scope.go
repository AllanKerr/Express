package cmd

import "github.com/spf13/cobra"

var clientScopeCmd = &cobra.Command{
	Use:   "scope",
	Short: "Add and remove scopes from existing clients.",
}

func init() {
	clientCmd.AddCommand(clientScopeCmd)
}