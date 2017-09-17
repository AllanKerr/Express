package main

import (
	"github.com/gocql/gocql"
	"github.com/sirupsen/logrus"
)

type CQLDatastore struct {
	session *gocql.Session
}

func NewCQLDatastore(host string, keyspace string) (*CQLDatastore, error) {

	logrus.WithFields(logrus.Fields{
		"host": host,
		"keyspace" : keyspace,
	}).Info("Creating new CQL datastore.")

	cluster := gocql.NewCluster(host)
	cluster.Keyspace = keyspace

	session, err := cluster.CreateSession()
	if err != nil {
		logrus.WithField("", err).Fatal("Failed to create a new CQL datastore.")
		return nil, err
	}
	ds := new(CQLDatastore)
	ds.session = session
	return ds, nil
}

func (ds CQLDatastore)GetSession() interface{} {
	return ds.session
}

func (ds CQLDatastore)Close() {
	ds.session.Close()
}
