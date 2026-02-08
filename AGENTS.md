# 🤖 AI Agent & Developer Instructions

## 🗺️ Universal Context
**mcpenetes** is the universal configuration manager for the Model Context Protocol (MCP) ecosystem. It supports 30+ clients (IDEs, CLIs, Desktop Apps) and provides a unified interface (CLI & Web UI) to manage their configurations.

*   **Primary Goal:** Ensure complete, robust, and well-documented support for all MCP clients.
*   **Vision:** See `VISION.md` for the ultimate goal and design philosophy.
*   **Architecture:** See `HANDOFF.md` for current implementation details.

## 🛠️ Workflow Protocols
1.  **Autonomous Execution:** Plan, execute, verify, and commit autonomously. Do not stop for confirmation unless blocked.
2.  **Documentation First:** Every feature must have UI documentation (tooltips/help tab) and markdown documentation (`README.md`, `Help`).
3.  **Versioning:**
    *   Increment `internal/version/VERSION` for every significant change.
    *   Update `CHANGELOG.md` to reflect the new version and changes.
    *   Git commit message should reference the change.
4.  **Submodules:** Check for submodules (`git submodule status`) and update them if present. Maintain a dashboard/list of submodules in the UI.
5.  **Testing:**
    *   Backend: `go test ./...`
    *   Frontend: Write temporary Playwright scripts (`verify_ui.py`) to visually verify changes.
    *   System Health: Run `go run main.go doctor` to ensure environment integrity.

## 📂 Project Structure
*   `cmd/`: CLI command definitions (cobra).
*   `internal/client/`: Client registry (`registry.go`) and custom client logic (`custom.go`).
*   `internal/config/`: Configuration structs and persistence logic.
*   `internal/core/`: Business logic (Backup, Restore, Import, Apply).
*   `internal/doctor/`: System health checks.
*   `internal/ui/`: Web UI server and static assets (`static/index.html`).
*   `internal/version/`: Version string source of truth.

## 💡 Code Guidelines
*   **Safety:** Use `hujson` for JSON parsing (supports comments). Always backup before overwriting user files.
*   **Concurrency:** Use `sync.Mutex` for shared resources (e.g., logging buffer).
*   **Cross-Platform:** Use `filepath.Join` and OS-specific base paths (`BaseHome`, `BaseAppData`).
*   **Error Handling:** Return meaningful errors and log them appropriately.

## 🔄 Release Procedure
1.  Run tests: `go test ./...`
2.  Verify UI (if changed).
3.  Update `internal/version/VERSION`.
4.  Update `CHANGELOG.md`.
5.  Commit: `feat: description` or `fix: description`.
6.  Push.
