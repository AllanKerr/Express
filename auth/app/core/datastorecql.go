package core

import (
	"github.com/gocql/gocql"
	"github.com/sirupsen/logrus"
	"time"
)

type CQLDataStore struct {
	session *gocql.Session
}

func newCluster(host string, keyspace string) *gocql.ClusterConfig {

	logrus.WithFields(logrus.Fields{
		"host": host,
		"keyspace" : keyspace,
	}).Info("Creating new CQL session.")

	cluster := gocql.NewCluster(host)
	cluster.Keyspace = keyspace
	return cluster
}

func NewCQLDataStoreRetry(host string, keyspace string, interval int) *CQLDataStore {

    cluster := newCluster(host, keyspace)
	session, err := cluster.CreateSession()
	for err != nil {
		logrus.WithField("error", err).Error("Failed to create a new CQL datastore.")
		time.Sleep(time.Duration(interval) * time.Second)
		session, err = cluster.CreateSession()
	}
	ds := new(CQLDataStore)
	ds.session = session
	return ds
}

func NewCQLDataStore(host string, keyspace string) (*CQLDataStore, error) {

	cluster := newCluster(host, keyspace)
	session, err := cluster.CreateSession()
	if err != nil {
		logrus.WithField("error", err).Fatal("Failed to create a new CQL datastore.")
		return nil, err
	}
	ds := new(CQLDataStore)
	ds.session = session
	return ds, nil
}

func (ds CQLDataStore) GetSession() interface{} {
	return ds.session
}

func (ds CQLDataStore) Close() {
	ds.session.Close()
}
