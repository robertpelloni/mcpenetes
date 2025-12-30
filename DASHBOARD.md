# 📊 Project Dashboard

## ℹ️ Project Overview
**Name:** mcpenetes
**Version:** 1.4.1
**Description:** A configuration manager for Model Context Protocol (MCP) servers, supporting over 30 clients including VS Code, Claude, and various CLIs. Includes a Web UI for management, log viewing, and registry editing.

## 📂 Project Structure

```
.
├── cmd/                    # CLI command implementations (root, ui, search, proxy, etc.)
├── internal/
│   ├── client/             # Client registry and detection logic (registry.go)
│   ├── config/             # Configuration structs and file handling
│   ├── core/               # Core business logic (Manager)
│   ├── doctor/             # System diagnostic checks
│   ├── log/                # Logging utilities
│   ├── proxy/              # Proxy wrapper logic for capturing server logs
│   ├── registry/           # Registry client logic
│   │   └── manager/        # Registry management logic (Add/Update/Remove)
│   ├── search/             # Search functionality logic
│   ├── translator/         # Logic to translate/apply configs to clients (JSONC support)
│   ├── ui/                 # Web server (server.go) and embedded frontend
│   │   └── static/         # HTML/CSS/JS assets
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
- **Binaries:** Output to `bin/` (ignored by git)
