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

// Recovery and Failover Operations

// Recover performs automatic recovery for a failed instance
func (c *Client) Recover(ctx context.Context, instanceKey inst.InstanceKey, candidateKey *inst.InstanceKey) (*inst.Instance, error) {
	var path string
	if candidateKey != nil {
		path = fmt.Sprintf("/recover/%s/%d/%s/%d", instanceKey.Hostname, instanceKey.Port, candidateKey.Hostname, candidateKey.Port)
	} else {
		path = fmt.Sprintf("/recover/%s/%d", instanceKey.Hostname, instanceKey.Port)
	}

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

// GracefulMasterTakeover performs a graceful master takeover
func (c *Client) GracefulMasterTakeover(ctx context.Context, clusterName string, designatedKey *inst.InstanceKey) (*inst.Instance, error) {
	var path string
	if designatedKey != nil {
		path = fmt.Sprintf("/graceful-master-takeover/%s/%s/%d",
			url.QueryEscape(clusterName), designatedKey.Hostname, designatedKey.Port)
	} else {
		path = fmt.Sprintf("/graceful-master-takeover/%s", url.QueryEscape(clusterName))
	}

	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	_, err = c.parseAPIResponse(resp, &result)
	if err != nil {
		return nil, err
	}

	// Extract SuccessorKey from result
	if successorData, ok := result["SuccessorKey"].(map[string]interface{}); ok {
		hostname, _ := successorData["Hostname"].(string)
		port, _ := successorData["Port"].(float64)

		successorKey := inst.InstanceKey{
			Hostname: hostname,
			Port:     int(port),
		}

		// Get the full instance details
		return c.GetInstance(ctx, successorKey)
	}

	return nil, fmt.Errorf("unexpected response format")
}

// GracefulMasterTakeoverAuto performs automatic graceful master takeover
func (c *Client) GracefulMasterTakeoverAuto(ctx context.Context, clusterName string, designatedKey *inst.InstanceKey) (*inst.Instance, error) {
	var path string
	if designatedKey != nil {
		path = fmt.Sprintf("/graceful-master-takeover-auto/%s/%s/%d",
			url.QueryEscape(clusterName), designatedKey.Hostname, designatedKey.Port)
	} else {
		path = fmt.Sprintf("/graceful-master-takeover-auto/%s", url.QueryEscape(clusterName))
	}

	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	_, err = c.parseAPIResponse(resp, &result)
	if err != nil {
		return nil, err
	}

	// Extract SuccessorKey from result
	if successorData, ok := result["SuccessorKey"].(map[string]interface{}); ok {
		hostname, _ := successorData["Hostname"].(string)
		port, _ := successorData["Port"].(float64)

		successorKey := inst.InstanceKey{
			Hostname: hostname,
			Port:     int(port),
		}

		// Get the full instance details
		return c.GetInstance(ctx, successorKey)
	}

	return nil, fmt.Errorf("unexpected response format")
}

// ForceMasterFailover forces a master failover
func (c *Client) ForceMasterFailover(ctx context.Context, clusterName string) (*inst.Instance, error) {
	path := fmt.Sprintf("/force-master-failover/%s", url.QueryEscape(clusterName))
	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	_, err = c.parseAPIResponse(resp, &result)
	if err != nil {
		return nil, err
	}

	// Extract SuccessorKey from result
	if successorData, ok := result["SuccessorKey"].(map[string]interface{}); ok {
		hostname, _ := successorData["Hostname"].(string)
		port, _ := successorData["Port"].(float64)

		successorKey := inst.InstanceKey{
			Hostname: hostname,
			Port:     int(port),
		}

		// Get the full instance details
		return c.GetInstance(ctx, successorKey)
	}

	return nil, fmt.Errorf("unexpected response format")
}

// ForceMasterTakeover forces a master takeover to a specific instance
func (c *Client) ForceMasterTakeover(ctx context.Context, clusterName string, designatedKey inst.InstanceKey) (*inst.Instance, error) {
	path := fmt.Sprintf("/force-master-takeover/%s/%s/%d",
		url.QueryEscape(clusterName), designatedKey.Hostname, designatedKey.Port)

	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	_, err = c.parseAPIResponse(resp, &result)
	if err != nil {
		return nil, err
	}

	// Extract SuccessorKey from result
	if successorData, ok := result["SuccessorKey"].(map[string]interface{}); ok {
		hostname, _ := successorData["Hostname"].(string)
		port, _ := successorData["Port"].(float64)

		successorKey := inst.InstanceKey{
			Hostname: hostname,
			Port:     int(port),
		}

		// Get the full instance details
		return c.GetInstance(ctx, successorKey)
	}

	return nil, fmt.Errorf("unexpected response format")
}

// GetReplicationAnalysis gets replication analysis for all clusters
func (c *Client) GetReplicationAnalysis(ctx context.Context) ([]interface{}, error) {
	resp, err := c.doRequest(ctx, "GET", "/replication-analysis", nil)
	if err != nil {
		return nil, err
	}

	var analysis []interface{}
	_, err = c.parseAPIResponse(resp, &analysis)
	if err != nil {
		return nil, err
	}

	return analysis, nil
}

// GetReplicationAnalysisForCluster gets replication analysis for a specific cluster
func (c *Client) GetReplicationAnalysisForCluster(ctx context.Context, clusterName string) ([]interface{}, error) {
	path := fmt.Sprintf("/replication-analysis/%s", url.QueryEscape(clusterName))
	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var analysis []interface{}
	_, err = c.parseAPIResponse(resp, &analysis)
	if err != nil {
		return nil, err
	}

	return analysis, nil
}

// AcknowledgeClusterRecoveries acknowledges recoveries for a cluster
func (c *Client) AcknowledgeClusterRecoveries(ctx context.Context, clusterName string, comment string) error {
	path := fmt.Sprintf("/ack-recovery/cluster/%s", url.QueryEscape(clusterName))
	if comment != "" {
		path += "?comment=" + url.QueryEscape(comment)
	}

	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return err
	}

	_, err = c.parseAPIResponse(resp, nil)
	return err
}

