package search

import (
	"fmt"

	"github.com/tuannvm/mcpenetes/internal/config"
	"github.com/tuannvm/mcpenetes/internal/registry"
)

// AddServerToMCPConfig adds a server to the mcp.json configuration.
// It attempts to construct a basic configuration.
// If configOverride is provided, it uses that instead of the default.
func AddServerToMCPConfig(serverID string, serverData *registry.ServerData, configOverride *config.MCPServer) error {
	mcpCfg, err := config.LoadMCPConfig()
	if err != nil {
		return fmt.Errorf("failed to load mcp config: %w", err)
	}

	// Check if already exists?
	// If configOverride is provided (from UI install wizard), we assume update/overwrite is intentional
	// If not provided (from CLI search), we check existence to be safe?
	// The CLI flow logs success if added. Let's stick to overwrite or error if exists?
	// The original implementation errored. Let's keep that for safety unless forced?
	// Actually, for "Install" workflow, overwrite might be expected if user confirms.
	// But let's error if no override provided and key exists.

	if configOverride == nil {
		if _, exists := mcpCfg.MCPServers[serverID]; exists {
			return fmt.Errorf("server '%s' already exists in mcp.json", serverID)
		}
	}

	var newServer config.MCPServer

	if configOverride != nil {
		newServer = *configOverride
	} else {
		// Construct default server config
		newServer = config.MCPServer{
			Command: "npx",
			Args:    []string{"-y", serverID}, // Optimistic guess
			Env:     make(map[string]string),
			Disabled: false,
		}
	}

	// Add to config
	mcpCfg.MCPServers[serverID] = newServer

	// Save
	if err := config.SaveMCPConfig(mcpCfg); err != nil {
		return fmt.Errorf("failed to save mcp config: %w", err)
	}

	return nil
}
