package client

import (
	"context"
	"fmt"
)

// Cluster Operations

// GetCluster retrieves the cluster topology
func (c *Client) GetCluster(ctx context.Context, clusterHint string) ([]Instance, error) {
	path := fmt.Sprintf("/api/cluster/%s", clusterHint)
	return c.executeTopologyOperationMulti(ctx, path)
}

// GetClusterByAlias retrieves cluster by alias
func (c *Client) GetClusterByAlias(ctx context.Context, clusterAlias string) ([]Instance, error) {
	path := fmt.Sprintf("/api/cluster/alias/%s", clusterAlias)
	return c.executeTopologyOperationMulti(ctx, path)
}

// GetClusterByInstance retrieves the cluster of a specific instance
func (c *Client) GetClusterByInstance(ctx context.Context, instanceKey InstanceKey) ([]Instance, error) {
	path := fmt.Sprintf("/api/cluster/instance/%s/%d",
		instanceKey.Hostname, instanceKey.Port)
	return c.executeTopologyOperationMulti(ctx, path)
}

// GetClusterInfo retrieves information about a cluster
func (c *Client) GetClusterInfo(ctx context.Context, clusterHint string) (*ClusterInfo, error) {
	path := fmt.Sprintf("/api/cluster-info/%s", clusterHint)
	var response APIResponse
	if err := c.getJSON(ctx, path, &response); err != nil {
		return nil, err
	}

	if response.Code != OK {
		return nil, fmt.Errorf("API error: %s", response.Message)
	}

	return convertToClusterInfo(response.Details)
}

// GetClusterInfoByAlias retrieves cluster info by alias
func (c *Client) GetClusterInfoByAlias(ctx context.Context, clusterAlias string) (*ClusterInfo, error) {
	path := fmt.Sprintf("/api/cluster-info/alias/%s", clusterAlias)
	var response APIResponse
	if err := c.getJSON(ctx, path, &response); err != nil {
		return nil, err
	}

	if response.Code != OK {
		return nil, fmt.Errorf("API error: %s", response.Message)
	}

	return convertToClusterInfo(response.Details)
}

// GetClusterOSCReplicas retrieves OSC replicas for a cluster
func (c *Client) GetClusterOSCReplicas(ctx context.Context, clusterHint string) ([]Instance, error) {
	path := fmt.Sprintf("/api/cluster-osc-replicas/%s", clusterHint)
	return c.executeTopologyOperationMulti(ctx, path)
}

// SetClusterAlias sets an alias for a cluster
func (c *Client) SetClusterAlias(ctx context.Context, clusterName, alias string) error {
	path := fmt.Sprintf("/api/set-cluster-alias/%s", clusterName)
	// This endpoint expects the alias as a parameter
	// Using POST with alias in body or as query param
	var response APIResponse
	if err := c.postJSON(ctx, path, map[string]string{"alias": alias}, &response); err != nil {
		return err
	}

	if response.Code != OK {
		return fmt.Errorf("API error: %s", response.Message)
	}

	return nil
}

// GetClusters retrieves all known clusters
func (c *Client) GetClusters(ctx context.Context, ) ([]string, error) {
	path := "/api/clusters"
	var response APIResponse
	if err := c.getJSON(ctx, path, &response); err != nil {
		return nil, err
	}

	if response.Code != OK {
		return nil, fmt.Errorf("API error: %s", response.Message)
	}

	return convertToStringSlice(response.Details)
}

// GetClustersInfo retrieves information about all clusters
func (c *Client) GetClustersInfo(ctx context.Context, ) ([]ClusterInfo, error) {
	path := "/api/clusters-info"
	var response APIResponse
	if err := c.getJSON(ctx, path, &response); err != nil {
		return nil, err
	}

	if response.Code != OK {
		return nil, fmt.Errorf("API error: %s", response.Message)
	}

	return convertToClusterInfoSlice(response.Details)
}

// GetClusterMaster retrieves the master of a cluster
func (c *Client) GetClusterMaster(ctx context.Context, clusterHint string) (*Instance, error) {
	path := fmt.Sprintf("/api/master/%s", clusterHint)
	return c.executeTopologyOperation(ctx, path)
}

// GetAllMasters retrieves all cluster masters
func (c *Client) GetAllMasters(ctx context.Context, ) ([]Instance, error) {
	path := "/api/masters"
	return c.executeTopologyOperationMulti(ctx, path)
}

// ReloadClusterAlias reloads cluster alias configuration
func (c *Client) ReloadClusterAlias(ctx context.Context, ) error {
	path := "/api/reload-cluster-alias"
	var response APIResponse
	if err := c.getJSON(ctx, path, &response); err != nil {
		return err
	}

	if response.Code != OK {
		return fmt.Errorf("API error: %s", response.Message)
	}

	return nil
}
