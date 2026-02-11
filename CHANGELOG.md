# Changelog

All notable changes to this project will be documented in this file.

## [1.4.0] - 2025-05-27
### Added
- **Cloud Sync:** Sync configurations across devices using GitHub Gists.
- **Server Testing:** "Test" button in dashboard to verify MCP server connectivity (JSON-RPC handshake).
- **Global Environment Variables:** Define env vars once in Settings to inject into all servers.
- **Dynamic Project Structure:** System tab now visualizes the project directory tree and submodule status.
- **User Manual:** Comprehensive documentation in `docs/MANUAL.md`.
- **Developer Docs:** Consolidated instructions in `AGENTS.md`.

### Changed
- **Config:** Standardized configuration path to `~/.config/mcpetes` across all operating systems.
- **Backup:** Implemented automatic retention policy to prevent backup bloat.
- **UI:** Massive enhancements to Dashboard, Settings, and Help tabs.

## [1.3.2] (2025-05-27)
- **Settings:** Added Backup Path and Retention Policy settings.
- **Clients:** Added "Known Clients" list and manual configuration support.
- **Search:** Added dynamic install forms based on registry `EnvSchema`.

## [1.3.1] (2025-05-27)
- **UI:** Added Logs viewer, Backup management, and Registry editing.
- **Core:** Implemented Delete Backup and Clear Logs.

## [1.3.0] (2025-05-27)
- **ui:** add web-based user interface (`mcpenetes ui`) for managing configurations.

## [1.2.0] (2025-05-27)
- **core:** backup/restore logic.

## [1.1.0] (2025-05-27)
- **cli:** initial CLI commands.
