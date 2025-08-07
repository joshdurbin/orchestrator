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

	"github.com/openark/orchestrator/go/inst"
)

// Topology Manipulation Operations

// Relocate moves a replica to replicate under another instance
func (c *Client) Relocate(ctx context.Context, instanceKey inst.InstanceKey, belowKey inst.InstanceKey) (*inst.Instance, error) {
	path := fmt.Sprintf("/relocate/%s/%d/%s/%d",
		instanceKey.Hostname, instanceKey.Port, belowKey.Hostname, belowKey.Port)

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

// RelocateReplicas moves all replicas of an instance under another instance
func (c *Client) RelocateReplicas(ctx context.Context, instanceKey inst.InstanceKey, belowKey inst.InstanceKey) ([]inst.Instance, error) {
	path := fmt.Sprintf("/relocate-replicas/%s/%d/%s/%d",
		instanceKey.Hostname, instanceKey.Port, belowKey.Hostname, belowKey.Port)

	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var replicas []inst.Instance
	_, err = c.parseAPIResponse(resp, &replicas)
	if err != nil {
		return nil, err
	}

	return replicas, nil
}

// MoveUp moves a replica one level up in the topology
func (c *Client) MoveUp(ctx context.Context, instanceKey inst.InstanceKey) (*inst.Instance, error) {
	path := fmt.Sprintf("/move-up/%s/%d", instanceKey.Hostname, instanceKey.Port)
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

// MoveUpReplicas moves all replicas of an instance one level up in the topology
func (c *Client) MoveUpReplicas(ctx context.Context, instanceKey inst.InstanceKey) ([]inst.Instance, error) {
	path := fmt.Sprintf("/move-up-replicas/%s/%d", instanceKey.Hostname, instanceKey.Port)
	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var replicas []inst.Instance
	_, err = c.parseAPIResponse(resp, &replicas)
	if err != nil {
		return nil, err
	}

	return replicas, nil
}

// MoveBelow moves a replica beneath a sibling replica
func (c *Client) MoveBelow(ctx context.Context, instanceKey inst.InstanceKey, siblingKey inst.InstanceKey) (*inst.Instance, error) {
	path := fmt.Sprintf("/move-below/%s/%d/%s/%d",
		instanceKey.Hostname, instanceKey.Port, siblingKey.Hostname, siblingKey.Port)

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

// MakeCoMaster makes an instance a co-master with its current master
func (c *Client) MakeCoMaster(ctx context.Context, instanceKey inst.InstanceKey) (*inst.Instance, error) {
	path := fmt.Sprintf("/make-co-master/%s/%d", instanceKey.Hostname, instanceKey.Port)
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

// TakeMaster makes an instance the master of its current master
func (c *Client) TakeMaster(ctx context.Context, instanceKey inst.InstanceKey) (*inst.Instance, error) {
	path := fmt.Sprintf("/take-master/%s/%d", instanceKey.Hostname, instanceKey.Port)
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

// TakeSiblings makes an instance take all of its sibling replicas as its own replicas
func (c *Client) TakeSiblings(ctx context.Context, instanceKey inst.InstanceKey) (*inst.Instance, error) {
	path := fmt.Sprintf("/take-siblings/%s/%d", instanceKey.Hostname, instanceKey.Port)
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

// GTID-based operations

// MoveGTID moves a replica using GTID
func (c *Client) MoveGTID(ctx context.Context, instanceKey inst.InstanceKey, belowKey inst.InstanceKey) (*inst.Instance, error) {
	path := fmt.Sprintf("/move-gtid/%s/%d/%s/%d",
		instanceKey.Hostname, instanceKey.Port, belowKey.Hostname, belowKey.Port)

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

// MoveReplicasGTID moves all replicas of an instance using GTID
func (c *Client) MoveReplicasGTID(ctx context.Context, instanceKey inst.InstanceKey, belowKey inst.InstanceKey) ([]inst.Instance, error) {
	path := fmt.Sprintf("/move-replicas-gtid/%s/%d/%s/%d",
		instanceKey.Hostname, instanceKey.Port, belowKey.Hostname, belowKey.Port)

	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var replicas []inst.Instance
	_, err = c.parseAPIResponse(resp, &replicas)
	if err != nil {
		return nil, err
	}

	return replicas, nil
}

// Pseudo-GTID operations

// MatchBelow matches a replica below another instance using Pseudo-GTID
func (c *Client) MatchBelow(ctx context.Context, instanceKey inst.InstanceKey, belowKey inst.InstanceKey) (*inst.Instance, error) {
	path := fmt.Sprintf("/match/%s/%d/%s/%d",
		instanceKey.Hostname, instanceKey.Port, belowKey.Hostname, belowKey.Port)

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

// MatchUp matches a replica up one level using Pseudo-GTID
func (c *Client) MatchUp(ctx context.Context, instanceKey inst.InstanceKey) (*inst.Instance, error) {
	path := fmt.Sprintf("/match-up/%s/%d", instanceKey.Hostname, instanceKey.Port)
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

// MatchUpReplicas matches all replicas of an instance up one level using Pseudo-GTID
func (c *Client) MatchUpReplicas(ctx context.Context, instanceKey inst.InstanceKey) ([]inst.Instance, error) {
	path := fmt.Sprintf("/match-up-replicas/%s/%d", instanceKey.Hostname, instanceKey.Port)
	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var replicas []inst.Instance
	_, err = c.parseAPIResponse(resp, &replicas)
	if err != nil {
		return nil, err
	}

	return replicas, nil
}

// MultiMatchReplicas matches multiple replicas using Pseudo-GTID
func (c *Client) MultiMatchReplicas(ctx context.Context, instanceKey inst.InstanceKey, belowKey inst.InstanceKey) ([]inst.Instance, error) {
	path := fmt.Sprintf("/match-replicas/%s/%d/%s/%d",
		instanceKey.Hostname, instanceKey.Port, belowKey.Hostname, belowKey.Port)

	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var replicas []inst.Instance
	_, err = c.parseAPIResponse(resp, &replicas)
	if err != nil {
		return nil, err
	}

	return replicas, nil
}

// Repoint operations (for binlog servers)

// Repoint repositions a replica under another master with exact coordinates
func (c *Client) Repoint(ctx context.Context, instanceKey inst.InstanceKey, belowKey inst.InstanceKey) (*inst.Instance, error) {
	path := fmt.Sprintf("/repoint/%s/%d/%s/%d",
		instanceKey.Hostname, instanceKey.Port, belowKey.Hostname, belowKey.Port)

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

// RepointReplicas repoints all replicas of an instance
func (c *Client) RepointReplicas(ctx context.Context, instanceKey inst.InstanceKey) ([]inst.Instance, error) {
	path := fmt.Sprintf("/repoint-replicas/%s/%d", instanceKey.Hostname, instanceKey.Port)
	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var replicas []inst.Instance
	_, err = c.parseAPIResponse(resp, &replicas)
	if err != nil {
		return nil, err
	}

	return replicas, nil
}

// Utility functions for checking replication compatibility

// CanReplicateFrom checks if an instance can replicate from another
func (c *Client) CanReplicateFrom(ctx context.Context, instanceKey inst.InstanceKey, masterKey inst.InstanceKey) (bool, error) {
	path := fmt.Sprintf("/can-replicate-from/%s/%d/%s/%d",
		instanceKey.Hostname, instanceKey.Port, masterKey.Hostname, masterKey.Port)

	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return false, err
	}

	apiResp, err := c.parseAPIResponse(resp, nil)
	if err != nil {
		return false, err
	}

	// Check if the message contains "true"
	return apiResp.Message == "true", nil
}

// CanReplicateFromGTID checks if an instance can replicate from another using GTID
func (c *Client) CanReplicateFromGTID(ctx context.Context, instanceKey inst.InstanceKey, masterKey inst.InstanceKey) (bool, error) {
	path := fmt.Sprintf("/can-replicate-from-gtid/%s/%d/%s/%d",
		instanceKey.Hostname, instanceKey.Port, masterKey.Hostname, masterKey.Port)

	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return false, err
	}

	apiResp, err := c.parseAPIResponse(resp, nil)
	if err != nil {
		return false, err
	}

	// Check if the message contains "true"
	return apiResp.Message == "true", nil
}
