package core

type DataStore interface {
	CreateTable(object string) error
	CreateSchema(schema Schema) error
	GetSession()  interface{}
	Close()
}
