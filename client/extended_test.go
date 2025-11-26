package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Extended test suite for improved coverage

// Topology Operations Tests

func TestRelocateBelow(t *testing.T) {
	ctx := context.Background()
	mockResponse := APIResponse{
		Code:    OK,
		Message: "Success",
		Details: map[string]interface{}{
			"Key": map[string]interface{}{
				"Hostname": "replica1.example.com",
				"Port":     3306,
			},
			"MasterKey": map[string]interface{}{
				"Hostname": "new-master.example.com",
				"Port":     3306,
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/relocate-below/replica1.example.com/3306/new-master.example.com/3306" {
			t.Errorf("Expected specific path, got '%s'", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	config := Config{Host: server.URL[7:]}
	client := NewClient(config)

	instanceKey := InstanceKey{"replica1.example.com", 3306}
	belowKey := InstanceKey{"new-master.example.com", 3306}
	
	instance, err := client.RelocateBelow(ctx, instanceKey, belowKey)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if instance == nil {
		t.Fatal("Expected instance to be returned")
	}
}

func TestRegroupReplicas(t *testing.T) {
	ctx := context.Background()
	mockResponse := APIResponse{
		Code:    OK,
		Message: "Success",
		Details: []map[string]interface{}{
			{
				"Key": map[string]interface{}{
					"Hostname": "replica1.example.com",
					"Port":     3306,
				},
			},
			{
				"Key": map[string]interface{}{
					"Hostname": "replica2.example.com",
					"Port":     3306,
				},
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/regroup-replicas/master.example.com/3306" {
			t.Errorf("Expected specific path, got '%s'", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	config := Config{Host: server.URL[7:]}
	client := NewClient(config)

	instanceKey := InstanceKey{"master.example.com", 3306}
	
	instances, err := client.RegroupReplicas(ctx, instanceKey)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(instances) != 2 {
		t.Errorf("Expected 2 instances, got %d", len(instances))
	}
}

// Replication Control Tests

func TestStartReplica(t *testing.T) {
	ctx := context.Background()
	mockResponse := APIResponse{
		Code:    OK,
		Message: "Success",
		Details: map[string]interface{}{
			"Key": map[string]interface{}{
				"Hostname": "replica.example.com",
				"Port":     3306,
			},
			"Slave_SQL_Running": true,
			"Slave_IO_Running":  true,
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	config := Config{Host: server.URL[7:]}
	client := NewClient(config)

	instanceKey := InstanceKey{"replica.example.com", 3306}
	
	instance, err := client.StartReplica(ctx, instanceKey)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if instance == nil {
		t.Fatal("Expected instance to be returned")
	}
}

func TestDelayReplication(t *testing.T) {
	ctx := context.Background()
	mockResponse := APIResponse{
		Code:    OK,
		Message: "Success",
		Details: map[string]interface{}{
			"Key": map[string]interface{}{
				"Hostname": "replica.example.com",
				"Port":     3306,
			},
			"SQLDelay": 3600,
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/delay-replication/replica.example.com/3306/3600" {
			t.Errorf("Expected path with delay seconds, got '%s'", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	config := Config{Host: server.URL[7:]}
	client := NewClient(config)

	instanceKey := InstanceKey{"replica.example.com", 3306}
	
	instance, err := client.DelayReplication(ctx, instanceKey, 3600)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if instance == nil {
		t.Fatal("Expected instance to be returned")
	}
}

// Recovery Operations Tests

func TestGetReplicationAnalysis(t *testing.T) {
	ctx := context.Background()
	mockResponse := APIResponse{
		Code:    OK,
		Message: "Success",
		Details: []map[string]interface{}{
			{
				"AnalyzedInstanceKey": map[string]interface{}{
					"Hostname": "failed-master.example.com",
					"Port":     3306,
				},
				"Analysis":   "DeadMaster",
				"IsActionableRecovery": true,
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/replication-analysis" {
			t.Errorf("Expected /api/replication-analysis, got '%s'", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	config := Config{Host: server.URL[7:]}
	client := NewClient(config)
	
	analyses, err := client.GetReplicationAnalysis(ctx)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(analyses) != 1 {
		t.Errorf("Expected 1 analysis, got %d", len(analyses))
	}
}

func TestRecoverInstance(t *testing.T) {
	ctx := context.Background()
	mockResponse := APIResponse{
		Code:    OK,
		Message: "Success",
		Details: map[string]interface{}{
			"Id":  int64(123),
			"UID": "recovery-uid-123",
			"AnalysisEntry": map[string]interface{}{
				"Analysis": "DeadMaster",
			},
			"IsSuccessful": true,
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/recover/failed-master.example.com/3306" {
			t.Errorf("Expected recovery path, got '%s'", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	config := Config{Host: server.URL[7:]}
	client := NewClient(config)

	instanceKey := InstanceKey{"failed-master.example.com", 3306}
	
	recovery, err := client.RecoverInstance(ctx, instanceKey)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if recovery == nil {
		t.Fatal("Expected recovery to be returned")
	}

	if recovery.Id != 123 {
		t.Errorf("Expected recovery ID 123, got %d", recovery.Id)
	}
}

func TestGracefulMasterTakeoverCluster(t *testing.T) {
	ctx := context.Background()
	mockResponse := APIResponse{
		Code:    OK,
		Message: "Success",
		Details: map[string]interface{}{
			"Id":  int64(456),
			"UID": "takeover-uid-456",
			"SuccessorKey": map[string]interface{}{
				"Hostname": "new-master.example.com",
				"Port":     3306,
			},
			"IsSuccessful": true,
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/graceful-master-takeover/my-cluster" {
			t.Errorf("Expected takeover path, got '%s'", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	config := Config{Host: server.URL[7:]}
	client := NewClient(config)
	
	recovery, err := client.GracefulMasterTakeoverCluster(ctx, "my-cluster")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if recovery == nil {
		t.Fatal("Expected recovery to be returned")
	}
}

// Agent Operations Tests

func TestGetAgents(t *testing.T) {
	ctx := context.Background()
	mockResponse := APIResponse{
		Code:    OK,
		Message: "Success",
		Details: []map[string]interface{}{
			{
				"Hostname": "agent1.example.com",
				"Port":     3001,
			},
			{
				"Hostname": "agent2.example.com",
				"Port":     3001,
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/agents" {
			t.Errorf("Expected /api/agents, got '%s'", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	config := Config{Host: server.URL[7:]}
	client := NewClient(config)
	
	agents, err := client.GetAgents(ctx)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(agents) != 2 {
		t.Errorf("Expected 2 agents, got %d", len(agents))
	}
}

func TestAgentSeed(t *testing.T) {
	ctx := context.Background()
	mockResponse := APIResponse{
		Code:    OK,
		Message: "Success",
		Details: map[string]interface{}{
			"SeedId":         int64(789),
			"TargetHostname": "target.example.com",
			"SourceHostname": "source.example.com",
			"IsComplete":     false,
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/agent-seed/target.example.com/source.example.com" {
			t.Errorf("Expected seed path, got '%s'", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	config := Config{Host: server.URL[7:]}
	client := NewClient(config)
	
	seed, err := client.AgentSeed(ctx, "target.example.com", "source.example.com")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if seed == nil {
		t.Fatal("Expected seed to be returned")
	}

	if seed.SeedId != 789 {
		t.Errorf("Expected seed ID 789, got %d", seed.SeedId)
	}
}

// Cluster Operations Tests

func TestGetClusterInfo(t *testing.T) {
	ctx := context.Background()
	mockResponse := APIResponse{
		Code:    OK,
		Message: "Success",
		Details: map[string]interface{}{
			"ClusterName":                  "my-cluster",
			"ClusterAlias":                 "prod-cluster",
			"CountInstances":               10,
			"HasAutomatedMasterRecovery":   true,
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/cluster-info/my-cluster" {
			t.Errorf("Expected cluster-info path, got '%s'", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	config := Config{Host: server.URL[7:]}
	client := NewClient(config)
	
	info, err := client.GetClusterInfo(ctx, "my-cluster")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if info == nil {
		t.Fatal("Expected cluster info to be returned")
	}

	if info.ClusterName != "my-cluster" {
		t.Errorf("Expected cluster name 'my-cluster', got '%s'", info.ClusterName)
	}
}

func TestGetClusterMaster(t *testing.T) {
	ctx := context.Background()
	mockResponse := APIResponse{
		Code:    OK,
		Message: "Success",
		Details: map[string]interface{}{
			"Key": map[string]interface{}{
				"Hostname": "master.example.com",
				"Port":     3306,
			},
			"ReadOnly": false,
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/master/my-cluster" {
			t.Errorf("Expected master path, got '%s'", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	config := Config{Host: server.URL[7:]}
	client := NewClient(config)
	
	master, err := client.GetClusterMaster(ctx, "my-cluster")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if master == nil {
		t.Fatal("Expected master instance to be returned")
	}

	if master.ReadOnly {
		t.Error("Expected master to be writable")
	}
}

// Downtime Operations Tests

func TestBeginDowntime(t *testing.T) {
	ctx := context.Background()
	mockResponse := APIResponse{
		Code:    OK,
		Message: "Success",
		Details: map[string]interface{}{
			"Key": map[string]interface{}{
				"Hostname": "db1.example.com",
				"Port":     3306,
			},
			"IsDowntimed":    true,
			"DowntimeOwner":  "admin",
			"DowntimeReason": "maintenance",
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := "/api/begin-downtime/db1.example.com/3306/admin/maintenance"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path '%s', got '%s'", expectedPath, r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	config := Config{Host: server.URL[7:]}
	client := NewClient(config)

	instanceKey := InstanceKey{"db1.example.com", 3306}
	
	instance, err := client.BeginDowntime(ctx, instanceKey, "admin", "maintenance")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if instance == nil {
		t.Fatal("Expected instance to be returned")
	}

	if !instance.IsDowntimed {
		t.Error("Expected instance to be downtimed")
	}
}

func TestGetDowntimedInstances(t *testing.T) {
	ctx := context.Background()
	mockResponse := APIResponse{
		Code:    OK,
		Message: "Success",
		Details: []map[string]interface{}{
			{
				"Key": map[string]interface{}{
					"Hostname": "db1.example.com",
					"Port":     3306,
				},
				"IsDowntimed": true,
			},
			{
				"Key": map[string]interface{}{
					"Hostname": "db2.example.com",
					"Port":     3306,
				},
				"IsDowntimed": true,
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/downtimed" {
			t.Errorf("Expected /api/downtimed, got '%s'", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	config := Config{Host: server.URL[7:]}
	client := NewClient(config)
	
	instances, err := client.GetDowntimedInstances(ctx)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(instances) != 2 {
		t.Errorf("Expected 2 downtimed instances, got %d", len(instances))
	}
}

// Tagging Operations Tests

func TestTagInstance(t *testing.T) {
	ctx := context.Background()
	mockResponse := APIResponse{
		Code:    OK,
		Message: "Success",
		Details: map[string]interface{}{
			"Key": map[string]interface{}{
				"Hostname": "db1.example.com",
				"Port":     3306,
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := "/api/tag/db1.example.com/3306/environment/production"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path '%s', got '%s'", expectedPath, r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	config := Config{Host: server.URL[7:]}
	client := NewClient(config)

	instanceKey := InstanceKey{"db1.example.com", 3306}
	
	instance, err := client.TagInstance(ctx, instanceKey, "environment", "production")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if instance == nil {
		t.Fatal("Expected instance to be returned")
	}
}

func TestGetInstanceTags(t *testing.T) {
	ctx := context.Background()
	mockResponse := APIResponse{
		Code:    OK,
		Message: "Success",
		Details: []map[string]interface{}{
			{
				"TagName":  "environment",
				"TagValue": "production",
			},
			{
				"TagName":  "datacenter",
				"TagValue": "us-west-2",
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/tags/db1.example.com/3306" {
			t.Errorf("Expected tags path, got '%s'", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	config := Config{Host: server.URL[7:]}
	client := NewClient(config)

	instanceKey := InstanceKey{"db1.example.com", 3306}
	
	tags, err := client.GetInstanceTags(ctx, instanceKey)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(tags) != 2 {
		t.Errorf("Expected 2 tags, got %d", len(tags))
	}
}

// GTID Operations Tests

func TestEnableGTID(t *testing.T) {
	ctx := context.Background()
	mockResponse := APIResponse{
		Code:    OK,
		Message: "Success",
		Details: map[string]interface{}{
			"Key": map[string]interface{}{
				"Hostname": "db1.example.com",
				"Port":     3306,
			},
			"GTIDMode":        "ON",
			"UsingOracleGTID": true,
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/enable-gtid/db1.example.com/3306" {
			t.Errorf("Expected enable-gtid path, got '%s'", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	config := Config{Host: server.URL[7:]}
	client := NewClient(config)

	instanceKey := InstanceKey{"db1.example.com", 3306}
	
	instance, err := client.EnableGTID(ctx, instanceKey)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if instance == nil {
		t.Fatal("Expected instance to be returned")
	}
}

func TestLocateErrantGTID(t *testing.T) {
	ctx := context.Background()
	mockResponse := APIResponse{
		Code:    OK,
		Message: "Success",
		Details: map[string]interface{}{
			"Key": map[string]interface{}{
				"Hostname": "db1.example.com",
				"Port":     3306,
			},
			"GtidErrant": "uuid:1-5",
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/locate-gtid-errant/db1.example.com/3306" {
			t.Errorf("Expected locate-gtid-errant path, got '%s'", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	config := Config{Host: server.URL[7:]}
	client := NewClient(config)

	instanceKey := InstanceKey{"db1.example.com", 3306}
	
	instance, err := client.LocateErrantGTID(ctx, instanceKey)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if instance == nil {
		t.Fatal("Expected instance to be returned")
	}
}

// Raft Operations Tests

func TestGetRaftState(t *testing.T) {
	ctx := context.Background()
	mockResponse := APIResponse{
		Code:    OK,
		Message: "Success",
		Details: map[string]interface{}{
			"Leader":     "node1:10008",
			"Peer":       "node2:10008",
			"Peers":      []string{"node1:10008", "node2:10008", "node3:10008"},
			"IsLeader":   true,
			"IsFollower": false,
			"State":      "Leader",
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/raft-state" {
			t.Errorf("Expected /api/raft-state, got '%s'", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	config := Config{Host: server.URL[7:]}
	client := NewClient(config)
	
	state, err := client.GetRaftState(ctx)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if state == nil {
		t.Fatal("Expected raft state to be returned")
	}

	if !state.IsLeader {
		t.Error("Expected node to be leader")
	}

	if len(state.Peers) != 3 {
		t.Errorf("Expected 3 peers, got %d", len(state.Peers))
	}
}

func TestAddRaftPeer(t *testing.T) {
	ctx := context.Background()
	mockResponse := APIResponse{
		Code:    OK,
		Message: "Success",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/raft-add-peer/node4:10008" {
			t.Errorf("Expected raft-add-peer path, got '%s'", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	config := Config{Host: server.URL[7:]}
	client := NewClient(config)
	
	err := client.AddRaftPeer(ctx, "node4:10008")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}

// Monitoring Operations Tests

func TestGetDiscoveryMetricsRaw(t *testing.T) {
	ctx := context.Background()
	mockResponse := APIResponse{
		Code:    OK,
		Message: "Success",
		Details: []map[string]interface{}{
			{
				"Timestamp": "2024-01-01T10:00:00Z",
				"InstanceKey": map[string]interface{}{
					"Hostname": "db1.example.com",
					"Port":     3306,
				},
				"DurationMillis": int64(150),
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/discovery-metrics-raw/60" {
			t.Errorf("Expected discovery-metrics-raw path, got '%s'", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	config := Config{Host: server.URL[7:]}
	client := NewClient(config)
	
	metrics, err := client.GetDiscoveryMetricsRaw(ctx, 60)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(metrics) != 1 {
		t.Errorf("Expected 1 metric, got %d", len(metrics))
	}
}

// Pool Operations Tests

func TestGetClusterPoolInstances(t *testing.T) {
	ctx := context.Background()
	mockResponse := APIResponse{
		Code:    OK,
		Message: "Success",
		Details: map[string]interface{}{
			"read-pool": []map[string]interface{}{
				{
					"Hostname": "replica1.example.com",
					"Port":     3306,
				},
				{
					"Hostname": "replica2.example.com",
					"Port":     3306,
				},
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/cluster-pool-instances/my-cluster" {
			t.Errorf("Expected cluster-pool-instances path, got '%s'", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	config := Config{Host: server.URL[7:]}
	client := NewClient(config)
	
	pools, err := client.GetClusterPoolInstances(ctx, "my-cluster")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(pools) == 0 {
		t.Error("Expected at least one pool")
	}
}

// Audit Operations Tests

func TestGetAudit(t *testing.T) {
	ctx := context.Background()
	mockResponse := APIResponse{
		Code:    OK,
		Message: "Success",
		Details: []map[string]interface{}{
			{
				"AuditId":        int64(1),
				"AuditTimestamp": "2024-01-01T10:00:00Z",
				"AuditType":      "recover",
				"Message":        "Recovery completed successfully",
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/audit/0" {
			t.Errorf("Expected audit path with page, got '%s'", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	config := Config{Host: server.URL[7:]}
	client := NewClient(config)
	
	entries, err := client.GetAudit(ctx, 0)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(entries) != 1 {
		t.Errorf("Expected 1 audit entry, got %d", len(entries))
	}
}

// Error Handling Tests

func TestAPIErrorResponse(t *testing.T) {
	ctx := context.Background()
	mockResponse := APIResponse{
		Code:    ERROR,
		Message: "Instance not found",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	config := Config{Host: server.URL[7:]}
	client := NewClient(config)

	instanceKey := InstanceKey{"nonexistent.example.com", 3306}
	
	_, err := client.GetInstance(ctx, instanceKey)
	if err == nil {
		t.Fatal("Expected error, got none")
	}

	if err.Error() != "API error: Instance not found" {
		t.Errorf("Unexpected error message: %s", err.Error())
	}
}

func TestHTTPErrorResponse(t *testing.T) {
	ctx := context.Background()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	config := Config{Host: server.URL[7:]}
	client := NewClient(config)

	instanceKey := InstanceKey{"db1.example.com", 3306}
	
	_, err := client.GetInstance(ctx, instanceKey)
	if err == nil {
		t.Fatal("Expected error, got none")
	}
}

func TestNetworkError(t *testing.T) {
	ctx := context.Background()

	// Use a non-existent server
	config := Config{
		Host:    "localhost:99999",
		Timeout: 100, // 100ms timeout
	}
	client := NewClient(config)

	instanceKey := InstanceKey{"db1.example.com", 3306}
	
	_, err := client.GetInstance(ctx, instanceKey)
	if err == nil {
		t.Fatal("Expected network error, got none")
	}
}


