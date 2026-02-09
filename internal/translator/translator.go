package translator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/tailscale/hujson"
	"github.com/tuannvm/mcpenetes/internal/client"
	"github.com/tuannvm/mcpenetes/internal/config"
	"github.com/tuannvm/mcpenetes/internal/util"
	"gopkg.in/yaml.v3"
)

// Translator handles backing up and translating MCP configs for clients.
type Translator struct {
	AppConfig *config.Config
	MCPConfig *config.MCPConfig
}

// NewTranslator creates a new Translator instance.
func NewTranslator(appCfg *config.Config, mcpCfg *config.MCPConfig) *Translator {
	return &Translator{
		AppConfig: appCfg,
		MCPConfig: mcpCfg,
	}
}

// getBackupDir returns the backup directory for a client (internal helper)
func (t *Translator) getBackupDir(clientName string) string {
	// For now, it seems all backups go to a single configured path
	// But `DeleteBackup` implementation implied getting it.
	// Re-implementing logic from BackupClientConfig

	// Assuming t.AppConfig.Backups.Path is the root
	path, _ := util.ExpandPath(t.AppConfig.Backups.Path)
	return path
}

// BackupClientConfig creates a timestamped backup of a client's configuration file.
func (t *Translator) BackupClientConfig(clientName string, clientConf config.Client) (string, error) {
	backupDir, err := util.ExpandPath(t.AppConfig.Backups.Path)
	if err != nil {
		return "", fmt.Errorf("failed to expand backup path '%s': %w", t.AppConfig.Backups.Path, err)
	}

	clientConfigPath, err := util.ExpandPath(clientConf.ConfigPath)
	if err != nil {
		return "", fmt.Errorf("failed to expand client config path '%s' for %s: %w", clientConf.ConfigPath, clientName, err)
	}

	// Ensure the main backup directory exists
	if err := os.MkdirAll(backupDir, 0750); err != nil {
		return "", fmt.Errorf("failed to create backup directory '%s': %w", backupDir, err)
	}

	// Check if source file exists
	srcInfo, err := os.Stat(clientConfigPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Source config doesn't exist, nothing to back up
			return "", nil // Not an error, just nothing to do
		}
		return "", fmt.Errorf("failed to stat source config file '%s': %w", clientConfigPath, err)
	}
	if srcInfo.IsDir() {
		return "", fmt.Errorf("source config path '%s' is a directory, not a file", clientConfigPath)
	}

	// Create timestamped backup filename
	timestamp := time.Now().Format("20060102-150405") // YYYYMMDD-HHMMSS
	backupFileName := fmt.Sprintf("%s-%s%s", clientName, timestamp, filepath.Ext(clientConfigPath))
	backupFilePath := filepath.Join(backupDir, backupFileName)

	// Open source file
	srcFile, err := os.Open(clientConfigPath)
	if err != nil {
		return "", fmt.Errorf("failed to open source config file '%s': %w", clientConfigPath, err)
	}
	defer func() {
		_ = srcFile.Close()
	}()

	// Create destination backup file
	dstFile, err := os.Create(backupFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to create backup file '%s': %w", backupFilePath, err)
	}
	defer func() {
		_ = dstFile.Close()
	}()

	// Copy content
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return "", fmt.Errorf("failed to copy config to backup file '%s': %w", backupFilePath, err)
	}

	fmt.Printf("  Backed up '%s' to '%s'\n", clientConfigPath, backupFilePath)

	// TODO: Implement backup retention logic here or separately

	return backupFilePath, nil
}

// DeleteBackup removes a specific backup file
func (t *Translator) DeleteBackup(clientName, filename string) error {
	backupDir := t.getBackupDir(clientName)
	path := filepath.Join(backupDir, filename)

	// Security check: ensure path is within backup dir
	if filepath.Dir(path) != backupDir {
		return fmt.Errorf("invalid backup path")
	}

	return os.Remove(path)
}

