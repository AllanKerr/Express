package cmd

import (
	"github.com/spf13/cobra"
	"app/core"
	"app/oauth2"
	"fmt"
	"os"
)

var cmdClientScopeAdd = &cobra.Command{
	Use:   "add <id> <scope>",
	Short: "Add a scope to an existing client",
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {

		// start a new CQL session
		databaseUrl := os.Getenv("DATABASE_URL")
		ds, err := core.NewCQLDataStore(databaseUrl, "authorization", 3)
		if err != nil {
			fmt.Errorf("failed to create data store session %g", err)
		}
		defer ds.Close()

		adapter := oauth2.NewDataStoreAdapter(ds, config.GetHasher())
		c, err := adapter.GetClient(nil, args[0])
		if err != nil {
			fmt.Errorf("failed to get client: %g", err)
		}
		client, ok := c.(*oauth2.Client)
		if !ok {
			fmt.Println("unexpected client type")
		}
		client.AppendScope(args[1])
		err = adapter.UpdateClient(client, "scopes")
		if err != nil {
			fmt.Errorf("failed to update client: %g", err)
		}
		fmt.Println("added scope: " + args[1])
	},
}

func init() {
	clientScopeCmd.AddCommand(cmdClientScopeAdd)
}