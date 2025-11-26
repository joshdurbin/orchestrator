package client

import (
	"context"
	"fmt"
)

// Replication Control Operations

// StartReplica starts replication on an instance
func (c *Client) StartReplica(ctx context.Context, instanceKey InstanceKey) (*Instance, error) {
	path := fmt.Sprintf("/api/start-replica/%s/%d",
		instanceKey.Hostname, instanceKey.Port)
	return c.executeTopologyOperation(ctx, path)
}

// RestartReplica restarts replication on an instance
func (c *Client) RestartReplica(ctx context.Context, instanceKey InstanceKey) (*Instance, error) {
	path := fmt.Sprintf("/api/restart-replica/%s/%d",
		instanceKey.Hostname, instanceKey.Port)
	return c.executeTopologyOperation(ctx, path)
}

// StopReplica stops replication on an instance
func (c *Client) StopReplica(ctx context.Context, instanceKey InstanceKey) (*Instance, error) {
	path := fmt.Sprintf("/api/stop-replica/%s/%d",
		instanceKey.Hostname, instanceKey.Port)
	return c.executeTopologyOperation(ctx, path)
}

// StopReplicaNice stops replication gracefully on an instance
func (c *Client) StopReplicaNice(ctx context.Context, instanceKey InstanceKey) (*Instance, error) {
	path := fmt.Sprintf("/api/stop-replica-nice/%s/%d",
		instanceKey.Hostname, instanceKey.Port)
	return c.executeTopologyOperation(ctx, path)
}

// ResetReplica resets replication on an instance
func (c *Client) ResetReplica(ctx context.Context, instanceKey InstanceKey) (*Instance, error) {
	path := fmt.Sprintf("/api/reset-replica/%s/%d",
		instanceKey.Hostname, instanceKey.Port)
	return c.executeTopologyOperation(ctx, path)
}

// DetachReplica detaches a replica from its master
func (c *Client) DetachReplica(ctx context.Context, instanceKey InstanceKey) (*Instance, error) {
	path := fmt.Sprintf("/api/detach-replica/%s/%d",
		instanceKey.Hostname, instanceKey.Port)
	return c.executeTopologyOperation(ctx, path)
}

// ReattachReplica reattaches a replica to its master
func (c *Client) ReattachReplica(ctx context.Context, instanceKey InstanceKey) (*Instance, error) {
	path := fmt.Sprintf("/api/reattach-replica/%s/%d",
		instanceKey.Hostname, instanceKey.Port)
	return c.executeTopologyOperation(ctx, path)
}

// DetachReplicaMasterHost detaches a replica from master (alternative endpoint)
func (c *Client) DetachReplicaMasterHost(ctx context.Context, instanceKey InstanceKey) (*Instance, error) {
	path := fmt.Sprintf("/api/detach-replica-master-host/%s/%d",
		instanceKey.Hostname, instanceKey.Port)
	return c.executeTopologyOperation(ctx, path)
}

// ReattachReplicaMasterHost reattaches a replica to master (alternative endpoint)
func (c *Client) ReattachReplicaMasterHost(ctx context.Context, instanceKey InstanceKey) (*Instance, error) {
	path := fmt.Sprintf("/api/reattach-replica-master-host/%s/%d",
		instanceKey.Hostname, instanceKey.Port)
	return c.executeTopologyOperation(ctx, path)
}

// SkipQuery skips a problematic query on a replica
func (c *Client) SkipQuery(ctx context.Context, instanceKey InstanceKey) (*Instance, error) {
	path := fmt.Sprintf("/api/skip-query/%s/%d",
		instanceKey.Hostname, instanceKey.Port)
	return c.executeTopologyOperation(ctx, path)
}

// FlushBinaryLogs flushes binary logs on an instance
func (c *Client) FlushBinaryLogs(ctx context.Context, instanceKey InstanceKey) (*Instance, error) {
	path := fmt.Sprintf("/api/flush-binary-logs/%s/%d",
		instanceKey.Hostname, instanceKey.Port)
	return c.executeTopologyOperation(ctx, path)
}

// PurgeBinaryLogs purges binary logs up to a specific log file
func (c *Client) PurgeBinaryLogs(ctx context.Context, instanceKey InstanceKey, logFile string) (*Instance, error) {
	path := fmt.Sprintf("/api/purge-binary-logs/%s/%d/%s",
		instanceKey.Hostname, instanceKey.Port, logFile)
	return c.executeTopologyOperation(ctx, path)
}

// RestartReplicaStatements returns the statements needed to restart replication
func (c *Client) RestartReplicaStatements(ctx context.Context, instanceKey InstanceKey) ([]string, error) {
	path := fmt.Sprintf("/api/restart-replica-statements/%s/%d",
		instanceKey.Hostname, instanceKey.Port)
	var response APIResponse
	if err := c.getJSON(ctx, path, &response); err != nil {
		return nil, err
	}

	if response.Code != OK {
		return nil, fmt.Errorf("API error: %s", response.Message)
	}

	return convertToStringSlice(response.Details)
}

// DelayReplication adds a replication delay to an instance
func (c *Client) DelayReplication(ctx context.Context, instanceKey InstanceKey, seconds int) (*Instance, error) {
	path := fmt.Sprintf("/api/delay-replication/%s/%d/%d",
		instanceKey.Hostname, instanceKey.Port, seconds)
	return c.executeTopologyOperation(ctx, path)
}

// Semi-Sync Replication Operations

// EnableSemiSyncMaster enables semi-sync replication on master
func (c *Client) EnableSemiSyncMaster(ctx context.Context, instanceKey InstanceKey) (*Instance, error) {
	path := fmt.Sprintf("/api/enable-semi-sync-master/%s/%d",
		instanceKey.Hostname, instanceKey.Port)
	return c.executeTopologyOperation(ctx, path)
}

// DisableSemiSyncMaster disables semi-sync replication on master
func (c *Client) DisableSemiSyncMaster(ctx context.Context, instanceKey InstanceKey) (*Instance, error) {
	path := fmt.Sprintf("/api/disable-semi-sync-master/%s/%d",
		instanceKey.Hostname, instanceKey.Port)
	return c.executeTopologyOperation(ctx, path)
}

// EnableSemiSyncReplica enables semi-sync replication on replica
func (c *Client) EnableSemiSyncReplica(ctx context.Context, instanceKey InstanceKey) (*Instance, error) {
	path := fmt.Sprintf("/api/enable-semi-sync-replica/%s/%d",
		instanceKey.Hostname, instanceKey.Port)
	return c.executeTopologyOperation(ctx, path)
}

// DisableSemiSyncReplica disables semi-sync replication on replica
func (c *Client) DisableSemiSyncReplica(ctx context.Context, instanceKey InstanceKey) (*Instance, error) {
	path := fmt.Sprintf("/api/disable-semi-sync-replica/%s/%d",
		instanceKey.Hostname, instanceKey.Port)
	return c.executeTopologyOperation(ctx, path)
}
