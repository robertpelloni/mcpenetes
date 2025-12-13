package core

import (
	"fmt"

	"github.com/tuannvm/mcpenetes/internal/config"
	"github.com/tuannvm/mcpenetes/internal/translator"
)

// Manager orchestrates the application of MCP configurations to clients.
type Manager struct {
	Config    *config.Config
	MCPConfig *config.MCPConfig
	Trans     *translator.Translator
}

// NewManager creates a new Manager instance.
func NewManager(cfg *config.Config, mcpCfg *config.MCPConfig) *Manager {
	return &Manager{
		Config:    cfg,
		MCPConfig: mcpCfg,
		Trans:     translator.NewTranslator(cfg, mcpCfg),
	}
}

// ApplyResult holds the result of an apply operation for a single client.
type ApplyResult struct {
	ClientName string
	Success    bool
	BackupPath string
	Error      error
}

// ApplyToClient applies the current MCP configuration to a specific client.
// It handles backup, translation/application of all servers, and cleanup of obsolete servers.
func (m *Manager) ApplyToClient(clientName string, clientConf config.Client) ApplyResult {
	res := ApplyResult{ClientName: clientName, Success: true}

	// 1. Backup
	backupPath, err := m.Trans.BackupClientConfig(clientName, clientConf)
	if err != nil {
		res.Success = false
		res.Error = fmt.Errorf("backup failed: %w", err)
		return res
	}
	res.BackupPath = backupPath

	// 2. Apply Servers
	for serverName, serverConf := range m.MCPConfig.MCPServers {
		err := m.Trans.TranslateAndApply(clientName, clientConf, serverConf)
		if err != nil {
			res.Success = false
			res.Error = fmt.Errorf("failed to apply server %s: %w", serverName, err)
			return res
		}
	}

	// 3. Clean Obsolete
	err = m.Trans.RemoveClientServers(clientName, clientConf)
	if err != nil {
		res.Success = false
		res.Error = fmt.Errorf("failed to remove obsolete servers: %w", err)
		return res
	}

	return res
}
