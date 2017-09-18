package main

import (
	"github.com/sirupsen/logrus"
	"services/oauth2"
	"services/core"
)

func main() {

	ds, err := core.NewCQLDatastore("cassandra-0.cassandra:9042", "default")
	if err != nil {
		logrus.WithError(err).Fatal("Failed to create CQL datastore.")
	}
	app := core.NewApp(ds)
	oauth2.NewHTTPController(app)
	app.Start(8080)
}
