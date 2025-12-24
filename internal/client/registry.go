package client

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"gopkg.in/yaml.v3"
)

// ConfigFormatEnum defines the supported configuration formats
type ConfigFormatEnum string

const (
	FormatClaudeDesktop ConfigFormatEnum = "claude-desktop" // {"mcpServers": {...}}
	FormatVSCode        ConfigFormatEnum = "vscode"         // {"mcp": {"servers": {...}}} or {"mcp.servers": {...}}
	FormatSimpleJSON    ConfigFormatEnum = "simple-json"    // {"mcpServers": {...}} (Standard MCP)
	FormatYAML          ConfigFormatEnum = "yaml"           // YAML format
	FormatTOML          ConfigFormatEnum = "toml"           // TOML format
	FormatContinue      ConfigFormatEnum = "continue"       // Continue.dev config.json structure
)

// BaseDirEnum defines where the config is relative to
type BaseDirEnum string

const (
	BaseHome        BaseDirEnum = "home"
	BaseAppData     BaseDirEnum = "appdata"     // Windows %APPDATA%
	BaseUserProfile BaseDirEnum = "userprofile" // Windows %USERPROFILE%
)

// PathDefinition defines a path strategy for a specific OS
type PathDefinition struct {
	Base BaseDirEnum
	Path string // Relative path from the base
}

// ClientDefinition defines the metadata and paths for a tool
type ClientDefinition struct {
	ID           string
	Name         string
	ConfigFormat ConfigFormatEnum
	// ConfigKey is an optional field to override the default JSON key
	// e.g. "openctx.providers" instead of "mcp.servers" for VSCode format
	ConfigKey string
	// Map of OS to list of potential config paths
	// supported OS keys: "windows", "darwin", "linux"
	Paths map[string][]PathDefinition
}

