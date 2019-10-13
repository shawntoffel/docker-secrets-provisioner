package provider

// Provider is a secrets provider
type Provider interface {
	Name() string
	Secret(id string) (string, error)
}
