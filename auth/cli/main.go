package main

import (
	"github.com/spf13/cobra"
	"cli/client"
)

var RootCmd = &cobra.Command{
	Use:   "cli",
}

func main() {

	RootCmd.AddCommand(client.ClientCmd)
	RootCmd.Execute()
}