// Registry holds the list of known clients
var Registry = []ClientDefinition{
	// --- Desktop IDEs ---
	{
		ID:           "claude-desktop",
		Name:         "Claude Desktop",
		ConfigFormat: FormatClaudeDesktop,
		Paths: map[string][]PathDefinition{
			"darwin": {
				{Base: BaseHome, Path: filepath.Join("Library", "Application Support", "Claude", "claude_desktop_config.json")},
			},
			"windows": {
				{Base: BaseAppData, Path: filepath.Join("Claude", "claude_desktop_config.json")},
			},
			"linux": {
				{Base: BaseHome, Path: filepath.Join(".config", "Claude", "claude_desktop_config.json")},
			},
		},
	},
	{
		ID:           "cursor",
		Name:         "Cursor",
		ConfigFormat: FormatSimpleJSON,
		Paths: map[string][]PathDefinition{
			"darwin": {
				{Base: BaseHome, Path: filepath.Join(".cursor", "mcp.json")},
				{Base: BaseHome, Path: filepath.Join("Library", "Application Support", "Cursor", "User", "mcp.json")},
			},
			"windows": {
				{Base: BaseAppData, Path: filepath.Join("Cursor", "User", "mcp.json")},
			},
			"linux": {
				{Base: BaseHome, Path: filepath.Join(".cursor", "mcp.json")},
			},
		},
	},
	{
		ID:           "windsurf",
		Name:         "Windsurf",
		ConfigFormat: FormatSimpleJSON,
		Paths: map[string][]PathDefinition{
			"darwin": {
				{Base: BaseHome, Path: filepath.Join(".codeium", "windsurf", "mcp_config.json")},
			},
			"windows": {
				{Base: BaseUserProfile, Path: filepath.Join(".codeium", "windsurf", "mcp_config.json")},
			},
			"linux": {
				{Base: BaseHome, Path: filepath.Join(".codeium", "windsurf", "mcp_config.json")},
			},
		},
	},
	{
		ID:           "vscode",
		Name:         "VS Code",
		ConfigFormat: FormatVSCode,
		Paths: map[string][]PathDefinition{
			"darwin": {
				{Base: BaseHome, Path: filepath.Join("Library", "Application Support", "Code", "User", "settings.json")},
			},
			"windows": {
				{Base: BaseAppData, Path: filepath.Join("Code", "User", "settings.json")},
			},
			"linux": {
				{Base: BaseHome, Path: filepath.Join(".config", "Code", "User", "settings.json")},
			},
		},
	},
	{
		ID:           "vscode-insiders",
		Name:         "VS Code Insiders",
		ConfigFormat: FormatVSCode,
		Paths: map[string][]PathDefinition{
			"darwin": {
				{Base: BaseHome, Path: filepath.Join("Library", "Application Support", "Code - Insiders", "User", "settings.json")},
			},
			"windows": {
				{Base: BaseAppData, Path: filepath.Join("Code - Insiders", "User", "settings.json")},
			},
			"linux": {
				{Base: BaseHome, Path: filepath.Join(".config", "Code - Insiders", "User", "settings.json")},
			},
		},
	},
	{
		ID:           "zed",
		Name:         "Zed",
		ConfigFormat: FormatSimpleJSON,
		Paths: map[string][]PathDefinition{
			"darwin": {
				{Base: BaseHome, Path: filepath.Join(".config", "zed", "settings.json")},
			},
			"windows": {
				{Base: BaseAppData, Path: filepath.Join("Zed", "settings.json")},
			},
			"linux": {
				{Base: BaseHome, Path: filepath.Join(".config", "zed", "settings.json")},
			},
		},
	},
	{
		ID:           "trae",
		Name:         "Trae",
		ConfigFormat: FormatSimpleJSON, // Assuming standard format, need to verify docs/user info if available
		Paths: map[string][]PathDefinition{
			"darwin": {
				{Base: BaseHome, Path: filepath.Join("Library", "Application Support", "Trae", "User", "globalStorage", "mcp.json")}, // Guess based on Electron/VSCode forks
			},
			"windows": {
				{Base: BaseAppData, Path: filepath.Join("Trae", "User", "globalStorage", "mcp.json")}, // Guess
			},
			"linux": {
				{Base: BaseHome, Path: filepath.Join(".config", "Trae", "User", "globalStorage", "mcp.json")}, // Guess
			},
		},
	},
	{
		ID:           "amazon-q",
		Name:         "Amazon Q (CodeWhisperer)",
		ConfigFormat: FormatSimpleJSON,
		Paths: map[string][]PathDefinition{
			"darwin": {
				{Base: BaseHome, Path: filepath.Join(".aws", "amazonq", "mcp.json")},
			},
			"windows": {
				{Base: BaseUserProfile, Path: filepath.Join(".aws", "amazonq", "mcp.json")},
			},
			"linux": {
				{Base: BaseHome, Path: filepath.Join(".aws", "amazonq", "mcp.json")},
			},
		},
	},
	{
		ID:           "jetbrains-junie",
		Name:         "JetBrains (Junie)",
		ConfigFormat: FormatSimpleJSON,
		Paths: map[string][]PathDefinition{
			"darwin": {
				{Base: BaseHome, Path: filepath.Join(".junie", "mcp", "mcp.json")},
			},
			"windows": {
				{Base: BaseHome, Path: filepath.Join(".junie", "mcp", "mcp.json")}, // Note: Docs say ~/.junie even on Windows, need to verify if it respects %USERPROFILE% (which BaseHome maps to on detection)
			},
			"linux": {
				{Base: BaseHome, Path: filepath.Join(".junie", "mcp", "mcp.json")},
			},
		},
	},

	// --- VSCode Extensions / "Autonomous Agents" ---
	{
		ID:           "cody",
		Name:         "Cody (Sourcegraph)",
		ConfigFormat: FormatVSCode,
		ConfigKey:    "openctx.providers",
		Paths: map[string][]PathDefinition{
			"darwin": {
				{Base: BaseHome, Path: filepath.Join("Library", "Application Support", "Code", "User", "settings.json")},
			},
			"windows": {
				{Base: BaseAppData, Path: filepath.Join("Code", "User", "settings.json")},
			},
			"linux": {
				{Base: BaseHome, Path: filepath.Join(".config", "Code", "User", "settings.json")},
			},
		},
	},
	{
		ID:           "cline",
		Name:         "Cline",
		ConfigFormat: FormatSimpleJSON,
		Paths: map[string][]PathDefinition{
			"darwin": {
				{Base: BaseHome, Path: filepath.Join("Library", "Application Support", "Code", "User", "globalStorage", "saoudrizwan.claude-dev", "settings", "cline_mcp_settings.json")},
			},
			"windows": {
				{Base: BaseAppData, Path: filepath.Join("Code", "User", "globalStorage", "saoudrizwan.claude-dev", "settings", "cline_mcp_settings.json")},
			},
			"linux": {
				{Base: BaseHome, Path: filepath.Join(".config", "Code", "User", "globalStorage", "saoudrizwan.claude-dev", "settings", "cline_mcp_settings.json")},
			},
		},
	},
	{
		ID:           "roo-code",
		Name:         "Roo Code",
		ConfigFormat: FormatSimpleJSON,
		Paths: map[string][]PathDefinition{
			"darwin": {
				{Base: BaseHome, Path: filepath.Join("Library", "Application Support", "Code", "User", "globalStorage", "rooveterinaryinc.roo-cline", "settings", "cline_mcp_settings.json")},
			},
			"windows": {
				{Base: BaseAppData, Path: filepath.Join("Code", "User", "globalStorage", "rooveterinaryinc.roo-cline", "settings", "cline_mcp_settings.json")},
			},
			"linux": {
				{Base: BaseHome, Path: filepath.Join(".config", "Code", "User", "globalStorage", "rooveterinaryinc.roo-cline", "settings", "cline_mcp_settings.json")},
			},
		},
	},
	{
		ID:           "continue",
		Name:         "Continue",
		ConfigFormat: FormatContinue,
		Paths: map[string][]PathDefinition{
			"darwin": {
				{Base: BaseHome, Path: filepath.Join(".continue", "config.json")},
			},
			"windows": {
				{Base: BaseUserProfile, Path: filepath.Join(".continue", "config.json")},
			},
			"linux": {
				{Base: BaseHome, Path: filepath.Join(".continue", "config.json")},
			},
		},
	},

	// --- Desktop Apps ---
	{
		ID:           "lm-studio",
		Name:         "LM Studio",
		ConfigFormat: FormatSimpleJSON,
		Paths: map[string][]PathDefinition{
			"darwin": {
				{Base: BaseHome, Path: filepath.Join(".lmstudio", "mcp.json")},
			},
			"windows": {
				{Base: BaseUserProfile, Path: filepath.Join(".lmstudio", "mcp.json")},
			},
			"linux": {
				{Base: BaseHome, Path: filepath.Join(".lmstudio", "mcp.json")},
			},
		},
	},
	{
		ID:           "anythingllm",
		Name:         "AnythingLLM",
		ConfigFormat: FormatSimpleJSON,
		Paths: map[string][]PathDefinition{
			"darwin": {
				{Base: BaseHome, Path: filepath.Join("Library", "Application Support", "anythingllm-desktop", "storage", "plugins", "anythingllm_mcp_servers.json")},
			},
			"windows": {
				{Base: BaseAppData, Path: filepath.Join("anythingllm-desktop", "storage", "plugins", "anythingllm_mcp_servers.json")},
			},
			"linux": {
				{Base: BaseHome, Path: filepath.Join(".config", "anythingllm-desktop", "storage", "plugins", "anythingllm_mcp_servers.json")},
			},
		},
	},
	{
		ID:           "tabby",
		Name:         "Tabby",
		ConfigFormat: FormatTOML, // Tabby uses config.toml
		Paths: map[string][]PathDefinition{
			"darwin": {
				{Base: BaseHome, Path: filepath.Join(".tabby-client", "agent", "config.toml")}, // Typical agent config
			},
			"windows": {
				{Base: BaseAppData, Path: filepath.Join("Tabby", "config.toml")}, // Check specific location, often user profile or appdata
				{Base: BaseUserProfile, Path: filepath.Join(".tabby-client", "agent", "config.toml")},
			},
			"linux": {
				{Base: BaseHome, Path: filepath.Join(".tabby-client", "agent", "config.toml")},
			},
		},
	},
	{
		ID:           "librechat",
		Name:         "LibreChat",
		ConfigFormat: FormatYAML,
		Paths: map[string][]PathDefinition{
			"darwin": {
				{Base: BaseHome, Path: filepath.Join("librechat.yaml")}, // Often in project root or home
				{Base: BaseHome, Path: filepath.Join(".librechat", "librechat.yaml")},
			},
			"windows": {
				{Base: BaseUserProfile, Path: filepath.Join("librechat.yaml")},
				{Base: BaseUserProfile, Path: filepath.Join(".librechat", "librechat.yaml")},
			},
			"linux": {
				{Base: BaseHome, Path: filepath.Join("librechat.yaml")},
				{Base: BaseHome, Path: filepath.Join(".librechat", "librechat.yaml")},
			},
		},
	},

	// --- CLIs ---
	{
		ID:           "goose",
		Name:         "Goose CLI",
		ConfigFormat: FormatYAML,
		Paths: map[string][]PathDefinition{
			"darwin": {
				{Base: BaseHome, Path: filepath.Join(".config", "goose", "config.yaml")},
			},
			"windows": {
				{Base: BaseAppData, Path: filepath.Join("Block", "goose", "config", "config.yaml")},
			},
			"linux": {
				{Base: BaseHome, Path: filepath.Join(".config", "goose", "config.yaml")},
			},
		},
	},
	{
		ID:           "mistral-vibe",
		Name:         "Mistral Vibe",
		ConfigFormat: FormatTOML,
		Paths: map[string][]PathDefinition{
			"darwin": {
				{Base: BaseHome, Path: filepath.Join(".vibe", "config.toml")},
			},
			"windows": {
				{Base: BaseUserProfile, Path: filepath.Join(".vibe", "config.toml")},
			},
			"linux": {
				{Base: BaseHome, Path: filepath.Join(".vibe", "config.toml")},
			},
		},
	},
	{
		ID:           "code-cli",
		Name:         "Code CLI (Codex)",
		ConfigFormat: FormatSimpleJSON,
		Paths: map[string][]PathDefinition{
			"darwin": {
				{Base: BaseHome, Path: filepath.Join(".config", "code-cli", "mcp.json")}, // Guess based on convention
			},
			"windows": {
				{Base: BaseAppData, Path: filepath.Join("code-cli", "mcp.json")},
			},
			"linux": {
				{Base: BaseHome, Path: filepath.Join(".config", "code-cli", "mcp.json")},
			},
		},
	},
	{
		ID:           "grok-cli",
		Name:         "Grok CLI",
		ConfigFormat: FormatSimpleJSON,
		Paths: map[string][]PathDefinition{
			"darwin": {
				{Base: BaseHome, Path: filepath.Join(".grok", "config.json")}, // Typical CLI convention
			},
			"windows": {
				{Base: BaseUserProfile, Path: filepath.Join(".grok", "config.json")},
			},
			"linux": {
				{Base: BaseHome, Path: filepath.Join(".grok", "config.json")},
			},
		},
	},
	{
		ID:           "open-interpreter",
		Name:         "Open Interpreter",
		ConfigFormat: FormatYAML, // Often uses YAML for profiles
		Paths: map[string][]PathDefinition{
			"darwin": {
				{Base: BaseHome, Path: filepath.Join(".config", "open-interpreter", "config.yaml")},
			},
			"windows": {
				{Base: BaseAppData, Path: filepath.Join("Open Interpreter", "config.yaml")},
			},
			"linux": {
				{Base: BaseHome, Path: filepath.Join(".config", "open-interpreter", "config.yaml")},
			},
		},
	},
	{
		ID:           "factory-cli",
		Name:         "Factory CLI",
		ConfigFormat: FormatSimpleJSON,
		Paths: map[string][]PathDefinition{
			"darwin": {
				{Base: BaseHome, Path: filepath.Join(".factory", "config.json")},
			},
			"windows": {
				{Base: BaseUserProfile, Path: filepath.Join(".factory", "config.json")},
			},
			"linux": {
				{Base: BaseHome, Path: filepath.Join(".factory", "config.json")},
			},
		},
	},
	{
		ID:           "aider",
		Name:         "Aider",
		ConfigFormat: FormatYAML,
		Paths: map[string][]PathDefinition{
			"darwin": {
				{Base: BaseHome, Path: filepath.Join(".aider.conf.yml")},
			},
			"windows": {
				{Base: BaseUserProfile, Path: filepath.Join(".aider.conf.yml")},
			},
			"linux": {
				{Base: BaseHome, Path: filepath.Join(".aider.conf.yml")},
			},
		},
	},
	{
		ID:           "pearai",
		Name:         "PearAI",
		ConfigFormat: FormatSimpleJSON, // VSCode fork, uses settings.json
		Paths: map[string][]PathDefinition{
			"darwin": {
				{Base: BaseHome, Path: filepath.Join("Library", "Application Support", "PearAI", "User", "settings.json")},
			},
			"windows": {
				{Base: BaseAppData, Path: filepath.Join("PearAI", "User", "settings.json")},
			},
			"linux": {
				{Base: BaseHome, Path: filepath.Join(".config", "PearAI", "User", "settings.json")},
			},
		},
	},
	{
		ID:           "void",
		Name:         "Void",
		ConfigFormat: FormatVSCode, // VSCode fork
		Paths: map[string][]PathDefinition{
			"darwin": {
				{Base: BaseHome, Path: filepath.Join("Library", "Application Support", "Void", "User", "settings.json")},
			},
			"windows": {
				{Base: BaseAppData, Path: filepath.Join("Void", "User", "settings.json")},
			},
			"linux": {
				{Base: BaseHome, Path: filepath.Join(".config", "Void", "User", "settings.json")},
			},
		},
	},
	{
		ID:           "melty",
		Name:         "Melty",
		ConfigFormat: FormatVSCode, // VSCode fork
		Paths: map[string][]PathDefinition{
			"darwin": {
				{Base: BaseHome, Path: filepath.Join("Library", "Application Support", "Melty", "User", "settings.json")},
			},
			"windows": {
				{Base: BaseAppData, Path: filepath.Join("Melty", "User", "settings.json")},
			},
			"linux": {
				{Base: BaseHome, Path: filepath.Join(".config", "Melty", "User", "settings.json")},
			},
		},
	},
	{
		ID:           "codebuddy",
		Name:         "CodeBuddy",
		ConfigFormat: FormatVSCode, // VSCode fork/extension
		Paths: map[string][]PathDefinition{
			"darwin": {
				{Base: BaseHome, Path: filepath.Join(".codebuddy", "settings.json")},
			},
			"windows": {
				{Base: BaseUserProfile, Path: filepath.Join(".codebuddy", "settings.json")},
			},
			"linux": {
				{Base: BaseHome, Path: filepath.Join(".codebuddy", "settings.json")},
			},
		},
	},
	{
		ID:           "kiro",
		Name:         "Kiro",
		ConfigFormat: FormatSimpleJSON,
		Paths: map[string][]PathDefinition{
			"darwin": {
				{Base: BaseHome, Path: filepath.Join(".kiro", "settings", "mcp.json")},
			},
			"windows": {
				{Base: BaseUserProfile, Path: filepath.Join(".kiro", "settings", "mcp.json")},
			},
			"linux": {
				{Base: BaseHome, Path: filepath.Join(".kiro", "settings", "mcp.json")},
			},
		},
	},
	{
		ID:           "jan",
		Name:         "Jan",
		ConfigFormat: FormatSimpleJSON, // Guessing simple JSON for now, might be in assistant.json or settings.json
		Paths: map[string][]PathDefinition{
			"darwin": {
				{Base: BaseHome, Path: filepath.Join("Library", "Application Support", "Jan", "data", "settings.json")},
				{Base: BaseHome, Path: filepath.Join("jan", "settings.json")},
			},
			"windows": {
				{Base: BaseAppData, Path: filepath.Join("Jan", "data", "settings.json")},
				{Base: BaseUserProfile, Path: filepath.Join("jan", "settings.json")},
			},
			"linux": {
				{Base: BaseHome, Path: filepath.Join(".config", "Jan", "data", "settings.json")},
				{Base: BaseHome, Path: filepath.Join("jan", "settings.json")},
			},
		},
	},
	{
		ID:           "warp",
		Name:         "Warp",
		ConfigFormat: FormatSimpleJSON,
		Paths: map[string][]PathDefinition{
			"darwin": {
				{Base: BaseHome, Path: filepath.Join(".local", "state", "warp-terminal", "mcp")},
			},
			"windows": {
				{Base: BaseUserProfile, Path: filepath.Join(".local", "state", "warp-terminal", "mcp")},
				{Base: BaseAppData, Path: filepath.Join("Warp", "mcp.json")}, // Fallback guess
			},
			"linux": {
				{Base: BaseHome, Path: filepath.Join(".local", "state", "warp-terminal", "mcp")},
			},
		},
	},
	{
		ID:           "llm-cli",
		Name:         "LLM CLI (Simon Willison)",
		ConfigFormat: FormatSimpleJSON,
		Paths: map[string][]PathDefinition{
			"darwin": {
				{Base: BaseHome, Path: filepath.Join(".llm-tools-mcp", "mcp.json")},
				{Base: BaseHome, Path: filepath.Join(".config", "io.datasette.llm", "mcp.json")},
			},
			"windows": {
				{Base: BaseUserProfile, Path: filepath.Join(".llm-tools-mcp", "mcp.json")},
				{Base: BaseAppData, Path: filepath.Join("io.datasette.llm", "mcp.json")},
			},
			"linux": {
				{Base: BaseHome, Path: filepath.Join(".llm-tools-mcp", "mcp.json")},
				{Base: BaseHome, Path: filepath.Join(".config", "io.datasette.llm", "mcp.json")},
			},
		},
	},
	{
		ID:           "claude-code",
		Name:         "Claude Code CLI",
		ConfigFormat: FormatSimpleJSON, // ~/.claude.json uses {"mcpServers": ...}
		Paths: map[string][]PathDefinition{
			"darwin": {
				{Base: BaseHome, Path: filepath.Join(".claude.json")},
			},
			"windows": {
				{Base: BaseUserProfile, Path: filepath.Join(".claude.json")},
			},
			"linux": {
				{Base: BaseHome, Path: filepath.Join(".claude.json")},
			},
		},
	},
	{
		ID:           "boltai",
		Name:         "BoltAI",
		ConfigFormat: FormatSimpleJSON,
		Paths: map[string][]PathDefinition{
			"darwin": {
				{Base: BaseHome, Path: filepath.Join("Library", "Application Support", "BoltAI", "mcp.json")},
			},
			"windows": {
				{Base: BaseAppData, Path: filepath.Join("BoltAI", "mcp.json")}, // Standard assumption for Electron/similar apps
			},
			// Linux support for BoltAI unknown/unlikely
		},
	},
}

