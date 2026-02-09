package client

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"gopkg.in/yaml.v3"
)

var customClientsMutex sync.Mutex

// GetCustomRegistryPath returns the path to clients.yaml
func GetCustomRegistryPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	configDir := filepath.Join(homeDir, ".config", "mcpetes")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", err
	}
	return filepath.Join(configDir, UserRegistryFile), nil
}

// LoadCustomClients reads the user-defined clients from clients.yaml
func LoadCustomClients() ([]ClientDefinition, error) {
	path, err := GetCustomRegistryPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return []ClientDefinition{}, nil
	}
	if err != nil {
		return nil, err
	}

	var clients []ClientDefinition
	if err := yaml.Unmarshal(data, &clients); err != nil {
		return nil, fmt.Errorf("failed to parse clients.yaml: %w", err)
	}
	return clients, nil
}

// SaveCustomClients writes the list of clients to clients.yaml
func SaveCustomClients(clients []ClientDefinition) error {
	customClientsMutex.Lock()
	defer customClientsMutex.Unlock()

	path, err := GetCustomRegistryPath()
	if err != nil {
		return err
	}

	data, err := yaml.Marshal(clients)
	if err != nil {
		return fmt.Errorf("failed to marshal clients: %w", err)
	}

	return os.WriteFile(path, data, 0644)
}

// AddCustomClient adds a new client definition to the custom registry
func AddCustomClient(client ClientDefinition) error {
	clients, err := LoadCustomClients()
	if err != nil {
		return err
	}

	// Check for duplicates
	for _, c := range clients {
		if c.ID == client.ID {
			return fmt.Errorf("client with ID '%s' already exists", client.ID)
		}
	}

	clients = append(clients, client)
	return SaveCustomClients(clients)
}

// RemoveCustomClient removes a client definition by ID
func RemoveCustomClient(id string) error {
	clients, err := LoadCustomClients()
	if err != nil {
		return err
	}

	newClients := []ClientDefinition{}
	found := false
	for _, c := range clients {
		if c.ID == id {
			found = true
			continue
		}
		newClients = append(newClients, c)
	}

	if !found {
		return fmt.Errorf("client with ID '%s' not found", id)
	}

	return SaveCustomClients(newClients)
}
