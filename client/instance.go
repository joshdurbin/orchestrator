package client

import (
	"context"
	"fmt"
)

// Instance Management Operations

// GetInstance retrieves information about a specific instance
func (c *Client) GetInstance(ctx context.Context, key InstanceKey) (*Instance, error) {
	path := fmt.Sprintf("/api/instance/%s/%d", key.Hostname, key.Port)
	var response APIResponse
	if err := c.getJSON(ctx, path, &response); err != nil {
		return nil, err
	}

	if response.Code != OK {
		return nil, fmt.Errorf("API error: %s", response.Message)
	}

	// The Details field contains the Instance
	if response.Details == nil {
		return nil, fmt.Errorf("no instance data in response")
	}

	// Convert the Details interface{} to Instance
	instance, err := convertToInstance(response.Details)
	if err != nil {
		return nil, fmt.Errorf("failed to convert response: %w", err)
	}

	return instance, nil
}

// GetInstanceReplicas retrieves the list of replicas for a given instance
func (c *Client) GetInstanceReplicas(ctx context.Context, key InstanceKey) ([]Instance, error) {
	path := fmt.Sprintf("/api/instance-replicas/%s/%d", key.Hostname, key.Port)
	var response APIResponse
	if err := c.getJSON(ctx, path, &response); err != nil {
		return nil, err
	}

	if response.Code != OK {
		return nil, fmt.Errorf("API error: %s", response.Message)
	}

	return convertToInstanceSlice(response.Details)
}

// DiscoverInstance synchronously discovers and returns an instance
func (c *Client) DiscoverInstance(ctx context.Context, key InstanceKey) (*Instance, error) {
	path := fmt.Sprintf("/api/discover/%s/%d", key.Hostname, key.Port)
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

// AsyncDiscoverInstance asynchronously discovers an instance
func (c *Client) AsyncDiscoverInstance(ctx context.Context, key InstanceKey) (*Instance, error) {
	path := fmt.Sprintf("/api/async-discover/%s/%d", key.Hostname, key.Port)
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

// RefreshInstance refreshes the information about an instance
func (c *Client) RefreshInstance(ctx context.Context, key InstanceKey) (*Instance, error) {
	path := fmt.Sprintf("/api/refresh/%s/%d", key.Hostname, key.Port)
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

// ForgetInstance removes an instance from orchestrator's tracking
func (c *Client) ForgetInstance(ctx context.Context, key InstanceKey) (*Instance, error) {
	path := fmt.Sprintf("/api/forget/%s/%d", key.Hostname, key.Port)
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

// ForgetCluster removes all instances in a cluster from orchestrator's tracking
func (c *Client) ForgetCluster(ctx context.Context, clusterHint string) ([]Instance, error) {
	path := fmt.Sprintf("/api/forget-cluster/%s", clusterHint)
	var response APIResponse
	if err := c.getJSON(ctx, path, &response); err != nil {
		return nil, err
	}

	if response.Code != OK {
		return nil, fmt.Errorf("API error: %s", response.Message)
	}

	return convertToInstanceSlice(response.Details)
}

// GetAllInstances retrieves all known instances
func (c *Client) GetAllInstances(ctx context.Context) ([]Instance, error) {
	path := "/api/all-instances"
	var response APIResponse
	if err := c.getJSON(ctx, path, &response); err != nil {
		return nil, err
	}

	if response.Code != OK {
		return nil, fmt.Errorf("API error: %s", response.Message)
	}

	return convertToInstanceSlice(response.Details)
}

// ResolveInstance resolves hostname/port for an instance
func (c *Client) ResolveInstance(ctx context.Context, key InstanceKey) (*Instance, error) {
	path := fmt.Sprintf("/api/resolve/%s/%d", key.Hostname, key.Port)
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

// SetReadOnly sets an instance to read-only mode
func (c *Client) SetReadOnly(ctx context.Context, key InstanceKey) (*Instance, error) {
	path := fmt.Sprintf("/api/set-read-only/%s/%d", key.Hostname, key.Port)
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

// SetWriteable sets an instance to writeable (read-write) mode
func (c *Client) SetWriteable(ctx context.Context, key InstanceKey) (*Instance, error) {
	path := fmt.Sprintf("/api/set-writeable/%s/%d", key.Hostname, key.Port)
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

// KillQuery kills a running query on an instance
func (c *Client) KillQuery(ctx context.Context, key InstanceKey, processID int64) (*Instance, error) {
	path := fmt.Sprintf("/api/kill-query/%s/%d/%d", key.Hostname, key.Port, processID)
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

// SearchInstances searches for instances matching a search string
func (c *Client) SearchInstances(ctx context.Context, searchString string) ([]Instance, error) {
	path := fmt.Sprintf("/api/search/%s", searchString)
	var response APIResponse
	if err := c.getJSON(ctx, path, &response); err != nil {
		return nil, err
	}

	if response.Code != OK {
		return nil, fmt.Errorf("API error: %s", response.Message)
	}

	return convertToInstanceSlice(response.Details)
}

// SearchAllInstances searches all instances (empty search)
func (c *Client) SearchAllInstances(ctx context.Context) ([]Instance, error) {
	path := "/api/search"
	var response APIResponse
	if err := c.getJSON(ctx, path, &response); err != nil {
		return nil, err
	}

	if response.Code != OK {
		return nil, fmt.Errorf("API error: %s", response.Message)
	}

	return convertToInstanceSlice(response.Details)
}

// BulkInstances performs a bulk retrieval of instances
func (c *Client) BulkInstances(ctx context.Context, instanceKeys []InstanceKey) ([]Instance, error) {
	path := "/api/bulk-instances"

	// POST the instance keys as JSON
	var response APIResponse
	if err := c.postJSON(ctx, path, instanceKeys, &response); err != nil {
		return nil, err
	}

	if response.Code != OK {
		return nil, fmt.Errorf("API error: %s", response.Message)
	}

	return convertToInstanceSlice(response.Details)
}

// GetProblems retrieves all instances with topology problems
func (c *Client) GetProblems(ctx context.Context) ([]Instance, error) {
	path := "/api/problems"
	var response APIResponse
	if err := c.getJSON(ctx, path, &response); err != nil {
		return nil, err
	}

	if response.Code != OK {
		return nil, fmt.Errorf("API error: %s", response.Message)
	}

	return convertToInstanceSlice(response.Details)
}

// GetClusterProblems retrieves instances with topology problems for a specific cluster
func (c *Client) GetClusterProblems(ctx context.Context, clusterName string) ([]Instance, error) {
	path := fmt.Sprintf("/api/problems/%s", clusterName)
	var response APIResponse
	if err := c.getJSON(ctx, path, &response); err != nil {
		return nil, err
	}

	if response.Code != OK {
		return nil, fmt.Errorf("API error: %s", response.Message)
	}

	return convertToInstanceSlice(response.Details)
}
