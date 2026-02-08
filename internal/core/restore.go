package core

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/tuannvm/mcpenetes/internal/util"
)

// BackupFile represents a backup file for a client
type BackupFile struct {
	Name      string    `json:"name"`
	Path      string    `json:"path"`
	Timestamp time.Time `json:"timestamp"`
}

// ListBackups returns a map of client names to their backup files, sorted by timestamp descending
func (m *Manager) ListBackups() (map[string][]BackupFile, error) {
	backupDir, err := util.ExpandPath(m.Config.Backups.Path)
	if err != nil {
		return nil, fmt.Errorf("error expanding backup path: %w", err)
	}

	entries, err := os.ReadDir(backupDir)
	if err != nil {
		if os.IsNotExist(err) {
			return make(map[string][]BackupFile), nil
		}
		return nil, fmt.Errorf("error reading backup directory: %w", err)
	}

	backups := make(map[string][]BackupFile)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		fileName := entry.Name()
		// Format: <clientName>-<timestamp>.<ext>
		// Timestamp: YYYYMMDD-HHMMSS (15 chars)
		// We need to robustly parse this.

		// Find the last occurrence of "-" before the extension?
		// Actually the timestamp is fixed length 15 chars.
		// Let's assume the format is strictly clientName-YYYYMMDD-HHMMSS.ext

		ext := filepath.Ext(fileName)
		nameWithoutExt := strings.TrimSuffix(fileName, ext)

		// Extract timestamp from end (15 chars)
		if len(nameWithoutExt) <= 16 { // at least 1 char name + hyphen + 15 char timestamp
			continue
		}

		timestampStr := nameWithoutExt[len(nameWithoutExt)-15:]
		if nameWithoutExt[len(nameWithoutExt)-16] != '-' {
			continue // Expected hyphen before timestamp
		}

		clientName := nameWithoutExt[:len(nameWithoutExt)-16]

		// Parse timestamp
		ts, err := time.Parse("20060102-150405", timestampStr)
		if err != nil {
			// If strict parsing fails, fallback to file mod time or skip
			ts = info.ModTime()
		}

		backups[clientName] = append(backups[clientName], BackupFile{
			Name:      fileName,
			Path:      filepath.Join(backupDir, fileName),
			Timestamp: ts,
		})
	}

	// Sort backups by timestamp descending
	for clientName := range backups {
		clientBackups := backups[clientName]
		sort.Slice(clientBackups, func(i, j int) bool {
			return clientBackups[i].Timestamp.After(clientBackups[j].Timestamp)
		})
		backups[clientName] = clientBackups
	}

	return backups, nil
}

// RestoreClient restores a specific backup file for a client
func (m *Manager) RestoreClient(clientName, backupFileName string) error {
	clientConf, ok := m.Config.Clients[clientName]
	if !ok {
		return fmt.Errorf("client '%s' not found in configuration", clientName)
	}

	backupDir, err := util.ExpandPath(m.Config.Backups.Path)
	if err != nil {
		return fmt.Errorf("error expanding backup path: %w", err)
	}

	backupPath := filepath.Join(backupDir, backupFileName)

	// Verify backup file exists
	if _, err := os.Stat(backupPath); err != nil {
		return fmt.Errorf("backup file '%s' not found: %w", backupFileName, err)
	}

	clientConfigPath, err := util.ExpandPath(clientConf.ConfigPath)
	if err != nil {
		return fmt.Errorf("error expanding client config path: %w", err)
	}

	// Perform copy
	return copyFile(backupPath, clientConfigPath)
}

// RestoreAllLatest restores the latest backup for every client
func (m *Manager) RestoreAllLatest() (map[string]string, map[string]error) {
	backups, err := m.ListBackups()
	if err != nil {
		return nil, map[string]error{"all": err}
	}

	restored := make(map[string]string)
	errors := make(map[string]error)

	for clientName, clientBackups := range backups {
		if len(clientBackups) == 0 {
			continue
		}

		latest := clientBackups[0]
		if err := m.RestoreClient(clientName, latest.Name); err != nil {
			errors[clientName] = err
		} else {
			restored[clientName] = latest.Name
		}
	}

	return restored, errors
}

// copyFile copies a file from src to dst.
func copyFile(src, dst string) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	// Ensure destination directory exists
	dstDir := filepath.Dir(dst)
	if err := os.MkdirAll(dstDir, 0750); err != nil {
		return fmt.Errorf("failed to create destination directory '%s': %w", dstDir, err)
	}

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	return err
}
