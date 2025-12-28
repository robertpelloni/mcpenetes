package version

import (
	_ "embed"
	"strings"
)

//go:embed VERSION
var version string

// Version is the current version of the application
var Version = strings.TrimSpace(version)
