package translator_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/tuannvm/mcpenetes/internal/config"
	"github.com/tuannvm/mcpenetes/internal/translator"
)

// TestTranslateAndApply_JSONC tests that JSON with comments (VSCode style)
// is parsed correctly and not corrupted.
func TestTranslateAndApply_JSONC(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "settings.json")

	// 1. Create a VSCode-style config with comments
	initialContent := `
{
	// This is a comment
	"editor.fontSize": 14,
	"mcp": {
		"servers": {
			"existing-server": {
				"command": "echo"
			}
		}
	}
}
`
	err := os.WriteFile(configPath, []byte(initialContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write initial config: %v", err)
	}

	// 2. Setup Translator
	appCfg := &config.Config{
		Backups: config.BackupConfig{Path: filepath.Join(tmpDir, "backups")},
	}
	mcpCfg := &config.MCPConfig{
		MCPServers: map[string]config.MCPServer{
			"new-server": {Command: "node", Args: []string{"server.js"}},
		},
	}

	tr := translator.NewTranslator(appCfg, mcpCfg)

	// 3. Apply a new server
	clientConf := config.Client{
		ConfigPath: configPath,
		Type:       "vscode",
	}
	serverConf := config.MCPServer{
		Command: "node",
		Args:    []string{"server.js"},
	}

	err = tr.TranslateAndApply("vscode-test", clientConf, serverConf)
	if err != nil {
		t.Fatalf("TranslateAndApply failed: %v", err)
	}

	// 4. Verify results
	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read back config: %v", err)
	}

	var resultMap map[string]interface{}
	err = json.Unmarshal(content, &resultMap)
	if err != nil {
		t.Fatalf("Failed to parse result JSON: %v", err)
	}

	// Check if "editor.fontSize" is preserved
	if val, ok := resultMap["editor.fontSize"]; !ok || val.(float64) != 14 {
		t.Errorf("Existing setting 'editor.fontSize' was lost or corrupted")
	}

	// Check if new server was added
	mcp := resultMap["mcp"].(map[string]interface{})
	servers := mcp["servers"].(map[string]interface{})

	if _, ok := servers["new-server"]; !ok {
		t.Errorf("New server 'new-server' was not added")
	}
	if _, ok := servers["existing-server"]; !ok {
		t.Errorf("Existing server 'existing-server' was lost")
	}
}

// TestTranslateAndApply_InvalidJSON verifies that processing aborts on invalid JSON
func TestTranslateAndApply_InvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "broken.json")

	// 1. Create broken JSON
	err := os.WriteFile(configPath, []byte(`{ "unclosed_string": "oops `), 0644)
	if err != nil {
		t.Fatalf("Failed to write broken config: %v", err)
	}

	tr := translator.NewTranslator(&config.Config{}, &config.MCPConfig{})
	clientConf := config.Client{ConfigPath: configPath, Type: "simple-json"}
	serverConf := config.MCPServer{Command: "echo"}

	// 2. Apply should fail
	err = tr.TranslateAndApply("broken-client", clientConf, serverConf)
	if err == nil {
		t.Error("Expected error when parsing invalid JSON, but got nil")
	} else {
		t.Logf("Got expected error: %v", err)
	}
}
