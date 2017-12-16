package client

import (
	"github.com/spf13/cobra"
	"app/core"
	"app/oauth2"
	"github.com/scylladb/gocqlx/qb"
	"github.com/scylladb/gocqlx"
	"errors"
	"strconv"
	"github.com/sirupsen/logrus"
	"github.com/gocql/gocql"
	"fmt"
)

var cmdCreate = &cobra.Command{
	Use:   "create",
	Short: "Creates a new oauth2 client",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 2 {
			if _, err := strconv.ParseBool(args[1]); err != nil {
				return errors.New("expected boolean public")
			}
		} else if len(args) != 1 {
			return errors.New("requires owner")
		}
		return nil
	},
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

		// parse args
		var public bool
		if len(args) == 1 {
			public = true
		} else {
			public, _ = strconv.ParseBool(args[1])
		}
		c, err := oauth2.NewClient(args[0], public)

		// create insert query
		stmt, names := qb.Insert("default.clients").
			Columns("id", "secret_hash", "redirect_uris", "grant_types", "response_types", "scopes", "public").
			ToCql()

		// bind the new client to be inserted
		q := gocqlx.Query(session.Query(stmt), names).BindStruct(c)
		if err := q.ExecRelease(); err != nil {
			logrus.WithField("error", err).Fatal("insert failed")
		}
		fmt.Println("client_id: " + c.Id + "\nclient_secret: " + c.Secret)
	},
}

func init() {
	ClientCmd.AddCommand(cmdCreate)
}