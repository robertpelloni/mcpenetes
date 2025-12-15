package search

import (
	"fmt"

	"github.com/tuannvm/mcpenetes/internal/config"
	"github.com/tuannvm/mcpenetes/internal/registry"
)

// AddServerToMCPConfig adds a server to the mcp.json configuration.
// It attempts to construct a basic configuration.
func AddServerToMCPConfig(serverID string, serverData *registry.ServerData) error {
	mcpCfg, err := config.LoadMCPConfig()
	if err != nil {
		return fmt.Errorf("failed to load mcp config: %w", err)
	}

	// Check if already exists
	if _, exists := mcpCfg.MCPServers[serverID]; exists {
		// Already exists, maybe update? For now, let's just log and return
		return fmt.Errorf("server '%s' already exists in mcp.json", serverID)
	}

	// Construct new server config
	// Since we don't have the exact 'command' and 'args' from the registry (it only gives repo URL),
	// we will create a placeholder or "manual configuration required" entry.
	// Users often use 'npx' or 'uvx' for these.

	newServer := config.MCPServer{
		// Default to npx if it looks like a package, but we can't be sure.
		// Safe bet: disabled = true, and a helpful comment in args/command?
		// But we can't put comments in JSON.
		// Let's assume 'npx -y <package-name>' if the name looks like a package?
		// Often serverID is the package name.
		Command: "npx",
		Args:    []string{"-y", serverID}, // Optimistic guess
		Env:     make(map[string]string),
		Disabled: false, // Enable by default so users see it, even if it fails
	}

	// If repository URL is available, maybe we can be smarter?
	// For now, simple npx fallback.

	// Add to config
	mcpCfg.MCPServers[serverID] = newServer

	// Save
	if err := config.SaveMCPConfig(mcpCfg); err != nil {
		return fmt.Errorf("failed to save mcp config: %w", err)
	}

	return nil
}
