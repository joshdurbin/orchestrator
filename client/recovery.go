package client

import (
	"context"
	"encoding/json"
	"fmt"
)

// Recovery and Failure Handling Operations

// Replication Analysis

// GetReplicationAnalysis retrieves replication analysis for all clusters
func (c *Client) GetReplicationAnalysis(ctx context.Context, ) ([]ReplicationAnalysis, error) {
	path := "/api/replication-analysis"
	var response APIResponse
	if err := c.getJSON(ctx, path, &response); err != nil {
		return nil, err
	}

	if response.Code != OK {
		return nil, fmt.Errorf("API error: %s", response.Message)
	}

	return convertToReplicationAnalysisSlice(response.Details)
}

// GetClusterReplicationAnalysis retrieves replication analysis for a specific cluster
func (c *Client) GetClusterReplicationAnalysis(ctx context.Context, clusterName string) ([]ReplicationAnalysis, error) {
	path := fmt.Sprintf("/api/replication-analysis/%s", clusterName)
	var response APIResponse
	if err := c.getJSON(ctx, path, &response); err != nil {
		return nil, err
	}

	if response.Code != OK {
		return nil, fmt.Errorf("API error: %s", response.Message)
	}

	return convertToReplicationAnalysisSlice(response.Details)
}

// GetInstanceReplicationAnalysis retrieves replication analysis for a specific instance
func (c *Client) GetInstanceReplicationAnalysis(ctx context.Context, instanceKey InstanceKey) (*ReplicationAnalysis, error) {
	path := fmt.Sprintf("/api/replication-analysis/instance/%s/%d",
		instanceKey.Hostname, instanceKey.Port)
	var response APIResponse
	if err := c.getJSON(ctx, path, &response); err != nil {
		return nil, err
	}

	if response.Code != OK {
		return nil, fmt.Errorf("API error: %s", response.Message)
	}

	return convertToReplicationAnalysis(response.Details)
}

// GetReplicationAnalysisChangelog retrieves the changelog of replication analysis
func (c *Client) GetReplicationAnalysisChangelog(ctx context.Context, ) ([]ReplicationAnalysis, error) {
	path := "/api/replication-analysis-changelog"
	var response APIResponse
	if err := c.getJSON(ctx, path, &response); err != nil {
		return nil, err
	}

	if response.Code != OK {
		return nil, fmt.Errorf("API error: %s", response.Message)
	}

	return convertToReplicationAnalysisSlice(response.Details)
}

// Recovery Operations

// RecoverInstance initiates recovery for a failed instance
func (c *Client) RecoverInstance(ctx context.Context, instanceKey InstanceKey) (*TopologyRecovery, error) {
	path := fmt.Sprintf("/api/recover/%s/%d",
		instanceKey.Hostname, instanceKey.Port)
	return c.executeRecoveryOperation(ctx, path)
}

// RecoverInstanceWithCandidate initiates recovery with a specific candidate
func (c *Client) RecoverInstanceWithCandidate(ctx context.Context, instanceKey, candidateKey InstanceKey) (*TopologyRecovery, error) {
	path := fmt.Sprintf("/api/recover/%s/%d/%s/%d",
		instanceKey.Hostname, instanceKey.Port,
		candidateKey.Hostname, candidateKey.Port)
	return c.executeRecoveryOperation(ctx, path)
}

// RecoverLite initiates a lite recovery for an instance
func (c *Client) RecoverLite(ctx context.Context, instanceKey InstanceKey) (*TopologyRecovery, error) {
	path := fmt.Sprintf("/api/recover-lite/%s/%d",
		instanceKey.Hostname, instanceKey.Port)
	return c.executeRecoveryOperation(ctx, path)
}

// RecoverLiteWithCandidate initiates a lite recovery with a specific candidate
func (c *Client) RecoverLiteWithCandidate(ctx context.Context, instanceKey, candidateKey InstanceKey) (*TopologyRecovery, error) {
	path := fmt.Sprintf("/api/recover-lite/%s/%d/%s/%d",
		instanceKey.Hostname, instanceKey.Port,
		candidateKey.Hostname, candidateKey.Port)
	return c.executeRecoveryOperation(ctx, path)
}

// Graceful Master Takeover

// GracefulMasterTakeoverInstance performs a graceful master takeover from an instance
func (c *Client) GracefulMasterTakeoverInstance(ctx context.Context, instanceKey InstanceKey) (*TopologyRecovery, error) {
	path := fmt.Sprintf("/api/graceful-master-takeover/%s/%d",
		instanceKey.Hostname, instanceKey.Port)
	return c.executeRecoveryOperation(ctx, path)
}

