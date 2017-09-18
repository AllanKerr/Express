package oauth2

import "services/core"

type DatastoreAdapter struct {
	datastore core.Datastore
}

func NewDatastoreAdapter(ds core.Datastore) *DatastoreAdapter {
	adapter := new(DatastoreAdapter)
	adapter.datastore = ds
	return adapter
}