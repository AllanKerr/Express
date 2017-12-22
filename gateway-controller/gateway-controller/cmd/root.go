package cmd

import (
	"github.com/spf13/cobra"
	"gateway-controller/cmd/handlers"
	"gateway-controller/kube"
	"fmt"
)

var handler *handlers.CommandHandler

var RootCmd = &cobra.Command{
	Use:   "app",
}

func init() {

	c, err := kube.NewDefaultClient()
	if err != nil {
		fmt.Println(err)
		return
	}
	handler = handlers.NewCommandHandler(c)
}