package main

import (
	"app/core"
	"app/oauth2"
	"os"
	"log"
	"github.com/sirupsen/logrus"
)

func main() {

	// Load config environment variables
	secret := os.Getenv("SYSTEM_SECRET")
	if secret == "" {
		log.Fatal("system secret must be 32 characters long, received: " + secret)
	}

	ds := core.NewCQLDataStoreRetry("cassandra-0.cassandra:9042", "default", 5)
	app := core.NewApp(ds, true, logrus.DebugLevel)

	ctrl := oauth2.NewController(app, secret)

	clientId := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	err := ctrl.CreateRootClient(clientId, clientSecret)
	if err != nil {
		logrus.WithField("error", err).Error("failed to create root client")
	}
	app.Start(8080)
}
