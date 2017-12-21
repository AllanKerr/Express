package server

import (
	"app/oauth2"
	"app/core"
	"github.com/sirupsen/logrus"
)

func RunHost(config *oauth2.Config) {

	ds := core.NewCQLDataStoreRetry("cassandra-0.cassandra:9042", "default", 5)
	app := core.NewApp(ds, true, logrus.DebugLevel)
	oauth2.NewController(app, config)
	app.Start(8080)
}
