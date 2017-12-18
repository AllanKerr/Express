package util

import (
	"github.com/gocql/gocql"
	"github.com/sirupsen/logrus"
	"app/core"
	"github.com/scylladb/gocqlx/qb"
	"github.com/scylladb/gocqlx"
	"app/oauth2"
)

var CassandraKeySpace string = "default"

func CreateCassandraSession() (*gocql.Session, error) {

	ds, err := core.NewCQLDataStore("cassandra-0.cassandra:9042", "default")
	if err != nil {
		logrus.WithField("error", err).Fatal("failed to create data store session")
	}
	session, ok := ds.GetSession().(*gocql.Session)
	if !ok {
		logrus.WithField("error", err).Fatal("invalid session: expected CQL data store")
	}
	return session, nil
}

func GetClient(session *gocql.Session, id string) (*oauth2.Client, error) {

	stmt, names := qb.Select("default.clients").
		Where(qb.Eq("id")).
		ToCql()

	q := gocqlx.Query(session.Query(stmt), names).BindMap(qb.M{
		"id": id,
	})

	var c oauth2.Client
	if err := gocqlx.Get(&c, q.Query); err != nil {
		return nil, err
	}
	return &c, nil
}