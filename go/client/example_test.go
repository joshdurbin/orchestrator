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

package client_test

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/openark/orchestrator/go/client"
	"github.com/openark/orchestrator/go/inst"
)

func ExampleClient_basic() {
	// Create a client with single endpoint
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

	fmt.Printf("Discovered instance: %s:%d\n", instance.Key.Hostname, instance.Key.Port)
}

func ExampleClient_multipleEndpoints() {
	// Create a client with multiple endpoints for automatic leader detection
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

	ctx := context.Background()

	// Get all clusters
	clusters, err := c.GetClusters(ctx)
	if err != nil {
		log.Fatal("Failed to get clusters:", err)
	}

	for _, cluster := range clusters {
		fmt.Printf("Cluster: %s\n", cluster)
	}
}

func ExampleClient_topologyOperations() {
	config := &client.Config{
		BaseURL: "http://orchestrator.example.com:3000",
	}

	c, err := client.NewClient(config)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	// Get cluster instances
	instances, err := c.GetCluster(ctx, "main-cluster")
	if err != nil {
		log.Fatal("Failed to get cluster:", err)
	}

	fmt.Printf("Found %d instances in cluster\n", len(instances))

	// Get topology visualization
	topology, err := c.GetTopology(ctx, "main-cluster")
	if err != nil {
		log.Fatal("Failed to get topology:", err)
	}

	fmt.Println("Topology:")
	fmt.Println(topology)
}

func ExampleClient_failover() {
	config := &client.Config{
		BaseURL:  "http://orchestrator.example.com:3000",
		Username: "admin",
		Password: "secret",
	}

	c, err := client.NewClient(config)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	// Perform graceful master takeover
	newMaster, err := c.GracefulMasterTakeover(ctx, "main-cluster", nil)
	if err != nil {
		log.Fatal("Failover failed:", err)
	}

	fmt.Printf("New master: %s:%d\n", newMaster.Key.Hostname, newMaster.Key.Port)

	// Check replication analysis
	analysis, err := c.GetReplicationAnalysis(ctx)
	if err != nil {
		log.Fatal("Failed to get analysis:", err)
	}

	fmt.Printf("Found %d analysis items\n", len(analysis))
}

func ExampleClient_maintenance() {
	config := &client.Config{
		BaseURL: "http://orchestrator.example.com:3000",
	}

	c, err := client.NewClient(config)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	instanceKey := inst.InstanceKey{
		Hostname: "mysql1.example.com",
		Port:     3306,
	}

	// Begin maintenance
	err = c.BeginMaintenance(ctx, instanceKey, "admin", "Scheduled maintenance")
	if err != nil {
		log.Fatal("Failed to begin maintenance:", err)
	}

	fmt.Println("Maintenance mode activated")

	// Check if in maintenance
	inMaintenance, err := c.InMaintenance(ctx, instanceKey)
	if err != nil {
		log.Fatal("Failed to check maintenance status:", err)
	}

	fmt.Printf("In maintenance: %v\n", inMaintenance)

	// End maintenance
	err = c.EndMaintenance(ctx, instanceKey)
	if err != nil {
		log.Fatal("Failed to end maintenance:", err)
	}

	fmt.Println("Maintenance mode deactivated")
}

func ExampleClient_tags() {
	config := &client.Config{
		BaseURL: "http://orchestrator.example.com:3000",
	}

	c, err := client.NewClient(config)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	instanceKey := inst.InstanceKey{
		Hostname: "mysql1.example.com",
		Port:     3306,
	}

	// Set tags
	err = c.SetTag(ctx, instanceKey, "environment", "production")
	if err != nil {
		log.Fatal("Failed to set tag:", err)
	}

	err = c.SetTag(ctx, instanceKey, "role", "master")
	if err != nil {
		log.Fatal("Failed to set tag:", err)
	}

	// Get tags
	tags, err := c.GetTags(ctx, instanceKey)
	if err != nil {
		log.Fatal("Failed to get tags:", err)
	}

	fmt.Printf("Tags: %v\n", tags)

	// Get tagged instances
	prodInstances, err := c.GetTagged(ctx, "environment", "production")
	if err != nil {
		log.Fatal("Failed to get tagged instances:", err)
	}

	fmt.Printf("Found %d production instances\n", len(prodInstances))
}

func ExampleClient_replication() {
	config := &client.Config{
		BaseURL: "http://orchestrator.example.com:3000",
	}

	c, err := client.NewClient(config)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	instanceKey := inst.InstanceKey{
		Hostname: "mysql-replica1.example.com",
		Port:     3306,
	}

	// Stop replication
	instance, err := c.StopReplication(ctx, instanceKey)
	if err != nil {
		log.Fatal("Failed to stop replication:", err)
	}

	fmt.Printf("Stopped replication on %s:%d\n", instance.Key.Hostname, instance.Key.Port)

	// Start replication
	instance, err = c.StartReplication(ctx, instanceKey)
	if err != nil {
		log.Fatal("Failed to start replication:", err)
	}

	fmt.Printf("Started replication on %s:%d\n", instance.Key.Hostname, instance.Key.Port)
}

func ExampleClient_topologyModification() {
	config := &client.Config{
		BaseURL: "http://orchestrator.example.com:3000",
	}

	c, err := client.NewClient(config)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	replica := inst.InstanceKey{
		Hostname: "mysql-replica1.example.com",
		Port:     3306,
	}

	newMaster := inst.InstanceKey{
		Hostname: "mysql-master2.example.com",
		Port:     3306,
	}

	// Relocate replica under a new master
	instance, err := c.Relocate(ctx, replica, newMaster)
	if err != nil {
		log.Fatal("Failed to relocate:", err)
	}

	fmt.Printf("Relocated %s:%d under %s:%d\n",
		replica.Hostname, replica.Port, newMaster.Hostname, newMaster.Port)

	// Move replica up one level
	instance, err = c.MoveUp(ctx, replica)
	if err != nil {
		log.Fatal("Failed to move up:", err)
	}

	fmt.Printf("Moved up %s:%d\n", instance.Key.Hostname, instance.Key.Port)
}
