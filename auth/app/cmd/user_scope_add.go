package cmd

import (
	"github.com/spf13/cobra"
	"app/core"
	"app/oauth2"
	"fmt"
)

var cmdUserScopeAdd = &cobra.Command{
	Use:   "add <username> <scope>",
	Short: "Add a scope to an existing user",
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {

		// start a new CQL session
		ds, err := core.NewCQLDataStore("cassandra-0.cassandra:9042", "default")
		if err != nil {
			fmt.Errorf("failed to create data store session %g", err)
		}
		defer ds.Close()

		adapter := oauth2.NewDataStoreAdapter(ds, config.GetHasher())
		user, err := adapter.GetUser(nil, args[0])
		if err != nil {
			fmt.Errorf("failed to get user: %g", err)
		}
		user.AppendScope(args[1])

		err = adapter.UpdateUser(user, "scopes")
		if err != nil {
			fmt.Errorf("failed to update user: %g", err)
		}
		fmt.Println("added scope: " + args[1])
	},
}

func init() {
	userScopeCmd.AddCommand(cmdUserScopeAdd)
}