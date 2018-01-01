package cmd

import (
	"testing"
	"os"
	"app/core"
	"app/oauth2"
	"github.com/ory/fosite"
	"github.com/pkg/errors"
)

// create a new test user
func createTestUser(username string) error {

	databaseUrl := os.Getenv("DATABASE_URL")
	ds, err := core.NewCQLDataStore(databaseUrl, "default", 1)
	if err != nil {
		return err
	}
	defer ds.Close()

	adapter := oauth2.NewDataStoreAdapter(ds, config.GetHasher())

	user := oauth2.NewUser(username, "password")
	return adapter.CreateUser(user)
}

// load the specified user from the data store
func getTestUser(username string) (oauth2.User, error) {

	databaseUrl := os.Getenv("DATABASE_URL")
	ds, err := core.NewCQLDataStore(databaseUrl, "default", 1)
	if err != nil {
		return nil, err
	}
	defer ds.Close()

	adapter := oauth2.NewDataStoreAdapter(ds, config.GetHasher())
	return adapter.GetUser(nil, username)
}

// test creating a client and adding a scope to the client
func TestClient_User(t *testing.T) {

	name := "cmdtest"

	if err := createTestUser(name); err != nil {
		t.Errorf("Error creating test user: %v", err)
	}
	if err := addUserScope(name, "test_scope"); err != nil {
		t.Errorf("Error adding scope: %v", err)
	}

	user, err := getTestUser(name)
	if err != nil {
		t.Errorf("Error verifying user scopes: %v", err)
	}
	if !user.GetScopes().Has("test_scope") {
		t.Errorf("Scope not added")
	}

	// test adding a scope to a user that doesn't exist
	if err := addUserScope("non_existent", "test_scope"); err == nil || errors.Cause(err) != fosite.ErrNotFound  {
		t.Errorf("Error adding scope for non-existent user: %v", err)
	}
}