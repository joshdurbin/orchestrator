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

package client

import (
	"context"
	"fmt"
	"net/url"

	"github.com/openark/orchestrator/go/inst"
)

// Instance operations

// Discover discovers and reads information about a MySQL instance
func (c *Client) Discover(ctx context.Context, instanceKey inst.InstanceKey) (*inst.Instance, error) {
	path := fmt.Sprintf("/discover/%s/%d", instanceKey.Hostname, instanceKey.Port)
	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var instance inst.Instance
	apiResp, err := c.parseAPIResponse(resp, &instance)
	if err != nil {
		return nil, err
	}

	if apiResp.Code != "OK" {
		return nil, fmt.Errorf("discover failed: %s", apiResp.Message)
	}

	return &instance, nil
}

// AsyncDiscover initiates asynchronous discovery of a MySQL instance
func (c *Client) AsyncDiscover(ctx context.Context, instanceKey inst.InstanceKey) error {
	path := fmt.Sprintf("/async-discover/%s/%d", instanceKey.Hostname, instanceKey.Port)
	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return err
	}

	_, err = c.parseAPIResponse(resp, nil)
	return err
}

// GetInstance retrieves information about a specific instance
func (c *Client) GetInstance(ctx context.Context, instanceKey inst.InstanceKey) (*inst.Instance, error) {
	path := fmt.Sprintf("/instance/%s/%d", instanceKey.Hostname, instanceKey.Port)
	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var instance inst.Instance
	if err := c.parseResponse(resp, &instance); err != nil {
		return nil, err
	}

	return &instance, nil
}

// GetInstanceReplicas gets the list of replicas for a given instance
func (c *Client) GetInstanceReplicas(ctx context.Context, instanceKey inst.InstanceKey) ([]inst.Instance, error) {
	path := fmt.Sprintf("/instance-replicas/%s/%d", instanceKey.Hostname, instanceKey.Port)
	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var replicas []inst.Instance
	if err := c.parseResponse(resp, &replicas); err != nil {
		return nil, err
	}

	return replicas, nil
}

// Forget removes an instance from orchestrator's tracking
func (c *Client) Forget(ctx context.Context, instanceKey inst.InstanceKey) error {
	path := fmt.Sprintf("/forget/%s/%d", instanceKey.Hostname, instanceKey.Port)
	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return err
	}

	_, err = c.parseAPIResponse(resp, nil)
	return err
}

// Search searches for instances matching the given substring
func (c *Client) Search(ctx context.Context, searchString string) ([]inst.Instance, error) {
	path := fmt.Sprintf("/search/%s", url.QueryEscape(searchString))
	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var instances []inst.Instance
	if err := c.parseResponse(resp, &instances); err != nil {
		return nil, err
	}

	return instances, nil
}

// Cluster operations

// GetCluster retrieves all instances in a cluster
func (c *Client) GetCluster(ctx context.Context, clusterName string) ([]inst.Instance, error) {
	path := fmt.Sprintf("/cluster/%s", url.QueryEscape(clusterName))
	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var instances []inst.Instance
	if err := c.parseResponse(resp, &instances); err != nil {
		return nil, err
	}

	return instances, nil
}

// GetClusterByInstance retrieves all instances in the same cluster as the given instance
func (c *Client) GetClusterByInstance(ctx context.Context, instanceKey inst.InstanceKey) ([]inst.Instance, error) {
	path := fmt.Sprintf("/cluster/instance/%s/%d", instanceKey.Hostname, instanceKey.Port)
	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var instances []inst.Instance
	if err := c.parseResponse(resp, &instances); err != nil {
		return nil, err
	}

	return instances, nil
}

// GetClusterAlias retrieves a cluster by its alias
func (c *Client) GetClusterAlias(ctx context.Context, alias string) ([]inst.Instance, error) {
	path := fmt.Sprintf("/cluster/alias/%s", url.QueryEscape(alias))
	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var instances []inst.Instance
	if err := c.parseResponse(resp, &instances); err != nil {
		return nil, err
	}

	return instances, nil
}

