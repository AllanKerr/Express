package client

import (
	"github.com/spf13/cobra"
	"fmt"
	"strings"
)

var cmdCreate = &cobra.Command{
	Use:   "create",
	Short: "Creates a new oauth2 client",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Print: " + strings.Join(args, " "))
	},
}

func init() {
	ClientCmd.AddCommand(cmdCreate)
}