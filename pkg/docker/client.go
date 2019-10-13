package docker

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

type errorResponse struct {
	Message string
}

// Client is a Docker Engine client
type Client struct {
	host       string
	version    string
	httpClient *http.Client
}

// NewClient creates a new Docker client
func NewClient(host string, version string) Client {
	return NewClientWithHTTPClient(&http.Client{}, host, version)
}

// NewClientFromEnv creates a new Docker client with values from environment variables
func NewClientFromEnv() Client {
	return NewClientWithHTTPClient(
		&http.Client{},
		os.Getenv("DOCKER_HOST"),
		os.Getenv("DOCKER_API_VERSION"))
}

// NewClientWithHTTPClient creates a new Docker client with the provided http client
func NewClientWithHTTPClient(client *http.Client, host string, version string) Client {
	return Client{
		httpClient: client,
		host:       host,
		version:    version,
	}
}

func (c Client) create(endpoint string, data interface{}, response interface{}) error {
	body, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("could not marshal request to json: %s", err.Error())
	}

	resp, err := c.httpClient.Post(c.buildRequestURL(endpoint), "application/json", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("error sending request: %s", err.Error())
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		message, err := c.decodeErrorResponse(resp.Body)
		if err != nil {
			return err
		}

		return errors.New(message)
	}

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return fmt.Errorf("could not decode response: %s", err.Error())
	}

	return nil
}

func (c Client) buildRequestURL(endpoint string) string {
	if c.version != "" {
		return c.host + "/" + c.version + endpoint
	}

	return c.host + endpoint
}

func (c Client) decodeErrorResponse(body io.ReadCloser) (string, error) {
	errorResponse := errorResponse{}

	err := json.NewDecoder(body).Decode(&errorResponse)
	if err != nil {
		return "", fmt.Errorf("could not decode error response: %s", err.Error())
	}

	return "docker error: " + errorResponse.Message, nil
}
