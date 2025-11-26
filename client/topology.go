package client

import (
	"context"
	"fmt"
)

// Topology Management Operations

// Smart Relocation Operations

// Relocate relocates an instance below another instance
func (c *Client) Relocate(ctx context.Context, instanceKey, belowKey InstanceKey) (*Instance, error) {
	path := fmt.Sprintf("/api/relocate/%s/%d/%s/%d",
		instanceKey.Hostname, instanceKey.Port,
		belowKey.Hostname, belowKey.Port)
	return c.executeTopologyOperation(ctx, path)
}

// RelocateBelow explicitly relocates an instance below another
func (c *Client) RelocateBelow(ctx context.Context, instanceKey, belowKey InstanceKey) (*Instance, error) {
	path := fmt.Sprintf("/api/relocate-below/%s/%d/%s/%d",
		instanceKey.Hostname, instanceKey.Port,
		belowKey.Hostname, belowKey.Port)
	return c.executeTopologyOperation(ctx, path)
}

// RelocateReplicas relocates all replicas of an instance below another instance
func (c *Client) RelocateReplicas(ctx context.Context, instanceKey, belowKey InstanceKey) ([]Instance, error) {
	path := fmt.Sprintf("/api/relocate-replicas/%s/%d/%s/%d",
		instanceKey.Hostname, instanceKey.Port,
		belowKey.Hostname, belowKey.Port)
	return c.executeTopologyOperationMulti(ctx, path)
}

// RegroupReplicas regroups replicas of an instance
func (c *Client) RegroupReplicas(ctx context.Context, instanceKey InstanceKey) ([]Instance, error) {
	path := fmt.Sprintf("/api/regroup-replicas/%s/%d",
		instanceKey.Hostname, instanceKey.Port)
	return c.executeTopologyOperationMulti(ctx, path)
}

// File Position (Binlog) Based Relocation

// MoveUp moves an instance one level up in the topology
func (c *Client) MoveUp(ctx context.Context, instanceKey InstanceKey) (*Instance, error) {
	path := fmt.Sprintf("/api/move-up/%s/%d",
		instanceKey.Hostname, instanceKey.Port)
	return c.executeTopologyOperation(ctx, path)
}

// MoveUpReplicas moves all replicas of an instance one level up
func (c *Client) MoveUpReplicas(ctx context.Context, instanceKey InstanceKey) ([]Instance, error) {
	path := fmt.Sprintf("/api/move-up-replicas/%s/%d",
		instanceKey.Hostname, instanceKey.Port)
	return c.executeTopologyOperationMulti(ctx, path)
}

// MoveBelow moves an instance below a sibling
func (c *Client) MoveBelow(ctx context.Context, instanceKey, siblingKey InstanceKey) (*Instance, error) {
	path := fmt.Sprintf("/api/move-below/%s/%d/%s/%d",
		instanceKey.Hostname, instanceKey.Port,
		siblingKey.Hostname, siblingKey.Port)
	return c.executeTopologyOperation(ctx, path)
}

// MoveEquivalent moves an instance to an equivalent position below another instance
func (c *Client) MoveEquivalent(ctx context.Context, instanceKey, belowKey InstanceKey) (*Instance, error) {
	path := fmt.Sprintf("/api/move-equivalent/%s/%d/%s/%d",
		instanceKey.Hostname, instanceKey.Port,
		belowKey.Hostname, belowKey.Port)
	return c.executeTopologyOperation(ctx, path)
}

// Repoint repoints an instance to replicate from a different master
func (c *Client) Repoint(ctx context.Context, instanceKey, belowKey InstanceKey) (*Instance, error) {
	path := fmt.Sprintf("/api/repoint/%s/%d/%s/%d",
		instanceKey.Hostname, instanceKey.Port,
		belowKey.Hostname, belowKey.Port)
	return c.executeTopologyOperation(ctx, path)
}

// RepointReplicas repoints all replicas of an instance
func (c *Client) RepointReplicas(ctx context.Context, instanceKey InstanceKey) ([]Instance, error) {
	path := fmt.Sprintf("/api/repoint-replicas/%s/%d",
		instanceKey.Hostname, instanceKey.Port)
	return c.executeTopologyOperationMulti(ctx, path)
}

// MakeCoMaster promotes an instance to co-master
func (c *Client) MakeCoMaster(ctx context.Context, instanceKey InstanceKey) (*Instance, error) {
	path := fmt.Sprintf("/api/make-co-master/%s/%d",
		instanceKey.Hostname, instanceKey.Port)
	return c.executeTopologyOperation(ctx, path)
}

// TakeSiblings makes an instance the master of its siblings
func (c *Client) TakeSiblings(ctx context.Context, instanceKey InstanceKey) (*Instance, error) {
	path := fmt.Sprintf("/api/take-siblings/%s/%d",
		instanceKey.Hostname, instanceKey.Port)
	return c.executeTopologyOperation(ctx, path)
}

