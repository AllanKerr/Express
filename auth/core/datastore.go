package core

type Datastore interface {
	GetSession()  interface{}
	Close()
}
