package main

type Datastore interface {
	GetSession()  interface{}
	Close()
}
