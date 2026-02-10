# 📘 mcpenetes User Manual

## Table of Contents
1.  [Introduction](#introduction)
2.  [Installation](#installation)
3.  [Quick Start](#quick-start)
4.  [Web UI Dashboard](#web-ui-dashboard)
    *   [Dashboard](#dashboard)
    *   [Search & Install](#search--install)
    *   [Registries](#registries)
    *   [Clients](#clients)
    *   [Backups](#backups)
    *   [Logs](#logs)
    *   [System](#system)
    *   [Settings](#settings)
5.  [CLI Reference](#cli-reference)
6.  [Advanced Configuration](#advanced-configuration)
7.  [Troubleshooting](#troubleshooting)

---

## Introduction

**mcpenetes** is a universal configuration manager for the Model Context Protocol (MCP) ecosystem. It allows you to centrally manage MCP servers and automatically sync their configurations across multiple supported clients, such as VS Code, Claude Desktop, Cursor, and many others.

Key features:
*   **Centralized Management:** Define servers once, apply everywhere.
*   **Web UI:** A comprehensive visual dashboard for all operations.
*   **Safety:** Automatic backups and configuration validation.
*   **Extensibility:** Support for custom clients and registries.

## Installation

### From Source
```bash
git clone https://github.com/tuannvm/mcpenetes.git
cd mcpenetes
make build
# Binary is at ./bin/mcpenetes
```

### Using Go Install
```bash
go install github.com/tuannvm/mcpenetes@latest
```

## Quick Start

1.  **Start the Web UI:**
    ```bash
    mcpenetes ui
    ```
    This opens the dashboard at `http://localhost:3000`.

2.  **Find a Server:**
    Go to the **Search** tab and type a keyword (e.g., "filesystem"). Click **Install** on a server.

3.  **Apply Configuration:**
    Go to the **Dashboard** tab. Select the clients you want to configure (checked by default if detected) and click **Apply Configuration**.

## Web UI Dashboard

### Dashboard
The command center.
*   **Detected Clients:** Lists all supported tools found on your system. Uncheck any you don't want to manage.
*   **Configured MCP Servers:** Lists currently selected servers.
    *   **Inspect:** View the raw command line or JSON.
    *   **Test:** (⚡) Run a connectivity test (ping) to ensure the server starts and responds to MCP handshake.
    *   **Edit:** Modify configuration (command, args, env).
    *   **Delete:** Remove the server from your configuration.
*   **Apply:** Writes the configuration to all selected clients.

### Search & Install
Find servers from configured registries (Glama, Smithery, etc.).
*   **Search:** Type keywords to filter servers.
*   **Install:** Click to add a server. If the server requires environment variables (e.g., API keys), a form will appear for you to fill them in.

### Registries
Manage where `mcpenetes` looks for servers.
*   **Add:** Add a new registry URL (e.g., a private company registry).
*   **Edit/Remove:** Modify existing sources.

### Clients
Manage supported tools.
*   **Known Clients:** A list of all 30+ built-in supported clients (VS Code, Claude, etc.) and their detection status.
*   **Configure:** If a client is installed but not detected (e.g., non-standard path), click **Configure** to manually set its configuration file path.
*   **Custom Clients:** Define completely new tools by specifying their config file path and format.

### Backups
Safety first.
*   Every time you click **Apply**, a backup of the client's previous config is created.
*   **Restore:** Revert to a previous state instantly.
*   **Delete:** Remove old backups. (Automatic retention policy is also available in Settings).

### Logs
View internal application logs for debugging `mcpenetes` itself. Useful if an "Apply" operation fails.

### System
**Doctor:** Runs diagnostic checks to ensure your system is healthy (permissions, dependencies like `npx` or `python`).
**Info:** Displays build version and environment details.

### Settings
Global configuration.
*   **Backup Storage Path:** Where backups are saved.
*   **Retention Policy:** How many backups to keep per client before deleting the oldest ones.
*   **Global Environment Variables:** Define environment variables (e.g., `OPENAI_API_KEY`) here to inject them into **all** server configurations automatically. Server-specific variables override these.

## CLI Reference

*   `mcpenetes ui`: Start the Web Dashboard.
*   `mcpenetes search`: Interactive CLI search.
*   `mcpenetes apply`: Apply configuration to all clients.
*   `mcpenetes load`: Import configuration from clipboard.
*   `mcpenetes restore`: Restore latest backups.
*   `mcpenetes doctor`: Run health checks.

## Advanced Configuration

### Configuration Files
Located in `~/.config/mcpetes/` (or `%USERPROFILE%\.config\mcpetes\` on Windows).
*   `config.yaml`: Main settings, registries, and selected servers.
*   `clients.yaml`: Custom client definitions.
*   `mcp.json`: The raw MCP server configuration.

### Custom Clients via YAML
You can manually add clients to `clients.yaml`:
```yaml
- id: my-tool
  name: My Tool
  config_format: simple-json # or vscode, yaml, toml
  paths:
    darwin:
      - base: home
        path: .config/mytool/config.json
```

## Troubleshooting

*   **Server not starting?** Use the **Test** button in the Dashboard to check for errors. Ensure you have the runtime installed (e.g., Node.js for `npx`, Python for `uvx`/`python`).
*   **Client not detected?** Check the **Clients** tab and try manually configuring the path.
*   **Permission denied?** Run `mcpenetes doctor` to check file permissions.
