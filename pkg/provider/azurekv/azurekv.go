package azurekv

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// APIVersion is the Azure Key Vault API version
var APIVersion = "7.0"

// AzureKV is an Azure Key Vault secrets provider
type AzureKV struct {
	httpClient     *http.Client
	clientID       string
	clientSecret   string
	tenantID       string
	subscriptionID string
}

type secretResponse struct {
	ID          string `json:"id"`
	ContentType string `json:"contentType"`
	Value       string `json:"value"`
}

type keyVaultError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
type keyVaultErrorResponse struct {
	Error keyVaultError `json:"error"`
}

// NewProviderFromEnv creates a new AzureKV provider using values from environment variables
func NewProviderFromEnv() AzureKV {
	return AzureKV{
		httpClient:     &http.Client{Timeout: 5 * time.Second},
		clientID:       os.Getenv("ARM_CLIENT_ID"),
		clientSecret:   os.Getenv("ARM_CLIENT_SECRET"),
		tenantID:       os.Getenv("ARM_TENANT_ID"),
		subscriptionID: os.Getenv("ARM_SUBSCRIPTION_ID"),
	}
}

// Name returns the provider name
func (a AzureKV) Name() string {
	return "AzureKV"
}

// Secret returns the requested secret
func (a AzureKV) Secret(id string) (string, error) {
	accessToken, err := a.accessToken()
	if err != nil {
		return "", fmt.Errorf("could not get access token: %s", err.Error())
	}

	response := secretResponse{}

	httpRequest, err := http.NewRequest("GET", id+"?api-version="+APIVersion, nil)
	if err != nil {
		return "", fmt.Errorf("could not create http request: %s", err.Error())
	}

	httpRequest.Header.Add("Authorization", "Bearer "+accessToken)

	httpResponse, err := a.httpClient.Do(httpRequest)
	if err != nil {
		return "", fmt.Errorf("error sending request: %s", err.Error())
	}

	if httpResponse.StatusCode != http.StatusOK {
		message, err := a.decodeKeyVaultErrorResponse(httpResponse.Body)
		if err != nil {
			return "", err
		}

		return "", errors.New(message)
	}

	err = json.NewDecoder(httpResponse.Body).Decode(&response)
	if err != nil {
		return "", fmt.Errorf("could not decode response: %s", err.Error())
	}

	return response.Value, nil
}

func (a AzureKV) decodeKeyVaultErrorResponse(body io.ReadCloser) (string, error) {
	errorResponse := keyVaultErrorResponse{}

	err := json.NewDecoder(body).Decode(&errorResponse)
	if err != nil {
		return "", fmt.Errorf("could not decode error response: %s", err.Error())
	}

	return "azurekv: " + errorResponse.Error.Message, nil
}
