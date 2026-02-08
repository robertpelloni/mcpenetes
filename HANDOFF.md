# 🧙‍♂️ mcpenetes Handoff Documentation

## 📅 Session Summary
**Date:** 2025-05-27
**Status:** Major refactor completed. Web UI enhanced with comprehensive features. Version 1.3.1.

This session focused on transforming `mcpenetes` from a simple CLI tool into a robust configuration manager with a full-featured Web UI. We implemented backup/restore, custom client management, system diagnostics, and extensive documentation.

## 🏗️ Architectural Changes

### 1. Client Registry (`internal/client`)
We moved away from hardcoded detection logic in `util` to a data-driven **Registry** approach.
-   **`Registry`**: A slice of `ClientDefinition` structs defining the Tool ID, Name, Config Format, OS-specific paths, and optional `ConfigKey` overrides.
-   **User-Defined Registry**: The tool now automatically loads custom client definitions from `~/.config/mcpenetes/clients.yaml` (or equivalent on Windows), allowing users to support new tools without waiting for a release.
-   **`custom.go`**: Logic to programmatically add/remove entries from `clients.yaml`.

### 2. Core Logic (`internal/core`)
-   **`Manager`**: Encapsulates the logic for `ApplyToClient` (Backup -> Translate -> Apply -> Cleanup).
-   Decoupled from `cobra` CLI commands to allow reuse by the Web UI.

### 3. Web UI (`internal/ui`, `cmd/ui.go`)
-   **Embedded Server**: Uses Go `embed` to serve a static HTML/JS frontend.
-   **Features**:
    -   **Dashboard**: View detected clients and configured servers. Inspect server commands, edit configurations, or delete servers.
    -   **Search & Install**: Find new servers from configured registries and install them with one click.
    -   **Clients**: Manage custom client definitions.
    -   **Backups**: View a history of configuration backups for each client and restore them if needed.
    -   **Logs**: Real-time application log viewer.
    -   **System**: View version info, build details, and project structure.
    -   **Import Config**: Easily import an existing `mcpServers` JSON configuration by pasting it into the UI.
    -   **Help**: Built-in documentation and troubleshooting guides.
-   **API Endpoints**: Comprehensive set of endpoints for all UI features.

### 4. Robust Translation (`internal/translator`)
-   **JSONC Support**: Integrated `github.com/tailscale/hujson` to safely parse VS Code `settings.json` files containing comments.
-   **Safety**: The translator now aborts if the existing config file cannot be parsed, preventing data loss.
-   **Custom Keys**: Supports injecting MCP configurations into custom JSON keys (e.g., `openctx.providers` for Cody) via the `ConfigKey` property.

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
| `continue` | Continue | Custom | VS Code Extension |
| `pearai` | PearAI | JSON | VS Code Fork |
| `void` | Void | VSCode-JSON | VS Code Fork |
| `melty` | Melty | VSCode-JSON | VS Code Fork |
| `codebuddy` | CodeBuddy | VSCode-JSON | VS Code Fork |
| `kiro` | Kiro | JSON | |
| `codegpt` | CodeGPT | JSON | VS Code Extension |
| `cody` | Cody (Sourcegraph) | VSCode-JSON | Targets `openctx.providers` |
| `5ire` | 5ire | JSON | |
| `lm-studio` | LM Studio | JSON | |
| `anythingllm` | AnythingLLM | JSON | |
| `tabby` | Tabby | TOML | |
| `librechat` | LibreChat | YAML | |
| `jan` | Jan | JSON | |
| `boltai` | BoltAI | JSON | |
| `goose` | Goose CLI | YAML | |
| `mistral-vibe` | Mistral Vibe | TOML | |
| `code-cli` | Code CLI (Codex) | JSON | |
| `grok-cli` | Grok CLI | JSON | |
| `open-interpreter`| Open Interpreter | YAML | |
| `factory-cli` | Factory CLI | JSON | |
| `aider` | Aider | YAML | |
| `warp` | Warp Terminal | JSON | |
| `llm-cli` | LLM CLI | JSON | Simon Willison's tool |
| `claude-code` | Claude Code CLI | JSON | `~/.claude.json` |
| `amazon-q` | Amazon Q | JSON | CodeWhisperer |

## 🧠 Findings & Decisions

1.  **Windows Path Priority**: The registry prioritizes Windows paths (AppData/UserProfile) as requested, but falls back to Home for cross-platform compatibility.
2.  **Detection Heuristic**: We detect clients by checking for the *config file first*. If missing, we check for the *parent directory*. This allows us to configure tools that are installed but haven't generated a config file yet (fresh installs).
3.  **Search Workflow**: The CLI `search` command previously only updated `config.yaml` (legacy list). We refactored it to update `mcp.json` directly with a default `npx` configuration, making the "Search -> Apply" workflow functional.
4.  **Race Conditions**: When multiple clients (e.g., VS Code and Cody) target the same file (`settings.json`), sequential processing is enforced to avoid data corruption.
5.  **Logging Safety**: Implemented a `sync.Mutex` in the logging buffer to prevent data races during concurrent web requests.

## 🚀 Future Roadmap

1.  **Cloud Sync**: Synchronize configurations across machines via GitHub Gists or a cloud backend.
2.  **Plugin System**: Allow third-party plugins to extend functionality (e.g., custom translators).
3.  **MCP Server Marketplace**: Integration with a centralized registry for one-click installation of community servers.
4.  **Docker Integration**: Native support for managing Docker-based MCP servers.

## 📝 Memories
-   The project uses `github.com/tailscale/hujson` to parse JSONC.
-   `internal/client` is the source of truth for tool detection.
-   `internal/core` handles the business logic.
-   The UI runs on port 3000 by default (`mcpenetes ui`).
