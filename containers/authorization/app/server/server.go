package server

import (
	"app/oauth2"
	"app/core"
	"github.com/sirupsen/logrus"
	"os"
)

// Server used to initialize and run the authorization service.
type Server struct {
	app *core.App
	authController *oauth2.HTTPController
}

// Create the database schema from the schemas directory
func CreateSchema(ds core.DataStore) error {
	schema, err := core.NewCqlSchema("schemas")
	if err != nil {
		return err
	}
	return ds.CreateSchema(schema);
}

// Initialize the authorization server by creating the database
// session, initializing the database's schema, and adding the
// HTTP endpoints.
func Initialize(config *oauth2.Config) *Server {

	databaseUrl := os.Getenv("DATABASE_URL")
	ds := core.NewCQLDataStoreRetry(databaseUrl, "authorization", 1, 5)

	if err := CreateSchema(ds); err != nil {
		logrus.WithField("error", err).Error("Failed to create schema.")
	}

	app := core.NewApp(ds, true, logrus.DebugLevel)

	return &Server{
		app: app,
		authController: oauth2.NewController(app, config),
	}
}

// Run the authorization server by starting HTTP listening.
// This function should not return.
func (server *Server) Run() {
	server.app.Start(8080)
}