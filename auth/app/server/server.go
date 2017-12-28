package server

import (
	"app/oauth2"
	"app/core"
	"github.com/sirupsen/logrus"
	"os"
)

func RunHost(config *oauth2.Config) {

	databaseUrl := os.Getenv("DATABASE_URL")
	ds := core.NewCQLDataStoreRetry(databaseUrl, "default", 5)
	app := core.NewApp(ds, true, logrus.DebugLevel)
	oauth2.NewController(app, config)
	app.Start(8080)
}
