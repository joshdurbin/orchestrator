package client

import (
	"context"
	"encoding/json"
	"fmt"
)

// Miscellaneous Operations (Pools, KV Store, Audit, Health, Hostname Management, etc.)

// Pool Operations

// SubmitPoolInstances submits instances to a pool
func (c *Client) SubmitPoolInstances(ctx context.Context, pool string) error {
	path := fmt.Sprintf("/api/submit-pool-instances/%s", pool)
	var response APIResponse
	if err := c.getJSON(ctx, path, &response); err != nil {
		return err
	}

	if response.Code != OK {
		return fmt.Errorf("API error: %s", response.Message)
	}

	return nil
}

// GetClusterPoolInstances retrieves pool instances for a cluster
func (c *Client) GetClusterPoolInstances(ctx context.Context, clusterName string) (PoolInstancesMap, error) {
	path := fmt.Sprintf("/api/cluster-pool-instances/%s", clusterName)
	var response APIResponse
	if err := c.getJSON(ctx, path, &response); err != nil {
		return nil, err
	}

	if response.Code != OK {
		return nil, fmt.Errorf("API error: %s", response.Message)
	}

	return convertToPoolInstancesMap(response.Details)
}

// GetClusterPoolInstancesForPool retrieves instances for a specific pool in a cluster
func (c *Client) GetClusterPoolInstancesForPool(ctx context.Context, clusterName, pool string) ([]InstanceKey, error) {
	path := fmt.Sprintf("/api/cluster-pool-instances/%s/%s", clusterName, pool)
	var response APIResponse
	if err := c.getJSON(ctx, path, &response); err != nil {
		return nil, err
	}

	if response.Code != OK {
		return nil, fmt.Errorf("API error: %s", response.Message)
	}

	// Convert string slice to InstanceKey slice
	strSlice, err := convertToStringSlice(response.Details)
	if err != nil {
		return nil, err
	}

	keys := make([]InstanceKey, 0, len(strSlice))
	for _, keyStr := range strSlice {
		key, err := ParseInstanceKey(keyStr)
		if err != nil {
			continue
		}
		keys = append(keys, key)
	}

	return keys, nil
}

// GetHeuristicClusterPoolInstances retrieves heuristic pool instances for a cluster
func (c *Client) GetHeuristicClusterPoolInstances(ctx context.Context, clusterName string) (PoolInstancesMap, error) {
	path := fmt.Sprintf("/api/heuristic-cluster-pool-instances/%s", clusterName)
	var response APIResponse
	if err := c.getJSON(ctx, path, &response); err != nil {
		return nil, err
	}

	if response.Code != OK {
		return nil, fmt.Errorf("API error: %s", response.Message)
	}

	return convertToPoolInstancesMap(response.Details)
}

// GetHeuristicClusterPoolInstancesForPool retrieves heuristic instances for a specific pool
func (c *Client) GetHeuristicClusterPoolInstancesForPool(ctx context.Context, clusterName, pool string) ([]InstanceKey, error) {
	path := fmt.Sprintf("/api/heuristic-cluster-pool-instances/%s/%s", clusterName, pool)
	var response APIResponse
	if err := c.getJSON(ctx, path, &response); err != nil {
		return nil, err
	}

	if response.Code != OK {
		return nil, fmt.Errorf("API error: %s", response.Message)
	}

	strSlice, err := convertToStringSlice(response.Details)
	if err != nil {
		return nil, err
	}

	keys := make([]InstanceKey, 0, len(strSlice))
	for _, keyStr := range strSlice {
		key, err := ParseInstanceKey(keyStr)
		if err != nil {
			continue
		}
		keys = append(keys, key)
	}

	return keys, nil
}

// GetHeuristicClusterPoolLag retrieves heuristic pool lag for a cluster
func (c *Client) GetHeuristicClusterPoolLag(ctx context.Context, clusterName string) (map[string]int64, error) {
	path := fmt.Sprintf("/api/heuristic-cluster-pool-lag/%s", clusterName)
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

	var lagMap map[string]int64
	if err := json.Unmarshal(jsonData, &lagMap); err != nil {
		return nil, fmt.Errorf("failed to unmarshal lag map: %w", err)
	}

	return lagMap, nil
}

