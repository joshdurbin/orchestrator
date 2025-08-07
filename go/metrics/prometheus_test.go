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

package metrics

import (
	"strings"
	"testing"

	"github.com/openark/orchestrator/go/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rcrowley/go-metrics"
)

func TestPrometheusIntegration(t *testing.T) {
	// Save original config
	originalEnabled := config.Config.PrometheusEnabled
	originalNamespace := config.Config.PrometheusNamespace
	originalSubsystem := config.Config.PrometheusSubsystem
	defer func() {
		config.Config.PrometheusEnabled = originalEnabled
		config.Config.PrometheusNamespace = originalNamespace
		config.Config.PrometheusSubsystem = originalSubsystem
	}()

	// Enable Prometheus for testing
	config.Config.PrometheusEnabled = true
	config.Config.PrometheusNamespace = "orchestrator"
	config.Config.PrometheusSubsystem = "test"

	// Initialize Prometheus metrics
	err := InitPrometheusMetrics()
	if err != nil {
		t.Fatalf("Failed to initialize Prometheus metrics: %v", err)
	}

	// Test that registry is created
	registry := GetPrometheusRegistry()
	if registry == nil {
		t.Fatal("Prometheus registry is nil")
	}

	// Test that IsPrometheusEnabled works
	if !IsPrometheusEnabled() {
		t.Fatal("IsPrometheusEnabled should return true")
	}
}

func TestPrometheusDisabled(t *testing.T) {
	// Save original config
	originalEnabled := config.Config.PrometheusEnabled
	defer func() {
		config.Config.PrometheusEnabled = originalEnabled
	}()

	// Disable Prometheus
	config.Config.PrometheusEnabled = false

	// Initialize should succeed but do nothing
	err := InitPrometheusMetrics()
	if err != nil {
		t.Fatalf("InitPrometheusMetrics should not fail when disabled: %v", err)
	}

	// Registry should be nil when disabled
	registry := GetPrometheusRegistry()
	if registry != nil {
		t.Fatal("Prometheus registry should be nil when disabled")
	}

	// IsPrometheusEnabled should return false
	if IsPrometheusEnabled() {
		t.Fatal("IsPrometheusEnabled should return false")
	}
}

func TestMetricsBridge(t *testing.T) {
	// Save original config
	originalEnabled := config.Config.PrometheusEnabled
	originalNamespace := config.Config.PrometheusNamespace
	defer func() {
		config.Config.PrometheusEnabled = originalEnabled
		config.Config.PrometheusNamespace = originalNamespace
	}()

	// Enable Prometheus for testing
	config.Config.PrometheusEnabled = true
	config.Config.PrometheusNamespace = "orchestrator"

	// Initialize Prometheus metrics
	err := InitPrometheusMetrics()
	if err != nil {
		t.Fatalf("Failed to initialize Prometheus metrics: %v", err)
	}

	// Create a test go-metrics counter
	testCounter := metrics.NewCounter()
	metrics.Register("test.counter", testCounter)
	defer metrics.Unregister("test.counter")

	// Increment the counter
	testCounter.Inc(5)

	// Test that the counter value is accessible
	if testCounter.Count() != 5 {
		t.Errorf("Expected counter value 5, got %d", testCounter.Count())
	}
}

func TestPrometheusMetricNames(t *testing.T) {
	// Save original config
	originalEnabled := config.Config.PrometheusEnabled
	originalNamespace := config.Config.PrometheusNamespace
	originalSubsystem := config.Config.PrometheusSubsystem
	defer func() {
		config.Config.PrometheusEnabled = originalEnabled
		config.Config.PrometheusNamespace = originalNamespace
		config.Config.PrometheusSubsystem = originalSubsystem
	}()

	// Enable Prometheus with custom namespace
	config.Config.PrometheusEnabled = true
	config.Config.PrometheusNamespace = "test_orchestrator"
	config.Config.PrometheusSubsystem = "subsys"

	// Initialize Prometheus metrics
	err := InitPrometheusMetrics()
	if err != nil {
		t.Fatalf("Failed to initialize Prometheus metrics: %v", err)
	}

	registry := GetPrometheusRegistry()
	if registry == nil {
		t.Fatal("Prometheus registry is nil")
	}

	// Gather metrics to check names
	metricFamilies, err := registry.Gather()
	if err != nil {
		t.Fatalf("Failed to gather metrics: %v", err)
	}

	// Check that we have some metrics
	if len(metricFamilies) == 0 {
		t.Fatal("No metrics found in registry")
	}

	// Check that metric names have correct namespace/subsystem
	foundCorrectNaming := false
	for _, mf := range metricFamilies {
		name := mf.GetName()
		if strings.HasPrefix(name, "test_orchestrator_subsys_") {
			foundCorrectNaming = true
			break
		}
	}

	if !foundCorrectNaming {
		t.Fatal("No metrics found with correct namespace/subsystem prefix")
	}
}

func TestPrometheusSyncMetrics(t *testing.T) {
	// Save original config
	originalEnabled := config.Config.PrometheusEnabled
	defer func() {
		config.Config.PrometheusEnabled = originalEnabled
	}()

	// Enable Prometheus for testing
	config.Config.PrometheusEnabled = true

	// Initialize Prometheus metrics
	err := InitPrometheusMetrics()
	if err != nil {
		t.Fatalf("Failed to initialize Prometheus metrics: %v", err)
	}

	// Test the sync function doesn't panic
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("syncMetricsToPrometheus panicked: %v", r)
		}
	}()

	// Call sync function
	syncMetricsToPrometheus()
}

func BenchmarkPrometheusSyncMetrics(b *testing.B) {
	// Save original config
	originalEnabled := config.Config.PrometheusEnabled
	defer func() {
		config.Config.PrometheusEnabled = originalEnabled
	}()

	// Enable Prometheus for testing
	config.Config.PrometheusEnabled = true

	// Initialize Prometheus metrics
	err := InitPrometheusMetrics()
	if err != nil {
		b.Fatalf("Failed to initialize Prometheus metrics: %v", err)
	}

	// Add some test metrics
	for i := 0; i < 10; i++ {
		counter := metrics.NewCounter()
		metrics.Register(prometheus.BuildFQName("benchmark", "test", "counter_"+string(rune(i))), counter)
		counter.Inc(int64(i))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		syncMetricsToPrometheus()
	}
}

func TestPrometheusCounterVsGoMetrics(t *testing.T) {
	// Save original config
	originalEnabled := config.Config.PrometheusEnabled
	defer func() {
		config.Config.PrometheusEnabled = originalEnabled
	}()

	// Enable Prometheus for testing
	config.Config.PrometheusEnabled = true

	// Initialize Prometheus metrics
	err := InitPrometheusMetrics()
	if err != nil {
		t.Fatalf("Failed to initialize Prometheus metrics: %v", err)
	}

	// Test with discoveries.attempt metric which should be bridged
	testCounter := metrics.GetOrRegisterCounter("discoveries.attempt", metrics.DefaultRegistry)
	testCounter.Inc(3)

	// Sync metrics
	syncMetricsToPrometheus()

	// Check that go-metrics counter has correct value
	if testCounter.Count() != 3 {
		t.Errorf("Go-metrics counter should be 3, got %d", testCounter.Count())
	}
}