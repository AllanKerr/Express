package handlers

import (
	apiv1 "k8s.io/api/core/v1"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"gateway-controller/kube"
	"fmt"
)

func RequiresUpdate(flags *pflag.FlagSet, updater kube.ObjectUpdater) bool {

	for _, mod := range updater.GetModifiers() {
		if flags.Changed(mod) {
			return true
		}
	}
	return false
}

func (ch *CommandHandler) updateDeployment(namespace string, name string, flags *pflag.FlagSet) {

	updater := kube.NewDeploymentUpdater(ch.Client, namespace)
	if !RequiresUpdate(flags, updater) {
		return
	}

	var imagePtr *string
	image, _ := flags.GetString("image")
	if flags.Changed("image") {
		imagePtr = &image
	}

	update := &kube.ContainerUpdate{
		Image: imagePtr,
	}
	if err := updater.Update(name, update); err != nil {
		fmt.Printf("Unable to update deployment: %v\n", err.Error())
	} else {
		fmt.Println("Updated deployment.")
	}
}

func (ch *CommandHandler) updateAutoscaler(namespace string, name string, flags *pflag.FlagSet) {

	updater := kube.NewAutoscalerUpdater(ch.Client, namespace)
	if !RequiresUpdate(flags, updater) {
		return
	}

	var minPtr *int32
	var maxPtr *int32
	min, _ := flags.GetInt32("min")
	max, _ := flags.GetInt32("max")
	if flags.Changed("min") {
		minPtr = &min
	}
	if flags.Changed("max") {
		maxPtr = &max
	}

	update := &kube.AutoscalerUpdate {
		MinReplicas: minPtr,
		MaxReplicas: maxPtr,
	}
	if err := updater.Update(name, update); err != nil {
		fmt.Printf("Unable to update autoscaler: %v\n", err.Error())
	} else {
		fmt.Println("Updated autoscaler.")
	}
}

func (ch *CommandHandler) Update(command *cobra.Command, args []string) {

	name := args[0]
	namespace := apiv1.NamespaceDefault
	flags := command.Flags()

	ch.updateDeployment(namespace, name, flags)
	ch.updateAutoscaler(namespace, name, flags)


}