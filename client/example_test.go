package client_test

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/openark/orchestrator/client"
)

// Example demonstrating basic client usage
func Example_basic() {
	ctx := context.Background()
	// Create a new client
	cfg := client.Config{
		Host:     "localhost:3000",
		Username: "admin",
		Password: "secret",
	}

	c := client.NewClient(cfg)

	// Get all clusters
	clusters, err := c.GetClusters(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Found %d clusters\n", len(clusters))
}

// Example demonstrating instance management
func Example_instanceManagement() {
	ctx := context.Background()
	cfg := client.Config{
		Host: "localhost:3000",
	}

	c := client.NewClient(cfg)

	instanceKey := client.InstanceKey{
		Hostname: "db1.example.com",
		Port:     3306,
	}

	// Discover instance
	instance, err := c.DiscoverInstance(ctx, instanceKey)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Instance: %s, Version: %s\n", instanceKey, instance.Version)

	// Set read-only
	_, err = c.SetReadOnly(ctx, instanceKey)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Instance %s set to read-only\n", instanceKey)
}

// Example demonstrating topology operations
func Example_topologyOperations() {
	ctx := context.Background()
	cfg := client.Config{
		Host: "localhost:3000",
	}

	c := client.NewClient(cfg)

	replicaKey := client.InstanceKey{
		Hostname: "replica1.example.com",
		Port:     3306,
	}

	newMasterKey := client.InstanceKey{
		Hostname: "master2.example.com",
		Port:     3306,
	}

	// Relocate replica under new master
	_, err := c.Relocate(ctx, replicaKey, newMasterKey)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Relocated %s below %s\n", replicaKey, newMasterKey)
}

// Example demonstrating maintenance windows
func Example_maintenance() {
	ctx := context.Background()
	cfg := client.Config{
		Host: "localhost:3000",
	}

	c := client.NewClient(cfg)

	instanceKey := client.InstanceKey{
		Hostname: "db1.example.com",
		Port:     3306,
	}

	// Begin maintenance for 1 hour
	maintenance, err := c.BeginMaintenanceWithDuration(ctx, 
		instanceKey,
		"admin",
		"Planned maintenance",
		time.Hour,
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Maintenance window started: %d\n", maintenance.MaintenanceId)

	// Check if in maintenance
	inMaintenance, err := c.InMaintenance(ctx, instanceKey)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("In maintenance: %v\n", inMaintenance)

	// End maintenance
	_, err = c.EndMaintenanceByInstance(ctx, instanceKey)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Maintenance ended")
}

// Example demonstrating recovery operations
func Example_recovery() {
	ctx := context.Background()
	cfg := client.Config{
		Host: "localhost:3000",
	}

	c := client.NewClient(cfg)

	// Get replication analysis
	analyses, err := c.GetReplicationAnalysis(ctx)
	if err != nil {
		log.Fatal(err)
	}

	for _, analysis := range analyses {
		fmt.Printf("Cluster: %s\n", analysis.ClusterDetails.ClusterName)
		fmt.Printf("Analysis: %s\n", analysis.Analysis)
		fmt.Printf("Actionable: %v\n", analysis.IsActionableRecovery)
	}

	// If recovery is needed
	instanceKey := client.InstanceKey{
		Hostname: "failed-master.example.com",
		Port:     3306,
	}

	recovery, err := c.RecoverInstance(ctx, instanceKey)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Recovery initiated: %s\n", recovery.UID)
}

// Example demonstrating graceful master takeover
func Example_gracefulTakeover() {
	ctx := context.Background()
	cfg := client.Config{
		Host: "localhost:3000",
	}

	c := client.NewClient(cfg)

	clusterName := "production-cluster"

	// Perform graceful master takeover
	recovery, err := c.GracefulMasterTakeoverCluster(ctx, clusterName)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Takeover completed for cluster %s\n", clusterName)
	fmt.Printf("New master: %v\n", recovery.SuccessorKey)
}

// Example demonstrating tagging
func Example_tagging() {
	ctx := context.Background()
	cfg := client.Config{
		Host: "localhost:3000",
	}

	c := client.NewClient(cfg)

	instanceKey := client.InstanceKey{
		Hostname: "db1.example.com",
		Port:     3306,
	}

	// Tag instance
	_, err := c.TagInstance(ctx, instanceKey, "environment", "production")
	if err != nil {
		log.Fatal(err)
	}

	_, err = c.TagInstance(ctx, instanceKey, "role", "primary")
	if err != nil {
		log.Fatal(err)
	}

	// Get tags
	tags, err := c.GetInstanceTags(ctx, instanceKey)
	if err != nil {
		log.Fatal(err)
	}

	for _, tag := range tags {
		fmt.Printf("%s: %s\n", tag.TagName, tag.TagValue)
	}
}

// Example demonstrating cluster operations
func Example_clusterOperations() {
	ctx := context.Background()
	cfg := client.Config{
		Host: "localhost:3000",
	}

	c := client.NewClient(cfg)

	// Get all clusters
	clusters, err := c.GetClusters(ctx)
	if err != nil {
		log.Fatal(err)
	}

	for _, clusterName := range clusters {
		// Get cluster info
		info, err := c.GetClusterInfo(ctx, clusterName)
		if err != nil {
			continue
		}

		fmt.Printf("Cluster: %s\n", info.ClusterName)
		fmt.Printf("  Alias: %s\n", info.ClusterAlias)
		fmt.Printf("  Instances: %d\n", info.CountInstances)
		fmt.Printf("  Auto Recovery: %v\n", info.HasAutomatedMasterRecovery)

		// Get cluster master
		master, err := c.GetClusterMaster(ctx, clusterName)
		if err != nil {
			continue
		}

		fmt.Printf("  Master: %s\n", master.Key)
	}
}

// Example demonstrating monitoring and metrics
func Example_monitoring() {
	ctx := context.Background()
	cfg := client.Config{
		Host: "localhost:3000",
	}

	c := client.NewClient(cfg)

	// Get discovery metrics for last 60 seconds
	metrics, err := c.GetDiscoveryMetricsAggregated(ctx, 60)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Discovery metrics (last 60s): %d entries\n", len(metrics))

	for _, metric := range metrics {
		fmt.Printf("Instance: %s, Duration: %dms\n",
			metric.InstanceKey,
			metric.DurationMillis)
	}

	// Get backend query metrics
	queryMetrics, err := c.GetBackendQueryMetricsAggregated(ctx, 60)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Query metrics: %d entries\n", len(queryMetrics))
}

// Example demonstrating HTTPS with custom client
func Example_httpsWithCustomClient() {
	ctx := context.Background()
	cfg := client.Config{
		Host:     "orchestrator.example.com:443",
		Username: "admin",
		Password: "secret",
		UseHTTPS: true,
		Timeout:  60 * time.Second,
	}

	c := client.NewClient(cfg)

	// Health check
	err := c.Health(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Orchestrator is healthy")
}

// Example demonstrating Raft operations
func Example_raftOperations() {
	ctx := context.Background()
	cfg := client.Config{
		Host: "localhost:3000",
	}

	c := client.NewClient(cfg)

	// Get Raft state
	state, err := c.GetRaftState(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Raft Leader: %s\n", state.Leader)
	fmt.Printf("Is Leader: %v\n", state.IsLeader)

	// Get Raft peers
	peers, err := c.GetRaftPeers(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Raft Peers: %v\n", peers)

	// Get Raft health
	health, err := c.GetRaftHealth(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Raft Healthy: %v\n", health.Healthy)
	if !health.Healthy {
		fmt.Printf("Reason: %s\n", health.Reason)
	}
}

// Example demonstrating agent operations
func Example_agentOperations() {
	ctx := context.Background()
	cfg := client.Config{
		Host: "localhost:3000",
	}

	c := client.NewClient(cfg)

	// Get all agents
	agents, err := c.GetAgents(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Found %d agents\n", len(agents))

	for _, agent := range agents {
		fmt.Printf("Agent: %s:%d\n", agent.Hostname, agent.Port)
		fmt.Printf("  Available snapshots: %d\n", len(agent.AvailableLocalSnapshots))
		fmt.Printf("  Disk space ratio: %.2f\n", agent.AvailableDiskSpaceRatio)
	}

	// Get specific agent
	agent, err := c.GetAgent(ctx, "db1.example.com")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Agent details for %s\n", agent.Hostname)

	// Seed from source to target
	seed, err := c.AgentSeed(ctx, "target.example.com", "source.example.com")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Seed initiated: ID=%d\n", seed.SeedId)

	// Monitor seed progress
	activeSeeds, err := c.GetAgentActiveSeeds(ctx, "target.example.com")
	if err != nil {
		log.Fatal(err)
	}

	for _, activeSeed := range activeSeeds {
		fmt.Printf("Active seed %d: %s -> %s, Complete: %v\n",
			activeSeed.SeedId,
			activeSeed.SourceHostname,
			activeSeed.TargetHostname,
			activeSeed.IsComplete)
	}
}

// Example demonstrating downtime management
func Example_downtimeManagement() {
	ctx := context.Background()
	cfg := client.Config{
		Host: "localhost:3000",
	}

	c := client.NewClient(cfg)

	instanceKey := client.InstanceKey{
		Hostname: "db1.example.com",
		Port:     3306,
	}

	// Begin downtime with 2 hour duration
	instance, err := c.BeginDowntimeWithDuration(ctx,
		instanceKey,
		"sre-team",
		"Scheduled hardware upgrade",
		2*time.Hour,
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Downtime started for %s\n", instanceKey)
	fmt.Printf("Owner: %s\n", instance.DowntimeOwner)
	fmt.Printf("Reason: %s\n", instance.DowntimeReason)

	// Get all downtimed instances
	downtimed, err := c.GetDowntimedInstances(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Total downtimed instances: %d\n", len(downtimed))

	// End downtime when done
	_, err = c.EndDowntime(ctx, instanceKey)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Downtime ended for %s\n", instanceKey)
}

// Example demonstrating replication control
func Example_replicationControl() {
	ctx := context.Background()
	cfg := client.Config{
		Host: "localhost:3000",
	}

	c := client.NewClient(cfg)

	replicaKey := client.InstanceKey{
		Hostname: "replica1.example.com",
		Port:     3306,
	}

	// Stop replication
	instance, err := c.StopReplica(ctx, replicaKey)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Stopped replication on %s\n", replicaKey)
	fmt.Printf("SQL Thread Running: %v\n", instance.Slave_SQL_Running)
	fmt.Printf("IO Thread Running: %v\n", instance.Slave_IO_Running)

	// Add replication delay (for testing)
	_, err = c.DelayReplication(ctx, replicaKey, 60) // 60 second delay
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Added 60 second delay to %s\n", replicaKey)

	// Start replication
	instance, err = c.StartReplica(ctx, replicaKey)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Started replication on %s\n", replicaKey)

	// Enable semi-sync replication
	_, err = c.EnableSemiSyncReplica(ctx, replicaKey)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Enabled semi-sync replication on %s\n", replicaKey)
}

// Example demonstrating GTID operations
func Example_gtidOperations() {
	ctx := context.Background()
	cfg := client.Config{
		Host: "localhost:3000",
	}

	c := client.NewClient(cfg)

	instanceKey := client.InstanceKey{
		Hostname: "db1.example.com",
		Port:     3306,
	}

	// Enable GTID mode
	instance, err := c.EnableGTID(ctx, instanceKey)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("GTID enabled on %s\n", instanceKey)
	fmt.Printf("GTID Mode: %s\n", instance.GTIDMode)
	fmt.Printf("Using Oracle GTID: %v\n", instance.UsingOracleGTID)

	// Check for errant GTIDs
	instance, err = c.LocateErrantGTID(ctx, instanceKey)
	if err != nil {
		log.Fatal(err)
	}

	if instance.GtidErrant != "" {
		fmt.Printf("Errant GTIDs found: %s\n", instance.GtidErrant)
		
		// Inject empty transactions to resolve
		_, err = c.GTIDErrantInjectEmpty(ctx, instanceKey)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Injected empty transactions to resolve errant GTIDs")
	} else {
		fmt.Println("No errant GTIDs found")
	}
}

// Example demonstrating audit trail
func Example_auditTrail() {
	ctx := context.Background()
	cfg := client.Config{
		Host: "localhost:3000",
	}

	c := client.NewClient(cfg)

	// Get recent audit entries (page 0)
	entries, err := c.GetAudit(ctx, 0)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Recent audit entries: %d\n", len(entries))

	for _, entry := range entries {
		fmt.Printf("[%s] %s: %s - %s\n",
			entry.AuditTimestamp,
			entry.AuditType,
			entry.AuditInstanceKey,
			entry.Message)
	}

	// Get audit for specific instance
	instanceKey := client.InstanceKey{
		Hostname: "db1.example.com",
		Port:     3306,
	}

	instanceEntries, err := c.GetAuditForInstance(ctx, instanceKey, 0)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Audit entries for %s: %d\n", instanceKey, len(instanceEntries))
}
