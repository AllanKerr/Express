package handlers

import (
	"github.com/spf13/cobra"
	"gateway-controller/kube"
	apiv1 "k8s.io/api/core/v1"
	"fmt"
)

func (ch *CommandHandler) Teardown(command *cobra.Command, args []string) {

	name := args[0]
	namespace := apiv1.NamespaceDefault

	fmt.Println("Starting teardown...")
	if err := kube.NewIngressTransaction(ch.Client, namespace).Rollback(name); err != nil {
		fmt.Printf("Unable to delete ingresses: %v\n", err)
	} else {
		fmt.Printf("Deleted ingresses: %v\n", name)
	}
	if err := kube.NewServiceTransaction(ch.Client, namespace).Rollback(name); err != nil {
		fmt.Printf("Unable to delete service: %v\n", err)
	} else {
		fmt.Printf("Deleted service: %v\n", name)
	}
	if err := kube.NewAutoscalerTransaction(ch.Client, namespace).Rollback(name); err != nil {
		fmt.Printf("Unable to delete autoscaler: %v\n", err)
	} else {
		fmt.Printf("Deleted autoscaler: %v\n", name)
	}
	if err := kube.NewDeploymentTransaction(ch.Client, namespace).Rollback(name); err != nil {
		fmt.Printf("Unable to delete deployment: %v\n", err)
	} else {
		fmt.Printf("Deleted deployment: %v\n", name)
	}
	fmt.Println("Teardown complete")
}