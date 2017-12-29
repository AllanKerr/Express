package core

import (
	"github.com/gocql/gocql"
	"github.com/sirupsen/logrus"
	"time"
	"fmt"
)

type CQLDataStore struct {
	session *gocql.Session
}

func createCluster(host string) *gocql.ClusterConfig {

	logrus.WithField("host", host).Info("Creating new CQL session.")

	cluster := gocql.NewCluster(host)
	cluster.Consistency = gocql.Quorum
	return cluster
}

func createTable(s *gocql.Session, table string) error {
	if err := s.Query(table).RetryPolicy(nil).Exec(); err != nil {
		return err
	}
	return nil
}

func createKeyspace(cluster *gocql.ClusterConfig, keyspace string, replicationFactor int) error {

	c := *cluster
	c.Keyspace = "system"
	session, err := c.CreateSession()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
			"keyspace": keyspace,
		}).Info("Failed to create keyspace.")
		return err
	}
	defer session.Close()

	return createTable(session, fmt.Sprintf(`CREATE KEYSPACE IF NOT EXISTS "%s" WITH replication = {
		'class' : 'SimpleStrategy',
		'replication_factor' : %d
	}`, keyspace, replicationFactor))
}

func createSession(cluster *gocql.ClusterConfig, keyspace string, replicationFactor int) (*gocql.Session, error) {

	if err := createKeyspace(cluster, keyspace, replicationFactor); err != nil {
		return nil, err
	}
	cluster.Keyspace = keyspace
	return cluster.CreateSession()
}

func NewCQLDataStoreRetry(host string, keyspace string, replicationFactor int, interval int) *CQLDataStore {

    cluster := createCluster(host)
	session, err := createSession(cluster, keyspace, replicationFactor)
	for err != nil {
		logrus.WithField("error", err).Error("Failed to create a new CQL datastore.")
		time.Sleep(time.Duration(interval) * time.Second)
		session, err = createSession(cluster, keyspace, replicationFactor)
	}
	ds := new(CQLDataStore)
	ds.session = session
	return ds
}

func NewCQLDataStore(host string, keyspace string, replicationFactor int) (*CQLDataStore, error) {

	cluster := createCluster(host)
	session, err := createSession(cluster, keyspace, replicationFactor)
	if err != nil {
		logrus.WithField("error", err).Fatal("Failed to create a new CQL datastore.")
		return nil, err
	}
	ds := new(CQLDataStore)
	ds.session = session
	return ds, nil
}

func (ds *CQLDataStore) CreateTable(schema string) error {
	return createTable( ds.session, schema)
}

func (ds CQLDataStore) GetSession() interface{} {
	return ds.session
}

func (ds CQLDataStore) Close() {
	ds.session.Close()
}
