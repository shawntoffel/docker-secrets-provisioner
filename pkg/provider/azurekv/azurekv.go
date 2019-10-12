package azurekv

// AzureKV is an Azure Key Vault secrets provider
type AzureKV struct {
}

// NewProvider creates a new AzureKV provider
func NewProvider() AzureKV {
	return AzureKV{}
}

// Name returns the provider name
func (azurekv AzureKV) Name() string {
	return "AzureKV"
}

// Secret returns the requested secret
func (azurekv AzureKV) Secret(name string, version string) (string, error) {
	return "", nil
}
