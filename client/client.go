// Package client provides a Go SDK for interacting with the orchestrator REST API.
package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Client is the orchestrator API client
type Client struct {
	baseURL    string
	httpClient *http.Client
	username   string
	password   string
	urlPrefix  string
}

// Config holds the configuration for creating a new client
type Config struct {
	// Host is the orchestrator host (e.g., "localhost:3000")
	Host string
	// Username for authentication
	Username string
	// Password for authentication
	Password string
	// UseHTTPS enables HTTPS instead of HTTP
	UseHTTPS bool
	// URLPrefix is the URL prefix (e.g., "/orchestrator")
	URLPrefix string
	// HTTPClient allows providing a custom HTTP client
	HTTPClient *http.Client
	// Timeout for HTTP requests (default: 30 seconds)
	Timeout time.Duration
}

// NewClient creates a new orchestrator API client
func NewClient(config Config) *Client {
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}

	httpClient := config.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: config.Timeout,
		}
	}

	scheme := "http"
	if config.UseHTTPS {
		scheme = "https"
	}

	baseURL := fmt.Sprintf("%s://%s", scheme, config.Host)

	// Normalize URL prefix
	urlPrefix := strings.TrimSuffix(config.URLPrefix, "/")
	if urlPrefix != "" && !strings.HasPrefix(urlPrefix, "/") {
		urlPrefix = "/" + urlPrefix
	}

	return &Client{
		baseURL:    baseURL,
		httpClient: httpClient,
		username:   config.Username,
		password:   config.Password,
		urlPrefix:  urlPrefix,
	}
}

// doRequest executes an HTTP request and handles the response
func (c *Client) doRequest(ctx context.Context, method, path string, body io.Reader) (*http.Response, error) {
	url := c.baseURL + c.urlPrefix + path

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if c.username != "" || c.password != "" {
		req.SetBasicAuth(c.username, c.password)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	return resp, nil
}

// get executes a GET request
func (c *Client) get(ctx context.Context, path string) (*http.Response, error) {
	return c.doRequest(ctx, "GET", path, nil)
}

// post executes a POST request
func (c *Client) post(ctx context.Context, path string, body interface{}) (*http.Response, error) {
	var bodyReader io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonData)
	}
	return c.doRequest(ctx, "POST", path, bodyReader)
}

// parseResponse parses an API response into the target structure
func (c *Client) parseResponse(resp *http.Response, target interface{}) error {
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	if target != nil {
		if err := json.Unmarshal(body, target); err != nil {
			return fmt.Errorf("failed to parse response: %w (body: %s)", err, string(body))
		}
	}

	return nil
}

// getJSON performs a GET request and parses the JSON response
func (c *Client) getJSON(ctx context.Context, path string, target interface{}) error {
	resp, err := c.get(ctx, path)
	if err != nil {
		return err
	}
	return c.parseResponse(resp, target)
}

// postJSON performs a POST request and parses the JSON response
func (c *Client) postJSON(ctx context.Context, path string, body, target interface{}) error {
	resp, err := c.post(ctx, path, body)
	if err != nil {
		return err
	}
	return c.parseResponse(resp, target)
}

// getPlainText performs a GET request and returns plain text response
func (c *Client) getPlainText(ctx context.Context, path string) (string, error) {
	resp, err := c.get(ctx, path)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	return string(body), nil
}

// buildPath builds an API path with parameters
func buildPath(pathTemplate string, params ...interface{}) string {
	parts := make([]string, len(params))
	for i, p := range params {
		parts[i] = url.PathEscape(fmt.Sprint(p))
	}

	// Replace placeholders with actual values
	result := pathTemplate
	for _, part := range parts {
		result = strings.Replace(result, "{}", part, 1)
	}

	return result
}

// InstanceKey represents a unique identifier for a database instance
type InstanceKey struct {
	Hostname string
	Port     int
}

// String returns the string representation of an instance key
func (key InstanceKey) String() string {
	return fmt.Sprintf("%s:%d", key.Hostname, key.Port)
}

// ParseInstanceKey parses a host:port string into an InstanceKey
func ParseInstanceKey(hostPort string) (InstanceKey, error) {
	parts := strings.Split(hostPort, ":")
	if len(parts) != 2 {
		return InstanceKey{}, fmt.Errorf("invalid instance key format: %s", hostPort)
	}

	var port int
	_, err := fmt.Sscanf(parts[1], "%d", &port)
	if err != nil {
		return InstanceKey{}, fmt.Errorf("invalid port in instance key: %s", parts[1])
	}

	return InstanceKey{
		Hostname: parts[0],
		Port:     port,
	}, nil
}
