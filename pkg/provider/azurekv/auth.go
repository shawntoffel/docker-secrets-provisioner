package azurekv

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// AuthBaseURL is the Azure OAuth2 base url
var AuthBaseURL = "https://login.microsoftonline.com/"

// TokenEndpoint is the Azure OAuth2 token endpoint
var TokenEndpoint = "/oauth2/token"

type tokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   string `json:"expires_in"`
	ExpiresOn   string `json:"expires_on"`
	Resource    string `json:"resource"`
}

type tokenErrorResponse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

func (a AzureKV) accessToken() (string, error) {
	response := tokenResponse{}

	httpResponse, err := a.httpClient.PostForm(AuthBaseURL+a.tenantID+TokenEndpoint, url.Values{
		"grant_type":    {"client_credentials"},
		"client_id":     {a.clientID},
		"client_secret": {a.clientSecret},
		"resource":      {"https://vault.azure.net"},
	})
	if err != nil {
		return "", fmt.Errorf("error sending request: %s", err.Error())
	}

	if httpResponse.StatusCode != http.StatusOK {
		message, err := a.decodeTokenErrorResponse(httpResponse.Body)
		if err != nil {
			return "", err
		}

		return "", errors.New(message)
	}

	err = json.NewDecoder(httpResponse.Body).Decode(&response)
	if err != nil {
		return "", fmt.Errorf("could not decode response: %s", err.Error())
	}

	return response.AccessToken, nil
}

func (a AzureKV) decodeTokenErrorResponse(body io.ReadCloser) (string, error) {
	errorResponse := tokenErrorResponse{}

	err := json.NewDecoder(body).Decode(&errorResponse)
	if err != nil {
		return "", fmt.Errorf("could not decode error response: %s", err.Error())
	}

	return "Azure token error: " + errorResponse.Error + ": " + errorResponse.ErrorDescription, nil
}
