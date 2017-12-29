package core

import (
	"github.com/gocql/gocql"
	"github.com/sirupsen/logrus"
	"time"
	"fmt"
)

const SchemaTimeout = 5*time.Second

type CQLDataStore struct {
	session *gocql.Session
	cluster *gocql.ClusterConfig
}

func createCluster(host string) *gocql.ClusterConfig {

	logrus.WithField("host", host).Info("Creating new CQL session.")

	cluster := gocql.NewCluster(host)
	cluster.Consistency = gocql.Quorum
	return cluster
}

func createTable(s *gocql.Session, table string) error {

	q := s.Query(table)
	defer q.Release()

	if err := q.Exec(); err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
			"table": table,
		}).Error("Failed to create table.")
		return err
	}
	return nil
}

func createKeyspace(cluster *gocql.ClusterConfig, keyspace string, replicationFactor int) error {

	c := *cluster
	c.Keyspace = "system"
	c.Timeout = SchemaTimeout
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

func createSession(cluster *gocql.ClusterConfig) (*gocql.Session, error) {
	return cluster.CreateSession()
}

func createSessionAndKeyspace(cluster *gocql.ClusterConfig, keyspace string, replicationFactor int) (*gocql.Session, error) {

	if err := createKeyspace(cluster, keyspace, replicationFactor); err != nil {
		return nil, err
	}
	cluster.Keyspace = keyspace
	return createSession(cluster)
}

func NewCQLDataStoreRetry(host string, keyspace string, replicationFactor int, interval int) *CQLDataStore {

    cluster := createCluster(host)
	session, err := createSessionAndKeyspace(cluster, keyspace, replicationFactor)
	for err != nil {
		logrus.WithField("error", err).Error("Failed to create a new CQL datastore.")
		time.Sleep(time.Duration(interval) * time.Second)
		session, err = createSessionAndKeyspace(cluster, keyspace, replicationFactor)
	}
	return &CQLDataStore {
		session: session,
		cluster: cluster,
	}
}

func NewCQLDataStore(host string, keyspace string, replicationFactor int) (*CQLDataStore, error) {

	cluster := createCluster(host)
	session, err := createSessionAndKeyspace(cluster, keyspace, replicationFactor)
	if err != nil {
		logrus.WithField("error", err).Fatal("Failed to create a new CQL datastore.")
		return nil, err
	}
	return &CQLDataStore {
		session: session,
		cluster: cluster,
	}, nil
}

func (ds *CQLDataStore) CreateTable(object string) error {

	// Create session with increased timeout for schema creation
	ds.cluster.Timeout = SchemaTimeout
	session, err := createSession(ds.cluster)
	if err != nil {
		return err
	}
	defer session.Close()

	return createTable(session, object)
}

func (ds *CQLDataStore) CreateSchema(schema Schema) error {

	// Create session with increased timeout for schema creation
	ds.cluster.Timeout = SchemaTimeout
	session, err := createSession(ds.cluster)
	if err != nil {
		return err
	}
	defer session.Close()

	for _, object := range schema.GetObjects() {
		if err := createTable(session, object); err != nil {
			return err
		}
	}
	return nil
}

func (ds CQLDataStore) GetSession() interface{} {
	return ds.session
}

func (ds CQLDataStore) Close() {
	ds.session.Close()
}