// TakeMaster makes an instance the master of its current master
func (c *Client) TakeMaster(ctx context.Context, instanceKey InstanceKey) (*Instance, error) {
	path := fmt.Sprintf("/api/take-master/%s/%d",
		instanceKey.Hostname, instanceKey.Port)
	return c.executeTopologyOperation(ctx, path)
}

// MasterEquivalent finds the equivalent master position
func (c *Client) MasterEquivalent(ctx context.Context, instanceKey InstanceKey, logFile string, logPos int64) (*Instance, error) {
	path := fmt.Sprintf("/api/master-equivalent/%s/%d/%s/%d",
		instanceKey.Hostname, instanceKey.Port, logFile, logPos)
	return c.executeTopologyOperation(ctx, path)
}

// Binlog Server Relocation

// RegroupReplicasBinlogServers regroups replicas using binlog servers
func (c *Client) RegroupReplicasBinlogServers(ctx context.Context, instanceKey InstanceKey) ([]Instance, error) {
	path := fmt.Sprintf("/api/regroup-replicas-bls/%s/%d",
		instanceKey.Hostname, instanceKey.Port)
	return c.executeTopologyOperationMulti(ctx, path)
}

// GTID-Based Relocation

// MoveBelowGTID moves an instance below another using GTID
func (c *Client) MoveBelowGTID(ctx context.Context, instanceKey, belowKey InstanceKey) (*Instance, error) {
	path := fmt.Sprintf("/api/move-below-gtid/%s/%d/%s/%d",
		instanceKey.Hostname, instanceKey.Port,
		belowKey.Hostname, belowKey.Port)
	return c.executeTopologyOperation(ctx, path)
}

// MoveReplicasGTID moves all replicas below another instance using GTID
func (c *Client) MoveReplicasGTID(ctx context.Context, instanceKey, belowKey InstanceKey) ([]Instance, error) {
	path := fmt.Sprintf("/api/move-replicas-gtid/%s/%d/%s/%d",
		instanceKey.Hostname, instanceKey.Port,
		belowKey.Hostname, belowKey.Port)
	return c.executeTopologyOperationMulti(ctx, path)
}

// RegroupReplicasGTID regroups replicas using GTID
func (c *Client) RegroupReplicasGTID(ctx context.Context, instanceKey InstanceKey) ([]Instance, error) {
	path := fmt.Sprintf("/api/regroup-replicas-gtid/%s/%d",
		instanceKey.Hostname, instanceKey.Port)
	return c.executeTopologyOperationMulti(ctx, path)
}

// Pseudo-GTID Based Relocation

// Match matches an instance below another using Pseudo-GTID
func (c *Client) Match(ctx context.Context, instanceKey, belowKey InstanceKey) (*Instance, error) {
	path := fmt.Sprintf("/api/match/%s/%d/%s/%d",
		instanceKey.Hostname, instanceKey.Port,
		belowKey.Hostname, belowKey.Port)
	return c.executeTopologyOperation(ctx, path)
}

// MatchBelow explicitly matches an instance below another using Pseudo-GTID
func (c *Client) MatchBelow(ctx context.Context, instanceKey, belowKey InstanceKey) (*Instance, error) {
	path := fmt.Sprintf("/api/match-below/%s/%d/%s/%d",
		instanceKey.Hostname, instanceKey.Port,
		belowKey.Hostname, belowKey.Port)
	return c.executeTopologyOperation(ctx, path)
}

// MatchUp matches an instance one level up using Pseudo-GTID
func (c *Client) MatchUp(ctx context.Context, instanceKey InstanceKey) (*Instance, error) {
	path := fmt.Sprintf("/api/match-up/%s/%d",
		instanceKey.Hostname, instanceKey.Port)
	return c.executeTopologyOperation(ctx, path)
}

// MatchReplicas multi-matches replicas below another instance using Pseudo-GTID
func (c *Client) MatchReplicas(ctx context.Context, instanceKey, belowKey InstanceKey) ([]Instance, error) {
	path := fmt.Sprintf("/api/match-replicas/%s/%d/%s/%d",
		instanceKey.Hostname, instanceKey.Port,
		belowKey.Hostname, belowKey.Port)
	return c.executeTopologyOperationMulti(ctx, path)
}

// MatchUpReplicas matches all replicas one level up using Pseudo-GTID
func (c *Client) MatchUpReplicas(ctx context.Context, instanceKey InstanceKey) ([]Instance, error) {
	path := fmt.Sprintf("/api/match-up-replicas/%s/%d",
		instanceKey.Hostname, instanceKey.Port)
	return c.executeTopologyOperationMulti(ctx, path)
}

// RegroupReplicasPseudoGTID regroups replicas using Pseudo-GTID
func (c *Client) RegroupReplicasPseudoGTID(ctx context.Context, instanceKey InstanceKey) ([]Instance, error) {
	path := fmt.Sprintf("/api/regroup-replicas-pgtid/%s/%d",
		instanceKey.Hostname, instanceKey.Port)
	return c.executeTopologyOperationMulti(ctx, path)
}

// Replication Analysis

