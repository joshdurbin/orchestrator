package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// MockServer creates a test HTTP server for testing
func createMockServer(t *testing.T, handler http.HandlerFunc) *httptest.Server {
	return httptest.NewServer(handler)
}

// Test client creation
func TestNewClient(t *testing.T) {
	config := Config{
		Host:     "localhost:3000",
		Username: "admin",
		Password: "secret",
		UseHTTPS: false,
	}

	client := NewClient(config)

	if client == nil {
		t.Fatal("Expected client to be created")
	}

	if client.baseURL != "http://localhost:3000" {
		t.Errorf("Expected baseURL to be 'http://localhost:3000', got '%s'", client.baseURL)
	}

	if client.username != "admin" {
		t.Errorf("Expected username to be 'admin', got '%s'", client.username)
	}
}

// Test client with HTTPS
func TestNewClientHTTPS(t *testing.T) {
	config := Config{
		Host:     "orchestrator.example.com:443",
		Username: "admin",
		Password: "secret",
		UseHTTPS: true,
	}

	client := NewClient(config)

	if client.baseURL != "https://orchestrator.example.com:443" {
		t.Errorf("Expected baseURL to be 'https://orchestrator.example.com:443', got '%s'", client.baseURL)
	}
}

// Test client with URL prefix
func TestNewClientWithURLPrefix(t *testing.T) {
	config := Config{
		Host:      "localhost:3000",
		URLPrefix: "/orchestrator",
	}

	client := NewClient(config)

	if client.urlPrefix != "/orchestrator" {
		t.Errorf("Expected urlPrefix to be '/orchestrator', got '%s'", client.urlPrefix)
	}
}

// Test InstanceKey parsing
func TestParseInstanceKey(t *testing.T) {
	tests := []struct {
		input       string
		expectedKey InstanceKey
		shouldError bool
	}{
		{"localhost:3306", InstanceKey{"localhost", 3306}, false},
		{"db1.example.com:3307", InstanceKey{"db1.example.com", 3307}, false},
		{"invalid", InstanceKey{}, true},
		{"host:abc", InstanceKey{}, true},
		{"host:3306:extra", InstanceKey{}, true},
	}

	for _, test := range tests {
		key, err := ParseInstanceKey(test.input)

		if test.shouldError {
			if err == nil {
				t.Errorf("Expected error for input '%s', but got none", test.input)
			}
		} else {
			if err != nil {
				t.Errorf("Unexpected error for input '%s': %v", test.input, err)
			}
			if key.Hostname != test.expectedKey.Hostname || key.Port != test.expectedKey.Port {
				t.Errorf("Expected key %v, got %v", test.expectedKey, key)
			}
		}
	}
}

// Test InstanceKey String method
func TestInstanceKeyString(t *testing.T) {
	key := InstanceKey{"localhost", 3306}
	expected := "localhost:3306"

	if key.String() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, key.String())
	}
}

// Test GetInstance
func TestGetInstance(t *testing.T) {
	ctx := context.Background()
	mockResponse := APIResponse{
		Code:    OK,
		Message: "Success",
		Details: map[string]interface{}{
			"Key": map[string]interface{}{
				"Hostname": "localhost",
				"Port":     3306,
			},
			"ServerID":    1,
			"Version":     "8.0.30",
			"ClusterName": "test-cluster",
			"ReadOnly":    false,
		},
	}

	server := createMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/instance/localhost/3306" {
			t.Errorf("Expected path '/api/instance/localhost/3306', got '%s'", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	})
	defer server.Close()

	config := Config{
		Host: server.URL[7:], // Remove "http://"
	}
	client := NewClient(config)

	instance, err := client.GetInstance(ctx, InstanceKey{"localhost", 3306})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if instance == nil {
		t.Fatal("Expected instance to be returned")
	}

	if instance.Key.Hostname != "localhost" || instance.Key.Port != 3306 {
		t.Errorf("Expected instance key localhost:3306, got %v", instance.Key)
	}
}

// Test GetClusters
func TestGetClusters(t *testing.T) {
	ctx := context.Background()
	mockResponse := APIResponse{
		Code:    OK,
		Message: "Success",
		Details: []string{"cluster1", "cluster2", "cluster3"},
	}

	server := createMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/clusters" {
			t.Errorf("Expected path '/api/clusters', got '%s'", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	})
	defer server.Close()

	config := Config{
		Host: server.URL[7:],
	}
	client := NewClient(config)

	clusters, err := client.GetClusters(ctx)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(clusters) != 3 {
		t.Errorf("Expected 3 clusters, got %d", len(clusters))
	}

	if clusters[0] != "cluster1" {
		t.Errorf("Expected first cluster to be 'cluster1', got '%s'", clusters[0])
	}
}

// Test Health check
func TestHealth(t *testing.T) {
	ctx := context.Background()
	mockResponse := APIResponse{
		Code:    OK,
		Message: "OK",
		Details: nil,
	}

	server := createMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/health" {
			t.Errorf("Expected path '/api/health', got '%s'", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	})
	defer server.Close()

	config := Config{
		Host: server.URL[7:],
	}
	client := NewClient(config)

	err := client.Health(ctx)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}

