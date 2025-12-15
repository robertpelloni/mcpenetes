# 🧙‍♂️ mcpenetes Handoff Documentation

## 📅 Session Summary
**Date:** 2025-05-27
**Status:** Major refactor completed. Web UI added. Extensive tool support added.

This session focused on transforming `mcpenetes` from a simple CLI tool into a comprehensive configuration manager for the Model Context Protocol (MCP) ecosystem. We addressed the user's need to support a vast array of AI tools (IDEs, CLIs, Desktop Apps) and provided a graphical interface.

## 🏗️ Architectural Changes

### 1. Client Registry (`internal/client`)
We moved away from hardcoded detection logic in `util` to a data-driven **Registry** approach.
-   **`Registry`**: A slice of `ClientDefinition` structs defining the Tool ID, Name, Config Format, and OS-specific paths.
-   **`ConfigFormatEnum`**: explicit support for:
    -   `simple-json`: Standard `{"mcpServers": {...}}` (Claude, Cursor, etc.)
    -   `vscode`: Nested `{"mcp": {"servers": {...}}}`
    -   `claude-desktop`: Specific handling for Claude Desktop.
    -   `yaml`: For tools like Goose CLI.
    -   `toml`: For tools like Mistral Vibe.
-   **`PathDefinition`**: Supports paths relative to `BaseHome`, `BaseAppData`, and `BaseUserProfile`.

### 2. Core Logic (`internal/core`)
-   **`Manager`**: Encapsulates the logic for `ApplyToClient` (Backup -> Translate -> Apply -> Cleanup).
-   Decoupled from `cobra` CLI commands to allow reuse by the Web UI.

### 3. Web UI (`internal/ui`, `cmd/ui.go`)
-   **Embedded Server**: Uses Go `embed` to serve a static HTML/JS frontend.
-   **API Endpoints**:
    -   `GET /api/data`: Returns clients, servers, and registries.
    -   `POST /api/apply`: Applies configs to selected clients.
    -   `POST /api/install`: Adds a server from the registry to `mcp.json` (defaults to `npx` execution).
    -   `POST /api/server/update`: Updates server config directly (Edit feature).
-   **Frontend**: Single-page dashboard using Pico.css.

### 4. Robust Translation (`internal/translator`)
-   **JSONC Support**: Integrated `github.com/tailscale/hujson` to safely parse VS Code `settings.json` files containing comments.
-   **Safety**: The translator now aborts if the existing config file cannot be parsed, preventing data loss.

## 🛠️ Supported Clients (Built-in)

The following clients are currently supported in `internal/client/registry.go`:

| ID | Name | Format | Notes |
| :--- | :--- | :--- | :--- |
| `claude-desktop` | Claude Desktop | JSON | |
| `cursor` | Cursor | JSON | |
| `windsurf` | Windsurf | JSON | |
| `vscode` | VS Code | VSCode-JSON | Supports `settings.json` with comments |
| `vscode-insiders`| VS Code Insiders | VSCode-JSON | |
| `zed` | Zed | JSON | |
| `trae` | Trae | JSON | |
| `jetbrains-junie`| JetBrains (Junie)| JSON | Detects `~/.junie/mcp/mcp.json` |
| `cline` | Cline | JSON | VS Code Extension |
| `roo-code` | Roo Code | JSON | VS Code Extension |
| `lm-studio` | LM Studio | JSON | |
| `anythingllm` | AnythingLLM | JSON | |
| `goose` | Goose CLI | YAML | |
| `mistral-vibe` | Mistral Vibe | TOML | |
| `code-cli` | Code CLI (Codex) | JSON | |
| `grok-cli` | Grok CLI | JSON | |
| `open-interpreter`| Open Interpreter | YAML | |
| `factory-cli` | Factory CLI | JSON | |

## 🧠 Findings & Decisions

1.  **Windows Path Priority**: The registry prioritizes Windows paths (AppData/UserProfile) as requested, but falls back to Home for cross-platform compatibility.
2.  **Detection Heuristic**: We detect clients by checking for the *config file first*. If missing, we check for the *parent directory*. This allows us to configure tools that are installed but haven't generated a config file yet (fresh installs).
3.  **Search Workflow**: The CLI `search` command previously only updated `config.yaml` (legacy list). We refactored it to update `mcp.json` directly with a default `npx` configuration, making the "Search -> Apply" workflow functional.

## 🚀 Future Roadmap

1.  **User-Defined Registry**: Implement loading of `clients.yaml` from the user's config directory. This enables users to add support for the "long tail" of tools (Adrenaline, Aider, etc.) without modifying the source code.
2.  **Expanded Built-in Support**: Continue adding definitions for popular tools from the user's CSV list (Aider, PearAI, Void).
3.  **UI Enhancements**:
    -   Log viewer for the MCP servers? (Hard since they run inside the clients).
    -   Better "Install" wizard that asks for `command` and `args` instead of defaulting to `npx`.

## 📝 Memories
-   The project uses `github.com/tailscale/hujson` to parse JSONC.
-   `internal/client` is the source of truth for tool detection.
-   `internal/core` handles the business logic.
-   The UI runs on port 3000 by default (`mcpenetes ui`).
