package client

import (
	"context"
	"encoding/json"
	"fmt"
)

// Monitoring and Metrics Operations

// Discovery Metrics

// GetDiscoveryMetricsRaw retrieves raw discovery metrics for the last N seconds
func (c *Client) GetDiscoveryMetricsRaw(ctx context.Context, seconds int) ([]DiscoveryMetric, error) {
	path := fmt.Sprintf("/api/discovery-metrics-raw/%d", seconds)
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

	var metrics []DiscoveryMetric
	if err := json.Unmarshal(jsonData, &metrics); err != nil {
		return nil, fmt.Errorf("failed to unmarshal discovery metrics: %w", err)
	}

	return metrics, nil
}

// GetDiscoveryMetricsAggregated retrieves aggregated discovery metrics for the last N seconds
func (c *Client) GetDiscoveryMetricsAggregated(ctx context.Context, seconds int) ([]DiscoveryMetric, error) {
	path := fmt.Sprintf("/api/discovery-metrics-aggregated/%d", seconds)
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

	var metrics []DiscoveryMetric
	if err := json.Unmarshal(jsonData, &metrics); err != nil {
		return nil, fmt.Errorf("failed to unmarshal discovery metrics: %w", err)
	}

	return metrics, nil
}

// Discovery Queue Metrics

// GetDiscoveryQueueMetricsRaw retrieves raw discovery queue metrics for the last N seconds
func (c *Client) GetDiscoveryQueueMetricsRaw(ctx context.Context, seconds int) ([]DiscoveryQueueMetric, error) {
	path := fmt.Sprintf("/api/discovery-queue-metrics-raw/%d", seconds)
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

	var metrics []DiscoveryQueueMetric
	if err := json.Unmarshal(jsonData, &metrics); err != nil {
		return nil, fmt.Errorf("failed to unmarshal queue metrics: %w", err)
	}

	return metrics, nil
}

// GetDiscoveryQueueMetricsAggregated retrieves aggregated discovery queue metrics for the last N seconds
func (c *Client) GetDiscoveryQueueMetricsAggregated(ctx context.Context, seconds int) ([]DiscoveryQueueMetric, error) {
	path := fmt.Sprintf("/api/discovery-queue-metrics-aggregated/%d", seconds)
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

	var metrics []DiscoveryQueueMetric
	if err := json.Unmarshal(jsonData, &metrics); err != nil {
		return nil, fmt.Errorf("failed to unmarshal queue metrics: %w", err)
	}

	return metrics, nil
}

// GetDiscoveryQueueMetricsRawByQueue retrieves raw queue metrics for a specific queue
func (c *Client) GetDiscoveryQueueMetricsRawByQueue(ctx context.Context, queue string, seconds int) ([]DiscoveryQueueMetric, error) {
	path := fmt.Sprintf("/api/discovery-queue-metrics-raw/%s/%d", queue, seconds)
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

	var metrics []DiscoveryQueueMetric
	if err := json.Unmarshal(jsonData, &metrics); err != nil {
		return nil, fmt.Errorf("failed to unmarshal queue metrics: %w", err)
	}

	return metrics, nil
}

// GetDiscoveryQueueMetricsAggregatedByQueue retrieves aggregated queue metrics for a specific queue
func (c *Client) GetDiscoveryQueueMetricsAggregatedByQueue(ctx context.Context, queue string, seconds int) ([]DiscoveryQueueMetric, error) {
	path := fmt.Sprintf("/api/discovery-queue-metrics-aggregated/%s/%d", queue, seconds)
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

	var metrics []DiscoveryQueueMetric
	if err := json.Unmarshal(jsonData, &metrics); err != nil {
		return nil, fmt.Errorf("failed to unmarshal queue metrics: %w", err)
	}

	return metrics, nil
}

// Backend Query Metrics

// GetBackendQueryMetricsRaw retrieves raw backend query metrics for the last N seconds
func (c *Client) GetBackendQueryMetricsRaw(ctx context.Context, seconds int) ([]BackendQueryMetric, error) {
	path := fmt.Sprintf("/api/backend-query-metrics-raw/%d", seconds)
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

	var metrics []BackendQueryMetric
	if err := json.Unmarshal(jsonData, &metrics); err != nil {
		return nil, fmt.Errorf("failed to unmarshal backend metrics: %w", err)
	}

	return metrics, nil
}

// GetBackendQueryMetricsAggregated retrieves aggregated backend query metrics for the last N seconds
func (c *Client) GetBackendQueryMetricsAggregated(ctx context.Context, seconds int) ([]BackendQueryMetric, error) {
	path := fmt.Sprintf("/api/backend-query-metrics-aggregated/%d", seconds)
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

	var metrics []BackendQueryMetric
	if err := json.Unmarshal(jsonData, &metrics); err != nil {
		return nil, fmt.Errorf("failed to unmarshal backend metrics: %w", err)
	}

	return metrics, nil
}

// Write Buffer Metrics

// GetWriteBufferMetricsRaw retrieves raw write buffer metrics for the last N seconds
func (c *Client) GetWriteBufferMetricsRaw(ctx context.Context, seconds int) ([]WriteBufferMetric, error) {
	path := fmt.Sprintf("/api/write-buffer-metrics-raw/%d", seconds)
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

	var metrics []WriteBufferMetric
	if err := json.Unmarshal(jsonData, &metrics); err != nil {
		return nil, fmt.Errorf("failed to unmarshal write buffer metrics: %w", err)
	}

	return metrics, nil
}

// GetWriteBufferMetricsAggregated retrieves aggregated write buffer metrics for the last N seconds
func (c *Client) GetWriteBufferMetricsAggregated(ctx context.Context, seconds int) ([]WriteBufferMetric, error) {
	path := fmt.Sprintf("/api/write-buffer-metrics-aggregated/%d", seconds)
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

	var metrics []WriteBufferMetric
	if err := json.Unmarshal(jsonData, &metrics); err != nil {
		return nil, fmt.Errorf("failed to unmarshal write buffer metrics: %w", err)
	}

	return metrics, nil
}
