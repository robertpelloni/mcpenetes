package util

import (
	"github.com/tuannvm/mcpenetes/internal/client"
	"github.com/tuannvm/mcpenetes/internal/config"
)

// DetectMCPClients automatically detects installed MCP-compatible clients
// and their configuration paths on the user's system
func DetectMCPClients() (map[string]config.Client, error) {
	detected, err := client.DetectClients()
	if err != nil {
		return nil, err
	}

	result := make(map[string]config.Client)
	for id, c := range detected {
		result[id] = config.Client{
			ConfigPath: c.ConfigPath,
			Type:       string(c.ConfigFormat),
			Key:        c.ConfigKey,
		}
	}

	return result, nil
}
