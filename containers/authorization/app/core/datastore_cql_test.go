package core

import (
	"testing"
	"os"
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx"
	"github.com/scylladb/gocqlx/qb"
	//"github.com/ory/fosite"
)

var datastore DataStore

func TestMain(m *testing.M) {

	databaseUrl := os.Getenv("DATABASE_URL")
	datastore = NewCQLDataStoreRetry(databaseUrl, "datastore_cql_test", 1, 5)
	retCode := m.Run()
	datastore.Close()
	os.Exit(retCode)
}

type test_item struct {
	TestPrimary string
	TestSecondary string
}

// test that create table creates a new table in the session's keyspace
func TestCQLDataStore_CreateTable(t *testing.T) {

	// drop it if it exists
	deleteTable := `DROP TABLE IF EXISTS "test_table";`
	if err := datastore.CreateTable(deleteTable); err != nil {
		t.Errorf("Error deleting table: %v", err)
	}

	table := `
	CREATE TABLE "test_table" (
    	test_primary text PRIMARY KEY,
    	test_secondary text,
    	test_collection set<text>,
	);
	`
	if err := datastore.CreateTable(table); err != nil {
		t.Errorf("Error creating table: %v", err)
	}

	session := datastore.GetSession().(*gocql.Session)

	// query for an item and ensure it is 'not found'
	stmt, names := qb.Select("test_table").Where(qb.Eq("test_primary")).ToCql()
	q := gocqlx.Query(session.Query(stmt), names).BindMap(qb.M{
		"test_primary": "not_found",
	})
	defer q.Release()

	var item interface{}
	if err := gocqlx.Get(&item, q.Query); err == nil || err != gocql.ErrNotFound {
		if err != nil {
			t.Errorf("Error creating table: %v", err)
		} else {
			t.Error("Error creating table")
		}
	}
}