// GracefulMasterTakeoverCluster performs a graceful master takeover for a cluster
func (c *Client) GracefulMasterTakeoverCluster(ctx context.Context, clusterHint string) (*TopologyRecovery, error) {
	path := fmt.Sprintf("/api/graceful-master-takeover/%s", clusterHint)
	return c.executeRecoveryOperation(ctx, path)
}

// GracefulMasterTakeoverInstanceWithDesignated performs a graceful takeover with designated master
func (c *Client) GracefulMasterTakeoverInstanceWithDesignated(ctx context.Context, instanceKey, designatedKey InstanceKey) (*TopologyRecovery, error) {
	path := fmt.Sprintf("/api/graceful-master-takeover/%s/%d/%s/%d",
		instanceKey.Hostname, instanceKey.Port,
		designatedKey.Hostname, designatedKey.Port)
	return c.executeRecoveryOperation(ctx, path)
}

// GracefulMasterTakeoverClusterWithDesignated performs a graceful cluster takeover with designated master
func (c *Client) GracefulMasterTakeoverClusterWithDesignated(ctx context.Context, clusterHint string, designatedKey InstanceKey) (*TopologyRecovery, error) {
	path := fmt.Sprintf("/api/graceful-master-takeover/%s/%s/%d",
		clusterHint, designatedKey.Hostname, designatedKey.Port)
	return c.executeRecoveryOperation(ctx, path)
}

// Auto Graceful Master Takeover

// GracefulMasterTakeoverAutoInstance performs auto graceful master takeover from an instance
func (c *Client) GracefulMasterTakeoverAutoInstance(ctx context.Context, instanceKey InstanceKey) (*TopologyRecovery, error) {
	path := fmt.Sprintf("/api/graceful-master-takeover-auto/%s/%d",
		instanceKey.Hostname, instanceKey.Port)
	return c.executeRecoveryOperation(ctx, path)
}

// GracefulMasterTakeoverAutoCluster performs auto graceful master takeover for a cluster
func (c *Client) GracefulMasterTakeoverAutoCluster(ctx context.Context, clusterHint string) (*TopologyRecovery, error) {
	path := fmt.Sprintf("/api/graceful-master-takeover-auto/%s", clusterHint)
	return c.executeRecoveryOperation(ctx, path)
}

// Force Master Failover

// ForceMasterFailoverInstance forces a master failover from an instance
func (c *Client) ForceMasterFailoverInstance(ctx context.Context, instanceKey InstanceKey) (*TopologyRecovery, error) {
	path := fmt.Sprintf("/api/force-master-failover/%s/%d",
		instanceKey.Hostname, instanceKey.Port)
	return c.executeRecoveryOperation(ctx, path)
}

// ForceMasterFailoverCluster forces a master failover for a cluster
func (c *Client) ForceMasterFailoverCluster(ctx context.Context, clusterHint string) (*TopologyRecovery, error) {
	path := fmt.Sprintf("/api/force-master-failover/%s", clusterHint)
	return c.executeRecoveryOperation(ctx, path)
}

// ForceMasterTakeoverCluster forces a master takeover for a cluster with designated master
func (c *Client) ForceMasterTakeoverCluster(ctx context.Context, clusterHint string, designatedKey InstanceKey) (*TopologyRecovery, error) {
	path := fmt.Sprintf("/api/force-master-takeover/%s/%s/%d",
		clusterHint, designatedKey.Hostname, designatedKey.Port)
	return c.executeRecoveryOperation(ctx, path)
}

// ForceMasterTakeoverInstance forces a master takeover from an instance with designated master
func (c *Client) ForceMasterTakeoverInstance(ctx context.Context, instanceKey, designatedKey InstanceKey) (*TopologyRecovery, error) {
	path := fmt.Sprintf("/api/force-master-takeover/%s/%d/%s/%d",
		instanceKey.Hostname, instanceKey.Port,
		designatedKey.Hostname, designatedKey.Port)
	return c.executeRecoveryOperation(ctx, path)
}

// Recovery Candidate Registration

// RegisterCandidate registers an instance as a recovery candidate
func (c *Client) RegisterCandidate(ctx context.Context, instanceKey InstanceKey, promotionRule CandidatePromotionRule) (*Instance, error) {
	path := fmt.Sprintf("/api/register-candidate/%s/%d/%s",
		instanceKey.Hostname, instanceKey.Port, promotionRule)
	return c.executeTopologyOperation(ctx, path)
}

// Recovery Filters

