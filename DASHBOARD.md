# 📊 Project Dashboard

## ℹ️ Project Overview
**Name:** mcpenetes
**Version:** 1.1.0
**Description:** A configuration manager for Model Context Protocol (MCP) servers, supporting over 30 clients including VS Code, Claude, and various CLIs.

## 📂 Project Structure

```
.
├── cmd/                    # CLI command implementations (root, ui, search, etc.)
├── internal/
│   ├── client/             # Client registry and detection logic (registry.go)
│   ├── config/             # Configuration structs and file handling
│   ├── core/               # Core business logic (Manager)
│   ├── doctor/             # System diagnostic checks
│   ├── search/             # Search functionality logic
│   ├── translator/         # Logic to translate/apply configs to clients
│   ├── ui/                 # Web server and frontend assets
│   ├── util/               # Utility functions
│   └── version/            # Version information
├── main.go                 # Entry point
└── ... (Docs & Configs)
```

## 🧩 Submodules
*No git submodules are currently used in this repository.*

## 🛠️ Build Information
- **Language:** Go 1.23+
- **Build System:** Makefile / `go build`
- **Frontend:** Embedded static HTML/JS (Pico.css)
