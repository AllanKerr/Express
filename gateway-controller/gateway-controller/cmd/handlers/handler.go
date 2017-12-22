package handlers

import "gateway-controller/kube"

type CommandHandler struct {
	client *kube.Client
}

func NewCommandHandler(client *kube.Client) *CommandHandler {
	return &CommandHandler{
		client,
	}
}