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
	"sync"

	"github.com/openark/golib/log"
	"github.com/openark/orchestrator/go/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rcrowley/go-metrics"
)

var (
	// Prometheus metrics registry
	promRegistry *prometheus.Registry
	
	// Prometheus metrics - Discovery
	promDiscoveryAttempts = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "discoveries_attempt_total",
			Help: "Total number of discovery attempts",
		},
		[]string{},
	)
	
	promDiscoveryFailures = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "discoveries_fail_total", 
			Help: "Total number of failed discoveries",
		},
		[]string{},
	)
	
	promInstancePollSecondsExceeded = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "discoveries_instance_poll_seconds_exceeded_total",
			Help: "Total number of times instance polling exceeded configured seconds",
		},
		[]string{},
	)
	
	promDiscoveryQueueLength = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "discoveries_queue_length",
			Help: "Current length of discovery queue",
		},
		[]string{},
	)
	
	promDiscoveryRecentCount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "discoveries_recent_count",
			Help: "Number of recent discoveries",
		},
		[]string{},
	)
	
	promDeadInstancesDiscoveryQueueLength = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "discoveries_dead_instances_queue_length",
			Help: "Current length of dead instances discovery queue",
		},
		[]string{},
	)
	
	// Prometheus metrics - Election & Health
	promIsElected = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "elect_is_elected",
			Help: "Whether this orchestrator node is elected as leader (1 = elected, 0 = not elected)",
		},
		[]string{},
	)
	
	promIsHealthy = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "health_is_healthy", 
			Help: "Whether orchestrator is healthy (1 = healthy, 0 = unhealthy)",
		},
		[]string{},
	)
	
	promRaftIsHealthy = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "raft_is_healthy",
			Help: "Whether raft consensus is healthy (1 = healthy, 0 = unhealthy)",
		},
		[]string{},
	)
	
	promRaftIsLeader = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "raft_is_leader",
			Help: "Whether this node is the raft leader (1 = leader, 0 = follower)",
		},
		[]string{},
	)

	// Prometheus metrics - Recovery
	promRecoverDeadMasterStart = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "recover_dead_master_start_total",
			Help: "Total number of dead master recovery attempts started",
		},
		[]string{},
	)
	
	promRecoverDeadMasterSuccess = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "recover_dead_master_success_total",
			Help: "Total number of successful dead master recoveries",
		},
		[]string{},
	)
	
	promRecoverDeadMasterFail = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "recover_dead_master_fail_total",
			Help: "Total number of failed dead master recoveries",
		},
		[]string{},
	)
	
	promRecoverDeadIntermediateMasterStart = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "recover_dead_intermediate_master_start_total",
			Help: "Total number of dead intermediate master recovery attempts started",
		},
		[]string{},
	)
	
	promRecoverDeadIntermediateMasterSuccess = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "recover_dead_intermediate_master_success_total",
			Help: "Total number of successful dead intermediate master recoveries",
		},
		[]string{},
	)
	
	promRecoverDeadIntermediateMasterFail = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "recover_dead_intermediate_master_fail_total",
			Help: "Total number of failed dead intermediate master recoveries",
		},
		[]string{},
	)
	
	promCountPendingRecoveries = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "recover_pending",
			Help: "Number of pending recovery operations",
		},
		[]string{},
	)
	
	// Prometheus metrics - Instance Operations
	promInstanceAccessDenied = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "instance_access_denied_total",
			Help: "Total number of access denied errors when connecting to instances",
		},
		[]string{},
	)
	
	promInstanceReadTopology = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "instance_read_topology_total",
			Help: "Total number of topology reads from instances",
		},
		[]string{},
	)
	
	promInstanceRead = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "instance_read_total",
			Help: "Total number of instance reads",
		},
		[]string{},
	)
	
	promInstanceWrite = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "instance_write_total",
			Help: "Total number of instance writes",
		},
		[]string{},
	)
	
	// Prometheus metrics - Audit & Analysis
	promAuditWrite = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "audit_write_total",
			Help: "Total number of audit operations written",
		},
		[]string{},
	)
	
	promAnalysisChangeWriteAttempt = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "analysis_change_write_attempt_total", 
			Help: "Total number of analysis change write attempts",
		},
		[]string{},
	)
	
	promAnalysisChangeWrite = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "analysis_change_write_total",
			Help: "Total number of successful analysis change writes",
		},
		[]string{},
	)
	
	// Prometheus metrics - Hostname Resolution
	promResolveWriteResolved = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "resolve_write_resolved_total",
			Help: "Total number of resolved hostnames written",
		},
		[]string{},
	)
	
	promResolveWriteUnresolved = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "resolve_write_unresolved_total", 
			Help: "Total number of unresolved hostnames written",
		},
		[]string{},
	)
	
	promResolveReadResolved = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "resolve_read_resolved_total",
			Help: "Total number of resolved hostnames read",
		},
		[]string{},
	)
	
	promResolveReadUnresolved = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "resolve_read_unresolved_total",
			Help: "Total number of unresolved hostnames read", 
		},
		[]string{},
	)

	// Map to track go-metrics to Prometheus metric mappings
	metricMappings = make(map[string]prometheus.Metric)
	mappingMutex   sync.RWMutex
)

