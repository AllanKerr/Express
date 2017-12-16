package core

import (
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"strconv"
	"net/http"
)

type App struct {
	router    *mux.Router
	datastore DataStore
}

func (app *App)GetDatastore() DataStore {
	return app.datastore
}

func NewApp(ds DataStore, readinessProbe bool) *App {
	if ds == nil {
		logrus.Fatal("Attempted to create new app with nil datastore.");
	}
	app := new(App)
	app.router = mux.NewRouter()
	app.datastore = ds

	if readinessProbe {
		app.AddEndpoint("/monitor/readiness", false, app.readinessProbe)
	}
	return app
}

func (app *App) readinessProbe(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func (app *App) AddEndpoint(path string, internal bool, f func(http.ResponseWriter, *http.Request)) *mux.Route {

	route := app.router.HandleFunc(path, f)
	if internal {
		route = route.Host("localhost")
	}
	return route
}

func (app *App) Start(port uint16) {
	defer app.datastore.Close()
	portStr := strconv.Itoa(int(port))
	logrus.WithField("port", port).Info("Starting routing.");
	logrus.Fatal(http.ListenAndServe(":" + portStr, app.router))
}
