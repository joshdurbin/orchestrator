package client

import (
	"context"
	"fmt"
	"time"
)

// Maintenance and Downtime Operations

// Maintenance Windows

// BeginMaintenance begins a maintenance window for an instance
func (c *Client) BeginMaintenance(ctx context.Context, instanceKey InstanceKey, owner, reason string) (*Maintenance, error) {
	path := fmt.Sprintf("/api/begin-maintenance/%s/%d/%s/%s",
		instanceKey.Hostname, instanceKey.Port, owner, reason)
	return c.executeMaintenanceOperation(ctx, path)
}

// BeginMaintenanceWithDuration begins a maintenance window with a specific duration
func (c *Client) BeginMaintenanceWithDuration(ctx context.Context, instanceKey InstanceKey, owner, reason string, duration time.Duration) (*Maintenance, error) {
	durationStr := fmt.Sprintf("%ds", int(duration.Seconds()))
	path := fmt.Sprintf("/api/begin-maintenance/%s/%d/%s/%s/%s",
		instanceKey.Hostname, instanceKey.Port, owner, reason, durationStr)
	return c.executeMaintenanceOperation(ctx, path)
}

// EndMaintenanceByInstance ends maintenance for a specific instance
func (c *Client) EndMaintenanceByInstance(ctx context.Context, instanceKey InstanceKey) (*Maintenance, error) {
	path := fmt.Sprintf("/api/end-maintenance/%s/%d",
		instanceKey.Hostname, instanceKey.Port)
	return c.executeMaintenanceOperation(ctx, path)
}

// EndMaintenanceByKey ends maintenance by maintenance key
func (c *Client) EndMaintenanceByKey(ctx context.Context, maintenanceKey uint) (*Maintenance, error) {
	path := fmt.Sprintf("/api/end-maintenance/%d", maintenanceKey)
	return c.executeMaintenanceOperation(ctx, path)
}

// InMaintenance checks if an instance is in maintenance
func (c *Client) InMaintenance(ctx context.Context, instanceKey InstanceKey) (bool, error) {
	path := fmt.Sprintf("/api/in-maintenance/%s/%d",
		instanceKey.Hostname, instanceKey.Port)
	var response APIResponse
	if err := c.getJSON(ctx, path, &response); err != nil {
		return false, err
	}

	if response.Code != OK {
		return false, fmt.Errorf("API error: %s", response.Message)
	}

	return convertToBool(response.Details)
}

// GetMaintenance retrieves all active maintenance windows
func (c *Client) GetMaintenance(ctx context.Context, ) ([]Maintenance, error) {
	path := "/api/maintenance"
	var response APIResponse
	if err := c.getJSON(ctx, path, &response); err != nil {
		return nil, err
	}

	if response.Code != OK {
		return nil, fmt.Errorf("API error: %s", response.Message)
	}

	return convertToMaintenanceSlice(response.Details)
}

// Downtime Windows

// BeginDowntime begins a downtime window for an instance
func (c *Client) BeginDowntime(ctx context.Context, instanceKey InstanceKey, owner, reason string) (*Instance, error) {
	path := fmt.Sprintf("/api/begin-downtime/%s/%d/%s/%s",
		instanceKey.Hostname, instanceKey.Port, owner, reason)
	return c.executeTopologyOperation(ctx, path)
}

// BeginDowntimeWithDuration begins a downtime window with a specific duration
func (c *Client) BeginDowntimeWithDuration(ctx context.Context, instanceKey InstanceKey, owner, reason string, duration time.Duration) (*Instance, error) {
	durationStr := fmt.Sprintf("%ds", int(duration.Seconds()))
	path := fmt.Sprintf("/api/begin-downtime/%s/%d/%s/%s/%s",
		instanceKey.Hostname, instanceKey.Port, owner, reason, durationStr)
	return c.executeTopologyOperation(ctx, path)
}

// EndDowntime ends a downtime window for an instance
func (c *Client) EndDowntime(ctx context.Context, instanceKey InstanceKey) (*Instance, error) {
	path := fmt.Sprintf("/api/end-downtime/%s/%d",
		instanceKey.Hostname, instanceKey.Port)
	return c.executeTopologyOperation(ctx, path)
}

// GetDowntimedInstances retrieves all downtimed instances
func (c *Client) GetDowntimedInstances(ctx context.Context, ) ([]Instance, error) {
	path := "/api/downtimed"
	return c.executeTopologyOperationMulti(ctx, path)
}

// GetDowntimedInstancesByCluster retrieves downtimed instances for a specific cluster
func (c *Client) GetDowntimedInstancesByCluster(ctx context.Context, clusterHint string) ([]Instance, error) {
	path := fmt.Sprintf("/api/downtimed/%s", clusterHint)
	return c.executeTopologyOperationMulti(ctx, path)
}

// Helper functions

// executeMaintenanceOperation executes a maintenance operation
func (c *Client) executeMaintenanceOperation(ctx context.Context, path string) (*Maintenance, error) {
	var response APIResponse
	if err := c.getJSON(ctx, path, &response); err != nil {
		return nil, err
	}

	if response.Code != OK {
		return nil, fmt.Errorf("API error: %s", response.Message)
	}

	return convertToMaintenance(response.Details)
}