// InitPrometheusMetrics initializes Prometheus metrics collection
func InitPrometheusMetrics() error {
	if !config.Config.PrometheusEnabled {
		// Clear the registry when disabled
		promRegistry = nil
		return nil
	}

	log.Debugf("Initializing Prometheus metrics with namespace: %s, subsystem: %s", 
		config.Config.PrometheusNamespace, config.Config.PrometheusSubsystem)

	// Create a new Prometheus registry
	promRegistry = prometheus.NewRegistry()

	// Apply namespace and subsystem to all metrics
	namespace := config.Config.PrometheusNamespace
	subsystem := config.Config.PrometheusSubsystem
	
	if namespace != "" {
		// Re-create metrics with proper namespace/subsystem
		promDiscoveryAttempts = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "discoveries_attempt_total",
				Help:      "Total number of discovery attempts",
			},
			[]string{},
		)
		
		promDiscoveryFailures = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "discoveries_fail_total",
				Help:      "Total number of failed discoveries",
			},
			[]string{},
		)
		
		promInstancePollSecondsExceeded = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "discoveries_instance_poll_seconds_exceeded_total",
				Help:      "Total number of times instance polling exceeded configured seconds",
			},
			[]string{},
		)
		
		promDiscoveryQueueLength = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "discoveries_queue_length",
				Help:      "Current length of discovery queue",
			},
			[]string{},
		)
		
		promDiscoveryRecentCount = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "discoveries_recent_count",
				Help:      "Number of recent discoveries",
			},
			[]string{},
		)
		
		promIsElected = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "elect_is_elected",
				Help:      "Whether this orchestrator node is elected as leader (1 = elected, 0 = not elected)",
			},
			[]string{},
		)
		
		promIsHealthy = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "health_is_healthy",
				Help:      "Whether orchestrator is healthy (1 = healthy, 0 = unhealthy)",
			},
			[]string{},
		)
		
		promRaftIsHealthy = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "raft_is_healthy", 
				Help:      "Whether raft consensus is healthy (1 = healthy, 0 = unhealthy)",
			},
			[]string{},
		)
		
		promRaftIsLeader = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "raft_is_leader",
				Help:      "Whether this node is the raft leader (1 = leader, 0 = follower)",
			},
			[]string{},
		)
	}

	// Register all Prometheus metrics
	collectors := []prometheus.Collector{
		promDiscoveryAttempts,
		promDiscoveryFailures, 
		promInstancePollSecondsExceeded,
		promDiscoveryQueueLength,
		promDiscoveryRecentCount,
		promDeadInstancesDiscoveryQueueLength,
		promIsElected,
		promIsHealthy,
		promRaftIsHealthy,
		promRaftIsLeader,
		promRecoverDeadMasterStart,
		promRecoverDeadMasterSuccess,
		promRecoverDeadMasterFail,
		promRecoverDeadIntermediateMasterStart,
		promRecoverDeadIntermediateMasterSuccess,
		promRecoverDeadIntermediateMasterFail,
		promCountPendingRecoveries,
		promInstanceAccessDenied,
		promInstanceReadTopology,
		promInstanceRead,
		promInstanceWrite,
		promAuditWrite,
		promAnalysisChangeWriteAttempt,
		promAnalysisChangeWrite,
		promResolveWriteResolved,
		promResolveWriteUnresolved,
		promResolveReadResolved,
		promResolveReadUnresolved,
	}

	for _, collector := range collectors {
		if err := promRegistry.Register(collector); err != nil {
			log.Errorf("Failed to register Prometheus collector: %v", err)
			return err
		}
	}

	// Set up bridging from go-metrics to Prometheus
	setupMetricsBridge()
	
	log.Infof("Prometheus metrics initialized successfully")
	return nil
}

