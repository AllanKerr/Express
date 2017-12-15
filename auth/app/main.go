package main

func main() {

	ds := NewCQLDatastoreRetry("cassandra-0.cassandra:9042", "default", 5)
	app := NewApp(ds, true)
	NewOAuth2Controller(app, "temp secret")
	app.Start(8080)
}
