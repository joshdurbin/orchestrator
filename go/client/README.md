# Orchestrator Go Client

The Orchestrator Go client provides a comprehensive programmatic interface to the Orchestrator REST API. It allows you to perform all orchestrator operations including topology discovery, replication management, failover operations, and cluster management from Go applications.

## Features

- **Complete API Coverage**: All Orchestrator REST endpoints are supported
- **Automatic Leader Detection**: Supports multiple endpoints with automatic leader discovery
- **Authentication**: HTTP Basic Auth and custom header authentication
- **Error Handling**: Comprehensive error handling with typed responses
- **Type Safety**: Uses existing Orchestrator types for consistency
- **Context Support**: All operations support Go contexts for timeouts and cancellation
- **Connection Management**: Automatic retry and connection management

## Installation

```bash
go get github.com/openark/orchestrator/go/client
```

## Quick Start

### Basic Usage

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/openark/orchestrator/go/client"
    "github.com/openark/orchestrator/go/inst"
)

func main() {
    // Create a client
    config := &client.Config{
        BaseURL: "http://orchestrator.example.com:3000",
        Timeout: 30 * time.Second,
    }

    c, err := client.NewClient(config)
    if err != nil {
        log.Fatal(err)
    }

    ctx := context.Background()

    // Check health
    if err := c.Health(ctx); err != nil {
        log.Fatal("Health check failed:", err)
    }

    // Discover an instance
    instanceKey := inst.InstanceKey{
        Hostname: "mysql1.example.com",
        Port:     3306,
    }

    instance, err := c.Discover(ctx, instanceKey)
    if err != nil {
        log.Fatal("Discovery failed:", err)
    }

    fmt.Printf("Discovered: %s:%d\n", instance.Key.Hostname, instance.Key.Port)
}
```

### High Availability Setup

```go
// Configure client for HA orchestrator setup
config := &client.Config{
    Endpoints: []string{
        "http://orchestrator1.example.com:3000/api",
        "http://orchestrator2.example.com:3000/api", 
        "http://orchestrator3.example.com:3000/api",
    },
    Username: "admin",
    Password: "secret",
    Timeout:  30 * time.Second,
}

c, err := client.NewClient(config)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Connected to leader: %s\n", c.GetLeader())
```

### Authentication

#### HTTP Basic Auth
```go
config := &client.Config{
    BaseURL:  "http://orchestrator.example.com:3000",
    Username: "admin",
    Password: "secret123",
}
```

#### Custom Headers
```go
config := &client.Config{
    BaseURL: "http://orchestrator.example.com:3000",
    Headers: map[string]string{
        "X-Auth-User": "admin@example.com",
        "X-Auth-Token": "abc123",
    },
}
```

## API Operations

### Instance Operations

```go
ctx := context.Background()
instanceKey := inst.InstanceKey{Hostname: "mysql1.example.com", Port: 3306}

// Discovery
instance, err := c.Discover(ctx, instanceKey)
instance, err := c.GetInstance(ctx, instanceKey)
replicas, err := c.GetInstanceReplicas(ctx, instanceKey)

// Search
instances, err := c.Search(ctx, "mysql1")

// Lifecycle
err = c.Forget(ctx, instanceKey)
```

### Cluster Operations

```go
// Get cluster information
clusters, err := c.GetClusters(ctx)
instances, err := c.GetCluster(ctx, "main-cluster")
master, err := c.GetClusterMaster(ctx, "main-cluster")
masters, err := c.GetMasters(ctx)

// Topology visualization
topology, err := c.GetTopology(ctx, "main-cluster")
fmt.Println(topology)
```

### Replication Control

```go
instanceKey := inst.InstanceKey{Hostname: "mysql1.example.com", Port: 3306}

// Basic replication control
instance, err := c.StartReplication(ctx, instanceKey)
instance, err := c.StopReplication(ctx, instanceKey) 
instance, err := c.RestartReplication(ctx, instanceKey)
instance, err := c.ResetReplication(ctx, instanceKey)

// Advanced operations
instance, err := c.SetReadOnly(ctx, instanceKey)
instance, err := c.SetWriteable(ctx, instanceKey)
instance, err := c.SkipQuery(ctx, instanceKey)
```

### Topology Manipulation

```go
replica := inst.InstanceKey{Hostname: "mysql-replica1.example.com", Port: 3306}
master := inst.InstanceKey{Hostname: "mysql-master2.example.com", Port: 3306}

// Move operations
instance, err := c.MoveUp(ctx, replica)
instances, err := c.MoveUpReplicas(ctx, replica)
instance, err := c.MoveBelow(ctx, replica, master)

// Relocate operations
instance, err := c.Relocate(ctx, replica, master)
instances, err := c.RelocateReplicas(ctx, replica, master)

