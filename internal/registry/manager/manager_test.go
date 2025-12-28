package manager

import (
	"path/filepath"
	"testing"

	"github.com/tuannvm/mcpenetes/internal/config"
)

// Helper to mock config loading/saving for tests
func setupTestConfig(t *testing.T) string {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	t.Setenv("HOME", tmpDir)
	// For Windows
	t.Setenv("USERPROFILE", tmpDir)

	return configPath
}

func TestRegistryOperations(t *testing.T) {
	setupTestConfig(t)

	// 1. Add Registry
	err := AddRegistry("TestRegistry", "https://test.com")
	if err != nil {
		t.Fatalf("Failed to add registry: %v", err)
	}

	// Verify it exists
	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	found := false
	for _, r := range cfg.Registries {
		if r.Name == "TestRegistry" {
			found = true
			if r.URL != "https://test.com" {
				t.Errorf("Registry URL incorrect")
			}
			break
		}
	}
	if !found {
		t.Errorf("Registry not added")
	}

	// 2. Update Registry
	err = UpdateRegistry("TestRegistry", "UpdatedRegistry", "https://updated.com")
	if err != nil {
		t.Fatalf("Failed to update registry: %v", err)
	}

	cfg, err = config.LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	found = false
	for _, r := range cfg.Registries {
		if r.Name == "UpdatedRegistry" {
			found = true
			if r.URL != "https://updated.com" {
				t.Errorf("Registry URL not updated")
			}
			break
		}
	}
	if !found {
		t.Errorf("Registry name update failed")
	}

	// 3. Remove Registry
	err = RemoveRegistry("UpdatedRegistry")
	if err != nil {
		t.Fatalf("Failed to remove registry: %v", err)
	}

	cfg, err = config.LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	for _, r := range cfg.Registries {
		if r.Name == "UpdatedRegistry" {
			t.Errorf("Registry was not removed")
		}
	}
}
