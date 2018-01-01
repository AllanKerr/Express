package cmd

import (
	"github.com/spf13/cobra"
	"app/core"
	"app/oauth2"
	"fmt"
	"os"
	"github.com/sirupsen/logrus"
)

// add a new scope to an existing user
func addUserScope(username string, scope string) error {

	// start a new CQL session
	databaseUrl := os.Getenv("DATABASE_URL")
	ds, err := core.NewCQLDataStore(databaseUrl, "default", 3)
	if err != nil {
		return err
	}
	defer ds.Close()

	// get the user
	adapter := oauth2.NewDataStoreAdapter(ds, config.GetHasher())
	user, err := adapter.GetUser(nil, username)
	if err != nil {
		logrus.WithField("error", err).Info("Error getting user")
		return err
	}

	// add scope and update
	user.AppendScope(scope)
	return adapter.UpdateUser(user, "scopes")
}

// cli command for adding scopes to users
var cmdUserScopeAdd = &cobra.Command{
	Use:   "add <username> <scope>",
	Short: "Add a scope to an existing user",
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {

		if err := addUserScope(args[0], args[1]); err != nil {
			fmt.Errorf("failed to add scope %v", err)
			return
		}
		fmt.Println("added scope: " + args[1])
	},
}

func init() {
	userScopeCmd.AddCommand(cmdUserScopeAdd)
}