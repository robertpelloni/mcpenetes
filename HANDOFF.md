# 🧙‍♂️ mcpenetes Handoff Documentation

## 📅 Session Summary
**Date:** 2025-05-27
**Status:** Major refactor completed. Web UI added. Extensive tool support added (30+ clients). User extensibility enabled.

This session focused on transforming `mcpenetes` from a simple CLI tool into a comprehensive configuration manager for the Model Context Protocol (MCP) ecosystem. We addressed the user's need to support a vast array of AI tools (IDEs, CLIs, Desktop Apps) and provided a graphical interface.

## 🏗️ Architectural Changes

### 1. Client Registry (`internal/client`)
We moved away from hardcoded detection logic in `util` to a data-driven **Registry** approach.
-   **`Registry`**: A slice of `ClientDefinition` structs defining the Tool ID, Name, Config Format, OS-specific paths, and optional `ConfigKey` overrides.
-   **User-Defined Registry**: The tool now automatically loads custom client definitions from `~/.config/mcpetes/clients.yaml` (or equivalent on Windows), allowing users to support new tools without waiting for a release.
-   **`ConfigFormatEnum`**: explicit support for:
    -   `simple-json`: Standard `{"mcpServers": {...}}` (Claude, Cursor, etc.)
    -   `vscode`: Nested `{"mcp": {"servers": {...}}}`. Can be overridden with `ConfigKey`.
    -   `claude-desktop`: Specific handling for Claude Desktop.
    -   `yaml`: For tools like Goose CLI.
    -   `toml`: For tools like Mistral Vibe.
    -   `continue`: For the Continue extension's nested array format.
-   **`PathDefinition`**: Supports paths relative to `BaseHome`, `BaseAppData`, and `BaseUserProfile`.

### 2. Core Logic (`internal/core`)
-   **`Manager`**: Encapsulates the logic for `ApplyToClient` (Backup -> Translate -> Apply -> Cleanup).
-   Decoupled from `cobra` CLI commands to allow reuse by the Web UI.

### 3. Web UI (`internal/ui`, `cmd/ui.go`)
-   **Embedded Server**: Uses Go `embed` to serve a static HTML/JS frontend.
-   **Safety**: Implements sequential processing of client updates (grouped by file path) to prevent race conditions when multiple logical clients modify the same physical configuration file.
-   **API Endpoints**:
    -   `GET /api/data`: Returns clients, servers, and registries.
    -   `POST /api/apply`: Applies configs to selected clients.
    -   `POST /api/install`: Adds a server from the registry to `mcp.json` (defaults to `npx` execution).
    -   `POST /api/server/update`: Updates server config directly (Edit feature).
    -   `POST /api/server/remove`: Removes a server configuration (Delete feature).
    -   `GET /api/doctor`: Runs system health checks.
-   **Frontend**: Single-page dashboard using Pico.css with features to:
    -   View detected clients and configured servers.
    -   Search specifically for MCP servers in registries.
    -   Install new servers with a customizable command wizard.
    -   **Edit** existing server configurations via a JSON modal.
    -   **Delete** server configurations.

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
| `cody` | Cody (Sourcegraph) | VSCode-JSON | Targets `openctx.providers` |
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

## 🚀 Future Roadmap

1.  **More Tool Support**: Keep expanding the registry as new tools emerge.
2.  **UI Enhancements**:
    -   Log viewer for the MCP servers? (Hard since they run inside the clients).
    -   Visual editor for `config.yaml` (Registries management).
3.  **Complex Configuration**: Support more advanced `openctx` provider mappings if Cody's requirements evolve.

## 📝 Memories
-   The project uses `github.com/tailscale/hujson` to parse JSONC.
-   `internal/client` is the source of truth for tool detection.
-   `internal/core` handles the business logic.
-   The UI runs on port 3000 by default (`mcpenetes ui`).
