package manager

import (
	"fmt"

	"github.com/tuannvm/mcpenetes/internal/config"
)

// AddRegistry adds a new registry to the configuration
func AddRegistry(name, url string) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Check for duplicates
	for _, reg := range cfg.Registries {
		if reg.Name == name {
			return fmt.Errorf("registry with name '%s' already exists", name)
		}
		if reg.URL == url {
			return fmt.Errorf("registry with URL '%s' already exists", url)
		}
	}

	cfg.Registries = append(cfg.Registries, config.Registry{
		Name: name,
		URL:  url,
	})

	if err := config.SaveConfig(cfg); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	return nil
}

// RemoveRegistry removes a registry from the configuration
func RemoveRegistry(name string) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	found := false
	newRegistries := []config.Registry{}
	for _, reg := range cfg.Registries {
		if reg.Name == name {
			found = true
			continue
		}
		newRegistries = append(newRegistries, reg)
	}

	if !found {
		return fmt.Errorf("registry '%s' not found", name)
	}

	cfg.Registries = newRegistries

	if err := config.SaveConfig(cfg); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	return nil
}

// UpdateRegistry updates an existing registry
func UpdateRegistry(originalName, newName, newURL string) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	found := false
	for i, reg := range cfg.Registries {
		if reg.Name == originalName {
			// Update fields
			cfg.Registries[i].Name = newName
			cfg.Registries[i].URL = newURL
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("registry '%s' not found", originalName)
	}

	if err := config.SaveConfig(cfg); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	return nil
}
