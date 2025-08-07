# Prometheus Metrics

Orchestrator supports native Prometheus metrics collection and export. This document describes how to configure and use Prometheus metrics with Orchestrator.

## Configuration

Prometheus metrics are disabled by default. To enable them, add the following configuration options to your `orchestrator.conf.json`:

```json
{
  "PrometheusEnabled": true,
  "PrometheusPath": "/metrics",
  "PrometheusNamespace": "orchestrator",
  "PrometheusSubsystem": ""
}
```

### Configuration Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `PrometheusEnabled` | boolean | `false` | Enable or disable Prometheus metrics collection |
| `PrometheusPath` | string | `"/metrics"` | HTTP path where Prometheus metrics will be exposed |
| `PrometheusNamespace` | string | `"orchestrator"` | Namespace prefix for all Prometheus metrics |
| `PrometheusSubsystem` | string | `""` | Optional subsystem prefix for metrics |

## Metrics Endpoint

When enabled, Prometheus metrics are available at the configured path (default: `/metrics`).

Example:
```bash
curl http://orchestrator:3000/metrics
```

## Available Metrics

Orchestrator exports the following Prometheus metrics:

### Discovery Metrics

- `orchestrator_discoveries_attempt_total` - Total number of discovery attempts
- `orchestrator_discoveries_fail_total` - Total number of failed discoveries  
- `orchestrator_discoveries_instance_poll_seconds_exceeded_total` - Discovery attempts that exceeded configured timeout
- `orchestrator_discoveries_queue_length` - Current length of discovery queue
- `orchestrator_discoveries_recent_count` - Number of recent discoveries
- `orchestrator_discoveries_dead_instances_queue_length` - Length of dead instances discovery queue

### Election & Health Metrics

- `orchestrator_elect_is_elected` - Whether this node is elected as leader (1 = elected, 0 = not elected)
- `orchestrator_health_is_healthy` - Whether orchestrator is healthy (1 = healthy, 0 = unhealthy)
- `orchestrator_raft_is_healthy` - Whether raft consensus is healthy (1 = healthy, 0 = unhealthy)
- `orchestrator_raft_is_leader` - Whether this node is the raft leader (1 = leader, 0 = follower)

### Recovery Metrics

- `orchestrator_recover_dead_master_start_total` - Dead master recovery attempts started
- `orchestrator_recover_dead_master_success_total` - Successful dead master recoveries
- `orchestrator_recover_dead_master_fail_total` - Failed dead master recoveries
- `orchestrator_recover_dead_intermediate_master_start_total` - Dead intermediate master recovery attempts
- `orchestrator_recover_dead_intermediate_master_success_total` - Successful dead intermediate master recoveries  
- `orchestrator_recover_dead_intermediate_master_fail_total` - Failed dead intermediate master recoveries
- `orchestrator_recover_pending` - Number of pending recovery operations

### Instance Operation Metrics

- `orchestrator_instance_access_denied_total` - Access denied errors when connecting to instances
- `orchestrator_instance_read_topology_total` - Topology reads from instances
- `orchestrator_instance_read_total` - Instance read operations
- `orchestrator_instance_write_total` - Instance write operations

### Audit & Analysis Metrics

- `orchestrator_audit_write_total` - Audit operations written
- `orchestrator_analysis_change_write_attempt_total` - Analysis change write attempts
- `orchestrator_analysis_change_write_total` - Successful analysis change writes

### Hostname Resolution Metrics

- `orchestrator_resolve_write_resolved_total` - Resolved hostnames written
- `orchestrator_resolve_write_unresolved_total` - Unresolved hostnames written
- `orchestrator_resolve_read_resolved_total` - Resolved hostnames read
- `orchestrator_resolve_read_unresolved_total` - Unresolved hostnames read

## Example Prometheus Configuration

Add the following to your `prometheus.yml`:

```yaml
scrape_configs:
  - job_name: 'orchestrator'
    static_configs:
      - targets: ['orchestrator:3000']
    metrics_path: /metrics
    scrape_interval: 30s
```

For multiple Orchestrator instances:

