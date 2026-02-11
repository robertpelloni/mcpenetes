package config

// Config represents the structure of config.yaml
type Config struct {
	Version    int               `yaml:"version"`
	Registries []Registry        `yaml:"registries"`
	MCPs       []string          `yaml:"mcps"`
	Clients    map[string]Client `yaml:"clients"`
	Backups    BackupConfig      `yaml:"backups"`
	GlobalEnv  map[string]string `yaml:"global_env,omitempty"`
	Sync       SyncConfig        `yaml:"sync,omitempty"`
}

// SyncConfig defines cloud synchronization settings
type SyncConfig struct {
	GitHubToken string `yaml:"github_token,omitempty"`
	GistID      string `yaml:"gist_id,omitempty"`
	LastSynced  string `yaml:"last_synced,omitempty"` // Timestamp
}

// Registry defines a registry endpoint
type Registry struct {
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
}

// Client defines a target client configuration location
type Client struct {
	ConfigPath string `yaml:"config_path"`
	// Type/Format of the client config, used for translation
	// If empty, defaults will be inferred
	Type string `yaml:"type,omitempty"`
	// Key overrides the default JSON key if set (e.g. "openctx.providers")
	Key string `yaml:"key,omitempty"`
}

// BackupConfig defines backup settings
type BackupConfig struct {
	Path      string `yaml:"path"`
	Retention int    `yaml:"retention"`
}

// MCPConfig represents the structure of mcp.json
type MCPConfig struct {
	MCPServers map[string]MCPServer `json:"mcpServers"`
}

// MCPServer defines the configuration for a single MCP server
// According to the schema, it must have either command or url,
// and can optionally have args and env
type MCPServer struct {
	Command     string            `json:"command,omitempty"`
	Args        []string          `json:"args,omitempty"`
	URL         string            `json:"url,omitempty"`
	Env         map[string]string `json:"env,omitempty"`
	Disabled    bool              `json:"disabled,omitempty"`
	AutoApprove []string          `json:"autoApprove,omitempty"`
}
