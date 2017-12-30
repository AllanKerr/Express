package handlers

import (
	"github.com/spf13/cobra"
	"fmt"
	"text/tabwriter"
	"os"
	apiv1 "k8s.io/api/core/v1"
	"gateway-controller/kube"
	"io"
)

// Print a list of applications in a tab delimited table
func printApplications(w io.Writer, apps []kube.Application) {

	fmt.Fprintln(w, "NAME\tPORT\tCREATION\t")
	for _, app := range apps {
		fmt.Fprintf(w, "%v\t", app.GetName())
		fmt.Fprintf(w, "%v\t", app.GetPort())
		fmt.Fprintf(w, "%v\t\n", app.GetCreationTimestamp())
	}
}

// List the set of deployed applications
func (ch *CommandHandler) List(command *cobra.Command, args []string) {

	applications, err := ch.Client.ListApplications(apiv1.NamespaceDefault)
	if err != nil {
		fmt.Printf("Unable to list applications: %v\n", err)
		return
	}
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 0, 4, ' ', 0)
	printApplications(w, applications)
	w.Flush()
}
