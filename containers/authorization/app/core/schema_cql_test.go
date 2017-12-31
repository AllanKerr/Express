package core

import (
	"testing"
	"io/ioutil"
	"os"
)


func createTestSchema(t *testing.T) (Schema, []string, error) {

	// create tmp directory
	name, err := ioutil.TempDir("", "schema_test")
	if err != nil {
		t.Errorf("Error creating temporary directory: %v", err)
	}

	var tables []string

	table1 := `
	CREATE TABLE IF NOT EXISTS "test_schema_table" (
    	key text PRIMARY KEY,
    	collection set<text>,
	);
	`
	tables = append(tables, table1)

	table2 := `
	CREATE TABLE IF NOT EXISTS "test_schema_table_2" (
    	key int PRIMARY KEY,
    	item text
	);
	`
	tables = append(tables, table2)

	// add schema files to the temporary directory
	ioutil.WriteFile(name + "/notpart.sql", []byte(""), os.ModeTemporary)
	ioutil.WriteFile(name + "/table1.cql", []byte(table1), os.ModeTemporary)
	ioutil.WriteFile(name + "/table2.cql", []byte(table2), os.ModeTemporary)

	schema, err := NewCqlSchema(name)
	return schema, tables, err
}

func TestNewCqlSchema(t *testing.T) {

	schema, tables, err := createTestSchema(t)
	if err != nil {
		t.Errorf("Error creating schema: %v", schema)
	}
	if len(schema.GetObjects()) != 2 {
		t.Errorf("Error, unexpected schema objects: %v", schema.GetObjects())
	}

	// test that all expected tables are found in the schema
	for _, obj := range schema.GetObjects() {
		found := false
		for _, table := range tables {
			found = found || table == obj
		}
		if !found {
			t.Errorf("Error, unexpected schema object: %v", obj)
		}
	}
}