package cmd

import (
	"testing"
	"os"
	"app/core"
	"app/server"
	"github.com/sirupsen/logrus"
)

func TestMain(m *testing.M) {

	os.Chdir("..")

	databaseUrl := os.Getenv("DATABASE_URL")
	ds := core.NewCQLDataStoreRetry(databaseUrl, "authorization", 1, 5)

	if err := server.CreateSchema(ds); err != nil {
		logrus.WithField("error", err).Error("Failed to create schema.")
	}
	retCode := m.Run()
	ds.Close()
	os.Exit(retCode)
}

// test creating a client and adding a scope to the client
func TestClient_Client(t *testing.T) {

	if err := createClient("test_client", "test_secret", false); err != nil {
		t.Errorf("Error creating client: %v", err)
	}
	if err := addClientScope("test_client", "testscope"); err != nil {
		t.Errorf("Error adding scope: %v", err)
	}
}