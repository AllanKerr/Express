package cmd

import "github.com/spf13/cobra"

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Create, update, and delete users.",
}

func init() {
	RootCmd.AddCommand(userCmd)
}