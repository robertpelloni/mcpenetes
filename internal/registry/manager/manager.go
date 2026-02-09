package manager

import (
	"fmt"

	"github.com/tuannvm/mcpenetes/internal/config"
)

// AddRegistry adds a new registry to the configuration
func AddRegistry(name, url string) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	// Check if already exists
	for _, reg := range cfg.Registries {
		if reg.Name == name {
			return fmt.Errorf("registry '%s' already exists", name)
		}
	}

	cfg.Registries = append(cfg.Registries, config.Registry{
		Name: name,
		URL:  url,
	})

	return config.SaveConfig(cfg)
}

// RemoveRegistry removes a registry from the configuration
func RemoveRegistry(name string) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
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
	return config.SaveConfig(cfg)
}

// UpdateRegistry updates an existing registry's URL
func UpdateRegistry(name, newURL string) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	found := false
	for i, reg := range cfg.Registries {
		if reg.Name == name {
			cfg.Registries[i].URL = newURL
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("registry '%s' not found", name)
	}

	return config.SaveConfig(cfg)
}