// DetectedClient represents a client found on the system
type DetectedClient struct {
	ID           string
	Name         string
	ConfigPath   string
	ConfigFormat ConfigFormatEnum
	ConfigKey    string
}

// UserRegistryFile is the path to the user-defined registry file
const UserRegistryFile = "clients.yaml"

// DetectClients scans the system for known clients
func DetectClients() (map[string]DetectedClient, error) {
	clients := make(map[string]DetectedClient)
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	appData := os.Getenv("APPDATA")
	userProfile := os.Getenv("USERPROFILE")

	// Pre-calculate base paths
	basePaths := map[BaseDirEnum]string{
		BaseHome:        homeDir,
		BaseAppData:     appData,
		BaseUserProfile: userProfile,
	}

	// Windows fallback for AppData
	if runtime.GOOS == "windows" {
		if basePaths[BaseAppData] == "" {
			basePaths[BaseAppData] = homeDir
		}
		if basePaths[BaseUserProfile] == "" {
			basePaths[BaseUserProfile] = homeDir
		}
	}

	// Combine built-in registry with user-defined registry
	registryToScan := Registry

	// Load user-defined registry
	configDir := filepath.Join(homeDir, ".config", "mcpetes")
	userRegPath := filepath.Join(configDir, UserRegistryFile)
	if data, err := os.ReadFile(userRegPath); err == nil {
		var userClients []ClientDefinition
		if err := yaml.Unmarshal(data, &userClients); err == nil {
			registryToScan = append(registryToScan, userClients...)
		} else {
			fmt.Printf("Warning: Failed to parse user registry at %s: %v\n", userRegPath, err)
		}
	}

	for _, def := range registryToScan {
		paths, ok := def.Paths[runtime.GOOS]
		if !ok {
			continue
		}

		for _, pathDef := range paths {
			basePath := basePaths[pathDef.Base]
			if basePath == "" {
				continue
			}

			fullPath := filepath.Join(basePath, pathDef.Path)

			// Check if file exists
			if _, err := os.Stat(fullPath); err == nil {
				clients[def.ID] = DetectedClient{
					ID:           def.ID,
					Name:         def.Name,
					ConfigPath:   fullPath,
					ConfigFormat: def.ConfigFormat,
					ConfigKey:    def.ConfigKey,
				}
				break // Found valid config file
			}

			// Fallback: Check if directory exists
			dirPath := filepath.Dir(fullPath)
			if _, err := os.Stat(dirPath); err == nil {
				// Directory exists, allow this client so we can create the config file
				clients[def.ID] = DetectedClient{
					ID:           def.ID,
					Name:         def.Name,
					ConfigPath:   fullPath,
					ConfigFormat: def.ConfigFormat,
					ConfigKey:    def.ConfigKey,
				}
				break
			}
		}
	}

	return clients, nil
}
