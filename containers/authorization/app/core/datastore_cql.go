package core

import (
	"github.com/gocql/gocql"
	"github.com/sirupsen/logrus"
	"time"
	"fmt"
)

// increased timeout for creating tables and keyspaces
const SchemaTimeout = 5*time.Second

// data store for accessing the Cassandra session and cluster
type CQLDataStore struct {
	session *gocql.Session
	cluster *gocql.ClusterConfig
}

// create a new cluster on the specified host
func createCluster(host string) *gocql.ClusterConfig {

	logrus.WithField("host", host).Info("Creating new CQL session.")

	cluster := gocql.NewCluster(host)
	cluster.NumConns = 20
	cluster.Consistency = gocql.Quorum
	return cluster
}

// create a new table with the specified session and CQL table query
// table may also be a query with no output like delete or update
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

// create a new keyspace on the specified cluster
func createKeyspace(cluster *gocql.ClusterConfig, keyspace string, replicationFactor int) error {

	c := *cluster
	c.Keyspace = "system"
	// increase timeout for keyspace creation to avoid timeout
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

	// create the keyspace
	return createTable(session, fmt.Sprintf(`CREATE KEYSPACE IF NOT EXISTS "%s" WITH replication = {
		'class' : 'SimpleStrategy',
		'replication_factor' : %d
	}`, keyspace, replicationFactor))
}

// create a session for interacting with the cluster
func createSession(cluster *gocql.ClusterConfig) (*gocql.Session, error) {
	return cluster.CreateSession()
}

// create session on the specified keyspace
// the keyspace will be created if it does not already exist
func createSessionAndKeyspace(cluster *gocql.ClusterConfig, keyspace string, replicationFactor int) (*gocql.Session, error) {

	if err := createKeyspace(cluster, keyspace, replicationFactor); err != nil {
		return nil, err
	}
	cluster.Keyspace = keyspace
	return createSession(cluster)
}

// create a new CQL data store for interacting with the database at the host path on the specified keyspace
// If creation fails then it will be retried every interval until it succeeds.
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

// Create a new CQL data store for interacting with the database at the host path on the specified keyspace.
// If creation fails, the error is returned.
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

// Create a new table in the current Cassandra session.
// Object may be a table, keyspace, update, or delete query.
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

// Create a new CQL schema by creating all objects in the
// schema for the current Cassandra session.
// If any of the creations fail, the error is immediately returned.
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

func (ds *CQLDataStore) GetSession() interface{} {
	return ds.session
}

func (ds *CQLDataStore) Close() {
	ds.session.Close()
}