// GetAutomatedRecoveryFilters retrieves automated recovery filters
func (c *Client) GetAutomatedRecoveryFilters(ctx context.Context, ) ([]AutomatedRecoveryFilter, error) {
	path := "/api/automated-recovery-filters"
	var response APIResponse
	if err := c.getJSON(ctx, path, &response); err != nil {
		return nil, err
	}

	if response.Code != OK {
		return nil, fmt.Errorf("API error: %s", response.Message)
	}

	// Convert to AutomatedRecoveryFilter slice
	jsonData, err := convertToStringSlice(response.Details)
	if err != nil {
		return nil, err
	}

	filters := make([]AutomatedRecoveryFilter, 0, len(jsonData))
	for _, pattern := range jsonData {
		filters = append(filters, AutomatedRecoveryFilter{Pattern: pattern})
	}

	return filters, nil
}

// Recovery Audit

// AuditFailureDetection retrieves failure detection audit log
func (c *Client) AuditFailureDetection(ctx context.Context, page int) ([]ReplicationAnalysis, error) {
	path := fmt.Sprintf("/api/audit-failure-detection/%d", page)
	var response APIResponse
	if err := c.getJSON(ctx, path, &response); err != nil {
		return nil, err
	}

	if response.Code != OK {
		return nil, fmt.Errorf("API error: %s", response.Message)
	}

	return convertToReplicationAnalysisSlice(response.Details)
}

// AuditFailureDetectionByID retrieves failure detection audit by ID
func (c *Client) AuditFailureDetectionByID(ctx context.Context, id int64) (*ReplicationAnalysis, error) {
	path := fmt.Sprintf("/api/audit-failure-detection/id/%d", id)
	var response APIResponse
	if err := c.getJSON(ctx, path, &response); err != nil {
		return nil, err
	}

	if response.Code != OK {
		return nil, fmt.Errorf("API error: %s", response.Message)
	}

	return convertToReplicationAnalysis(response.Details)
}

// AuditFailureDetectionByAlias retrieves failure detection audit by cluster alias
func (c *Client) AuditFailureDetectionByAlias(ctx context.Context, clusterAlias string, page int) ([]ReplicationAnalysis, error) {
	path := fmt.Sprintf("/api/audit-failure-detection/alias/%s/%d", clusterAlias, page)
	var response APIResponse
	if err := c.getJSON(ctx, path, &response); err != nil {
		return nil, err
	}

	if response.Code != OK {
		return nil, fmt.Errorf("API error: %s", response.Message)
	}

	return convertToReplicationAnalysisSlice(response.Details)
}

// AuditRecovery retrieves recovery audit log
func (c *Client) AuditRecovery(ctx context.Context, page int) ([]TopologyRecovery, error) {
	path := fmt.Sprintf("/api/audit-recovery/%d", page)
	var response APIResponse
	if err := c.getJSON(ctx, path, &response); err != nil {
		return nil, err
	}

	if response.Code != OK {
		return nil, fmt.Errorf("API error: %s", response.Message)
	}

	return convertToTopologyRecoverySlice(response.Details)
}

// AuditRecoveryByID retrieves recovery audit by ID
func (c *Client) AuditRecoveryByID(ctx context.Context, id int64) (*TopologyRecovery, error) {
	path := fmt.Sprintf("/api/audit-recovery/id/%d", id)
	var response APIResponse
	if err := c.getJSON(ctx, path, &response); err != nil {
		return nil, err
	}

	if response.Code != OK {
		return nil, fmt.Errorf("API error: %s", response.Message)
	}

	return convertToTopologyRecovery(response.Details)
}

// AuditRecoveryByUID retrieves recovery audit by UID
func (c *Client) AuditRecoveryByUID(ctx context.Context, uid string) ([]TopologyRecovery, error) {
	path := fmt.Sprintf("/api/audit-recovery/uid/%s", uid)
	var response APIResponse
	if err := c.getJSON(ctx, path, &response); err != nil {
		return nil, err
	}

	if response.Code != OK {
		return nil, fmt.Errorf("API error: %s", response.Message)
	}

	return convertToTopologyRecoverySlice(response.Details)
}

// AuditRecoveryByCluster retrieves recovery audit by cluster
func (c *Client) AuditRecoveryByCluster(ctx context.Context, clusterName string, page int) ([]TopologyRecovery, error) {
	path := fmt.Sprintf("/api/audit-recovery/cluster/%s/%d", clusterName, page)
	var response APIResponse
	if err := c.getJSON(ctx, path, &response); err != nil {
		return nil, err
	}

	if response.Code != OK {
		return nil, fmt.Errorf("API error: %s", response.Message)
	}

	return convertToTopologyRecoverySlice(response.Details)
}

