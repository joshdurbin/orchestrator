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
	"time"

	"github.com/openark/orchestrator/go/inst"
)

// Maintenance Operations

// BeginMaintenance begins maintenance mode for an instance
func (c *Client) BeginMaintenance(ctx context.Context, instanceKey inst.InstanceKey, owner, reason string) error {
	path := fmt.Sprintf("/begin-maintenance/%s/%d/%s/%s",
		instanceKey.Hostname, instanceKey.Port, url.QueryEscape(owner), url.QueryEscape(reason))

	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return err
	}

	_, err = c.parseAPIResponse(resp, nil)
	return err
}

// EndMaintenance ends maintenance mode for an instance
func (c *Client) EndMaintenance(ctx context.Context, instanceKey inst.InstanceKey) error {
	path := fmt.Sprintf("/end-maintenance/%s/%d", instanceKey.Hostname, instanceKey.Port)
	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return err
	}

	_, err = c.parseAPIResponse(resp, nil)
	return err
}

// InMaintenance checks if an instance is in maintenance mode
func (c *Client) InMaintenance(ctx context.Context, instanceKey inst.InstanceKey) (bool, error) {
	path := fmt.Sprintf("/in-maintenance/%s/%d", instanceKey.Hostname, instanceKey.Port)
	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return false, err
	}

	var result string
	_, err = c.parseAPIResponse(resp, &result)
	if err != nil {
		return false, err
	}

	return result == "true", nil
}

// GetMaintenance gets list of instances under maintenance
func (c *Client) GetMaintenance(ctx context.Context) ([]interface{}, error) {
	resp, err := c.doRequest(ctx, "GET", "/maintenance", nil)
	if err != nil {
		return nil, err
	}

	var maintenance []interface{}
	if err := c.parseResponse(resp, &maintenance); err != nil {
		return nil, err
	}

	return maintenance, nil
}

// Downtime Operations

// BeginDowntime begins downtime for an instance
func (c *Client) BeginDowntime(ctx context.Context, instanceKey inst.InstanceKey, owner, reason string, duration time.Duration) error {
	durationStr := ""
	if duration > 0 {
		durationStr = fmt.Sprintf("/%ds", int(duration.Seconds()))
	}

	path := fmt.Sprintf("/begin-downtime/%s/%d/%s/%s%s",
		instanceKey.Hostname, instanceKey.Port, url.QueryEscape(owner), url.QueryEscape(reason), durationStr)

	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return err
	}

	_, err = c.parseAPIResponse(resp, nil)
	return err
}

// EndDowntime ends downtime for an instance
func (c *Client) EndDowntime(ctx context.Context, instanceKey inst.InstanceKey) error {
	path := fmt.Sprintf("/end-downtime/%s/%d", instanceKey.Hostname, instanceKey.Port)
	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return err
	}

	_, err = c.parseAPIResponse(resp, nil)
	return err
}

// GetDowntimed gets list of downtimed instances
func (c *Client) GetDowntimed(ctx context.Context, clusterName string) ([]inst.Instance, error) {
	var path string
	if clusterName != "" {
		path = fmt.Sprintf("/downtimed/%s", url.QueryEscape(clusterName))
	} else {
		path = "/downtimed"
	}

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

// Tag Operations

// GetTags gets all tags for an instance
func (c *Client) GetTags(ctx context.Context, instanceKey inst.InstanceKey) ([]string, error) {
	path := fmt.Sprintf("/tags/%s/%d", instanceKey.Hostname, instanceKey.Port)
	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var tags []string
	if err := c.parseResponse(resp, &tags); err != nil {
		return nil, err
	}

	return tags, nil
}

// GetTagValue gets the value of a specific tag for an instance
func (c *Client) GetTagValue(ctx context.Context, instanceKey inst.InstanceKey, tagName string) (string, error) {
	var path string
	if tagName != "" {
		path = fmt.Sprintf("/tag-value/%s/%d/%s", instanceKey.Hostname, instanceKey.Port, url.QueryEscape(tagName))
	} else {
		path = fmt.Sprintf("/tag-value/%s/%d", instanceKey.Hostname, instanceKey.Port)
	}

	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return "", err
	}

	var value string
	if err := c.parseResponse(resp, &value); err != nil {
		return "", err
	}

	return value, nil
}

// SetTag sets a tag on an instance
func (c *Client) SetTag(ctx context.Context, instanceKey inst.InstanceKey, tagName, tagValue string) error {
	var path string
	if tagValue != "" {
		path = fmt.Sprintf("/tag/%s/%d/%s/%s",
			instanceKey.Hostname, instanceKey.Port, url.QueryEscape(tagName), url.QueryEscape(tagValue))
	} else {
		path = fmt.Sprintf("/tag/%s/%d?tag=%s",
			instanceKey.Hostname, instanceKey.Port, url.QueryEscape(tagName))
	}

	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return err
	}

	_, err = c.parseAPIResponse(resp, nil)
	return err
}

// RemoveTag removes a tag from an instance
func (c *Client) RemoveTag(ctx context.Context, instanceKey inst.InstanceKey, tagName string) error {
	var path string
	if tagName != "" {
		path = fmt.Sprintf("/untag/%s/%d/%s", instanceKey.Hostname, instanceKey.Port, url.QueryEscape(tagName))
	} else {
		path = fmt.Sprintf("/untag/%s/%d", instanceKey.Hostname, instanceKey.Port)
	}

	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return err
	}

	_, err = c.parseAPIResponse(resp, nil)
	return err
}

