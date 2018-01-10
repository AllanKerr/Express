package handlers

import "express-controller/kube"

type CommandHandler struct {
	Client kube.Client
}
