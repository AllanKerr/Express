package server

import (
	"app/oauth2"
	"app/core"
	"github.com/sirupsen/logrus"
	"os"
)

type Server struct {
	app *core.App
	authController *oauth2.HTTPController
}

func CreateSchema(ds core.DataStore) error {
	schema, err := core.NewCqlSchema("schemas")
	if err != nil {
		return err
	}
	return ds.CreateSchema(schema);
}

func Initialize(config *oauth2.Config) *Server {

	databaseUrl := os.Getenv("DATABASE_URL")
	ds := core.NewCQLDataStoreRetry(databaseUrl, "authorization", 3, 5)

	if err := CreateSchema(ds); err != nil {
		logrus.WithField("error", err).Error("Failed to create schema.")
	}

	app := core.NewApp(ds, true, logrus.DebugLevel)

	return &Server{
		app: app,
		authController: oauth2.NewController(app, config),
	}
}

func (server *Server) Run() {
	server.app.Start(8080)
}