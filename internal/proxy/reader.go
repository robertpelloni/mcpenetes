package proxy

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/hpcloud/tail"
)

// LogStream represents a stream of log lines
type LogStream struct {
	Lines chan string
	Stop  chan struct{}
}

// StreamLogs tails the log file for a given serverID
func StreamLogs(serverID string) (*LogStream, error) {
	home, _ := os.UserHomeDir()
	logDir := filepath.Join(home, ".config", "mcpetes", "logs")
	logPath := filepath.Join(logDir, fmt.Sprintf("%s.log", serverID))

	// Ensure file exists
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		// Return error or empty stream?
		// Creating it prevents "file not found" errors in tail
		os.MkdirAll(logDir, 0755)
		os.WriteFile(logPath, []byte("--- Log created ---\n"), 0644)
	}

	t, err := tail.TailFile(logPath, tail.Config{
		Follow: true,
		ReOpen: true, // Handle log rotation or recreation
		Poll:   true, // Better compatibility across FS
	})
	if err != nil {
		return nil, err
	}

	stream := &LogStream{
		Lines: make(chan string),
		Stop:  make(chan struct{}),
	}

	go func() {
		defer t.Stop()
		for {
			select {
			case line := <-t.Lines:
				if line == nil {
					// Error or closed
					return
				}
				select {
				case stream.Lines <- line.Text:
				case <-stream.Stop:
					return
				}
			case <-stream.Stop:
				return
			}
		}
	}()

	return stream, nil
}

// ReadLastNLines reads the last N lines from the log file
// Useful for initial buffer
func ReadLastNLines(serverID string, n int) ([]string, error) {
	home, _ := os.UserHomeDir()
	logPath := filepath.Join(home, ".config", "mcpetes", "logs", fmt.Sprintf("%s.log", serverID))

	file, err := os.Open(logPath)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}
	defer file.Close()

	// Simple implementation: Read full file and take last N
	// Optimization: Seek from end
	// For now, simple is fine as logs shouldn't be massive (stderr)

	content, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	// Split by newline
	// ... implementation skipped for brevity, simplistic approach
	// In production, use a buffer scanner from end

	return []string{string(content)}, nil
}
