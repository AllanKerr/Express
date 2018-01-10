package cmd

import (
	"github.com/spf13/cobra"
	"express-controller/cmd/handlers"
	"express-controller/kube"
	"fmt"
)

var handler = &handlers.CommandHandler{}

var RootCmd = &cobra.Command{
	Use:   "app",
}

func init() {

	c, err := kube.NewDefaultClient()
	if err != nil {
		fmt.Println(err)
		panic("error")
	}
	handler.Client = c

}
