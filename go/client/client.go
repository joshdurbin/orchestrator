/*
   Copyright 2025 Percona Inc.

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

// Package client provides a Go client for the Orchestrator REST API.
// It offers a programmatic interface to all Orchestrator operations including
// topology discovery, replication management, failover operations, and cluster management.
package client

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Client represents an Orchestrator API client
type Client struct {
	baseURL    string
	httpClient *http.Client
	username   string
	password   string
	headers    map[string]string
	leader     string
	endpoints  []string
}

// Config holds configuration for the client
type Config struct {
	// BaseURL is the primary Orchestrator API endpoint
	BaseURL string

	// Endpoints is a list of Orchestrator API endpoints for leader detection
	// If provided, the client will automatically detect the leader
	Endpoints []string

	// Username for HTTP basic authentication
	Username string

	// Password for HTTP basic authentication
	Password string

	// Headers contains additional HTTP headers to include in requests
	Headers map[string]string

	// Timeout for HTTP requests
	Timeout time.Duration

	// InsecureSkipVerify controls whether the client verifies the server's certificate chain
	InsecureSkipVerify bool
}

// APIResponse represents the standard API response format
type APIResponse struct {
	Code    string      `json:"Code"`
	Message string      `json:"Message"`
	Details interface{} `json:"Details"`
}

// NewClient creates a new Orchestrator API client
func NewClient(config *Config) (*Client, error) {
	if config == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	timeout := config.Timeout
	if timeout == 0 {
		timeout = 30 * time.Second
	}

	httpClient := &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: config.InsecureSkipVerify,
			},
		},
	}

	client := &Client{
		baseURL:    normalizeURL(config.BaseURL),
		httpClient: httpClient,
		username:   config.Username,
		password:   config.Password,
		headers:    config.Headers,
		endpoints:  normalizeEndpoints(config.Endpoints),
	}

	// If multiple endpoints provided, detect leader
	if len(client.endpoints) > 1 {
		if err := client.detectLeader(); err != nil {
			return nil, fmt.Errorf("failed to detect leader: %w", err)
		}
	} else if len(client.endpoints) == 1 {
		client.leader = client.endpoints[0]
	} else {
		client.leader = client.baseURL
	}

	return client, nil
}

// normalizeURL ensures the URL ends with /api
func normalizeURL(baseURL string) string {
	baseURL = strings.TrimSuffix(baseURL, "/")
	if !strings.HasSuffix(baseURL, "/api") {
		baseURL += "/api"
	}
	return baseURL
}

// normalizeEndpoints normalizes a list of endpoints
func normalizeEndpoints(endpoints []string) []string {
	var normalized []string
	for _, endpoint := range endpoints {
		normalized = append(normalized, normalizeURL(endpoint))
	}
	return normalized
}

// detectLeader finds the current leader among the configured endpoints
func (c *Client) detectLeader() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, endpoint := range c.endpoints {
		if c.isLeader(ctx, endpoint) {
			c.leader = endpoint
			return nil
		}
	}

	// If no direct leader found, try routed leader check
	for _, endpoint := range c.endpoints {
		if c.hasRoutedLeader(ctx, endpoint) {
			c.leader = endpoint
			return nil
		}
	}

	return fmt.Errorf("no leader found among endpoints")
}

// isLeader checks if the given endpoint is the leader
func (c *Client) isLeader(ctx context.Context, endpoint string) bool {
	req, err := http.NewRequestWithContext(ctx, "GET", endpoint+"/leader-check", nil)
	if err != nil {
		return false
	}
	c.setAuth(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

// hasRoutedLeader checks if the endpoint can route to the leader
func (c *Client) hasRoutedLeader(ctx context.Context, endpoint string) bool {
	req, err := http.NewRequestWithContext(ctx, "GET", endpoint+"/routed-leader-check", nil)
	if err != nil {
		return false
	}
	c.setAuth(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

// setAuth sets authentication headers on the request
func (c *Client) setAuth(req *http.Request) {
	if c.username != "" && c.password != "" {
		req.SetBasicAuth(c.username, c.password)
	}

	for key, value := range c.headers {
		req.Header.Set(key, value)
	}
}

// doRequest performs an HTTP request and returns the response
func (c *Client) doRequest(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
	var bodyReader io.Reader
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	url := c.leader + "/" + strings.TrimPrefix(path, "/")
	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	c.setAuth(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	return resp, nil
}

// parseResponse parses an API response into the provided structure
func (c *Client) parseResponse(resp *http.Response, result interface{}) error {
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		var apiResp APIResponse
		if err := json.Unmarshal(body, &apiResp); err == nil {
			return fmt.Errorf("API error: %s - %s", apiResp.Code, apiResp.Message)
		}
		return fmt.Errorf("HTTP error %d: %s", resp.StatusCode, string(body))
	}

	if result != nil {
		if err := json.Unmarshal(body, result); err != nil {
			return fmt.Errorf("failed to unmarshal response: %w", err)
		}
	}

	return nil
}

// parseAPIResponse parses a standard API response
func (c *Client) parseAPIResponse(resp *http.Response, result interface{}) (*APIResponse, error) {
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var apiResp APIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal API response: %w", err)
	}

	if apiResp.Code == "ERROR" {
		return &apiResp, fmt.Errorf("API error: %s", apiResp.Message)
	}

	if result != nil && apiResp.Details != nil {
		detailsBytes, err := json.Marshal(apiResp.Details)
		if err != nil {
			return &apiResp, fmt.Errorf("failed to marshal details: %w", err)
		}
		if err := json.Unmarshal(detailsBytes, result); err != nil {
			return &apiResp, fmt.Errorf("failed to unmarshal details: %w", err)
		}
	}

	return &apiResp, nil
}

// Health checks the health of the Orchestrator service
func (c *Client) Health(ctx context.Context) error {
	resp, err := c.doRequest(ctx, "GET", "/health", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("health check failed with status %d", resp.StatusCode)
	}

	return nil
}

// LeaderCheck verifies if the current endpoint is the leader
func (c *Client) LeaderCheck(ctx context.Context) (bool, error) {
	resp, err := c.doRequest(ctx, "GET", "/leader-check", nil)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK, nil
}

// GetLeader returns the current leader endpoint
func (c *Client) GetLeader() string {
	return c.leader
}

// RefreshLeader re-detects the leader among configured endpoints
func (c *Client) RefreshLeader() error {
	if len(c.endpoints) > 1 {
		return c.detectLeader()
	}
	return nil
}
