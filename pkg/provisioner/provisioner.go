package provisioner

import (
	"encoding/base64"
	"fmt"

	"github.com/shawntoffel/docker-secrets-provisioner/pkg/docker"
	"github.com/shawntoffel/docker-secrets-provisioner/pkg/provider"
)

// Provisioner is a Docker secret provisioner
type Provisioner struct {
	provider provider.Provider
	docker   docker.Client
}

// New creates a new Provisioner
func New(provider provider.Provider, dockerClient docker.Client) Provisioner {
	return Provisioner{provider: provider, docker: dockerClient}
}

// Provision provisions a docker secret from the provider
func (p Provisioner) Provision(sourceName string, sourceVersion string, targetName string) (string, error) {
	secret, err := p.provider.Secret(sourceName, sourceVersion)
	if err != nil {
		return "", fmt.Errorf("couldn't retrieve secret from provider: %s", err.Error())
	}

	encoded := base64.StdEncoding.EncodeToString([]byte(secret))

	resp, err := p.docker.CreateSecret(docker.CreateSecretRequest{Name: targetName, Data: encoded})
	if err != nil {
		return "", fmt.Errorf("couldn't create docker secret: %s", err.Error())
	}

	return resp.ID, nil
}
