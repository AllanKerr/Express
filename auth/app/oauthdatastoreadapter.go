package main

type DatastoreAdapter struct {
	datastore Datastore
}

func NewDatastoreAdapter(ds Datastore) *DatastoreAdapter {
	adapter := new(DatastoreAdapter)
	adapter.datastore = ds
	return adapter
}
