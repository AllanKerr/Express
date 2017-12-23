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

func (ch *CommandHandler) updateDeployment(name string, flags *pflag.FlagSet) {

	updater := kube.NewDeploymentUpdater(ch.Client, apiv1.NamespaceDefault)
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
		fmt.Println(err.Error())
	}
}

func (ch *CommandHandler) Update(command *cobra.Command, args []string) {

	name := args[0]
	flags := command.Flags()

	ch.updateDeployment(name, flags)


}