// GTID operations
instance, err := c.MoveGTID(ctx, replica, master)
instances, err := c.MoveReplicasGTID(ctx, replica, master)

// Pseudo-GTID operations
instance, err := c.MatchBelow(ctx, replica, master)
instance, err := c.MatchUp(ctx, replica)
```

### Failover and Recovery

```go
// Graceful operations
newMaster, err := c.GracefulMasterTakeover(ctx, "main-cluster", nil)
newMaster, err := c.GracefulMasterTakeoverAuto(ctx, "main-cluster", nil)

// Forced operations  
newMaster, err := c.ForceMasterFailover(ctx, "main-cluster")
designated := inst.InstanceKey{Hostname: "mysql2.example.com", Port: 3306}
newMaster, err := c.ForceMasterTakeover(ctx, "main-cluster", designated)

// Recovery
failed := inst.InstanceKey{Hostname: "mysql1.example.com", Port: 3306}
candidate := inst.InstanceKey{Hostname: "mysql2.example.com", Port: 3306}
recovered, err := c.Recover(ctx, failed, &candidate)

// Analysis
analysis, err := c.GetReplicationAnalysis(ctx)
```

### Maintenance and Downtime

```go
instanceKey := inst.InstanceKey{Hostname: "mysql1.example.com", Port: 3306}

// Maintenance
err = c.BeginMaintenance(ctx, instanceKey, "admin", "Scheduled maintenance")
inMaintenance, err := c.InMaintenance(ctx, instanceKey)
err = c.EndMaintenance(ctx, instanceKey)

// Downtime
duration := 1 * time.Hour
err = c.BeginDowntime(ctx, instanceKey, "admin", "OS update", duration)
err = c.EndDowntime(ctx, instanceKey)
```

### Tagging

```go
instanceKey := inst.InstanceKey{Hostname: "mysql1.example.com", Port: 3306}

// Set and get tags
err = c.SetTag(ctx, instanceKey, "environment", "production")
err = c.SetTag(ctx, instanceKey, "role", "master")
tags, err := c.GetTags(ctx, instanceKey)

// Query by tags
prodInstances, err := c.GetTagged(ctx, "environment", "production")

// Remove tags
err = c.RemoveTag(ctx, instanceKey, "role")
instances, err := c.RemoveTagFromAll(ctx, "environment", "staging")
```

## Error Handling

The client provides detailed error information:

```go
instance, err := c.Discover(ctx, instanceKey)
if err != nil {
    // Check for specific error types
    if strings.Contains(err.Error(), "not found") {
        // Handle instance not found
    } else if strings.Contains(err.Error(), "timeout") {
        // Handle timeout
    } else {
        // Handle other errors
        log.Printf("Discovery failed: %v", err)
    }
}
```

## Configuration Options

### Client Config

| Field | Type | Description |
|-------|------|-------------|
| `BaseURL` | `string` | Primary Orchestrator API endpoint |
| `Endpoints` | `[]string` | Multiple endpoints for HA setup |
| `Username` | `string` | HTTP Basic Auth username |
| `Password` | `string` | HTTP Basic Auth password |
| `Headers` | `map[string]string` | Custom HTTP headers |
| `Timeout` | `time.Duration` | Request timeout (default: 30s) |
| `InsecureSkipVerify` | `bool` | Skip TLS verification |

### Environment Variables

You can also configure the client using environment variables (similar to orchestrator-client):

```bash
export ORCHESTRATOR_API="http://orchestrator.example.com:3000/api"
export ORCHESTRATOR_AUTH_USER="admin"
export ORCHESTRATOR_AUTH_PASSWORD="secret"
```

## Testing

Run the tests:

```bash
go test ./...
```

Run benchmarks:

```bash
go test -bench=. ./...
```

## Migration from orchestrator-client

If you're migrating from the bash `orchestrator-client`, here are the equivalent operations:

| orchestrator-client | Go Client |
|-------------------|-----------|
| `orchestrator-client -c discover -i host:3306` | `c.Discover(ctx, instanceKey)` |
| `orchestrator-client -c clusters` | `c.GetClusters(ctx)` |
| `orchestrator-client -c topology -i cluster` | `c.GetTopology(ctx, cluster)` |
| `orchestrator-client -c relocate -i replica -d master` | `c.Relocate(ctx, replica, master)` |
| `orchestrator-client -c graceful-master-takeover -i cluster` | `c.GracefulMasterTakeover(ctx, cluster, nil)` |

## Contributing

Contributions are welcome! Please ensure:

1. All new API endpoints are covered
2. Tests are included for new functionality  
3. Documentation is updated
4. Error handling follows existing patterns

## License

Licensed under the Apache License, Version 2.0. See LICENSE for details.