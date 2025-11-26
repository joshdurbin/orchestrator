package client

import (
	"encoding/json"
	"fmt"
)

// convertToInstance converts an interface{} to an Instance
func convertToInstance(data interface{}) (*Instance, error) {
	if data == nil {
		return nil, fmt.Errorf("nil data")
	}

	// Marshal and unmarshal to convert map to struct
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data: %w", err)
	}

	var instance Instance
	if err := json.Unmarshal(jsonData, &instance); err != nil {
		return nil, fmt.Errorf("failed to unmarshal to Instance: %w", err)
	}

	return &instance, nil
}

// convertToInstanceSlice converts an interface{} to a slice of Instances
func convertToInstanceSlice(data interface{}) ([]Instance, error) {
	if data == nil {
		return []Instance{}, nil
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data: %w", err)
	}

	var instances []Instance
	if err := json.Unmarshal(jsonData, &instances); err != nil {
		return nil, fmt.Errorf("failed to unmarshal to []Instance: %w", err)
	}

	return instances, nil
}

// convertToClusterInfo converts an interface{} to a ClusterInfo
func convertToClusterInfo(data interface{}) (*ClusterInfo, error) {
	if data == nil {
		return nil, fmt.Errorf("nil data")
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data: %w", err)
	}

	var clusterInfo ClusterInfo
	if err := json.Unmarshal(jsonData, &clusterInfo); err != nil {
		return nil, fmt.Errorf("failed to unmarshal to ClusterInfo: %w", err)
	}

	return &clusterInfo, nil
}

// convertToClusterInfoSlice converts an interface{} to a slice of ClusterInfo
func convertToClusterInfoSlice(data interface{}) ([]ClusterInfo, error) {
	if data == nil {
		return []ClusterInfo{}, nil
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data: %w", err)
	}

	var clusterInfos []ClusterInfo
	if err := json.Unmarshal(jsonData, &clusterInfos); err != nil {
		return nil, fmt.Errorf("failed to unmarshal to []ClusterInfo: %w", err)
	}

	return clusterInfos, nil
}

// convertToReplicationAnalysis converts an interface{} to a ReplicationAnalysis
func convertToReplicationAnalysis(data interface{}) (*ReplicationAnalysis, error) {
	if data == nil {
		return nil, fmt.Errorf("nil data")
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data: %w", err)
	}

	var analysis ReplicationAnalysis
	if err := json.Unmarshal(jsonData, &analysis); err != nil {
		return nil, fmt.Errorf("failed to unmarshal to ReplicationAnalysis: %w", err)
	}

	return &analysis, nil
}

// convertToReplicationAnalysisSlice converts an interface{} to a slice of ReplicationAnalysis
func convertToReplicationAnalysisSlice(data interface{}) ([]ReplicationAnalysis, error) {
	if data == nil {
		return []ReplicationAnalysis{}, nil
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data: %w", err)
	}

	var analyses []ReplicationAnalysis
	if err := json.Unmarshal(jsonData, &analyses); err != nil {
		return nil, fmt.Errorf("failed to unmarshal to []ReplicationAnalysis: %w", err)
	}

	return analyses, nil
}

// convertToTopologyRecovery converts an interface{} to a TopologyRecovery
func convertToTopologyRecovery(data interface{}) (*TopologyRecovery, error) {
	if data == nil {
		return nil, fmt.Errorf("nil data")
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data: %w", err)
	}

	var recovery TopologyRecovery
	if err := json.Unmarshal(jsonData, &recovery); err != nil {
		return nil, fmt.Errorf("failed to unmarshal to TopologyRecovery: %w", err)
	}

	return &recovery, nil
}

// convertToTopologyRecoverySlice converts an interface{} to a slice of TopologyRecovery
func convertToTopologyRecoverySlice(data interface{}) ([]TopologyRecovery, error) {
	if data == nil {
		return []TopologyRecovery{}, nil
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data: %w", err)
	}

	var recoveries []TopologyRecovery
	if err := json.Unmarshal(jsonData, &recoveries); err != nil {
		return nil, fmt.Errorf("failed to unmarshal to []TopologyRecovery: %w", err)
	}

	return recoveries, nil
}

// convertToMaintenance converts an interface{} to a Maintenance
func convertToMaintenance(data interface{}) (*Maintenance, error) {
	if data == nil {
		return nil, fmt.Errorf("nil data")
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data: %w", err)
	}

	var maintenance Maintenance
	if err := json.Unmarshal(jsonData, &maintenance); err != nil {
		return nil, fmt.Errorf("failed to unmarshal to Maintenance: %w", err)
	}

	return &maintenance, nil
}

