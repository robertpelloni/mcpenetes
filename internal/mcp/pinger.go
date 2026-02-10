package mcp

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/tuannvm/mcpenetes/internal/config"
)

// PingResult contains the outcome of a server ping
type PingResult struct {
	Success bool          `json:"success"`
	Latency time.Duration `json:"latency"`
	Message string        `json:"message"`
	Version string        `json:"version,omitempty"`
}

// PingServer tests the connection to an MCP server by performing an initialization handshake.
// It returns the latency and any error encountered.
func PingServer(serverConf config.MCPServer) PingResult {
	if serverConf.Command == "" {
		return PingResult{Success: false, Message: "No command specified"}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, serverConf.Command, serverConf.Args...)

	// Set environment variables
	cmd.Env = append(cmd.Env, "PATH="+strings.Join([]string{
		"/usr/local/bin", "/usr/bin", "/bin", "/usr/sbin", "/sbin", // Standard paths
		// Add user paths if needed, e.g. for npx
	}, ":"))

	// Add configured env vars
	for k, v := range serverConf.Env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return PingResult{Success: false, Message: fmt.Sprintf("Failed to create stdin: %v", err)}
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return PingResult{Success: false, Message: fmt.Sprintf("Failed to create stdout: %v", err)}
	}

	// Capture stderr for debugging
	// cmd.Stderr = os.Stderr // Or buffer it

	startTime := time.Now()
	if err := cmd.Start(); err != nil {
		return PingResult{Success: false, Message: fmt.Sprintf("Failed to start server: %v", err)}
	}
	defer func() {
		_ = cmd.Process.Kill()
		_ = cmd.Wait()
	}()

	// Send initialize request
	// https://spec.modelcontextprotocol.io/specification/lifecycle/
	initReq := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "initialize",
		"params": map[string]interface{}{
			"protocolVersion": "2024-11-05",
			"capabilities":    map[string]interface{}{},
			"clientInfo": map[string]interface{}{
				"name":    "mcpenetes-pinger",
				"version": "1.0.0",
			},
		},
		"id": 1,
	}

	reqBytes, _ := json.Marshal(initReq)
	if _, err := stdin.Write(append(reqBytes, '\n')); err != nil {
		return PingResult{Success: false, Message: fmt.Sprintf("Failed to write to stdin: %v", err)}
	}

	// Read response
	scanner := bufio.NewScanner(stdout)

	// We expect a JSON-RPC response.
	// It might be preceded by logs (non-JSON lines).
	// We scan until we find a line that parses as the expected response ID.

	responseFound := false
	var serverVersion string

	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}

		var resp map[string]interface{}
		if err := json.Unmarshal([]byte(line), &resp); err == nil {
			// Check if it's our response
			if id, ok := resp["id"].(float64); ok && id == 1 {
				responseFound = true

				// Try to extract server info
				if result, ok := resp["result"].(map[string]interface{}); ok {
					if serverInfo, ok := result["serverInfo"].(map[string]interface{}); ok {
						if name, ok := serverInfo["name"].(string); ok {
							serverVersion = name
						}
						if ver, ok := serverInfo["version"].(string); ok {
							serverVersion += " " + ver
						}
					}
				}
				break
			}
		}
	}

	latency := time.Since(startTime)

	if !responseFound {
		if ctx.Err() == context.DeadlineExceeded {
			return PingResult{Success: false, Message: "Timeout waiting for initialization response"}
		}
		return PingResult{Success: false, Message: "Invalid or no response from server"}
	}

	return PingResult{
		Success: true,
		Latency: latency,
		Message: "Connection successful",
		Version: serverVersion,
	}
}
