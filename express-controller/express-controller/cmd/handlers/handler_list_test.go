package handlers

import (
	"testing"
	"express-controller/kube"
	"time"
)

// Mock client implementation for ListApplications
type mockClient struct {
	kube.Client
}

// Mock implementation for Application interface
type mockApplication struct {
	name string
	port int32
	creation time.Time
}

func (container *mockApplication) GetName() string {
	return container.name
}

func (container *mockApplication) GetPort() int32 {
	return container.port
}

func (container *mockApplication) GetCreationTimestamp() time.Time {
	return container.creation
}

func (client *mockClient) ListApplications(namespace string) ([]kube.Application, error) {

	var apps []kube.Application
	for i := 0; i < 5; i++ {
		apps = append(apps, &mockApplication{
			name: "abc",
			port: 21,
			creation: time.Now(),
		})
	}
	return apps, nil
}

// Test that no errors occur during list
func TestCommandHandler_List(t *testing.T) {

	handler := CommandHandler{
		Client: &mockClient{},
	}
	handler.List(nil, nil)
}
