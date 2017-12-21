package admin

import (
	"github.com/spf13/cobra"
	"app/core"
	"app/oauth2"
	"github.com/scylladb/gocqlx/qb"
	"github.com/scylladb/gocqlx"
	"github.com/sirupsen/logrus"
	"github.com/gocql/gocql"
	"fmt"
)

var cmdCreate = &cobra.Command{
	Use:   "create",
	Short: "Creates a new admin user",
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {

		// start a new CQL session
		ds, err := core.NewCQLDataStore("cassandra-0.cassandra:9042", "default")
		if err != nil {
			logrus.WithField("error", err).Fatal("failed to create data store session")
		}
		session, ok := ds.GetSession().(*gocql.Session)
		if !ok {
			logrus.WithField("error", err).Fatal("invalid session: expected CQL data store")
		}

		u, err := oauth2.NewAdmin(args[0], args[1])

		// create insert query
		stmt, names := qb.Insert("default.users").
			Columns("username", "password_hash").
			ToCql()

		// bind the new client to be inserted
		q := gocqlx.Query(session.Query(stmt), names).BindStruct(u)
		if err := q.ExecRelease(); err != nil {
			logrus.WithField("error", err).Fatal("insert failed")
		}
		fmt.Println("admin created: " + u.Username)
	},
}

func init() {
	AdminCmd.AddCommand(cmdCreate)
}