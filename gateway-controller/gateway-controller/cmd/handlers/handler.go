package handlers

import "gateway-controller/kube"

type CommandHandler struct {
	Client kube.Client
}