// CanReplicateFrom checks if an instance can replicate from another
func (c *Client) CanReplicateFrom(ctx context.Context, instanceKey, fromKey InstanceKey) (*Instance, error) {
	path := fmt.Sprintf("/api/can-replicate-from/%s/%d/%s/%d",
		instanceKey.Hostname, instanceKey.Port,
		fromKey.Hostname, fromKey.Port)
	return c.executeTopologyOperation(ctx, path)
}

// CanReplicateFromGTID checks if an instance can replicate from another using GTID
func (c *Client) CanReplicateFromGTID(ctx context.Context, instanceKey, fromKey InstanceKey) (*Instance, error) {
	path := fmt.Sprintf("/api/can-replicate-from-gtid/%s/%d/%s/%d",
		instanceKey.Hostname, instanceKey.Port,
		fromKey.Hostname, fromKey.Port)
	return c.executeTopologyOperation(ctx, path)
}

// LastPseudoGTID retrieves the last Pseudo-GTID entry for an instance
func (c *Client) LastPseudoGTID(ctx context.Context, instanceKey InstanceKey) (*Instance, error) {
	path := fmt.Sprintf("/api/last-pseudo-gtid/%s/%d",
		instanceKey.Hostname, instanceKey.Port)
	return c.executeTopologyOperation(ctx, path)
}

// Topology Views

// GetTopologyASCII retrieves an ASCII representation of cluster topology
func (c *Client) GetTopologyASCII(ctx context.Context, clusterHint string) (string, error) {
	path := fmt.Sprintf("/api/topology/%s", clusterHint)
	return c.getPlainText(ctx, path)
}

// GetTopologyASCIIFromInstance retrieves an ASCII topology starting from a specific instance
func (c *Client) GetTopologyASCIIFromInstance(ctx context.Context, instanceKey InstanceKey) (string, error) {
	path := fmt.Sprintf("/api/topology/%s/%d",
		instanceKey.Hostname, instanceKey.Port)
	return c.getPlainText(ctx, path)
}

// GetTopologyTabulated retrieves a tabulated representation of cluster topology
func (c *Client) GetTopologyTabulated(ctx context.Context, clusterHint string) (string, error) {
	path := fmt.Sprintf("/api/topology-tabulated/%s", clusterHint)
	return c.getPlainText(ctx, path)
}

// GetTopologyTabulatedFromInstance retrieves a tabulated topology from a specific instance
func (c *Client) GetTopologyTabulatedFromInstance(ctx context.Context, instanceKey InstanceKey) (string, error) {
	path := fmt.Sprintf("/api/topology-tabulated/%s/%d",
		instanceKey.Hostname, instanceKey.Port)
	return c.getPlainText(ctx, path)
}

// GetTopologyTags retrieves topology with tags for a cluster
func (c *Client) GetTopologyTags(ctx context.Context, clusterHint string) ([]Instance, error) {
	path := fmt.Sprintf("/api/topology-tags/%s", clusterHint)
	return c.executeTopologyOperationMulti(ctx, path)
}

// GetTopologyTagsFromInstance retrieves topology with tags from a specific instance
func (c *Client) GetTopologyTagsFromInstance(ctx context.Context, instanceKey InstanceKey) ([]Instance, error) {
	path := fmt.Sprintf("/api/topology-tags/%s/%d",
		instanceKey.Hostname, instanceKey.Port)
	return c.executeTopologyOperationMulti(ctx, path)
}

// SnapshotTopologies takes a snapshot of all topologies
func (c *Client) SnapshotTopologies(ctx context.Context, ) ([]InstanceKey, error) {
	path := "/api/snapshot-topologies"
	var response APIResponse
	if err := c.getJSON(ctx, path, &response); err != nil {
		return nil, err
	}

	if response.Code != OK {
		return nil, fmt.Errorf("API error: %s", response.Message)
	}

	// Convert to instance keys
	jsonData, err := convertToStringSlice(response.Details)
	if err != nil {
		return nil, err
	}

	keys := make([]InstanceKey, 0, len(jsonData))
	for _, keyStr := range jsonData {
		key, err := ParseInstanceKey(keyStr)
		if err != nil {
			continue
		}
		keys = append(keys, key)
	}

	return keys, nil
}

// Helper functions

// executeTopologyOperation executes a topology operation that returns a single instance
func (c *Client) executeTopologyOperation(ctx context.Context, path string) (*Instance, error) {
	var response APIResponse
	if err := c.getJSON(ctx, path, &response); err != nil {
		return nil, err
	}

	if response.Code != OK {
		return nil, fmt.Errorf("API error: %s", response.Message)
	}

	instance, err := convertToInstance(response.Details)
	if err != nil {
		return nil, fmt.Errorf("failed to convert response: %w", err)
	}

	return instance, nil
}

// executeTopologyOperationMulti executes a topology operation that returns multiple instances
func (c *Client) executeTopologyOperationMulti(ctx context.Context, path string) ([]Instance, error) {
	var response APIResponse
	if err := c.getJSON(ctx, path, &response); err != nil {
		return nil, err
	}

	if response.Code != OK {
		return nil, fmt.Errorf("API error: %s", response.Message)
	}

	return convertToInstanceSlice(response.Details)
}
