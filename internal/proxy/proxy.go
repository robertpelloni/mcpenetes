package proxy

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"
)

// Options defines the configuration for the proxy
type Options struct {
	Command   string
	Args      []string
	ServerID  string // Used for log file naming
	LogDir    string
}

// Run starts the proxy execution
func Run(opts Options) error {
	// 1. Setup Log File
	if opts.LogDir == "" {
		home, _ := os.UserHomeDir()
		opts.LogDir = filepath.Join(home, ".config", "mcpetes", "logs")
	}
	if err := os.MkdirAll(opts.LogDir, 0755); err != nil {
		return fmt.Errorf("failed to create log dir: %w", err)
	}

	// Use ServerID for filename if provided, otherwise a generic timestamped name
	// In a real usage, we might want a persistent name per server to tail it easily.
	logFileName := "proxy.log"
	if opts.ServerID != "" {
		logFileName = fmt.Sprintf("%s.log", opts.ServerID)
	}
	logPath := filepath.Join(opts.LogDir, logFileName)

	logFile, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}
	defer logFile.Close()

	// 2. Prepare Command
	cmd := exec.Command(opts.Command, opts.Args...)

	// Pipes
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	// 3. Start Command
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start command: %w", err)
	}

	// 4. Handle Signals
	// Forward signals to the child process
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigChan
		cmd.Process.Signal(sig)
	}()

	// 5. IO Forwarding
	var wg sync.WaitGroup
	wg.Add(3)

	// Stdin: OS Stdin -> Cmd Stdin
	go func() {
		defer wg.Done()
		defer stdin.Close()
		io.Copy(stdin, os.Stdin)
	}()

	// Stdout: Cmd Stdout -> OS Stdout (Transparent JSON-RPC)
	// We do NOT log stdout to keep the protocol clean and performant,
	// unless we implement a full inspector/parser.
	go func() {
		defer wg.Done()
		io.Copy(os.Stdout, stdout)
	}()

	// Stderr: Cmd Stderr -> OS Stderr AND Log File
	go func() {
		defer wg.Done()
		// MultiWriter to write to both original stderr and our log file
		writer := io.MultiWriter(os.Stderr, logFile)

		// Add a header to indicate new session
		fmt.Fprintf(logFile, "\n--- Session Started: %s ---\n", time.Now().Format(time.RFC3339))

		io.Copy(writer, stderr)
	}()

	// 6. Wait
	wg.Wait()
	return cmd.Wait()
}
