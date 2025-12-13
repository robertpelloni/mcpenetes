package client_test

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/tuannvm/mcpenetes/internal/client"
)

// TestDetectClients_Found verifies that detectClients finds a file if it exists.
// We use a temporary directory as HOME to simulate the environment.
func TestDetectClients_Found(t *testing.T) {
	// 1. Setup fake home directory
	tmpHome := t.TempDir()

	// Set environment variables for the test
	t.Setenv("HOME", tmpHome)
	t.Setenv("USERPROFILE", tmpHome) // Windows
	t.Setenv("APPDATA", filepath.Join(tmpHome, "AppData", "Roaming")) // Windows

	// 2. Create a dummy config file for a known client (e.g. Claude Desktop)
	// We need to know the relative path for the current OS from the Registry
	// Let's look up Claude Desktop in the registry to get the expected path
	var targetPath string
	found := false

	for _, def := range client.Registry {
		if def.ID == "claude-desktop" {
			paths, ok := def.Paths[runtime.GOOS]
			if ok && len(paths) > 0 {
				pathDef := paths[0]

				var basePath string
				switch pathDef.Base {
				case client.BaseHome:
					basePath = tmpHome
				case client.BaseAppData:
					basePath = os.Getenv("APPDATA")
				case client.BaseUserProfile:
					basePath = os.Getenv("USERPROFILE")
				}

				targetPath = filepath.Join(basePath, pathDef.Path)
				found = true
				break
			}
		}
	}

	if !found {
		t.Skip("Claude Desktop not supported on this OS, skipping test")
	}

	// Create the directory and file
	err := os.MkdirAll(filepath.Dir(targetPath), 0755)
	if err != nil {
		t.Fatalf("Failed to create test directories: %v", err)
	}
	err = os.WriteFile(targetPath, []byte("{}"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	// 3. Run DetectClients
	detected, err := client.DetectClients()
	if err != nil {
		t.Fatalf("DetectClients failed: %v", err)
	}

	// 4. Verify Claude Desktop is found
	if _, ok := detected["claude-desktop"]; !ok {
		t.Errorf("Expected to detect 'claude-desktop', but it was not found")
	} else {
		t.Logf("Successfully detected claude-desktop at %s", targetPath)
	}
}

// TestDetectClients_FallbackDirectory verifies that detectClients finds a client
// if the config file is missing but the parent directory exists.
func TestDetectClients_FallbackDirectory(t *testing.T) {
	// 1. Setup fake home directory
	tmpHome := t.TempDir()

	t.Setenv("HOME", tmpHome)
	t.Setenv("USERPROFILE", tmpHome)
	t.Setenv("APPDATA", filepath.Join(tmpHome, "AppData", "Roaming"))

	// 2. Create directory for a client (e.g. Windsurf) but NO config file
	// Lookup Windsurf path
	var targetPath string
	found := false

	for _, def := range client.Registry {
		if def.ID == "windsurf" {
			paths, ok := def.Paths[runtime.GOOS]
			if ok && len(paths) > 0 {
				pathDef := paths[0]
				var basePath string
				switch pathDef.Base {
				case client.BaseHome:
					basePath = tmpHome
				case client.BaseAppData:
					basePath = os.Getenv("APPDATA")
				case client.BaseUserProfile:
					basePath = os.Getenv("USERPROFILE")
				}
				targetPath = filepath.Join(basePath, pathDef.Path)
				found = true
				break
			}
		}
	}

	if !found {
		t.Skip("Windsurf not supported on this OS, skipping test")
	}

	// Create ONLY the directory
	err := os.MkdirAll(filepath.Dir(targetPath), 0755)
	if err != nil {
		t.Fatalf("Failed to create test directories: %v", err)
	}
	// DO NOT create the file

	// 3. Run DetectClients
	detected, err := client.DetectClients()
	if err != nil {
		t.Fatalf("DetectClients failed: %v", err)
	}

	// 4. Verify Windsurf is found via directory fallback
	if _, ok := detected["windsurf"]; !ok {
		t.Errorf("Expected to detect 'windsurf' via directory fallback, but it was not found")
	} else {
		t.Logf("Successfully detected windsurf (directory only) at %s", targetPath)
	}
}