// GetClusters retrieves a list of all cluster names
func (c *Client) GetClusters(ctx context.Context) ([]string, error) {
	resp, err := c.doRequest(ctx, "GET", "/clusters", nil)
	if err != nil {
		return nil, err
	}

	var clusters []string
	if err := c.parseResponse(resp, &clusters); err != nil {
		return nil, err
	}

	return clusters, nil
}

// GetClusterMaster retrieves the master instance of a cluster
func (c *Client) GetClusterMaster(ctx context.Context, clusterName string) (*inst.Instance, error) {
	path := fmt.Sprintf("/master/%s", url.QueryEscape(clusterName))
	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var master inst.Instance
	if err := c.parseResponse(resp, &master); err != nil {
		return nil, err
	}

	return &master, nil
}

// ForgetCluster removes all instances of a cluster from orchestrator's tracking
func (c *Client) ForgetCluster(ctx context.Context, clusterName string) error {
	path := fmt.Sprintf("/forget-cluster/%s", url.QueryEscape(clusterName))
	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return err
	}

	_, err = c.parseAPIResponse(resp, nil)
	return err
}

// GetAllInstances retrieves all instances known to orchestrator
func (c *Client) GetAllInstances(ctx context.Context) ([]inst.Instance, error) {
	resp, err := c.doRequest(ctx, "GET", "/all-instances", nil)
	if err != nil {
		return nil, err
	}

	var instances []inst.Instance
	if err := c.parseResponse(resp, &instances); err != nil {
		return nil, err
	}

	return instances, nil
}

// GetMasters retrieves all master instances (one per cluster)
func (c *Client) GetMasters(ctx context.Context) ([]inst.Instance, error) {
	resp, err := c.doRequest(ctx, "GET", "/masters", nil)
	if err != nil {
		return nil, err
	}

	var masters []inst.Instance
	if err := c.parseResponse(resp, &masters); err != nil {
		return nil, err
	}

	return masters, nil
}

// GetTopology retrieves ASCII topology representation
func (c *Client) GetTopology(ctx context.Context, clusterName string) (string, error) {
	path := fmt.Sprintf("/topology/%s", url.QueryEscape(clusterName))
	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return "", err
	}

	apiResp, err := c.parseAPIResponse(resp, nil)
	if err != nil {
		return "", err
	}

	if topology, ok := apiResp.Details.(string); ok {
		return topology, nil
	}

	return "", fmt.Errorf("unexpected response format")
}

// GetTopologyTabulated retrieves tabulated ASCII topology representation
func (c *Client) GetTopologyTabulated(ctx context.Context, clusterName string) (string, error) {
	path := fmt.Sprintf("/topology-tabulated/%s", url.QueryEscape(clusterName))
	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return "", err
	}

	apiResp, err := c.parseAPIResponse(resp, nil)
	if err != nil {
		return "", err
	}

	if topology, ok := apiResp.Details.(string); ok {
		return topology, nil
	}

	return "", fmt.Errorf("unexpected response format")
}

// Replication Control Operations

// StartReplication starts replication on an instance
func (c *Client) StartReplication(ctx context.Context, instanceKey inst.InstanceKey) (*inst.Instance, error) {
	path := fmt.Sprintf("/start-replica/%s/%d", instanceKey.Hostname, instanceKey.Port)
	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var instance inst.Instance
	_, err = c.parseAPIResponse(resp, &instance)
	if err != nil {
		return nil, err
	}

	return &instance, nil
}

// StopReplication stops replication on an instance
func (c *Client) StopReplication(ctx context.Context, instanceKey inst.InstanceKey) (*inst.Instance, error) {
	path := fmt.Sprintf("/stop-replica/%s/%d", instanceKey.Hostname, instanceKey.Port)
	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var instance inst.Instance
	_, err = c.parseAPIResponse(resp, &instance)
	if err != nil {
		return nil, err
	}

	return &instance, nil
}

