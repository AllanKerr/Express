package main

import (
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"strconv"
	"net/http"
)

type App struct {
	router *mux.Router
	datastore Datastore
}

func NewApp(ds Datastore) *App {
	if ds == nil {
		logrus.Fatal("Attempted to create new app with nil datastore.");
	}
	app := new(App)
	app.router = mux.NewRouter()
	app.datastore = ds
	return app
}

func (app *App) Start(port uint16) {
	defer app.datastore.Close()
	portStr := strconv.Itoa(int(port))
	logrus.WithField("port", port).Info("Starting routing.");
	logrus.Fatal(http.ListenAndServe(":" + portStr, app.router))
}
