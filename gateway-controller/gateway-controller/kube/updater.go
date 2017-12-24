package kube

type ObjectUpdater interface {
	GetModifiers() []string
	Update(name string, spec interface{}) error
}

