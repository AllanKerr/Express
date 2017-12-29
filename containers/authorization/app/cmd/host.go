package cmd

import (
	"github.com/spf13/cobra"
	"app/server"
)

// hostCmd represents the host command
var hostCmd = &cobra.Command{
	Use:   "host",
	Short: "Start the HTTP auth service",
	Run: func(cmd *cobra.Command, args []string) {
		server.RunHost(config)
	},
}

func init() {
	RootCmd.AddCommand(hostCmd)
}