// AcknowledgeInstanceRecoveries acknowledges recoveries for an instance
func (c *Client) AcknowledgeInstanceRecoveries(ctx context.Context, instanceKey inst.InstanceKey, comment string) error {
	path := fmt.Sprintf("/ack-recovery/instance/%s/%d", instanceKey.Hostname, instanceKey.Port)
	if comment != "" {
		path += "?comment=" + url.QueryEscape(comment)
	}

	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return err
	}

	_, err = c.parseAPIResponse(resp, nil)
	return err
}

// AcknowledgeAllRecoveries acknowledges all recoveries
func (c *Client) AcknowledgeAllRecoveries(ctx context.Context, comment string) error {
	path := "/ack-all-recoveries"
	if comment != "" {
		path += "?comment=" + url.QueryEscape(comment)
	}

	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return err
	}

	_, err = c.parseAPIResponse(resp, nil)
	return err
}

// DisableGlobalRecoveries disables global recoveries
func (c *Client) DisableGlobalRecoveries(ctx context.Context) error {
	resp, err := c.doRequest(ctx, "GET", "/disable-global-recoveries", nil)
	if err != nil {
		return err
	}

	_, err = c.parseAPIResponse(resp, nil)
	return err
}

// EnableGlobalRecoveries enables global recoveries
func (c *Client) EnableGlobalRecoveries(ctx context.Context) error {
	resp, err := c.doRequest(ctx, "GET", "/enable-global-recoveries", nil)
	if err != nil {
		return err
	}

	_, err = c.parseAPIResponse(resp, nil)
	return err
}

// CheckGlobalRecoveries checks the status of global recoveries
func (c *Client) CheckGlobalRecoveries(ctx context.Context) (bool, error) {
	resp, err := c.doRequest(ctx, "GET", "/check-global-recoveries", nil)
	if err != nil {
		return false, err
	}

	var result bool
	_, err = c.parseAPIResponse(resp, &result)
	if err != nil {
		return false, err
	}

	return result, nil
}

// GetBlockedRecoveries gets list of blocked recoveries
func (c *Client) GetBlockedRecoveries(ctx context.Context, clusterName string) ([]interface{}, error) {
	var path string
	if clusterName != "" {
		path = fmt.Sprintf("/blocked-recoveries/cluster/%s", url.QueryEscape(clusterName))
	} else {
		path = "/blocked-recoveries"
	}

	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var blocked []interface{}
	_, err = c.parseAPIResponse(resp, &blocked)
	if err != nil {
		return nil, err
	}

	return blocked, nil
}
