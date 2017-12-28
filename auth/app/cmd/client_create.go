package cmd

import (
	"github.com/spf13/cobra"
	"app/core"
	"app/oauth2"
	"fmt"
	"os"
)

var cmdClientCreate = &cobra.Command{
	Use:   "create <id> <secret>",
	Short: "Creates a new OAuth2 client",
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {

		// start a new CQL session
		databaseUrl := os.Getenv("DATABASE_URL")
		ds, err := core.NewCQLDataStore(databaseUrl, "default")
		if err != nil {
			fmt.Errorf("failed to create data store session %g", err)
		}
		defer ds.Close()

		public, _ := cmd.Flags().GetBool("public")
		client := oauth2.NewClient(args[0], args[1], public)

		adapter := oauth2.NewDataStoreAdapter(ds, config.GetHasher())
		if err := adapter.CreateClient(client); err != nil {
			fmt.Errorf("failed to create client %g", err)
		}
		fmt.Println("client_id: " + client.Id + "\nclient_secret: " + client.Secret)
	},
}

func init() {
	clientCmd.AddCommand(cmdClientCreate)
	cmdClientCreate.Flags().Bool("public", false, "Use to create a public client ")
}