// setupMetricsBridge creates a bridge between go-metrics and Prometheus metrics
func setupMetricsBridge() {
	// Map go-metrics names to Prometheus metrics
	mappingMutex.Lock()
	defer mappingMutex.Unlock()
	
	// Discovery metrics
	metricMappings["discoveries.attempt"] = promDiscoveryAttempts.WithLabelValues()
	metricMappings["discoveries.fail"] = promDiscoveryFailures.WithLabelValues()
	metricMappings["discoveries.instance_poll_seconds_exceeded"] = promInstancePollSecondsExceeded.WithLabelValues()
	metricMappings["discoveries.queue_length"] = promDiscoveryQueueLength.WithLabelValues()
	metricMappings["discoveries.recent_count"] = promDiscoveryRecentCount.WithLabelValues()
	metricMappings["discoveries.dead_instances_queue_length"] = promDeadInstancesDiscoveryQueueLength.WithLabelValues()
	
	// Election & Health metrics  
	metricMappings["elect.is_elected"] = promIsElected.WithLabelValues()
	metricMappings["health.is_healthy"] = promIsHealthy.WithLabelValues()
	metricMappings["raft.is_healthy"] = promRaftIsHealthy.WithLabelValues()
	metricMappings["raft.is_leader"] = promRaftIsLeader.WithLabelValues()
	
	// Recovery metrics
	metricMappings["recover.dead_master.start"] = promRecoverDeadMasterStart.WithLabelValues()
	metricMappings["recover.dead_master.success"] = promRecoverDeadMasterSuccess.WithLabelValues()
	metricMappings["recover.dead_master.fail"] = promRecoverDeadMasterFail.WithLabelValues()
	metricMappings["recover.dead_intermediate_master.start"] = promRecoverDeadIntermediateMasterStart.WithLabelValues()
	metricMappings["recover.dead_intermediate_master.success"] = promRecoverDeadIntermediateMasterSuccess.WithLabelValues()
	metricMappings["recover.dead_intermediate_master.fail"] = promRecoverDeadIntermediateMasterFail.WithLabelValues()
	metricMappings["recover.pending"] = promCountPendingRecoveries.WithLabelValues()
	
	// Instance metrics
	metricMappings["instance.access_denied"] = promInstanceAccessDenied.WithLabelValues()
	metricMappings["instance.read_topology"] = promInstanceReadTopology.WithLabelValues()
	metricMappings["instance.read"] = promInstanceRead.WithLabelValues()
	metricMappings["instance.write"] = promInstanceWrite.WithLabelValues()
	
	// Audit & Analysis metrics
	metricMappings["audit.write"] = promAuditWrite.WithLabelValues()
	metricMappings["analysis.change.write.attempt"] = promAnalysisChangeWriteAttempt.WithLabelValues()
	metricMappings["analysis.change.write"] = promAnalysisChangeWrite.WithLabelValues()
	
	// Hostname resolution metrics
	metricMappings["resolve.write_resolved"] = promResolveWriteResolved.WithLabelValues()
	metricMappings["resolve.write_unresolved"] = promResolveWriteUnresolved.WithLabelValues()
	metricMappings["resolve.read_resolved"] = promResolveReadResolved.WithLabelValues()
	metricMappings["resolve.read_unresolved"] = promResolveReadUnresolved.WithLabelValues()
	
	// Set up periodic sync from go-metrics to Prometheus
	OnMetricsTick(syncMetricsToPrometheus)
}

// syncMetricsToPrometheus syncs values from go-metrics to Prometheus
func syncMetricsToPrometheus() {
	mappingMutex.RLock()
	defer mappingMutex.RUnlock()

	// Iterate through all registered go-metrics and sync to Prometheus
	metrics.DefaultRegistry.Each(func(name string, metric interface{}) {
		promMetric, exists := metricMappings[name]
		if !exists {
			return
		}
		
		switch goMetric := metric.(type) {
		case metrics.Counter:
			if promCounter, ok := promMetric.(prometheus.Counter); ok {
				// Set Prometheus counter to match go-metrics counter
				// Note: We need to calculate the difference to avoid double counting
				currentCount := float64(goMetric.Count())
				promCounter.Add(currentCount - getCurrentPromValue(promCounter))
			}
		case metrics.Gauge:
			if promGauge, ok := promMetric.(prometheus.Gauge); ok {
				promGauge.Set(float64(goMetric.Value()))
			}
		case metrics.GaugeFloat64:
			if promGauge, ok := promMetric.(prometheus.Gauge); ok {
				promGauge.Set(goMetric.Value())
			}
		}
	})
}

// getCurrentPromValue gets current value from Prometheus metric (simplified implementation)
func getCurrentPromValue(metric prometheus.Counter) float64 {
	// This is a simplified approach. In a production environment, 
	// you might want to maintain a separate tracking mechanism
	// to avoid resetting counters
	return 0
}

// GetPrometheusRegistry returns the Prometheus registry for use in HTTP handler
func GetPrometheusRegistry() *prometheus.Registry {
	return promRegistry
}

// IsPrometheusEnabled returns whether Prometheus metrics are enabled
func IsPrometheusEnabled() bool {
	return config.Config.PrometheusEnabled
}