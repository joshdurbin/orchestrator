package client

import (
	"context"
	"fmt"
)

// Agent Operations

// GetAgents retrieves all known agents
func (c *Client) GetAgents(ctx context.Context, ) ([]Agent, error) {
	path := "/api/agents"
	var response APIResponse
	if err := c.getJSON(ctx, path, &response); err != nil {
		return nil, err
	}

	if response.Code != OK {
		return nil, fmt.Errorf("API error: %s", response.Message)
	}

	return convertToAgentSlice(response.Details)
}

// GetAgent retrieves information about a specific agent
func (c *Client) GetAgent(ctx context.Context, host string) (*Agent, error) {
	path := fmt.Sprintf("/api/agent/%s", host)
	var response APIResponse
	if err := c.getJSON(ctx, path, &response); err != nil {
		return nil, err
	}

	if response.Code != OK {
		return nil, fmt.Errorf("API error: %s", response.Message)
	}

	return convertToAgent(response.Details)
}

// AgentUnmount unmounts an agent's volume
func (c *Client) AgentUnmount(ctx context.Context, host string) error {
	path := fmt.Sprintf("/api/agent-umount/%s", host)
	return c.executeAgentCommand(ctx, path)
}

// AgentMount mounts an agent's logical volume
func (c *Client) AgentMount(ctx context.Context, host string) error {
	path := fmt.Sprintf("/api/agent-mount/%s", host)
	return c.executeAgentCommand(ctx, path)
}

// AgentCreateSnapshot creates a snapshot on an agent
func (c *Client) AgentCreateSnapshot(ctx context.Context, host string) error {
	path := fmt.Sprintf("/api/agent-create-snapshot/%s", host)
	return c.executeAgentCommand(ctx, path)
}

// AgentRemoveLV removes a logical volume from an agent
func (c *Client) AgentRemoveLV(ctx context.Context, host string) error {
	path := fmt.Sprintf("/api/agent-removelv/%s", host)
	return c.executeAgentCommand(ctx, path)
}

// AgentMySQLStop stops MySQL on an agent
func (c *Client) AgentMySQLStop(ctx context.Context, host string) error {
	path := fmt.Sprintf("/api/agent-mysql-stop/%s", host)
	return c.executeAgentCommand(ctx, path)
}

// AgentMySQLStart starts MySQL on an agent
func (c *Client) AgentMySQLStart(ctx context.Context, host string) error {
	path := fmt.Sprintf("/api/agent-mysql-start/%s", host)
	return c.executeAgentCommand(ctx, path)
}

// AgentSeed initiates a seed operation from source host to target host
func (c *Client) AgentSeed(ctx context.Context, targetHost, sourceHost string) (*AgentSeed, error) {
	path := fmt.Sprintf("/api/agent-seed/%s/%s", targetHost, sourceHost)
	var response APIResponse
	if err := c.getJSON(ctx, path, &response); err != nil {
		return nil, err
	}

	if response.Code != OK {
		return nil, fmt.Errorf("API error: %s", response.Message)
	}

	return convertToAgentSeed(response.Details)
}

// GetAgentActiveSeeds retrieves active seed operations for a host
func (c *Client) GetAgentActiveSeeds(ctx context.Context, host string) ([]AgentSeed, error) {
	path := fmt.Sprintf("/api/agent-active-seeds/%s", host)
	var response APIResponse
	if err := c.getJSON(ctx, path, &response); err != nil {
		return nil, err
	}

	if response.Code != OK {
		return nil, fmt.Errorf("API error: %s", response.Message)
	}

	return convertToAgentSeedSlice(response.Details)
}

// GetAgentRecentSeeds retrieves recent seed operations for a host
func (c *Client) GetAgentRecentSeeds(ctx context.Context, host string) ([]AgentSeed, error) {
	path := fmt.Sprintf("/api/agent-recent-seeds/%s", host)
	var response APIResponse
	if err := c.getJSON(ctx, path, &response); err != nil {
		return nil, err
	}

	if response.Code != OK {
		return nil, fmt.Errorf("API error: %s", response.Message)
	}

	return convertToAgentSeedSlice(response.Details)
}

// GetAgentSeedDetails retrieves details for a specific seed operation
func (c *Client) GetAgentSeedDetails(ctx context.Context, seedID int64) (*AgentSeed, error) {
	path := fmt.Sprintf("/api/agent-seed-details/%d", seedID)
	var response APIResponse
	if err := c.getJSON(ctx, path, &response); err != nil {
		return nil, err
	}

	if response.Code != OK {
		return nil, fmt.Errorf("API error: %s", response.Message)
	}

	return convertToAgentSeed(response.Details)
}

// GetAgentSeedStates retrieves states for a specific seed operation
func (c *Client) GetAgentSeedStates(ctx context.Context, seedID int64) ([]AgentSeedState, error) {
	path := fmt.Sprintf("/api/agent-seed-states/%d", seedID)
	var response APIResponse
	if err := c.getJSON(ctx, path, &response); err != nil {
		return nil, err
	}

	if response.Code != OK {
		return nil, fmt.Errorf("API error: %s", response.Message)
	}

	return convertToAgentSeedStateSlice(response.Details)
}

// AbortAgentSeed aborts a seed operation
func (c *Client) AbortAgentSeed(ctx context.Context, seedID int64) error {
	path := fmt.Sprintf("/api/agent-abort-seed/%d", seedID)
	return c.executeAgentCommand(ctx, path)
}

// AgentCustomCommand executes a custom command on an agent
func (c *Client) AgentCustomCommand(ctx context.Context, host, command string) error {
	path := fmt.Sprintf("/api/agent-custom-command/%s/%s", host, command)
	return c.executeAgentCommand(ctx, path)
}

// GetAllSeeds retrieves all seed operations
func (c *Client) GetAllSeeds(ctx context.Context, ) ([]AgentSeed, error) {
	path := "/api/seeds"
	var response APIResponse
	if err := c.getJSON(ctx, path, &response); err != nil {
		return nil, err
	}

	if response.Code != OK {
		return nil, fmt.Errorf("API error: %s", response.Message)
	}

	return convertToAgentSeedSlice(response.Details)
}

// Helper functions

// executeAgentCommand executes an agent command that doesn't return specific data
func (c *Client) executeAgentCommand(ctx context.Context, path string) error {
	var response APIResponse
	if err := c.getJSON(ctx, path, &response); err != nil {
		return err
	}

	if response.Code != OK {
		return fmt.Errorf("API error: %s", response.Message)
	}

	return nil
}
