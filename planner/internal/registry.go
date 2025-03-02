package internal

type Registry interface {
	NewRegistry() *Registry
	LoadRegistryFromFile(b []byte) (*Registry, error)
}