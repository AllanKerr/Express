package grant

import (
	"github.com/spf13/cobra"
	"github.com/scylladb/gocqlx/qb"
	"github.com/scylladb/gocqlx"
	"fmt"
	"cli/util"
)

var cmdAdd = &cobra.Command{
	Use:   "add",
	Short: "Adds a new grant type to the specified client",
	Args: cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {

		// Create CQL session
		session, err := util.CreateCassandraSession()
		if err != nil {
			return
		}

		// Get and verify client
		client, err := util.GetClient(session, args[0])
		if err != nil {
			fmt.Println("invalid client")
			return
		}
		if client.VerifyPassword(args[1]) != nil {
			fmt.Println("invalid client secret")
			return
		}
		client.AppendGrant(args[2])

		stmt, names := qb.Update(util.CassandraKeySpace + ".clients").
			Set("grant_types").
			Where(qb.Eq("id")).
			ToCql()

		q := gocqlx.Query(session.Query(stmt), names).BindStruct(client)
		if err := q.ExecRelease(); err != nil {
			fmt.Println(err)
		}
		fmt.Println("grant successfully added")
	},
}

func init() {
	ClientGrantCmd.AddCommand(cmdAdd)
}