// Test authentication
func TestAuthentication(t *testing.T) {
	ctx := context.Background()
	server := createMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok {
			t.Error("Expected basic auth to be present")
		}

		if username != "testuser" || password != "testpass" {
			t.Errorf("Expected credentials testuser:testpass, got %s:%s", username, password)
		}

		mockResponse := APIResponse{
			Code:    OK,
			Message: "Success",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	})
	defer server.Close()

	config := Config{
		Host:     server.URL[7:],
		Username: "testuser",
		Password: "testpass",
	}
	client := NewClient(config)

	err := client.Health(ctx)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}

// Test error handling
func TestErrorHandling(t *testing.T) {
	ctx := context.Background()
	mockResponse := APIResponse{
		Code:    ERROR,
		Message: "Instance not found",
		Details: nil,
	}

	server := createMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(mockResponse)
	})
	defer server.Close()

	config := Config{
		Host: server.URL[7:],
	}
	client := NewClient(config)

	_, err := client.GetInstance(ctx, InstanceKey{"nonexistent", 3306})
	if err == nil {
		t.Fatal("Expected error, got none")
	}
}

// Test SetReadOnly
func TestSetReadOnly(t *testing.T) {
	ctx := context.Background()
	mockResponse := APIResponse{
		Code:    OK,
		Message: "Success",
		Details: map[string]interface{}{
			"Key": map[string]interface{}{
				"Hostname": "localhost",
				"Port":     3306,
			},
			"ReadOnly": true,
		},
	}

	server := createMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/set-read-only/localhost/3306" {
			t.Errorf("Expected path '/api/set-read-only/localhost/3306', got '%s'", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	})
	defer server.Close()

	config := Config{
		Host: server.URL[7:],
	}
	client := NewClient(config)

	instance, err := client.SetReadOnly(ctx, InstanceKey{"localhost", 3306})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !instance.ReadOnly {
		t.Error("Expected instance to be read-only")
	}
}

// Test BeginMaintenance
func TestBeginMaintenance(t *testing.T) {
	ctx := context.Background()
	mockResponse := APIResponse{
		Code:    OK,
		Message: "Success",
		Details: map[string]interface{}{
			"MaintenanceId": uint(1),
			"Key": map[string]interface{}{
				"Hostname": "localhost",
				"Port":     3306,
			},
			"Owner":  "admin",
			"Reason": "testing",
		},
	}

	server := createMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		expectedPath := "/api/begin-maintenance/localhost/3306/admin/testing"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path '%s', got '%s'", expectedPath, r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	})
	defer server.Close()

	config := Config{
		Host: server.URL[7:],
	}
	client := NewClient(config)

	maintenance, err := client.BeginMaintenance(ctx, InstanceKey{"localhost", 3306}, "admin", "testing")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if maintenance.Owner != "admin" {
		t.Errorf("Expected owner to be 'admin', got '%s'", maintenance.Owner)
	}

	if maintenance.Reason != "testing" {
		t.Errorf("Expected reason to be 'testing', got '%s'", maintenance.Reason)
	}
}

// Test BeginMaintenanceWithDuration
func TestBeginMaintenanceWithDuration(t *testing.T) {
	ctx := context.Background()
	mockResponse := APIResponse{
		Code:    OK,
		Message: "Success",
		Details: map[string]interface{}{
			"MaintenanceId": uint(1),
			"Key": map[string]interface{}{
				"Hostname": "localhost",
				"Port":     3306,
			},
			"Owner":  "admin",
			"Reason": "testing",
		},
	}

	server := createMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		expectedPath := "/api/begin-maintenance/localhost/3306/admin/testing/3600s"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path '%s', got '%s'", expectedPath, r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	})
	defer server.Close()

	config := Config{
		Host: server.URL[7:],
	}
	client := NewClient(config)

	duration := time.Hour
	_, err := client.BeginMaintenanceWithDuration(ctx, InstanceKey{"localhost", 3306}, "admin", "testing", duration)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}

// Test custom timeout
func TestCustomTimeout(t *testing.T) {
	config := Config{
		Host:    "localhost:3000",
		Timeout: 5 * time.Second,
	}

	client := NewClient(config)

	if client.httpClient.Timeout != 5*time.Second {
		t.Errorf("Expected timeout to be 5s, got %v", client.httpClient.Timeout)
	}
}

// Test URL prefix normalization
func TestURLPrefixNormalization(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"orchestrator", "/orchestrator"},
		{"/orchestrator", "/orchestrator"},
		{"orchestrator/", "/orchestrator"},
		{"/orchestrator/", "/orchestrator"},
		{"", ""},
	}

	for _, test := range tests {
		config := Config{
			Host:      "localhost:3000",
			URLPrefix: test.input,
		}

		client := NewClient(config)

		if client.urlPrefix != test.expected {
			t.Errorf("For input '%s', expected urlPrefix '%s', got '%s'",
				test.input, test.expected, client.urlPrefix)
		}
	}
}
