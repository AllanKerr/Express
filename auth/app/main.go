package main

import (
	"github.com/sirupsen/logrus"
)

func main() {

	ds, err := NewCQLDatastore("cassandra-0.cassandra:9042", "default")
	if err != nil {
		logrus.WithError(err).Fatal("Failed to create CQL datastore.")
	}
	app := NewApp(ds)
	app.Start(8080)
}
