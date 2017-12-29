package core

type DataStore interface {
	CreateTable(schema string) error
	GetSession()  interface{}
	Close()
}
