package cmd

import (
	"github.com/spf13/cobra"
	"app/core"
	"app/oauth2"
	"fmt"
	"os"
)

func createClient(clientId string, clientSecret string, public bool) error {

	databaseUrl := os.Getenv("DATABASE_URL")
	ds, err := core.NewCQLDataStore(databaseUrl, "authorization", 3)
	if err != nil {
		return err
	}
	defer ds.Close()

	client := oauth2.NewClient(clientId, clientSecret, public)

	adapter := oauth2.NewDataStoreAdapter(ds, config.GetHasher())
	if err := adapter.CreateClient(client); err != nil {
		return err
	}
	return nil
}

var cmdClientCreate = &cobra.Command{
	Use:   "create <id> <secret>",
	Short: "Creates a new OAuth2 client",
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {

		public, _ := cmd.Flags().GetBool("public")
		if err := createClient(args[0], args[1], public); err != nil {
			fmt.Errorf("failed to create client %v", err)
			return
		}
		fmt.Println("client_id: " + args[0] + "\nclient_secret: " + args[1])
	},
}

func init() {
	clientCmd.AddCommand(cmdClientCreate)
	cmdClientCreate.Flags().Bool("public", false, "Use to create a public client ")
}