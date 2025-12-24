package doctor

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/tuannvm/mcpenetes/internal/config"
	"github.com/tuannvm/mcpenetes/internal/util"
)

// CheckResult represents the result of a single health check
type CheckResult struct {
	Name    string `json:"name"`
	Status  string `json:"status"` // "ok", "warning", "error"
	Message string `json:"message"`
}

// RunChecks performs all system health checks
func RunChecks() []CheckResult {
	var results []CheckResult

	results = append(results, checkConfigs()...)
	results = append(results, checkEnvironment()...)
	results = append(results, checkClients()...)

	return results
}

func checkConfigs() []CheckResult {
	var results []CheckResult

	// Check config.yaml
	_, err := config.LoadConfig()
	if err != nil {
		results = append(results, CheckResult{
			Name:    "Config File (config.yaml)",
			Status:  "error",
			Message: fmt.Sprintf("Failed to load: %v", err),
		})
	} else {
		results = append(results, CheckResult{
			Name:    "Config File (config.yaml)",
			Status:  "ok",
			Message: "Loaded successfully",
		})
	}

	// Check mcp.json
	_, err = config.LoadMCPConfig()
	if err != nil {
		results = append(results, CheckResult{
			Name:    "MCP Config (mcp.json)",
			Status:  "error",
			Message: fmt.Sprintf("Failed to load: %v", err),
		})
	} else {
		results = append(results, CheckResult{
			Name:    "MCP Config (mcp.json)",
			Status:  "ok",
			Message: "Loaded successfully",
		})
	}

	return results
}

func checkEnvironment() []CheckResult {
	var results []CheckResult
	tools := []string{"node", "npm", "npx", "python", "uv", "docker"}

	for _, tool := range tools {
		path, err := exec.LookPath(tool)
		if err != nil {
			status := "warning"
			msg := "Not found in PATH"

			// npx is critical for many MCP servers
			if tool == "npx" {
				status = "error"
				msg = "Required for running many MCP servers"
			}

			results = append(results, CheckResult{
				Name:    fmt.Sprintf("Tool: %s", tool),
				Status:  status,
				Message: msg,
			})
		} else {
			results = append(results, CheckResult{
				Name:    fmt.Sprintf("Tool: %s", tool),
				Status:  "ok",
				Message: fmt.Sprintf("Found at %s", path),
			})
		}
	}

	return results
}

func checkClients() []CheckResult {
	var results []CheckResult

	clients, err := util.DetectMCPClients()
	if err != nil {
		results = append(results, CheckResult{
			Name:    "Client Detection",
			Status:  "error",
			Message: fmt.Sprintf("Failed to detect clients: %v", err),
		})
		return results
	}

	if len(clients) == 0 {
		results = append(results, CheckResult{
			Name:    "Client Detection",
			Status:  "warning",
			Message: "No supported clients detected",
		})
	} else {
		results = append(results, CheckResult{
			Name:    "Client Detection",
			Status:  "ok",
			Message: fmt.Sprintf("Detected %d clients", len(clients)),
		})
	}

	for name, client := range clients {
		path, err := util.ExpandPath(client.ConfigPath)
		if err != nil {
			results = append(results, CheckResult{
				Name:    fmt.Sprintf("Client: %s", name),
				Status:  "error",
				Message: fmt.Sprintf("Path error: %v", err),
			})
			continue
		}

		// Check if file exists or is writable
		info, err := os.Stat(path)
		if err != nil {
			if os.IsNotExist(err) {
				// Check directory
				dir := filepath.Dir(path)
				if _, err := os.Stat(dir); err == nil {
					results = append(results, CheckResult{
						Name:    fmt.Sprintf("Client: %s", name),
						Status:  "ok",
						Message: "Config file missing, but directory exists (ready to create)",
					})
				} else {
					results = append(results, CheckResult{
						Name:    fmt.Sprintf("Client: %s", name),
						Status:  "warning",
						Message: "Config file and directory missing",
					})
				}
			} else {
				results = append(results, CheckResult{
					Name:    fmt.Sprintf("Client: %s", name),
					Status:  "error",
					Message: fmt.Sprintf("Access error: %v", err),
				})
			}
		} else {
			if info.Mode().Perm()&0200 == 0 {
				results = append(results, CheckResult{
					Name:    fmt.Sprintf("Client: %s", name),
					Status:  "warning",
					Message: "Config file exists but is read-only",
				})
			} else {
				results = append(results, CheckResult{
					Name:    fmt.Sprintf("Client: %s", name),
					Status:  "ok",
					Message: "Config file exists and is writable",
				})
			}
		}
	}

	return results
}
