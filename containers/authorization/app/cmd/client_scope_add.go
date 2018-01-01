package cmd

import (
	"github.com/spf13/cobra"
	"app/core"
	"app/oauth2"
	"fmt"
	"os"
)

func addClientScope(clientId string, scope string) error {

	// start a new CQL session
	databaseUrl := os.Getenv("DATABASE_URL")
	ds, err := core.NewCQLDataStore(databaseUrl, "authorization", 3)
	if err != nil {
		return err
	}
	defer ds.Close()

	adapter := oauth2.NewDataStoreAdapter(ds, config.GetHasher())
	c, err := adapter.GetClient(nil, clientId)
	if c == nil {}
	if err != nil {
		return err
	}
	client, ok := c.(*oauth2.Client)
	if !ok {
		return err
	}
	client.AppendScope(scope)
	err = adapter.UpdateClient(client, "scopes")
	if err != nil {
		return err
	}
	return nil
}

var cmdClientScopeAdd = &cobra.Command{
	Use:   "add <id> <scope>",
	Short: "Add a scope to an existing client",
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {

		if err := addClientScope(args[0], args[1]); err != nil {
			fmt.Println()
			return
		}
		fmt.Println("added scope: " + args[1])
	},
}

func init() {
	clientScopeCmd.AddCommand(cmdClientScopeAdd)
}