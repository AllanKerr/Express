package cmd

import (
	"github.com/spf13/cobra"
	"fmt"
)

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Run: func ( command *cobra.Command, args []string) {
		fmt.Println("APP")
	},
}

func init() {
	RootCmd.AddCommand(deployCmd)
}