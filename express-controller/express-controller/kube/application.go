package kube

import "time"

// Interface for a deployed application
// Abstracts over the individual deployment, service,
// load balancer, and Ingress configurations of the application
type Application interface {
	GetName() string
	GetPort() int32
	GetCreationTimestamp() time.Time
}

type DefaultApplication struct {
	name string
	port int32
	creation time.Time
}

func (container *DefaultApplication) GetName() string {
	return container.name
}

func (container *DefaultApplication) GetPort() int32 {
	return container.port
}

func (container *DefaultApplication) GetCreationTimestamp() time.Time {
	return container.creation
}
