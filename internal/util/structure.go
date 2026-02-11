package util

import (
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

// GenerateProjectStructure generates a tree-like string representation of the project
func GenerateProjectStructure(root string) (string, error) {
	// We delegate purely to generateTree to avoid redundancy,
	// passing the root name "." for the initial call if desired, or handle it inside.
	return generateTree(root, "", true)
}

func generateTree(path string, prefix string, isLast bool) (string, error) {
	stats, err := os.Stat(path)
	if err != nil {
		return "", err
	}

	// Check if directory
	if !stats.IsDir() {
		return "", nil
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return "", err
	}

	// Filter entries
	var visible []os.DirEntry
	for _, e := range entries {
		// Ignore hidden files, node_modules, bin folder, and the binary itself
		if strings.HasPrefix(e.Name(), ".") || e.Name() == "node_modules" || e.Name() == "bin" || e.Name() == "mcpenetes" {
			continue
		}
		visible = append(visible, e)
	}

	sort.Slice(visible, func(i, j int) bool {
		return visible[i].Name() < visible[j].Name()
	})

	var sb strings.Builder
	// If this is the root call (implied by empty prefix?), we assume "." was printed by caller or we start with content.
	// But let's just stick to the tree structure.

	if path == "." || prefix == "" {
		sb.WriteString(".\n")
	}

	for i, entry := range visible {
		isLastChild := i == len(visible)-1
		connector := "├── "
		if isLastChild {
			connector = "└── "
		}

		sb.WriteString(prefix + connector + entry.Name())

		// Add descriptions for known dirs
		if entry.IsDir() {
			desc := ""
			switch entry.Name() {
			case "cmd": desc = " # CLI commands"
			case "internal": desc = " # Internal packages"
			case "client": desc = " # Client registry"
			case "config": desc = " # Config management"
			case "core": desc = " # Core logic"
			case "doctor": desc = " # Health checks"
			case "log": desc = " # Logging"
			case "mcp": desc = " # MCP connectivity"
			case "registry": desc = " # Registry client"
			case "search": desc = " # Search logic"
			case "sync": desc = " # Cloud sync"
			case "translator": desc = " # Config translation"
			case "ui": desc = " # Web server & UI"
			case "util": desc = " # Utilities"
			case "version": desc = " # Version info"
			case "docs": desc = " # Documentation"
			}
			sb.WriteString(desc)
		}
		sb.WriteString("\n")

		if entry.IsDir() {
			newPrefix := prefix + "│   "
			if isLastChild {
				newPrefix = prefix + "    "
			}
			subTree, _ := generateTree(filepath.Join(path, entry.Name()), newPrefix, isLastChild)
			sb.WriteString(subTree)
		}
	}

	return sb.String(), nil
}

// GetSubmodules returns a list of git submodules status
func GetSubmodules() []string {
	// Try running git submodule status
	cmd := exec.Command("git", "submodule", "status")
	out, err := cmd.Output()
	if err != nil {
		// Not a git repo or no git installed
		return []string{}
	}

	lines := strings.Split(string(out), "\n")
	var result []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
