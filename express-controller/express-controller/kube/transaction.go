package kube

type Transaction interface {
	Execute(interface{}) error
	Rollback(name string) error
}
