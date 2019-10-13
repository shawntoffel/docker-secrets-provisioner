package docker

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
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

// NewClientFromEnv creates a new Docker client with values from environment variables
func NewClientFromEnv() (Client, error) {
	httpClient := &http.Client{}

	if os.Getenv("DOCKER_TLS_VERIFY") != "" {
		tlsConfig, err := loadTLSConfig()
		if err != nil {
			return Client{}, err
		}

		httpClient.Transport = &http.Transport{
			TLSClientConfig: tlsConfig,
		}
	}

	client := Client{
		httpClient: httpClient,
		host:       os.Getenv("DOCKER_HOST"),
		version:    os.Getenv("DOCKER_API_VERSION"),
	}

	return client, nil
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

func loadTLSConfig() (*tls.Config, error) {
	path := os.Getenv("DOCKER_CERT_PATH")

	cert, err := tls.LoadX509KeyPair(path+"/cert.pem", path+"/key.pem")
	if err != nil {
		return nil, fmt.Errorf("docker client: could not load key pair")
	}

	caCert, err := ioutil.ReadFile(path + "/ca.pem")
	if err != nil {
		return nil, fmt.Errorf("docker client: could not read CA file")
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
	}

	tlsConfig.BuildNameToCertificate()

	return tlsConfig, nil
}