// convertToMaintenanceSlice converts an interface{} to a slice of Maintenance
func convertToMaintenanceSlice(data interface{}) ([]Maintenance, error) {
	if data == nil {
		return []Maintenance{}, nil
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data: %w", err)
	}

	var maintenances []Maintenance
	if err := json.Unmarshal(jsonData, &maintenances); err != nil {
		return nil, fmt.Errorf("failed to unmarshal to []Maintenance: %w", err)
	}

	return maintenances, nil
}

// convertToAgent converts an interface{} to an Agent
func convertToAgent(data interface{}) (*Agent, error) {
	if data == nil {
		return nil, fmt.Errorf("nil data")
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data: %w", err)
	}

	var agent Agent
	if err := json.Unmarshal(jsonData, &agent); err != nil {
		return nil, fmt.Errorf("failed to unmarshal to Agent: %w", err)
	}

	return &agent, nil
}

// convertToAgentSlice converts an interface{} to a slice of Agents
func convertToAgentSlice(data interface{}) ([]Agent, error) {
	if data == nil {
		return []Agent{}, nil
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data: %w", err)
	}

	var agents []Agent
	if err := json.Unmarshal(jsonData, &agents); err != nil {
		return nil, fmt.Errorf("failed to unmarshal to []Agent: %w", err)
	}

	return agents, nil
}

// convertToAgentSeed converts an interface{} to an AgentSeed
func convertToAgentSeed(data interface{}) (*AgentSeed, error) {
	if data == nil {
		return nil, fmt.Errorf("nil data")
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data: %w", err)
	}

	var seed AgentSeed
	if err := json.Unmarshal(jsonData, &seed); err != nil {
		return nil, fmt.Errorf("failed to unmarshal to AgentSeed: %w", err)
	}

	return &seed, nil
}

// convertToAgentSeedSlice converts an interface{} to a slice of AgentSeeds
func convertToAgentSeedSlice(data interface{}) ([]AgentSeed, error) {
	if data == nil {
		return []AgentSeed{}, nil
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data: %w", err)
	}

	var seeds []AgentSeed
	if err := json.Unmarshal(jsonData, &seeds); err != nil {
		return nil, fmt.Errorf("failed to unmarshal to []AgentSeed: %w", err)
	}

	return seeds, nil
}

// convertToAgentSeedStateSlice converts an interface{} to a slice of AgentSeedStates
func convertToAgentSeedStateSlice(data interface{}) ([]AgentSeedState, error) {
	if data == nil {
		return []AgentSeedState{}, nil
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data: %w", err)
	}

	var states []AgentSeedState
	if err := json.Unmarshal(jsonData, &states); err != nil {
		return nil, fmt.Errorf("failed to unmarshal to []AgentSeedState: %w", err)
	}

	return states, nil
}

// convertToAuditEntrySlice converts an interface{} to a slice of AuditEntries
func convertToAuditEntrySlice(data interface{}) ([]AuditEntry, error) {
	if data == nil {
		return []AuditEntry{}, nil
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data: %w", err)
	}

	var entries []AuditEntry
	if err := json.Unmarshal(jsonData, &entries); err != nil {
		return nil, fmt.Errorf("failed to unmarshal to []AuditEntry: %w", err)
	}

	return entries, nil
}

// convertToRecoveryStepSlice converts an interface{} to a slice of RecoverySteps
func convertToRecoveryStepSlice(data interface{}) ([]RecoveryStep, error) {
	if data == nil {
		return []RecoveryStep{}, nil
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data: %w", err)
	}

	var steps []RecoveryStep
	if err := json.Unmarshal(jsonData, &steps); err != nil {
		return nil, fmt.Errorf("failed to unmarshal to []RecoveryStep: %w", err)
	}

	return steps, nil
}

// convertToStringSlice converts an interface{} to a slice of strings
func convertToStringSlice(data interface{}) ([]string, error) {
	if data == nil {
		return []string{}, nil
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data: %w", err)
	}

	var strings []string
	if err := json.Unmarshal(jsonData, &strings); err != nil {
		return nil, fmt.Errorf("failed to unmarshal to []string: %w", err)
	}

	return strings, nil
}

// convertToPoolInstancesMap converts an interface{} to a PoolInstancesMap
func convertToPoolInstancesMap(data interface{}) (PoolInstancesMap, error) {
	if data == nil {
		return PoolInstancesMap{}, nil
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data: %w", err)
	}

	var poolMap PoolInstancesMap
	if err := json.Unmarshal(jsonData, &poolMap); err != nil {
		return nil, fmt.Errorf("failed to unmarshal to PoolInstancesMap: %w", err)
	}

	return poolMap, nil
}

// convertToBool converts an interface{} to a bool
func convertToBool(data interface{}) (bool, error) {
	if data == nil {
		return false, fmt.Errorf("nil data")
	}

	if b, ok := data.(bool); ok {
		return b, nil
	}

	return false, fmt.Errorf("data is not a bool")
}
