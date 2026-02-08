package cmd

import (
	"os"
	"time"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
	"github.com/tuannvm/mcpenetes/internal/config"
	"github.com/tuannvm/mcpenetes/internal/core"
	"github.com/tuannvm/mcpenetes/internal/log"
)

// restoreCmd represents the restore command
var restoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "Restores client configurations from the latest backups.",
	Long: `Restores the configuration files for all defined clients 
from the most recent backup found in the backup directory specified in config.yaml.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("Executing restore command...")

		// 1. Load config
		cfg, err := config.LoadConfig()
		if err != nil {
			log.Fatal("Error loading config.yaml: %v", err)
		}

		mcpCfg, err := config.LoadMCPConfig()
		if err != nil {
			log.Fatal("Error loading mcp.json: %v", err)
		}

		manager := core.NewManager(cfg, mcpCfg)

		// 2. Perform restore
		log.Info("Restoring client configurations:")
		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Suffix = " Restoring..."
		s.Start()

		restored, errors := manager.RestoreAllLatest()

		s.Stop()

		// 3. Log results
		successCount := 0
		failureCount := 0

		if len(cfg.Clients) == 0 {
			log.Warn("No clients defined in configuration.")
			return
		}

		for clientName := range cfg.Clients {
			if backupFile, ok := restored[clientName]; ok {
				log.Success("- %s: Successfully restored from %s", clientName, backupFile)
				successCount++
			} else if err, ok := errors[clientName]; ok {
				log.Error("- %s: Failed restore - %v", clientName, err)
				failureCount++
			} else {
				// Not in restored or errors implies no backup found (or skipped)
				log.Warn("- %s: No backups found to restore.", clientName)
			}
		}

		// Also check for global errors (key "all")
		if err, ok := errors["all"]; ok {
			log.Fatal("Failed to list backups: %v", err)
		}

		log.Info("\nRestore finished.")
		if successCount > 0 {
			log.Success("Successfully restored %d clients.", successCount)
		}
		if failureCount > 0 {
			log.Error("Failed to restore %d clients.", failureCount)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(restoreCmd)
}
