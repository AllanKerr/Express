package main

func main() {

	ds := NewCQLDatastoreRetry("cassandra-0.cassandra:9042", "default", 5)
	app := NewApp(ds, true)
	/*oauth2.NewHTTPController(app)*/
	app.Start(8080)
}