// GetHeuristicClusterPoolLagForPool retrieves heuristic pool lag for a specific pool
func (c *Client) GetHeuristicClusterPoolLagForPool(ctx context.Context, clusterName, pool string) (int64, error) {
	path := fmt.Sprintf("/api/heuristic-cluster-pool-lag/%s/%s", clusterName, pool)
	var response APIResponse
	if err := c.getJSON(ctx, path, &response); err != nil {
		return 0, err
	}

	if response.Code != OK {
		return 0, fmt.Errorf("API error: %s", response.Message)
	}

	if response.Details == nil {
		return 0, nil
	}

	// Try to convert to int64
	switch v := response.Details.(type) {
	case float64:
		return int64(v), nil
	case int64:
		return v, nil
	case int:
		return int64(v), nil
	default:
		return 0, fmt.Errorf("unexpected type for lag value: %T", response.Details)
	}
}

// Key-Value Store Operations

// SubmitMastersToKVStores submits all masters to key-value stores
func (c *Client) SubmitMastersToKVStores(ctx context.Context, ) error {
	path := "/api/submit-masters-to-kv-stores"
	var response APIResponse
	if err := c.getJSON(ctx, path, &response); err != nil {
		return err
	}

	if response.Code != OK {
		return fmt.Errorf("API error: %s", response.Message)
	}

	return nil
}

// SubmitClusterMasterToKVStores submits a cluster's master to key-value stores
func (c *Client) SubmitClusterMasterToKVStores(ctx context.Context, clusterHint string) error {
	path := fmt.Sprintf("/api/submit-masters-to-kv-stores/%s", clusterHint)
	var response APIResponse
	if err := c.getJSON(ctx, path, &response); err != nil {
		return err
	}

	if response.Code != OK {
		return fmt.Errorf("API error: %s", response.Message)
	}

	return nil
}

// Audit Operations

// GetAudit retrieves the audit log
func (c *Client) GetAudit(ctx context.Context, page int) ([]AuditEntry, error) {
	path := fmt.Sprintf("/api/audit/%d", page)
	var response APIResponse
	if err := c.getJSON(ctx, path, &response); err != nil {
		return nil, err
	}

	if response.Code != OK {
		return nil, fmt.Errorf("API error: %s", response.Message)
	}

	return convertToAuditEntrySlice(response.Details)
}

// GetAuditForInstance retrieves the audit log for a specific instance
func (c *Client) GetAuditForInstance(ctx context.Context, instanceKey InstanceKey, page int) ([]AuditEntry, error) {
	path := fmt.Sprintf("/api/audit/instance/%s/%d/%d",
		instanceKey.Hostname, instanceKey.Port, page)
	var response APIResponse
	if err := c.getJSON(ctx, path, &response); err != nil {
		return nil, err
	}

	if response.Code != OK {
		return nil, fmt.Errorf("API error: %s", response.Message)
	}

	return convertToAuditEntrySlice(response.Details)
}

// Health and Status Operations

// GetHeaders retrieves request headers
func (c *Client) GetHeaders(ctx context.Context, ) (map[string][]string, error) {
	path := "/api/headers"
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

	var headers map[string][]string
	if err := json.Unmarshal(jsonData, &headers); err != nil {
		return nil, fmt.Errorf("failed to unmarshal headers: %w", err)
	}

	return headers, nil
}

// Health performs a health check
func (c *Client) Health(ctx context.Context, ) error {
	path := "/api/health"
	var response APIResponse
	if err := c.getJSON(ctx, path, &response); err != nil {
		return err
	}

	if response.Code != OK {
		return fmt.Errorf("health check failed: %s", response.Message)
	}

	return nil
}

// LBCheck performs a load balancer health check
func (c *Client) LBCheck(ctx context.Context, ) error {
	path := "/api/lb-check"
	var response APIResponse
	if err := c.getJSON(ctx, path, &response); err != nil {
		return err
	}

	if response.Code != OK {
		return fmt.Errorf("LB check failed: %s", response.Message)
	}

	return nil
}

// Ping performs a simple ping to the API
func (c *Client) Ping(ctx context.Context, ) error {
	path := "/api/_ping"
	var response APIResponse
	if err := c.getJSON(ctx, path, &response); err != nil {
		return err
	}

	if response.Code != OK {
		return fmt.Errorf("ping failed: %s", response.Message)
	}

	return nil
}

