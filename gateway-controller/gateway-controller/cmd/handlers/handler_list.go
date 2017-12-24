package handlers

import (
	"github.com/spf13/cobra"
	"fmt"
	"text/tabwriter"
	"os"
	apiv1 "k8s.io/api/core/v1"
)

func (ch *CommandHandler) List(command *cobra.Command, args []string) {

	services, err := ch.Client.ListServices(apiv1.NamespaceDefault)
	if err != nil {
		fmt.Printf("Unable to list services: %v\n", err)
	}

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 0, 4, ' ', 0)
	fmt.Fprintln(w, "NAME\tPORT\tCREATION\t")
	for _, service := range services {
		fmt.Fprintf(w, "%v\t", service.GetName())
		fmt.Fprintf(w, "%v\t", service.Spec.Ports[0].Port)
		fmt.Fprintf(w, "%v\t\n", service.GetCreationTimestamp().String())
	}
	w.Flush()
}