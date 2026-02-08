package core

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/tailscale/hujson"
	"github.com/tuannvm/mcpenetes/internal/config"
)

// ImportConfig imports MCP configuration from a JSON string and merges it into the current config.
// Returns the number of servers added/updated.
func (m *Manager) ImportConfig(jsonContent string) (int, error) {
	if strings.TrimSpace(jsonContent) == "" {
		return 0, fmt.Errorf("configuration content is empty")
	}

	// Use hujson to standardize potentially commented JSON
	standardized, err := hujson.Standardize([]byte(jsonContent))
	if err != nil {
		return 0, fmt.Errorf("failed to standardize JSON content: %w", err)
	}

	// Parse content as JSON map
	var importedData map[string]interface{}
	err = json.Unmarshal(standardized, &importedData)
	if err != nil {
		return 0, fmt.Errorf("failed to parse content as JSON: %w", err)
	}

	// Extract mcpServers
	mcpServersRaw, ok := importedData["mcpServers"]
	if !ok {
		return 0, fmt.Errorf("content does not contain 'mcpServers' key")
	}

	// Convert mcpServers back to JSON to parse into our struct
	mcpServersJSON, err := json.Marshal(mcpServersRaw)
	if err != nil {
		return 0, fmt.Errorf("failed to process mcpServers data: %w", err)
	}

	var tempConfig config.MCPConfig
	// Create a temporary wrapper to unmarshal just the servers
	wrapperJSON := fmt.Sprintf(`{"mcpServers": %s}`, string(mcpServersJSON))
	err = json.Unmarshal([]byte(wrapperJSON), &tempConfig)
	if err != nil {
		return 0, fmt.Errorf("failed to parse mcpServers config: %w", err)
	}

	// Load existing config (or use m.MCPConfig which should be up to date)
	// We should probably reload from disk to be safe, but m.MCPConfig is passed in.
	// Let's rely on m.MCPConfig but we must ensure it's saved.

	// Actually, Manager has *config.MCPConfig.
	if m.MCPConfig.MCPServers == nil {
		m.MCPConfig.MCPServers = make(map[string]config.MCPServer)
	}

	count := 0
	for name, server := range tempConfig.MCPServers {
		m.MCPConfig.MCPServers[name] = server
		count++
	}

	// Save the updated config
	err = config.SaveMCPConfig(m.MCPConfig)
	if err != nil {
		return 0, fmt.Errorf("failed to save merged configuration: %w", err)
	}

	return count, nil
}