// TranslateAndApply translates the selected MCP config and writes it to the client's path.
// If serverConf is nil, it will remove the server from the client's configuration.
func (t *Translator) TranslateAndApply(clientName string, clientConf config.Client, serverConf config.MCPServer) error {
	clientConfigPath, err := util.ExpandPath(clientConf.ConfigPath)
	if err != nil {
		return fmt.Errorf("failed to expand client config path '%s' for %s: %w", clientConf.ConfigPath, clientName, err)
	}

	fmt.Printf("  Translating config for %s ('%s')...\n", clientName, clientConfigPath)

	// Find the server ID (key) from the MCPConfig map by matching the server config
	serverID := ""
	for id, server := range t.MCPConfig.MCPServers {
		// Compare the relevant fields to find a match
		if server.Command == serverConf.Command &&
			server.URL == serverConf.URL &&
			fmt.Sprintf("%v", server.Args) == fmt.Sprintf("%v", serverConf.Args) &&
			fmt.Sprintf("%v", server.Env) == fmt.Sprintf("%v", serverConf.Env) {
			serverID = id
			break
		}
	}

	if serverID == "" {
		// Fallback: generate a server ID based on command or URL
		if serverConf.Command != "" {
			serverID = strings.Split(serverConf.Command, " ")[0]
		} else if serverConf.URL != "" {
			// Extract domain from URL
			parts := strings.Split(strings.TrimPrefix(strings.TrimPrefix(serverConf.URL, "https://"), "http://"), "/")
			if len(parts) > 0 {
				serverID = parts[0]
			} else {
				serverID = "mcp-server"
			}
		} else {
			serverID = "mcp-server"
		}
	}

	var outputData []byte
	formatType := client.ConfigFormatEnum(clientConf.Type)

	// Determine format type if not explicitly set
	if formatType == "" {
		ext := strings.ToLower(filepath.Ext(clientConfigPath))
		switch ext {
		case ".json":
			// Try to guess from client name or default to simple json
			if strings.Contains(clientName, "claude-desktop") {
				formatType = client.FormatClaudeDesktop
			} else if strings.Contains(clientName, "vscode") {
				formatType = client.FormatVSCode
			} else if strings.Contains(strings.ToLower(clientName), "continue") {
				formatType = client.FormatContinue
			} else {
				formatType = client.FormatSimpleJSON
			}
		case ".yaml", ".yml":
			formatType = client.FormatYAML
		case ".toml":
			formatType = client.FormatTOML
		default:
			return fmt.Errorf("unknown config format for client %s", clientName)
		}
	}

	// Helper function to read and parse JSON/JSONC safely
	parseJSONSafe := func(path string, v interface{}) error {
		data, err := os.ReadFile(path)
		if os.IsNotExist(err) {
			return nil // OK, just empty
		}
		if err != nil {
			return err
		}
		if len(data) == 0 {
			return nil
		}

		// Use hujson to standardize JSONC (comments, trailing commas) to standard JSON
		standardized, err := hujson.Standardize(data)
		if err != nil {
			// If we can't standardize it, it's likely invalid JSON/JSONC
			// To be safe, we abort to avoid overwriting a file we can't understand
			return fmt.Errorf("failed to parse existing config file (invalid JSON/JSONC): %w", err)
		}

		return json.Unmarshal(standardized, v)
	}

	// Prepare the server configuration based on client type
	switch formatType {
	case client.FormatClaudeDesktop:
		// Format: {"mcpServers": {"server-id": {...server config...}}}
		var claudeConfig map[string]interface{}

		if err := parseJSONSafe(clientConfigPath, &claudeConfig); err != nil {
			return err
		}

		if claudeConfig == nil {
			claudeConfig = make(map[string]interface{})
			claudeConfig["mcpServers"] = make(map[string]interface{})
		}

		// Check if mcpServers map exists
		mcpServers, ok := claudeConfig["mcpServers"].(map[string]interface{})
		if !ok {
			// Initialize or reset the mcpServers map if it doesn't exist or has wrong type
			mcpServers = make(map[string]interface{})
		}

		serverEntry := t.createServerMap(serverConf)
		mcpServers[serverID] = serverEntry
		claudeConfig["mcpServers"] = mcpServers

		outputData, err = json.MarshalIndent(claudeConfig, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal Claude Desktop config: %w", err)
		}

	case client.FormatSimpleJSON:
		// Format: {"mcpServers": {"server-id": {...}}} (Used by Cursor, Windsurf, Cline, etc.)
		var configMap map[string]interface{}

		if err := parseJSONSafe(clientConfigPath, &configMap); err != nil {
			return err
		}

		if configMap == nil {
			configMap = make(map[string]interface{})
		}

		// Ensure mcpServers exists
		if _, ok := configMap["mcpServers"]; !ok {
			configMap["mcpServers"] = make(map[string]interface{})
		}

		mcpServers, ok := configMap["mcpServers"].(map[string]interface{})
		if !ok {
			mcpServers = make(map[string]interface{})
		}

		serverEntry := t.createServerMap(serverConf)
		mcpServers[serverID] = serverEntry
		configMap["mcpServers"] = mcpServers

		outputData, err = json.MarshalIndent(configMap, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON config: %w", err)
		}

	case client.FormatVSCode:
		// Format: {"mcp": {"servers": {"server-id": {...}}}}
		// OR custom key if clientConf.Key is set
		var vscodeConfig map[string]interface{}

		if err := parseJSONSafe(clientConfigPath, &vscodeConfig); err != nil {
			return err
		}

		if vscodeConfig == nil {
			vscodeConfig = make(map[string]interface{})
		}

		// Determine the root key path
		// Default is "mcp.servers" (nested as "mcp": {"servers": ...})
		// If clientConf.Key is set (e.g. "openctx.providers"), use that.

		// For now, supporting specific cases based on key structure
		if clientConf.Key == "openctx.providers" {
			// Special handling for OpenCtx providers structure
			// "openctx.providers": { "https://...": { "mcp.provider.uri": "file:///..." } }
			// NOTE: This currently only supports the 'modelcontextprotocol' provider type from OpenCtx
			// and assumes we are configuring it to point to a local file/command?
			// Actually, OpenCtx MCP provider usually takes a URL or command.
			// Given the complexity of OpenCtx generic provider, we'll try to map standard MCP config to it if possible,
			// or simply inject it as a raw object if the user knows what they are doing.

			// For simplicity in this iteration, we will treat it as a map where we inject the server config
			// BUT OpenCtx providers usually expect a specific schema.
			// Let's assume for Cody we are injecting into "openctx.providers".
			// However, Cody/OpenCtx usually expects a map of Provider URLs to Config.

			// If we are just injecting a standard MCP server list, this might not work directly with OpenCtx
			// unless we are configuring the "mcp" provider specifically.

			// Fallback: If the user explicitly requested this format, we try to inject it as a key.
			var providers map[string]interface{}
			if existing, ok := vscodeConfig["openctx.providers"].(map[string]interface{}); ok {
				providers = existing
			} else {
				providers = make(map[string]interface{})
			}

			// We can't easily map a list of MCP servers to OpenCtx providers list without a specific adapter.
			// But if the requirement is just to support the key:
			// We will inject the server config under the serverID key in that map.
			serverEntry := t.createServerMap(serverConf)
			providers[serverID] = serverEntry
			vscodeConfig["openctx.providers"] = providers

		} else {
			// Standard VSCode MCP "mcp.servers" logic
			// Get or create mcp object
			var mcpObj map[string]interface{}
			if existingMcpObj, ok := vscodeConfig["mcp"].(map[string]interface{}); ok {
				mcpObj = existingMcpObj
			} else {
				mcpObj = make(map[string]interface{})
			}

			// Get or create servers object
			var mcpServers map[string]interface{}
			if existingServers, ok := mcpObj["servers"].(map[string]interface{}); ok {
				mcpServers = existingServers
			} else {
				mcpServers = make(map[string]interface{})
			}

			serverEntry := t.createServerMap(serverConf)
			// VSCode format explicitly needs env even if empty, usually
			if _, ok := serverEntry["env"]; !ok {
				serverEntry["env"] = make(map[string]string)
			}

			mcpServers[serverID] = serverEntry
			mcpObj["servers"] = mcpServers
			vscodeConfig["mcp"] = mcpObj
		}

		outputData, err = json.MarshalIndent(vscodeConfig, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal VS Code config: %w", err)
		}

	case client.FormatContinue:
		// Continue uses nested structure:
		// "experimental": { "modelContextProtocolServers": [ { "name": "...", "transport": { ... } } ] }
		var continueConfig map[string]interface{}
		if err := parseJSONSafe(clientConfigPath, &continueConfig); err != nil {
			return err
		}
		if continueConfig == nil {
			continueConfig = make(map[string]interface{})
		}

		// Ensure experimental object exists
		if _, ok := continueConfig["experimental"]; !ok {
			continueConfig["experimental"] = make(map[string]interface{})
		}
		experimental := continueConfig["experimental"].(map[string]interface{})

		// Ensure modelContextProtocolServers list exists
		if _, ok := experimental["modelContextProtocolServers"]; !ok {
			experimental["modelContextProtocolServers"] = []interface{}{}
		}

		// Convert existing list to a map for easy updating
		serversList, ok := experimental["modelContextProtocolServers"].([]interface{})
		if !ok {
			// If it's not a list, maybe reset it?
			serversList = []interface{}{}
		}

		// Check if server already exists in list and update it, or append
		found := false
		newServerEntry := map[string]interface{}{
			"name": serverID,
			"transport": map[string]interface{}{
				"type": "stdio", // Defaulting to stdio, check logic if http/sse is needed
				"command": serverConf.Command,
				"args": serverConf.Args,
				"env": serverConf.Env,
			},
		}

		if serverConf.URL != "" {
			// Handle HTTP/SSE transport mapping if applicable
			// Continue supports "type": "sse" with "url"
			// Assuming "url" implies SSE or HTTP
			newServerEntry["transport"] = map[string]interface{}{
				"type": "sse", // Simplification
				"url": serverConf.URL,
			}
		}

		updatedList := []interface{}{}
		for _, s := range serversList {
			sMap, ok := s.(map[string]interface{})
			if !ok { continue }
			if name, ok := sMap["name"].(string); ok && name == serverID {
				updatedList = append(updatedList, newServerEntry)
				found = true
			} else {
				updatedList = append(updatedList, s)
			}
		}

		if !found {
			updatedList = append(updatedList, newServerEntry)
		}

		experimental["modelContextProtocolServers"] = updatedList
		continueConfig["experimental"] = experimental

		outputData, err = json.MarshalIndent(continueConfig, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal Continue config: %w", err)
		}

	case client.FormatYAML:
		var yamlConfig map[string]interface{}

		existingFile, err := os.ReadFile(clientConfigPath)
		if err == nil && len(existingFile) > 0 {
			if err := yaml.Unmarshal(existingFile, &yamlConfig); err != nil {
				// Abort on invalid YAML
				return fmt.Errorf("failed to parse existing YAML config: %w", err)
			}
		}

		if yamlConfig == nil {
			yamlConfig = make(map[string]interface{})
		}

		if _, ok := yamlConfig["mcpServers"]; !ok {
			yamlConfig["mcpServers"] = make(map[string]interface{})
		}

		mcpServers, ok := yamlConfig["mcpServers"].(map[string]interface{})
		if !ok {
			mcpServers = make(map[string]interface{})
		}

		serverEntry := t.createServerMap(serverConf)
		mcpServers[serverID] = serverEntry
		yamlConfig["mcpServers"] = mcpServers

		outputData, err = yaml.Marshal(yamlConfig)
		if err != nil {
			return fmt.Errorf("failed to marshal YAML config: %w", err)
		}

	case client.FormatTOML:
		var tomlConfig map[string]interface{}

		existingFile, err := os.ReadFile(clientConfigPath)
		if err == nil && len(existingFile) > 0 {
			if err := toml.Unmarshal(existingFile, &tomlConfig); err != nil {
				// Abort on invalid TOML
				return fmt.Errorf("failed to parse existing TOML config: %w", err)
			}
		}

		if tomlConfig == nil {
			tomlConfig = make(map[string]interface{})
		}

		if _, ok := tomlConfig["mcpServers"]; !ok {
			tomlConfig["mcpServers"] = make(map[string]interface{})
		}

		mcpServers, ok := tomlConfig["mcpServers"].(map[string]interface{})
		if !ok {
			mcpServers = make(map[string]interface{})
		}

		serverEntry := t.createServerMap(serverConf)
		mcpServers[serverID] = serverEntry
		tomlConfig["mcpServers"] = mcpServers

		buf := new(bytes.Buffer)
		if err := toml.NewEncoder(buf).Encode(tomlConfig); err != nil {
			return fmt.Errorf("failed to marshal updated TOML config: %w", err)
		}
		outputData = buf.Bytes()

	default:
		return fmt.Errorf("unsupported config format '%s' for client %s", formatType, clientName)
	}

	// Ensure the target directory exists
	clientConfigDir := filepath.Dir(clientConfigPath)
	if err := os.MkdirAll(clientConfigDir, 0750); err != nil {
		return fmt.Errorf("failed to create directory '%s' for client %s: %w", clientConfigDir, clientName, err)
	}

	// Write the translated config file
	if err := os.WriteFile(clientConfigPath, outputData, 0644); err != nil {
		return fmt.Errorf("failed to write config file '%s' for client %s: %w", clientConfigPath, clientName, err)
	}

	fmt.Printf("  Successfully wrote config for %s to '%s'\n", clientName, clientConfigPath)
	return nil
}