// LeaderCheck checks if the current node is the leader
func (c *Client) LeaderCheck(ctx context.Context, ) error {
	path := "/api/leader-check"
	var response APIResponse
	if err := c.getJSON(ctx, path, &response); err != nil {
		return err
	}

	if response.Code != OK {
		return fmt.Errorf("leader check failed: %s", response.Message)
	}

	return nil
}

// LeaderCheckWithErrorCode checks if the current node is the leader with custom error code
func (c *Client) LeaderCheckWithErrorCode(ctx context.Context, errorStatusCode int) error {
	path := fmt.Sprintf("/api/leader-check/%d", errorStatusCode)
	var response APIResponse
	if err := c.getJSON(ctx, path, &response); err != nil {
		return err
	}

	if response.Code != OK {
		return fmt.Errorf("leader check failed: %s", response.Message)
	}

	return nil
}

// Status retrieves service status
func (c *Client) Status(ctx context.Context, ) (map[string]interface{}, error) {
	path := "/api/status"
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

	var status map[string]interface{}
	if err := json.Unmarshal(jsonData, &status); err != nil {
		return nil, fmt.Errorf("failed to unmarshal status: %w", err)
	}

	return status, nil
}

// Configuration Operations

// ReloadConfiguration reloads orchestrator configuration
func (c *Client) ReloadConfiguration(ctx context.Context, ) error {
	path := "/api/reload-configuration"
	var response APIResponse
	if err := c.getJSON(ctx, path, &response); err != nil {
		return err
	}

	if response.Code != OK {
		return fmt.Errorf("API error: %s", response.Message)
	}

	return nil
}

// Hostname Management

// GetHostnameResolveCache retrieves the hostname resolution cache
func (c *Client) GetHostnameResolveCache(ctx context.Context, ) ([]HostnameResolveCache, error) {
	path := "/api/hostname-resolve-cache"
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

	var cache []HostnameResolveCache
	if err := json.Unmarshal(jsonData, &cache); err != nil {
		return nil, fmt.Errorf("failed to unmarshal hostname cache: %w", err)
	}

	return cache, nil
}

// ResetHostnameResolveCache resets the hostname resolution cache
func (c *Client) ResetHostnameResolveCache(ctx context.Context, ) error {
	path := "/api/reset-hostname-resolve-cache"
	var response APIResponse
	if err := c.getJSON(ctx, path, &response); err != nil {
		return err
	}

	if response.Code != OK {
		return fmt.Errorf("API error: %s", response.Message)
	}

	return nil
}

// DeregisterHostnameUnresolve deregisters hostname unresolve mapping
func (c *Client) DeregisterHostnameUnresolve(ctx context.Context, instanceKey InstanceKey) error {
	path := fmt.Sprintf("/api/deregister-hostname-unresolve/%s/%d",
		instanceKey.Hostname, instanceKey.Port)
	var response APIResponse
	if err := c.getJSON(ctx, path, &response); err != nil {
		return err
	}

	if response.Code != OK {
		return fmt.Errorf("API error: %s", response.Message)
	}

	return nil
}

// RegisterHostnameUnresolve registers hostname unresolve mapping
func (c *Client) RegisterHostnameUnresolve(ctx context.Context, instanceKey InstanceKey, virtualName string) error {
	path := fmt.Sprintf("/api/register-hostname-unresolve/%s/%d/%s",
		instanceKey.Hostname, instanceKey.Port, virtualName)
	var response APIResponse
	if err := c.getJSON(ctx, path, &response); err != nil {
		return err
	}

	if response.Code != OK {
		return fmt.Errorf("API error: %s", response.Message)
	}

	return nil
}

// Bulk Operations

// BulkPromotionRules retrieves bulk promotion rules
func (c *Client) BulkPromotionRules(ctx context.Context, ) (map[string]CandidatePromotionRule, error) {
	path := "/api/bulk-promotion-rules"
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

	var rules map[string]CandidatePromotionRule
	if err := json.Unmarshal(jsonData, &rules); err != nil {
		return nil, fmt.Errorf("failed to unmarshal promotion rules: %w", err)
	}

	return rules, nil
}