// RemoveTagFromAll removes a tag from all instances that have it
func (c *Client) RemoveTagFromAll(ctx context.Context, tagName, tagValue string) ([]inst.Instance, error) {
	var path string
	if tagValue != "" {
		path = fmt.Sprintf("/untag-all/%s/%s", url.QueryEscape(tagName), url.QueryEscape(tagValue))
	} else {
		path = fmt.Sprintf("/untag-all?tag=%s", url.QueryEscape(tagName))
	}

	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var instances []inst.Instance
	_, err = c.parseAPIResponse(resp, &instances)
	if err != nil {
		return nil, err
	}

	return instances, nil
}

// GetTagged gets all instances with a specific tag
func (c *Client) GetTagged(ctx context.Context, tagName, tagValue string) ([]inst.Instance, error) {
	var path string
	if tagValue != "" {
		path = fmt.Sprintf("/tagged?tag=%s", url.QueryEscape(tagName+"="+tagValue))
	} else {
		path = fmt.Sprintf("/tagged?tag=%s", url.QueryEscape(tagName))
	}

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

// Candidate Registration

// RegisterCandidate registers an instance as a promotion candidate
func (c *Client) RegisterCandidate(ctx context.Context, instanceKey inst.InstanceKey, promotionRule string) error {
	path := fmt.Sprintf("/register-candidate/%s/%d/%s",
		instanceKey.Hostname, instanceKey.Port, url.QueryEscape(promotionRule))

	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return err
	}

	_, err = c.parseAPIResponse(resp, nil)
	return err
}

// Semi-sync Operations

// EnableSemiSyncMaster enables semi-sync on the master side
func (c *Client) EnableSemiSyncMaster(ctx context.Context, instanceKey inst.InstanceKey) (*inst.Instance, error) {
	path := fmt.Sprintf("/enable-semi-sync-master/%s/%d", instanceKey.Hostname, instanceKey.Port)
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

// DisableSemiSyncMaster disables semi-sync on the master side
func (c *Client) DisableSemiSyncMaster(ctx context.Context, instanceKey inst.InstanceKey) (*inst.Instance, error) {
	path := fmt.Sprintf("/disable-semi-sync-master/%s/%d", instanceKey.Hostname, instanceKey.Port)
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

// EnableSemiSyncReplica enables semi-sync on the replica side
func (c *Client) EnableSemiSyncReplica(ctx context.Context, instanceKey inst.InstanceKey) (*inst.Instance, error) {
	path := fmt.Sprintf("/enable-semi-sync-replica/%s/%d", instanceKey.Hostname, instanceKey.Port)
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

// DisableSemiSyncReplica disables semi-sync on the replica side
func (c *Client) DisableSemiSyncReplica(ctx context.Context, instanceKey inst.InstanceKey) (*inst.Instance, error) {
	path := fmt.Sprintf("/disable-semi-sync-replica/%s/%d", instanceKey.Hostname, instanceKey.Port)
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

// Pool Operations

// SubmitPoolInstances submits a list of instances for a pool
func (c *Client) SubmitPoolInstances(ctx context.Context, pool string, instances []inst.InstanceKey) error {
	instancesList := make([]string, len(instances))
	for i, instance := range instances {
		instancesList[i] = fmt.Sprintf("%s:%d", instance.Hostname, instance.Port)
	}
	instancesParam := url.QueryEscape(fmt.Sprintf("%v", instancesList))

	path := fmt.Sprintf("/submit-pool-instances/%s?instances=%s",
		url.QueryEscape(pool), instancesParam)

	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return err
	}

	_, err = c.parseAPIResponse(resp, nil)
	return err
}

// GetHeuristicClusterPoolInstances gets heuristic cluster pool instances
func (c *Client) GetHeuristicClusterPoolInstances(ctx context.Context, clusterName, pool string) ([]inst.Instance, error) {
	var path string
	if pool != "" {
		path = fmt.Sprintf("/heuristic-cluster-pool-instances/%s/%s",
			url.QueryEscape(clusterName), url.QueryEscape(pool))
	} else {
		path = fmt.Sprintf("/heuristic-cluster-pool-instances/%s", url.QueryEscape(clusterName))
	}

	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var instances []inst.Instance
	_, err = c.parseAPIResponse(resp, &instances)
	if err != nil {
		return nil, err
	}

	return instances, nil
}

// Hostname resolve operations

// RegisterHostnameUnresolve registers a hostname unresolve mapping
func (c *Client) RegisterHostnameUnresolve(ctx context.Context, instanceKey inst.InstanceKey, virtualname string) error {
	path := fmt.Sprintf("/register-hostname-unresolve/%s/%d/%s",
		instanceKey.Hostname, instanceKey.Port, url.QueryEscape(virtualname))

	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return err
	}

	_, err = c.parseAPIResponse(resp, nil)
	return err
}

// DeregisterHostnameUnresolve removes a hostname unresolve mapping
func (c *Client) DeregisterHostnameUnresolve(ctx context.Context, instanceKey inst.InstanceKey) error {
	path := fmt.Sprintf("/deregister-hostname-unresolve/%s/%d", instanceKey.Hostname, instanceKey.Port)
	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return err
	}

	_, err = c.parseAPIResponse(resp, nil)
	return err
}