// RestartReplication restarts replication on an instance
func (c *Client) RestartReplication(ctx context.Context, instanceKey inst.InstanceKey) (*inst.Instance, error) {
	path := fmt.Sprintf("/restart-replica/%s/%d", instanceKey.Hostname, instanceKey.Port)
	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var instance inst.Instance
	_, err = c.parseAPIResponse(resp, &instance)
	if err != nil {
		return nil, err
	}

	return &instance, nil
}

// ResetReplication resets replication on an instance
func (c *Client) ResetReplication(ctx context.Context, instanceKey inst.InstanceKey) (*inst.Instance, error) {
	path := fmt.Sprintf("/reset-replica/%s/%d", instanceKey.Hostname, instanceKey.Port)
	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var instance inst.Instance
	_, err = c.parseAPIResponse(resp, &instance)
	if err != nil {
		return nil, err
	}

	return &instance, nil
}

// SetReadOnly sets an instance to read-only mode
func (c *Client) SetReadOnly(ctx context.Context, instanceKey inst.InstanceKey) (*inst.Instance, error) {
	path := fmt.Sprintf("/set-read-only/%s/%d", instanceKey.Hostname, instanceKey.Port)
	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var instance inst.Instance
	_, err = c.parseAPIResponse(resp, &instance)
	if err != nil {
		return nil, err
	}

	return &instance, nil
}

// SetWriteable sets an instance to writeable mode
func (c *Client) SetWriteable(ctx context.Context, instanceKey inst.InstanceKey) (*inst.Instance, error) {
	path := fmt.Sprintf("/set-writeable/%s/%d", instanceKey.Hostname, instanceKey.Port)
	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var instance inst.Instance
	_, err = c.parseAPIResponse(resp, &instance)
	if err != nil {
		return nil, err
	}

	return &instance, nil
}

// DelayReplication sets replication delay on an instance
func (c *Client) DelayReplication(ctx context.Context, instanceKey inst.InstanceKey, seconds int) (*APIResponse, error) {
	path := fmt.Sprintf("/delay-replication/%s/%d/%d", instanceKey.Hostname, instanceKey.Port, seconds)
	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	return c.parseAPIResponse(resp, nil)
}

// SkipQuery skips a single query on a replica
func (c *Client) SkipQuery(ctx context.Context, instanceKey inst.InstanceKey) (*inst.Instance, error) {
	path := fmt.Sprintf("/skip-query/%s/%d", instanceKey.Hostname, instanceKey.Port)
	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var instance inst.Instance
	_, err = c.parseAPIResponse(resp, &instance)
	if err != nil {
		return nil, err
	}

	return &instance, nil
}

// PurgeBinaryLogs purges binary logs up to the specified log file
func (c *Client) PurgeBinaryLogs(ctx context.Context, instanceKey inst.InstanceKey, logFile string) (*inst.Instance, error) {
	path := fmt.Sprintf("/purge-binary-logs/%s/%d/%s", instanceKey.Hostname, instanceKey.Port, url.QueryEscape(logFile))
	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var instance inst.Instance
	_, err = c.parseAPIResponse(resp, &instance)
	if err != nil {
		return nil, err
	}

	return &instance, nil
}

// GetLastPseudoGTID retrieves the last pseudo-GTID entry for an instance
func (c *Client) GetLastPseudoGTID(ctx context.Context, instanceKey inst.InstanceKey) (string, error) {
	path := fmt.Sprintf("/last-pseudo-gtid/%s/%d", instanceKey.Hostname, instanceKey.Port)
	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return "", err
	}

	var result string
	_, err = c.parseAPIResponse(resp, &result)
	if err != nil {
		return "", err
	}

	return result, nil
}
