package client

import (
	"context"
	"encoding/json"
	"fmt"
)

// Raft Cluster Operations

// GrabElection attempts to become the Raft leader
func (c *Client) GrabElection(ctx context.Context, ) error {
	path := "/api/grab-election"
	var response APIResponse
	if err := c.getJSON(ctx, path, &response); err != nil {
		return err
	}

	if response.Code != OK {
		return fmt.Errorf("API error: %s", response.Message)
	}

	return nil
}

// AddRaftPeer adds a peer to the Raft cluster
func (c *Client) AddRaftPeer(ctx context.Context, addr string) error {
	path := fmt.Sprintf("/api/raft-add-peer/%s", addr)
	var response APIResponse
	if err := c.getJSON(ctx, path, &response); err != nil {
		return err
	}

	if response.Code != OK {
		return fmt.Errorf("API error: %s", response.Message)
	}

	return nil
}

// RemoveRaftPeer removes a peer from the Raft cluster
func (c *Client) RemoveRaftPeer(ctx context.Context, addr string) error {
	path := fmt.Sprintf("/api/raft-remove-peer/%s", addr)
	var response APIResponse
	if err := c.getJSON(ctx, path, &response); err != nil {
		return err
	}

	if response.Code != OK {
		return fmt.Errorf("API error: %s", response.Message)
	}

	return nil
}

// YieldRaft yields Raft leadership to a specific node
func (c *Client) YieldRaft(ctx context.Context, node string) error {
	path := fmt.Sprintf("/api/raft-yield/%s", node)
	var response APIResponse
	if err := c.getJSON(ctx, path, &response); err != nil {
		return err
	}

	if response.Code != OK {
		return fmt.Errorf("API error: %s", response.Message)
	}

	return nil
}

// YieldRaftHint yields Raft leadership with a hint
func (c *Client) YieldRaftHint(ctx context.Context, hint string) error {
	path := fmt.Sprintf("/api/raft-yield-hint/%s", hint)
	var response APIResponse
	if err := c.getJSON(ctx, path, &response); err != nil {
		return err
	}

	if response.Code != OK {
		return fmt.Errorf("API error: %s", response.Message)
	}

	return nil
}

// GetRaftPeers retrieves the list of Raft peers
func (c *Client) GetRaftPeers(ctx context.Context, ) ([]string, error) {
	path := "/api/raft-peers"
	var response APIResponse
	if err := c.getJSON(ctx, path, &response); err != nil {
		return nil, err
	}

	if response.Code != OK {
		return nil, fmt.Errorf("API error: %s", response.Message)
	}

	return convertToStringSlice(response.Details)
}

// GetRaftState retrieves the current Raft state
func (c *Client) GetRaftState(ctx context.Context, ) (*RaftState, error) {
	path := "/api/raft-state"
	var response APIResponse
	if err := c.getJSON(ctx, path, &response); err != nil {
		return nil, err
	}

	if response.Code != OK {
		return nil, fmt.Errorf("API error: %s", response.Message)
	}

	jsonData, err := json.Marshal(response.Details)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data: %w", err)
	}

	var state RaftState
	if err := json.Unmarshal(jsonData, &state); err != nil {
		return nil, fmt.Errorf("failed to unmarshal raft state: %w", err)
	}

	return &state, nil
}

// GetRaftLeader retrieves the current Raft leader
func (c *Client) GetRaftLeader(ctx context.Context, ) (string, error) {
	path := "/api/raft-leader"
	var response APIResponse
	if err := c.getJSON(ctx, path, &response); err != nil {
		return "", err
	}

	if response.Code != OK {
		return "", fmt.Errorf("API error: %s", response.Message)
	}

	if response.Details == nil {
		return "", nil
	}

	if str, ok := response.Details.(string); ok {
		return str, nil
	}

	return fmt.Sprint(response.Details), nil
}

// GetRaftHealth retrieves Raft health status
func (c *Client) GetRaftHealth(ctx context.Context, ) (*RaftMembershipHealth, error) {
	path := "/api/raft-health"
	var response APIResponse
	if err := c.getJSON(ctx, path, &response); err != nil {
		return nil, err
	}

	if response.Code != OK {
		return nil, fmt.Errorf("API error: %s", response.Message)
	}

	jsonData, err := json.Marshal(response.Details)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data: %w", err)
	}

	var health RaftMembershipHealth
	if err := json.Unmarshal(jsonData, &health); err != nil {
		return nil, fmt.Errorf("failed to unmarshal raft health: %w", err)
	}

	return &health, nil
}

// GetRaftStatus retrieves Raft status information
func (c *Client) GetRaftStatus(ctx context.Context, ) (map[string]interface{}, error) {
	path := "/api/raft-status"
	var response APIResponse
	if err := c.getJSON(ctx, path, &response); err != nil {
		return nil, err
	}

	if response.Code != OK {
		return nil, fmt.Errorf("API error: %s", response.Message)
	}

	if response.Details == nil {
		return map[string]interface{}{}, nil
	}

	jsonData, err := json.Marshal(response.Details)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data: %w", err)
	}

	var status map[string]interface{}
	if err := json.Unmarshal(jsonData, &status); err != nil {
		return nil, fmt.Errorf("failed to unmarshal raft status: %w", err)
	}

	return status, nil
}

// GetRaftSnapshot retrieves Raft snapshot information
func (c *Client) GetRaftSnapshot(ctx context.Context, ) ([]byte, error) {
	path := "/api/raft-snapshot"
	return c.getRawBytes(ctx, path)
}

// SubmitRaftFollowerHealthReport submits a follower health report
func (c *Client) SubmitRaftFollowerHealthReport(ctx context.Context, authToken, raftBind, raftAdvertise string) error {
	path := fmt.Sprintf("/api/raft-follower-health-report/%s/%s/%s",
		authToken, raftBind, raftAdvertise)
	var response APIResponse
	if err := c.getJSON(ctx, path, &response); err != nil {
		return err
	}

	if response.Code != OK {
		return fmt.Errorf("API error: %s", response.Message)
	}

	return nil
}

// Reelect triggers a reelection
func (c *Client) Reelect(ctx context.Context, ) error {
	path := "/api/reelect"
	var response APIResponse
	if err := c.getJSON(ctx, path, &response); err != nil {
		return err
	}

	if response.Code != OK {
		return fmt.Errorf("API error: %s", response.Message)
	}

	return nil
}

// Helper function for getting raw bytes
func (c *Client) getRawBytes(ctx context.Context, path string) ([]byte, error) {
	resp, err := c.get(ctx, path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := json.Marshal(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	return body, nil
}
