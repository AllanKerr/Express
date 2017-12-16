package main

import (
	"app/core"
	"app/oauth2"
)
func main() {
	ds := core.NewCQLDataStoreRetry("cassandra-0.cassandra:9042", "default", 5)
	app := core.NewApp(ds, true)
	oauth2.NewController(app, "temp secret")
	app.Start(8080)
}
