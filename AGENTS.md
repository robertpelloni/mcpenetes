# 🤖 AGENTS.md

> **Notice to AI Agents:** This file allows humans to give you instructions or tips for working with this code. The scope of this file is the entire directory tree rooted at the folder that contains it. You MUST obey instructions in this file.

## 🧠 Core Directives & Persona

1.  **Deep Planning Mode**: Before making changes, start a deep planning mode. Interact with the user to absolutely confirm expectations. Ask questions until you have zero doubt.
2.  **Autonomous Execution**: Once the plan is approved, proceed autonomously. You may complete a feature, commit, push, and continue to the next without stopping for confirmation unless blocked.
3.  **Frequent Commits**: Commit and push (`submit`) in between each major step or logical unit of work. Do not hoard changes.
4.  **Windows Priority**: When defining paths or OS-specific logic, prioritize **Windows** support, but ensure macOS and Linux are also handled where possible.

## 🛠️ Development Workflow

### 1. Versioning Strategy
**Every build/release must have a new version number.**
*   **Source of Truth**: `internal/version/version.go`
*   **Changelog**: `CHANGELOG.md`
*   **Procedure**:
    1.  Determine the new version (Semantic Versioning: Patch for bugfixes, Minor for features).
    2.  Update `const Version` in `internal/version/version.go`.
    3.  Add a new entry at the top of `CHANGELOG.md` following the existing format.

### 2. Testing
*   **Unit Tests**: Run `go test ./...` to verify changes.
*   **Safety Checks**: Verify that `hujson` is used for any JSON parsing that might involve comments (especially VS Code settings).
*   **Frontend**: If modifying `internal/ui`, ensure the embedded assets are correctly referenced.

### 3. Architecture Overview
The project is organized into modular packages to separate concerns:

*   **`cmd/`**: CLI entry points (`root`, `apply`, `search`, `ui`).
*   **`internal/client/`**: **The Brain**. Contains the `Registry` of supported tools (`registry.go`).
    *   **Rule**: Do not hardcode detection in `util`. Add new tools to `Registry`.
    *   **Formats**: Supports `FormatSimpleJSON`, `FormatVSCode` (nested), `FormatYAML`, `FormatTOML`, `FormatContinue`.
*   **`internal/core/`**: **The Logic**. `Manager` struct handles the `Backup -> Translate -> Apply` flow. Used by both CLI and UI.
*   **`internal/translator/`**: **The Worker**. Handles file parsing and writing.
    *   **Rule**: MUST use `hujson` for JSON. MUST abort if parsing fails to prevent data loss.
*   **`internal/ui/`**: **The Face**. Embedded Web UI server and static assets (`static/index.html`).

## 🔍 Key Implementation Details

*   **User-Defined Registry**: The tool loads `~/.config/mcpetes/clients.yaml` to allow users to add custom tools.
*   **Detection Heuristic**: Detects tools by config file existence OR parent directory existence (to support fresh installs).
*   **Search**: The `search` command adds a default configuration (usually `npx`) to `mcp.json`. The UI allows customizing this before install.

## 📝 Changelog Maintenance
When updating `CHANGELOG.md`:
*   Group changes by type: `Features`, `Bug Fixes`, `Refactoring`.
*   Include the scope (e.g., `**client:**`, `**ui:**`).
*   Keep it human-readable.

---
*You are resourceful. Use the tools at your disposal. Build magnificent software.*
