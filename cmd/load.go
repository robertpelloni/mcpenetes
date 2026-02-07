package cmd

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tuannvm/mcpenetes/internal/config"
	"github.com/tuannvm/mcpenetes/internal/core"
	"github.com/tuannvm/mcpenetes/internal/log"
)

// loadCmd represents the load command
var loadCmd = &cobra.Command{
	Use:   "load",
	Short: "Load MCP server configuration from clipboard",
	Long:  `Loads MCP server configuration from the clipboard and adds it to mcp.json`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("Reading configuration from clipboard...")

		// Get clipboard content
		clipboardContent, err := getClipboard()
		if err != nil {
			log.Fatal("Failed to read clipboard: %v", err)
			return
		}

		if clipboardContent == "" {
			log.Fatal("Clipboard is empty")
			return
		}

		// Load config
		cfg, err := config.LoadConfig()
		if err != nil {
			log.Fatal("Error loading config: %v", err)
		}

		mcpCfg, err := config.LoadMCPConfig()
		if err != nil {
			// If error is because the file doesn't exist, create a new config
			mcpCfg = &config.MCPConfig{
				MCPServers: make(map[string]config.MCPServer),
			}
		}

		manager := core.NewManager(cfg, mcpCfg)

		// Import
		count, err := manager.ImportConfig(clipboardContent)
		if err != nil {
			log.Fatal("Failed to import configuration: %v", err)
		}

		log.Success("Successfully loaded %d MCP servers from clipboard", count)
		log.Info("Run 'mcpenetes apply' to install these servers to your clients.")
	},
}

// getClipboard gets the content of the clipboard
func getClipboard() (string, error) {
	var cmd *exec.Cmd
	var out []byte
	var err error

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("pbpaste")
	case "linux":
		// Try xclip first, fallback to wl-paste?
		// Check if xclip is installed
		if _, err := exec.LookPath("xclip"); err == nil {
			cmd = exec.Command("xclip", "-selection", "clipboard", "-o")
		} else if _, err := exec.LookPath("wl-paste"); err == nil {
			cmd = exec.Command("wl-paste")
		} else {
			return "", fmt.Errorf("no clipboard utility found (install xclip or wl-clipboard)")
		}
	case "windows":
		cmd = exec.Command("powershell.exe", "-command", "Get-Clipboard")
	default:
		return "", fmt.Errorf("unsupported platform")
	}

	out, err = cmd.Output()
	if err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			if len(out) > 0 {
				return string(out), nil
			}
			return "", fmt.Errorf("clipboard command failed: %v", err)
		}
		return "", fmt.Errorf("failed to execute clipboard command: %v", err)
	}

	return strings.TrimSpace(string(out)), nil
}

func init() {
	rootCmd.AddCommand(loadCmd)
}