func (t *Translator) createServerMap(serverConf config.MCPServer) map[string]interface{} {
	serverEntry := make(map[string]interface{})

	if serverConf.Command != "" {
		serverEntry["command"] = serverConf.Command
	}
	if len(serverConf.Args) > 0 {
		serverEntry["args"] = serverConf.Args
	}
	if len(serverConf.Env) > 0 {
		serverEntry["env"] = serverConf.Env
	}
	if serverConf.URL != "" {
		serverEntry["url"] = serverConf.URL
	}
	if serverConf.Disabled {
		serverEntry["disabled"] = serverConf.Disabled
	}
	if len(serverConf.AutoApprove) > 0 {
		serverEntry["autoApprove"] = serverConf.AutoApprove
	}

	return serverEntry
}

// RemoveClientServers removes servers from client configurations that no longer exist in the main MCP configuration
func (t *Translator) RemoveClientServers(clientName string, clientConf config.Client) error {
	clientConfigPath, err := util.ExpandPath(clientConf.ConfigPath)
	if err != nil {
		return fmt.Errorf("failed to expand client config path '%s' for %s: %w", clientConf.ConfigPath, clientName, err)
	}

	// Check if client config file exists
	_, err = os.Stat(clientConfigPath)
	if os.IsNotExist(err) {
		// File doesn't exist, nothing to remove
		return nil
	} else if err != nil {
		return fmt.Errorf("failed to stat client config file '%s': %w", clientConfigPath, err)
	}

	// Read the client config file
	clientConfigData, err := os.ReadFile(clientConfigPath)
	if err != nil {
		return fmt.Errorf("failed to read client config file '%s': %w", clientConfigPath, err)
	}

	if len(clientConfigData) == 0 {
		return nil
	}

	formatType := client.ConfigFormatEnum(clientConf.Type)
	// Determine format type if not explicitly set
	if formatType == "" {
		ext := strings.ToLower(filepath.Ext(clientConfigPath))
		switch ext {
		case ".json":
			// Try to guess from client name or default to simple json
			if strings.Contains(clientName, "claude-desktop") {
				formatType = client.FormatClaudeDesktop
			} else if strings.Contains(clientName, "vscode") {
				formatType = client.FormatVSCode
			} else if strings.Contains(strings.ToLower(clientName), "continue") {
				formatType = client.FormatContinue
			} else {
				formatType = client.FormatSimpleJSON
			}
		case ".yaml", ".yml":
			formatType = client.FormatYAML
		case ".toml":
			formatType = client.FormatTOML
		default:
			// Fallback/Legacy logic if needed, or return error
		}
	}

	// Helper to parse JSON/JSONC
	parseJSON := func(v interface{}) error {
		standardized, err := hujson.Standardize(clientConfigData)
		if err != nil {
			return err
		}
		return json.Unmarshal(standardized, v)
	}

	changed := false

	switch formatType {
	case client.FormatClaudeDesktop, client.FormatSimpleJSON:
		var clientConfig map[string]interface{}
		if err := parseJSON(&clientConfig); err != nil {
			return fmt.Errorf("failed to parse client JSON config file '%s': %w", clientConfigPath, err)
		}

		if mcpServers, ok := clientConfig["mcpServers"].(map[string]interface{}); ok {
			if t.removeObsoleteServers(mcpServers) {
				changed = true
				clientConfig["mcpServers"] = mcpServers
			}
		}

		if changed {
			outputData, err := json.MarshalIndent(clientConfig, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal updated config: %w", err)
			}
			return os.WriteFile(clientConfigPath, outputData, 0644)
		}

	case client.FormatVSCode:
		var clientConfig map[string]interface{}
		if err := parseJSON(&clientConfig); err != nil {
			return fmt.Errorf("failed to parse client JSON config file '%s': %w", clientConfigPath, err)
		}

		if clientConf.Key == "openctx.providers" {
			if providers, ok := clientConfig["openctx.providers"].(map[string]interface{}); ok {
				if t.removeObsoleteServers(providers) {
					changed = true
					clientConfig["openctx.providers"] = providers
				}
			}
		} else {
			if mcpObj, ok := clientConfig["mcp"].(map[string]interface{}); ok {
				if servers, ok := mcpObj["servers"].(map[string]interface{}); ok {
					if t.removeObsoleteServers(servers) {
						changed = true
						mcpObj["servers"] = servers
						clientConfig["mcp"] = mcpObj
					}
				}
			}
		}

		if changed {
			outputData, err := json.MarshalIndent(clientConfig, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal updated config: %w", err)
			}
			return os.WriteFile(clientConfigPath, outputData, 0644)
		}

	case client.FormatContinue:
		// Not implementing removal for Continue format in this patch yet as it involves list filtering
		// leaving as TODO or simple log
		fmt.Println("Warning: Automatic removal of obsolete servers not yet implemented for Continue format")

	case client.FormatYAML:
		var clientConfig map[string]interface{}
		if err := yaml.Unmarshal(clientConfigData, &clientConfig); err != nil {
			return fmt.Errorf("failed to parse client YAML config: %w", err)
		}

		if mcpServers, ok := clientConfig["mcpServers"].(map[string]interface{}); ok {
			if t.removeObsoleteServers(mcpServers) {
				changed = true
				clientConfig["mcpServers"] = mcpServers
			}
		}

		if changed {
			outputData, err := yaml.Marshal(clientConfig)
			if err != nil {
				return fmt.Errorf("failed to marshal updated YAML config: %w", err)
			}
			return os.WriteFile(clientConfigPath, outputData, 0644)
		}

	case client.FormatTOML:
		var clientConfig map[string]interface{}
		if err := toml.Unmarshal(clientConfigData, &clientConfig); err != nil {
			return fmt.Errorf("failed to parse client TOML config: %w", err)
		}

		if mcpServers, ok := clientConfig["mcpServers"].(map[string]interface{}); ok {
			if t.removeObsoleteServers(mcpServers) {
				changed = true
				clientConfig["mcpServers"] = mcpServers
			}
		}

		if changed {
			buf := new(bytes.Buffer)
			if err := toml.NewEncoder(buf).Encode(clientConfig); err != nil {
				return fmt.Errorf("failed to marshal updated TOML config: %w", err)
			}
			return os.WriteFile(clientConfigPath, buf.Bytes(), 0644)
		}
	}

	return nil
}

// removeObsoleteServers removes server entries from a client config map that don't exist in the MCPConfig
// and returns whether any changes were made
func (t *Translator) removeObsoleteServers(servers map[string]interface{}) bool {
	if len(servers) == 0 {
		return false
	}

	changed := false
	for serverID := range servers {
		// Check if this server exists in the main MCP configuration
		if _, exists := t.MCPConfig.MCPServers[serverID]; !exists {
			delete(servers, serverID)
			fmt.Printf("  Removed obsolete server '%s' from client configuration\n", serverID)
			changed = true
		}
	}

	return changed
}
