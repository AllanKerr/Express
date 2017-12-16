package core

type DataStore interface {
	GetSession()  interface{}
	Close()
}
