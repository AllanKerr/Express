package main

import (
	"github.com/spf13/cobra"
	"cli/client"
	"cli/admin"
)

var RootCmd = &cobra.Command{
	Use:   "cli",
}

func main() {

	RootCmd.AddCommand(client.ClientCmd, admin.AdminCmd)
	RootCmd.Execute()
}
