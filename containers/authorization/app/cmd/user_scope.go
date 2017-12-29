package cmd

import "github.com/spf13/cobra"

var userScopeCmd = &cobra.Command{
	Use:   "scope",
	Short: "Add and remove scopes from existing users.",
}

func init() {
	userCmd.AddCommand(userScopeCmd)
}