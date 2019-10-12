package azurekv

// AzureKV is an Azure Key Vault secrets provider
type AzureKV struct {
}

// Secret returns the requested secret
func (azurekv AzureKV) Secret(name string, version string) (string, error) {
	return "", nil
}
