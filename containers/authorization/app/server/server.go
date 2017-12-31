package server

import (
	"app/oauth2"
	"app/core"
	"github.com/sirupsen/logrus"
	"os"
)

func CreateSchema(ds core.DataStore) error {
	schema, err := core.NewCqlSchema("schemas")
	if err != nil {
		return err
	}
	return ds.CreateSchema(schema);
}

func RunHost(config *oauth2.Config) {

	databaseUrl := os.Getenv("DATABASE_URL")
	ds := core.NewCQLDataStoreRetry(databaseUrl, "authorization", 3, 5)

	if err := CreateSchema(ds); err != nil {
		logrus.WithField("error", err).Error("Failed to create schema.")
	}

	app := core.NewApp(ds, true, logrus.DebugLevel)
	oauth2.NewController(app, config)
	app.Start(8080)
}
