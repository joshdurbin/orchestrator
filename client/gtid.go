package client

import (
	"context"
	"fmt"
)

// GTID Management Operations

// EnableGTID enables GTID replication on an instance
func (c *Client) EnableGTID(ctx context.Context, instanceKey InstanceKey) (*Instance, error) {
	path := fmt.Sprintf("/api/enable-gtid/%s/%d",
		instanceKey.Hostname, instanceKey.Port)
	return c.executeTopologyOperation(ctx, path)
}

// DisableGTID disables GTID replication on an instance
func (c *Client) DisableGTID(ctx context.Context, instanceKey InstanceKey) (*Instance, error) {
	path := fmt.Sprintf("/api/disable-gtid/%s/%d",
		instanceKey.Hostname, instanceKey.Port)
	return c.executeTopologyOperation(ctx, path)
}

// LocateErrantGTID locates errant GTIDs on an instance
func (c *Client) LocateErrantGTID(ctx context.Context, instanceKey InstanceKey) (*Instance, error) {
	path := fmt.Sprintf("/api/locate-gtid-errant/%s/%d",
		instanceKey.Hostname, instanceKey.Port)
	return c.executeTopologyOperation(ctx, path)
}

// GTIDErrantResetMaster resets master for errant GTIDs
func (c *Client) GTIDErrantResetMaster(ctx context.Context, instanceKey InstanceKey) (*Instance, error) {
	path := fmt.Sprintf("/api/gtid-errant-reset-master/%s/%d",
		instanceKey.Hostname, instanceKey.Port)
	return c.executeTopologyOperation(ctx, path)
}

// GTIDErrantInjectEmpty injects empty transaction for errant GTID
func (c *Client) GTIDErrantInjectEmpty(ctx context.Context, instanceKey InstanceKey) (*Instance, error) {
	path := fmt.Sprintf("/api/gtid-errant-inject-empty/%s/%d",
		instanceKey.Hostname, instanceKey.Port)
	return c.executeTopologyOperation(ctx, path)
}
