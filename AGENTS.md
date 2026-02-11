# 🤖 mcpenetes Developer & Agent Instructions

**mcpenetes** is the universal configuration manager for the Model Context Protocol (MCP) ecosystem. It supports 30+ clients (IDEs, CLIs, Desktop Apps) and provides a unified interface (CLI & Web UI) to manage their configurations.

## 🌟 Vision & Mission
*   **Mission:** Be the "Kubernetes for MCP" - a single control plane for managing MCP servers across diverse clients.
*   **Goal:** 100% feature completeness, robust error handling, comprehensive documentation, and a seamless user experience.
*   **Core Principles:** Universality, Safety (Backup first), Simplicity (Web UI), and Extensibility (Custom Clients).

## 📂 Architecture & Tech Stack
*   **Language:** Go (Golang) 1.22+
*   **UI:** Standard library `net/http` + Embedded `static/index.html` (HTML/CSS/JS). No external frontend frameworks.
*   **CLI:** `spf13/cobra`.
*   **Configuration:** YAML (`config.yaml`) + JSON/JSONC (`mcp.json`).
*   **Key Packages:**
    *   `internal/core`: Business logic (Manager, Backup, Restore, Import).
    *   `internal/client`: Client detection (`registry.go`) and custom client management (`custom.go`).
    *   `internal/mcp`: Server connectivity testing (`pinger.go`).
    *   `internal/sync`: Cloud synchronization via GitHub Gists (`gist.go`).
    *   `internal/ui`: Web server endpoints and static assets.

## 🛠️ Workflow Protocols (Strict)
1.  **Autonomous Execution:** Plan, execute, verify, and commit autonomously. Do not pause unless blocked.
2.  **Documentation First:** Update `MANUAL.md` and `README.md` *before* or *during* feature implementation. Ensure UI tooltips are present.
3.  **Versioning:**
    *   Update `internal/version/version.go` for every feature/fix.
    *   Update `CHANGELOG.md` with a detailed entry.
    *   Commit message format: `feat: description` or `fix: description`.
4.  **Testing:**
    *   Run `go test ./...` before every commit.
    *   Verify UI changes visually (using temporary Playwright scripts if needed).
    *   Run `mcpenetes doctor` to check health.
5.  **Submodules:** Check `git submodule status` and ensure submodules are synced if present.
6.  **Code Hygiene:**
    *   **NEVER commit binary artifacts** (add to `.gitignore`).
    *   Use `hujson` for parsing configs to support comments.
    *   Use `sync.Mutex` for shared resources.

## 🗺️ Feature Checklist (Current Status)
- [x] **Web UI:** Dashboard, Search, Clients, Backups, Logs, Settings, Sync, System, Help.
- [x] **Backend:** Apply Logic, Custom Clients, Registry Management, Backup Retention.
- [x] **Connectivity:** Server Ping (Test button).
- [x] **Config:** Global Environment Variables.
- [x] **Sync:** GitHub Gist (Push/Pull).
- [x] **Docs:** Comprehensive Manual (`docs/MANUAL.md`).

## 🔄 Release Procedure
1.  **Verify:** Run tests and UI checks.
2.  **Bump Version:** Increment `internal/version/version.go`.
3.  **Changelog:** Add entry to `CHANGELOG.md`.
4.  **Dashboard:** Update submodule dashboard (System tab covers this dynamically).
5.  **Commit & Push:** `git commit -am "chore: release vX.Y.Z ..."` -> `git push`.

## 🤖 Model-Specific Instructions
*   **Claude/Anthropic:** Focus on maintaining the "Vision" and "Architecture" alignment.
*   **GPT/OpenAI:** Focus on code correctness and error handling logic.
*   **Gemini/Google:** Focus on documentation depth and extensive analysis.
