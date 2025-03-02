package registry

type Registry interface {
	LoadRegistryFromFile(b []byte) (Registry, error)
	Get(s string) (any, error)
}