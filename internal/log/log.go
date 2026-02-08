package log

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/fatih/color"
)

// Predefine color functions for different log levels
var (
	InfoColor    = color.New(color.FgCyan)
	SuccessColor = color.New(color.FgGreen)
	WarnColor    = color.New(color.FgYellow)
	ErrorColor   = color.New(color.FgRed)
	DetailColor  = color.New(color.FgWhite) // For less important details
)

// Buffer Support
type LogEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Level     string    `json:"level"`
	Message   string    `json:"message"`
}

// Simple ring buffer implementation
const BufferSize = 1000

var (
	logEntries = make([]LogEntry, 0, BufferSize)
	logMutex   sync.Mutex
)

func addToBuffer(level, format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	entry := LogEntry{
		Timestamp: time.Now(),
		Level:     level,
		Message:   msg,
	}

	logMutex.Lock()
	defer logMutex.Unlock()

	if len(logEntries) >= BufferSize {
		logEntries = logEntries[1:]
	}
	logEntries = append(logEntries, entry)
}

func GetRecentLogs() []LogEntry {
	logMutex.Lock()
	defer logMutex.Unlock()
	// Return copy
	result := make([]LogEntry, len(logEntries))
	copy(result, logEntries)
	return result
}

// Info prints an informational message (cyan).
func Info(format string, a ...interface{}) {
	addToBuffer("INFO", format, a...)
	_, _ = InfoColor.Fprintf(os.Stdout, format+"\n", a...)
}

// Success prints a success message (green).
func Success(format string, a ...interface{}) {
	addToBuffer("SUCCESS", format, a...)
	_, _ = SuccessColor.Fprintf(os.Stdout, format+"\n", a...)
}

// Warn prints a warning message (yellow) to stderr.
func Warn(format string, a ...interface{}) {
	addToBuffer("WARN", format, a...)
	_, _ = WarnColor.Fprintf(os.Stderr, "Warning: "+format+"\n", a...)
}

// Error prints an error message (red) to stderr.
func Error(format string, a ...interface{}) {
	addToBuffer("ERROR", format, a...)
	_, _ = ErrorColor.Fprintf(os.Stderr, "Error: "+format+"\n", a...)
}

// Fatal prints an error message (red) to stderr and exits with status 1.
func Fatal(format string, a ...interface{}) {
	addToBuffer("FATAL", format, a...)
	Error(format, a...)
	os.Exit(1)
}

// Detail prints less important details (usually white/default).
func Detail(format string, a ...interface{}) {
	addToBuffer("DETAIL", format, a...)
	_, _ = DetailColor.Fprintf(os.Stdout, format+"\n", a...)
}

// Printf allows printing with a specific color.
func Printf(c *color.Color, format string, a ...interface{}) {
	// Not buffering generic printf as it might not be a log entry
	_, _ = c.Printf(format, a...)
}

// Fprintf allows printing to a specific writer with a specific color.
func Fprintf(w *os.File, c *color.Color, format string, a ...interface{}) {
	_, _ = c.Fprintf(w, format, a...)
}
