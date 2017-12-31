package core

import (
	"io/ioutil"
	"os"
	"strings"
)

// A CQL schema for creating a set of Cassandra objects
type CqlSchema struct {
	objects []string
}

// The array of objects in the schema
func (schema *CqlSchema) GetObjects() []string {
	return schema.objects
}

// The path for a file from its directory and file info
func getPath(dir string, file os.FileInfo) string {
	if !strings.HasSuffix(dir, "/") {
		dir += "/"
	}
	return dir + file.Name()
}

// Test whether a file is a schema file meaning it has a CQL extension
// and isn't a directory
func isObject(file os.FileInfo) bool {
	return !file.IsDir() && strings.HasSuffix(file.Name(), ".cql")
}

// Read the schema objects in the provided directory
func readObjects(dir string) ([]string, error) {

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var objects []string
	for _, file := range files {
		if isObject(file) {
			path := getPath(dir, file)
			object, err := ioutil.ReadFile(path)
			if err != nil {
				return nil, err
			}
			objects = append(objects, string(object))
		}
	}
	return objects, nil
}

// Create a new schema from the CQL files in the specified directory
func NewCqlSchema(dir string) (*CqlSchema, error) {

	objects, err := readObjects(dir)
	if err != nil {
		return nil, err
	}
	return &CqlSchema{
		objects: objects,
	}, nil
}