// AuditRecoveryByAlias retrieves recovery audit by cluster alias
func (c *Client) AuditRecoveryByAlias(ctx context.Context, clusterAlias string, page int) ([]TopologyRecovery, error) {
	path := fmt.Sprintf("/api/audit-recovery/alias/%s/%d", clusterAlias, page)
	var response APIResponse
	if err := c.getJSON(ctx, path, &response); err != nil {
		return nil, err
	}

	if response.Code != OK {
		return nil, fmt.Errorf("API error: %s", response.Message)
	}

	return convertToTopologyRecoverySlice(response.Details)
}

// AuditRecoverySteps retrieves recovery steps for a specific recovery
func (c *Client) AuditRecoverySteps(ctx context.Context, uid string) ([]RecoveryStep, error) {
	path := fmt.Sprintf("/api/audit-recovery-steps/%s", uid)
	var response APIResponse
	if err := c.getJSON(ctx, path, &response); err != nil {
		return nil, err
	}

	if response.Code != OK {
		return nil, fmt.Errorf("API error: %s", response.Message)
	}

	return convertToRecoveryStepSlice(response.Details)
}

// Active Recovery

// GetActiveClusterRecovery retrieves active recovery for a cluster
func (c *Client) GetActiveClusterRecovery(ctx context.Context, clusterName string) (*TopologyRecovery, error) {
	path := fmt.Sprintf("/api/active-cluster-recovery/%s", clusterName)
	var response APIResponse
	if err := c.getJSON(ctx, path, &response); err != nil {
		return nil, err
	}

	if response.Code != OK {
		return nil, fmt.Errorf("API error: %s", response.Message)
	}

	return convertToTopologyRecovery(response.Details)
}

// GetRecentlyActiveClusterRecovery retrieves recently active recovery for a cluster
func (c *Client) GetRecentlyActiveClusterRecovery(ctx context.Context, clusterName string) ([]TopologyRecovery, error) {
	path := fmt.Sprintf("/api/recently-active-cluster-recovery/%s", clusterName)
	var response APIResponse
	if err := c.getJSON(ctx, path, &response); err != nil {
		return nil, err
	}

	if response.Code != OK {
		return nil, fmt.Errorf("API error: %s", response.Message)
	}

	return convertToTopologyRecoverySlice(response.Details)
}

// GetRecentlyActiveInstanceRecovery retrieves recently active recovery for an instance
func (c *Client) GetRecentlyActiveInstanceRecovery(ctx context.Context, instanceKey InstanceKey) ([]TopologyRecovery, error) {
	path := fmt.Sprintf("/api/recently-active-instance-recovery/%s/%d",
		instanceKey.Hostname, instanceKey.Port)
	var response APIResponse
	if err := c.getJSON(ctx, path, &response); err != nil {
		return nil, err
	}

	if response.Code != OK {
		return nil, fmt.Errorf("API error: %s", response.Message)
	}

	return convertToTopologyRecoverySlice(response.Details)
}

// Recovery Acknowledgment

// AcknowledgeClusterRecovery acknowledges a cluster recovery
func (c *Client) AcknowledgeClusterRecovery(ctx context.Context, clusterHint, comment string) (*TopologyRecovery, error) {
	path := fmt.Sprintf("/api/ack-recovery/cluster/%s", clusterHint)
	return c.executeRecoveryOperationWithComment(ctx, path, comment)
}

// AcknowledgeClusterRecoveryByAlias acknowledges a cluster recovery by alias
func (c *Client) AcknowledgeClusterRecoveryByAlias(ctx context.Context, clusterAlias, comment string) (*TopologyRecovery, error) {
	path := fmt.Sprintf("/api/ack-recovery/cluster/alias/%s", clusterAlias)
	return c.executeRecoveryOperationWithComment(ctx, path, comment)
}

// AcknowledgeInstanceRecovery acknowledges an instance recovery
func (c *Client) AcknowledgeInstanceRecovery(ctx context.Context, instanceKey InstanceKey, comment string) (*TopologyRecovery, error) {
	path := fmt.Sprintf("/api/ack-recovery/instance/%s/%d",
		instanceKey.Hostname, instanceKey.Port)
	return c.executeRecoveryOperationWithComment(ctx, path, comment)
}

// AcknowledgeRecoveryByID acknowledges a recovery by ID
func (c *Client) AcknowledgeRecoveryByID(ctx context.Context, recoveryID int64, comment string) (*TopologyRecovery, error) {
	path := fmt.Sprintf("/api/ack-recovery/%d", recoveryID)
	return c.executeRecoveryOperationWithComment(ctx, path, comment)
}

