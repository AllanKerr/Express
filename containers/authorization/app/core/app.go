package core

import (
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"strconv"
	"net/http"
)

// app for managing the datastore and http endpoints
type App struct {
	router    *mux.Router
	datastore DataStore
}

func (app *App)GetDatastore() DataStore {
	return app.datastore
}

// create a new app with the specified data store and an http readiness probe
// endpoint at /monitor/readiness if set to true with the specified log level
func NewApp(ds DataStore, readinessProbe bool, level logrus.Level) *App {

	if ds == nil {
		logrus.Fatal("Attempted to create new app with nil data store.");
	}
	logrus.SetLevel(level)

	app := new(App)
	app.router = mux.NewRouter()
	app.datastore = ds

	// create readiness probe http endpoint
	if readinessProbe {
		app.AddEndpoint("/monitor/readiness", false, app.readinessProbe)
	}
	return app
}

// readiness probe handler
func (app *App) readinessProbe(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

// Add a new http endpoint to the app with the specified path.
// If the endpoint is internal then it only accepts requests from localhost
func (app *App) AddEndpoint(path string, internal bool, f func(http.ResponseWriter, *http.Request)) *mux.Route {

	route := app.router.HandleFunc(path, f)
	if internal {
		route = route.Host("localhost")
	}
	return route
}

// Start serving the added http endpoints and listenting for requests on the specified port
func (app *App) Start(port uint16) {
	defer app.datastore.Close()
	portStr := strconv.Itoa(int(port))
	logrus.WithField("port", port).Info("Starting routing.");
	logrus.Fatal(http.ListenAndServe(":" + portStr, app.router))
}