```yaml
scrape_configs:
  - job_name: 'orchestrator'
    static_configs:
      - targets: 
        - 'orchestrator1:3000'
        - 'orchestrator2:3000'  
        - 'orchestrator3:3000'
    metrics_path: /metrics
    scrape_interval: 30s
```

## Example Queries

### Basic Health Monitoring

```promql
# Check if orchestrator is healthy
orchestrator_health_is_healthy

# Check which node is the elected leader
orchestrator_elect_is_elected

# Check raft cluster health
orchestrator_raft_is_healthy
```

### Discovery Monitoring

```promql
# Discovery success rate
rate(orchestrator_discoveries_attempt_total[5m]) - rate(orchestrator_discoveries_fail_total[5m])

# Discovery queue length
orchestrator_discoveries_queue_length

# Failed discoveries per minute
rate(orchestrator_discoveries_fail_total[1m]) * 60
```

### Recovery Monitoring

```promql
# Master recovery success rate
rate(orchestrator_recover_dead_master_success_total[5m]) / 
rate(orchestrator_recover_dead_master_start_total[5m])

# Pending recoveries
orchestrator_recover_pending

# Recovery attempts by type
rate(orchestrator_recover_dead_master_start_total[5m])
```

## Alerting Rules

Example Prometheus alerting rules for Orchestrator:

```yaml
groups:
- name: orchestrator
  rules:
  - alert: OrchestratorDown
    expr: up{job="orchestrator"} == 0
    for: 1m
    labels:
      severity: critical
    annotations:
      summary: "Orchestrator instance is down"
      
  - alert: OrchestratorUnhealthy  
    expr: orchestrator_health_is_healthy == 0
    for: 2m
    labels:
      severity: warning
    annotations:
      summary: "Orchestrator reports unhealthy status"
      
  - alert: OrchestratorNoLeader
    expr: sum(orchestrator_elect_is_elected) == 0
    for: 30s
    labels:
      severity: critical
    annotations:
      summary: "No orchestrator leader elected"
      
  - alert: OrchestratorDiscoveryFailures
    expr: rate(orchestrator_discoveries_fail_total[5m]) > 0.1
    for: 2m
    labels:
      severity: warning
    annotations:
      summary: "High discovery failure rate"
      
  - alert: OrchestratorPendingRecoveries
    expr: orchestrator_recover_pending > 0
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: "Recovery operations are pending"
```

## Grafana Dashboard

Orchestrator includes a Grafana dashboard template at `resources/metrics/orchestrator-grafana.json`. This dashboard can be adapted for Prometheus by:

1. Changing the data source from Graphite to Prometheus
2. Converting metric names to Prometheus format
3. Updating queries to use PromQL syntax

## Compatibility with Graphite

Orchestrator can simultaneously export metrics to both Prometheus and Graphite. The existing Graphite configuration continues to work alongside the new Prometheus metrics.

## Performance Considerations

- Prometheus metrics collection has minimal performance overhead
- Metrics are synchronized periodically from the existing go-metrics registry  
- The `/metrics` endpoint is lightweight and safe to scrape frequently
- Consider the scrape interval based on your monitoring requirements (15-30s is typical)

## Troubleshooting

### Metrics Not Available

1. Verify `PrometheusEnabled` is set to `true` in configuration
2. Check that Orchestrator has restarted after configuration change
3. Verify the metrics endpoint is accessible: `curl http://orchestrator:3000/metrics`

### Missing Metrics

Some metrics may only appear after relevant events occur (e.g., recovery metrics only appear after recovery operations).

### High Cardinality

Orchestrator metrics use minimal labels to avoid high cardinality issues. Additional labels can be added at the Prometheus scrape level if needed.

## Migration from Graphite

If migrating from Graphite to Prometheus:

1. Enable Prometheus metrics while keeping Graphite enabled
2. Verify all required metrics are available in Prometheus
3. Update monitoring dashboards and alerts  
4. Disable Graphite once migration is complete

## Security Considerations

The `/metrics` endpoint exposes operational metrics but no sensitive data. However, you may want to:

- Restrict access to the metrics endpoint via firewall rules
- Use Orchestrator's authentication mechanisms if needed
- Monitor access to the metrics endpoint for security auditing