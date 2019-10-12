package provider

// Provider is a secrets provider
type Provider interface {
	Name() string
	Secret(name string, version string) (string, error)
}
