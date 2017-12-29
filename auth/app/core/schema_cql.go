package core

import (
	"io/ioutil"
	"os"
	"strings"
)

type CqlSchema struct {
	objects []string
}

func (schema *CqlSchema) GetObjects() []string {
	return schema.objects
}

func getPath(dir string, file os.FileInfo) string {
	if !strings.HasSuffix(dir, "/") {
		dir += "/"
	}
	return dir + file.Name()
}

func isObject(file os.FileInfo) bool {
	return !file.IsDir() && strings.HasSuffix(file.Name(), ".cql")
}

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

func NewCqlSchema(dir string) (*CqlSchema, error) {

	objects, err := readObjects(dir)
	if err != nil {
		return nil, err
	}
	return &CqlSchema{
		objects: objects,
	}, nil
}
