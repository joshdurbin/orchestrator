package client

import (
	"context"
	"encoding/json"
	"fmt"
)

// Tagging and Search Operations

// GetTaggedInstances retrieves all tagged instances
func (c *Client) GetTaggedInstances(ctx context.Context, ) ([]Instance, error) {
	path := "/api/tagged"
	return c.executeTopologyOperationMulti(ctx, path)
}

// GetInstanceTags retrieves tags for a specific instance
func (c *Client) GetInstanceTags(ctx context.Context, instanceKey InstanceKey) ([]Tag, error) {
	path := fmt.Sprintf("/api/tags/%s/%d",
		instanceKey.Hostname, instanceKey.Port)
	var response APIResponse
	if err := c.getJSON(ctx, path, &response); err != nil {
		return nil, err
	}

	if response.Code != OK {
		return nil, fmt.Errorf("API error: %s", response.Message)
	}

	// Convert to Tag slice
	jsonData, err := json.Marshal(response.Details)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data: %w", err)
	}

	var tags []Tag
	if err := json.Unmarshal(jsonData, &tags); err != nil {
		return nil, fmt.Errorf("failed to unmarshal tags: %w", err)
	}

	return tags, nil
}

// GetInstanceTagValue retrieves all tag values for an instance
func (c *Client) GetInstanceTagValue(ctx context.Context, instanceKey InstanceKey) (map[string]string, error) {
	path := fmt.Sprintf("/api/tag-value/%s/%d",
		instanceKey.Hostname, instanceKey.Port)
	var response APIResponse
	if err := c.getJSON(ctx, path, &response); err != nil {
		return nil, err
	}

	if response.Code != OK {
		return nil, fmt.Errorf("API error: %s", response.Message)
	}

	// Convert to map
	jsonData, err := json.Marshal(response.Details)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data: %w", err)
	}

	var tagMap map[string]string
	if err := json.Unmarshal(jsonData, &tagMap); err != nil {
		return nil, fmt.Errorf("failed to unmarshal tag values: %w", err)
	}

	return tagMap, nil
}

// GetInstanceSpecificTagValue retrieves a specific tag value for an instance
func (c *Client) GetInstanceSpecificTagValue(ctx context.Context, instanceKey InstanceKey, tagName string) (string, error) {
	path := fmt.Sprintf("/api/tag-value/%s/%d/%s",
		instanceKey.Hostname, instanceKey.Port, tagName)
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

// GetInstanceTag retrieves tags for an instance (list format)
func (c *Client) GetInstanceTag(ctx context.Context, instanceKey InstanceKey) ([]Tag, error) {
	path := fmt.Sprintf("/api/tag/%s/%d",
		instanceKey.Hostname, instanceKey.Port)
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

	var tags []Tag
	if err := json.Unmarshal(jsonData, &tags); err != nil {
		return nil, fmt.Errorf("failed to unmarshal tags: %w", err)
	}

	return tags, nil
}

// TagInstance sets a tag on an instance
func (c *Client) TagInstance(ctx context.Context, instanceKey InstanceKey, tagName, tagValue string) (*Instance, error) {
	path := fmt.Sprintf("/api/tag/%s/%d/%s/%s",
		instanceKey.Hostname, instanceKey.Port, tagName, tagValue)
	return c.executeTopologyOperation(ctx, path)
}

// UntagInstance removes all tags from an instance
func (c *Client) UntagInstance(ctx context.Context, instanceKey InstanceKey) (*Instance, error) {
	path := fmt.Sprintf("/api/untag/%s/%d",
		instanceKey.Hostname, instanceKey.Port)
	return c.executeTopologyOperation(ctx, path)
}

// UntagInstanceSpecific removes a specific tag from an instance
func (c *Client) UntagInstanceSpecific(ctx context.Context, instanceKey InstanceKey, tagName string) (*Instance, error) {
	path := fmt.Sprintf("/api/untag/%s/%d/%s",
		instanceKey.Hostname, instanceKey.Port, tagName)
	return c.executeTopologyOperation(ctx, path)
}

// UntagAll removes all instances with a specific tag
func (c *Client) UntagAll(ctx context.Context, ) ([]Instance, error) {
	path := "/api/untag-all"
	return c.executeTopologyOperationMulti(ctx, path)
}

// UntagAllSpecific removes a specific tag from all instances
func (c *Client) UntagAllSpecific(ctx context.Context, tagName, tagValue string) ([]Instance, error) {
	path := fmt.Sprintf("/api/untag-all/%s/%s", tagName, tagValue)
	return c.executeTopologyOperationMulti(ctx, path)
}
