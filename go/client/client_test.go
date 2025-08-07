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

package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/openark/orchestrator/go/inst"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name:    "nil config",
			config:  nil,
			wantErr: true,
		},
		{
			name: "valid config",
			config: &Config{
				BaseURL: "http://localhost:3000",
				Timeout: 10 * time.Second,
			},
			wantErr: false,
		},
		{
			name: "config with endpoints",
			config: &Config{
				Endpoints: []string{
					"http://localhost:3000",
					"http://localhost:3001",
				},
				Timeout: 10 * time.Second,
			},
			wantErr: true, // Will fail leader detection in test
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewClient(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"basic URL", "http://localhost:3000", "http://localhost:3000/api"},
		{"URL with trailing slash", "http://localhost:3000/", "http://localhost:3000/api"},
		{"URL already with api", "http://localhost:3000/api", "http://localhost:3000/api"},
		{"URL with api and slash", "http://localhost:3000/api/", "http://localhost:3000/api"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizeURL(tt.input)
			if result != tt.expected {
				t.Errorf("normalizeURL(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestClientHealth(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/health" {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	client := &Client{
		leader:     server.URL + "/api",
		httpClient: &http.Client{Timeout: 5 * time.Second},
	}

	ctx := context.Background()
	err := client.Health(ctx)
	if err != nil {
		t.Errorf("Health() failed: %v", err)
	}
}

func TestClientDiscover(t *testing.T) {
	// Create a test server that returns a mock instance
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/api/discover/") {
			response := APIResponse{
				Code:    "OK",
				Message: "Instance discovered",
				Details: map[string]interface{}{
					"Key": map[string]interface{}{
						"Hostname": "test-host",
						"Port":     3306,
					},
					"Version":    "8.0.25",
					"ReadOnly":   false,
					"ServerID":   1,
					"ServerUUID": "test-uuid",
				},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	client := &Client{
		leader:     server.URL + "/api",
		httpClient: &http.Client{Timeout: 5 * time.Second},
	}

	instanceKey := inst.InstanceKey{
		Hostname: "test-host",
		Port:     3306,
	}

	ctx := context.Background()
	instance, err := client.Discover(ctx, instanceKey)
	if err != nil {
		t.Errorf("Discover() failed: %v", err)
		return
	}

	if instance.Key.Hostname != "test-host" || instance.Key.Port != 3306 {
		t.Errorf("Expected instance key test-host:3306, got %s:%d",
			instance.Key.Hostname, instance.Key.Port)
	}
}

func TestClientGetInstance(t *testing.T) {
	// Create a test server that returns a mock instance
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/api/instance/") {
			instance := inst.Instance{
				Key: inst.InstanceKey{
					Hostname: "test-host",
					Port:     3306,
				},
				Version:    "8.0.25",
				ReadOnly:   false,
				ServerID:   1,
				ServerUUID: "test-uuid",
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(instance)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	client := &Client{
		leader:     server.URL + "/api",
		httpClient: &http.Client{Timeout: 5 * time.Second},
	}

	instanceKey := inst.InstanceKey{
		Hostname: "test-host",
		Port:     3306,
	}

	ctx := context.Background()
	instance, err := client.GetInstance(ctx, instanceKey)
	if err != nil {
		t.Errorf("GetInstance() failed: %v", err)
		return
	}

	if instance.Key.Hostname != "test-host" || instance.Key.Port != 3306 {
		t.Errorf("Expected instance key test-host:3306, got %s:%d",
			instance.Key.Hostname, instance.Key.Port)
	}

	if instance.Version != "8.0.25" {
		t.Errorf("Expected version 8.0.25, got %s", instance.Version)
	}
}

func TestClientGetClusters(t *testing.T) {
	// Create a test server that returns mock clusters
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/clusters" {
			clusters := []string{"cluster1", "cluster2", "cluster3"}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(clusters)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	client := &Client{
		leader:     server.URL + "/api",
		httpClient: &http.Client{Timeout: 5 * time.Second},
	}

	ctx := context.Background()
	clusters, err := client.GetClusters(ctx)
	if err != nil {
		t.Errorf("GetClusters() failed: %v", err)
		return
	}

	expected := []string{"cluster1", "cluster2", "cluster3"}
	if len(clusters) != len(expected) {
		t.Errorf("Expected %d clusters, got %d", len(expected), len(clusters))
		return
	}

	for i, cluster := range clusters {
		if cluster != expected[i] {
			t.Errorf("Expected cluster %s, got %s", expected[i], cluster)
		}
	}
}

func TestClientAPIError(t *testing.T) {
	// Create a test server that returns API errors
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := APIResponse{
			Code:    "ERROR",
			Message: "Instance not found",
			Details: nil,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK) // API returns 200 even for errors
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := &Client{
		leader:     server.URL + "/api",
		httpClient: &http.Client{Timeout: 5 * time.Second},
	}

	instanceKey := inst.InstanceKey{
		Hostname: "nonexistent",
		Port:     3306,
	}

	ctx := context.Background()
	_, err := client.Discover(ctx, instanceKey)
	if err == nil {
		t.Error("Expected error for nonexistent instance, got nil")
	}

	if !strings.Contains(err.Error(), "Instance not found") {
		t.Errorf("Expected 'Instance not found' in error, got: %v", err)
	}
}

func TestClientHTTPError(t *testing.T) {
	// Create a test server that returns HTTP errors
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	client := &Client{
		leader:     server.URL + "/api",
		httpClient: &http.Client{Timeout: 5 * time.Second},
	}

	ctx := context.Background()
	err := client.Health(ctx)
	if err == nil {
		t.Error("Expected error for HTTP 500, got nil")
	}

	if !strings.Contains(err.Error(), "500") {
		t.Errorf("Expected '500' in error, got: %v", err)
	}
}

func TestClientAuthentication(t *testing.T) {
	// Create a test server that checks authentication
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok || username != "testuser" || password != "testpass" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Test with correct credentials
	client := &Client{
		leader:     server.URL + "/api",
		httpClient: &http.Client{Timeout: 5 * time.Second},
		username:   "testuser",
		password:   "testpass",
	}

	ctx := context.Background()
	err := client.Health(ctx)
	if err != nil {
		t.Errorf("Health() with correct auth failed: %v", err)
	}

	// Test with wrong credentials
	clientWrong := &Client{
		leader:     server.URL + "/api",
		httpClient: &http.Client{Timeout: 5 * time.Second},
		username:   "wronguser",
		password:   "wrongpass",
	}

	err = clientWrong.Health(ctx)
	if err == nil {
		t.Error("Expected error with wrong credentials, got nil")
	}
}

func TestClientHeaders(t *testing.T) {
	// Create a test server that checks custom headers
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Custom-Header") != "test-value" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := &Client{
		leader:     server.URL + "/api",
		httpClient: &http.Client{Timeout: 5 * time.Second},
		headers: map[string]string{
			"X-Custom-Header": "test-value",
		},
	}

	ctx := context.Background()
	err := client.Health(ctx)
	if err != nil {
		t.Errorf("Health() with custom headers failed: %v", err)
	}
}

func TestClientTimeout(t *testing.T) {
	// Create a test server that delays response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Client with very short timeout
	client := &Client{
		leader:     server.URL + "/api",
		httpClient: &http.Client{Timeout: 10 * time.Millisecond},
	}

	ctx := context.Background()
	err := client.Health(ctx)
	if err == nil {
		t.Error("Expected timeout error, got nil")
	}
}

func BenchmarkClientHealth(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := &Client{
		leader:     server.URL + "/api",
		httpClient: &http.Client{Timeout: 5 * time.Second},
	}

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = client.Health(ctx)
	}
}
