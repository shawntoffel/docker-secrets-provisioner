package docker

// CreateSecretRequest contains data used to create a docker secret
type CreateSecretRequest struct {
	Name string
	Data string
}

// CreateSecretResponse contains data returned after secret creation
type CreateSecretResponse struct {
	ID string
}

// CreateSecret creates a Docker secret
func (c Client) CreateSecret(request CreateSecretRequest) (*CreateSecretResponse, error) {
	createSecretResponse := &CreateSecretResponse{}

	err := c.create("/secrets/create", request, createSecretResponse)
	if err != nil {
		return nil, err
	}

	return createSecretResponse, nil
}
