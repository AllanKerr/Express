package main

import (
	"github.com/gocql/gocql"
	"github.com/sirupsen/logrus"
	"time"
)

type CQLDatastore struct {
	session *gocql.Session
}

func newCluster(host string, keyspace string) *gocql.ClusterConfig {

	logrus.WithFields(logrus.Fields{
		"host": host,
		"keyspace" : keyspace,
	}).Info("Creating new CQL datastore.")

	cluster := gocql.NewCluster(host)
	cluster.Keyspace = keyspace
	return cluster
}

func NewCQLDatastoreRetry(host string, keyspace string, interval int) *CQLDatastore {

    cluster := newCluster(host, keyspace)
	session, err := cluster.CreateSession()
	for err != nil {
		logrus.WithField("error", err).Error("Failed to create a new CQL datastore.")
		time.Sleep(time.Duration(interval) * time.Second)
		session, err = cluster.CreateSession()
	}
	ds := new(CQLDatastore)
	ds.session = session
	return ds
}

func NewCQLDatastore(host string, keyspace string) (*CQLDatastore, error) {

	cluster := newCluster(host, keyspace)
	session, err := cluster.CreateSession()
	if err != nil {
		logrus.WithField("error", err).Fatal("Failed to create a new CQL datastore.")
		return nil, err
	}
	ds := new(CQLDatastore)
	ds.session = session
	return ds, nil
}

func (ds CQLDatastore) GetSession() interface{} {
	return ds.session
}

func (ds CQLDatastore) Close() {
	ds.session.Close()
}
