package provider

// Provider is a secrets provider
type Provider interface {
	Secret(name string, version string) (string, error)
}