// AcknowledgeRecoveryByUID acknowledges a recovery by UID
func (c *Client) AcknowledgeRecoveryByUID(ctx context.Context, uid, comment string) (*TopologyRecovery, error) {
	path := fmt.Sprintf("/api/ack-recovery/uid/%s", uid)
	return c.executeRecoveryOperationWithComment(ctx, path, comment)
}

// AcknowledgeAllRecoveries acknowledges all outstanding recoveries
func (c *Client) AcknowledgeAllRecoveries(ctx context.Context, comment string) ([]TopologyRecovery, error) {
	path := "/api/ack-all-recoveries"
	var response APIResponse
	body := map[string]string{"comment": comment}
	if err := c.postJSON(ctx, path, body, &response); err != nil {
		return nil, err
	}

	if response.Code != OK {
		return nil, fmt.Errorf("API error: %s", response.Message)
	}

	return convertToTopologyRecoverySlice(response.Details)
}

// Blocked Recoveries

// GetBlockedRecoveries retrieves all blocked recoveries
func (c *Client) GetBlockedRecoveries(ctx context.Context, ) ([]BlockedTopologyRecovery, error) {
	path := "/api/blocked-recoveries"
	var response APIResponse
	if err := c.getJSON(ctx, path, &response); err != nil {
		return nil, err
	}

	if response.Code != OK {
		return nil, fmt.Errorf("API error: %s", response.Message)
	}

	// Convert response
	jsonData, err := json.Marshal(response.Details)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data: %w", err)
	}

	var blocked []BlockedTopologyRecovery
	if err := json.Unmarshal(jsonData, &blocked); err != nil {
		return nil, fmt.Errorf("failed to unmarshal blocked recoveries: %w", err)
	}

	return blocked, nil
}

// GetBlockedRecoveriesByCluster retrieves blocked recoveries for a specific cluster
func (c *Client) GetBlockedRecoveriesByCluster(ctx context.Context, clusterName string) ([]BlockedTopologyRecovery, error) {
	path := fmt.Sprintf("/api/blocked-recoveries/cluster/%s", clusterName)
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

	var blocked []BlockedTopologyRecovery
	if err := json.Unmarshal(jsonData, &blocked); err != nil {
		return nil, fmt.Errorf("failed to unmarshal blocked recoveries: %w", err)
	}

	return blocked, nil
}

// Global Recovery Control

// DisableGlobalRecoveries disables global recovery
func (c *Client) DisableGlobalRecoveries(ctx context.Context, ) error {
	path := "/api/disable-global-recoveries"
	var response APIResponse
	if err := c.getJSON(ctx, path, &response); err != nil {
		return err
	}

	if response.Code != OK {
		return fmt.Errorf("API error: %s", response.Message)
	}

	return nil
}

// EnableGlobalRecoveries enables global recovery
func (c *Client) EnableGlobalRecoveries(ctx context.Context, ) error {
	path := "/api/enable-global-recoveries"
	var response APIResponse
	if err := c.getJSON(ctx, path, &response); err != nil {
		return err
	}

	if response.Code != OK {
		return fmt.Errorf("API error: %s", response.Message)
	}

	return nil
}

// CheckGlobalRecoveries checks if global recoveries are enabled
func (c *Client) CheckGlobalRecoveries(ctx context.Context, ) (bool, error) {
	path := "/api/check-global-recoveries"
	var response APIResponse
	if err := c.getJSON(ctx, path, &response); err != nil {
		return false, err
	}

	if response.Code != OK {
		return false, fmt.Errorf("API error: %s", response.Message)
	}

	return convertToBool(response.Details)
}

// Helper functions

// executeRecoveryOperation executes a recovery operation that returns TopologyRecovery
func (c *Client) executeRecoveryOperation(ctx context.Context, path string) (*TopologyRecovery, error) {
	var response APIResponse
	if err := c.getJSON(ctx, path, &response); err != nil {
		return nil, err
	}

	if response.Code != OK {
		return nil, fmt.Errorf("API error: %s", response.Message)
	}

	return convertToTopologyRecovery(response.Details)
}

// executeRecoveryOperationWithComment executes a recovery operation with a comment
func (c *Client) executeRecoveryOperationWithComment(ctx context.Context, path, comment string) (*TopologyRecovery, error) {
	var response APIResponse
	body := map[string]string{"comment": comment}
	if err := c.postJSON(ctx, path, body, &response); err != nil {
		return nil, err
	}

	if response.Code != OK {
		return nil, fmt.Errorf("API error: %s", response.Message)
	}

	return convertToTopologyRecovery(response.Details)
}
