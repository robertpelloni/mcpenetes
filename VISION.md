# 🌟 Project Vision: mcpenetes

## Mission Statement
**mcpenetes** aims to be the universal configuration manager for the Model Context Protocol (MCP) ecosystem. As AI coding assistants proliferate (VS Code, Cursor, Windsurf, Zed, etc.), managing MCP server configurations across multiple tools becomes a fragmented and tedious task. **mcpenetes** solves this by providing a single source of truth for MCP configurations and seamlessly applying them to all supported clients.

## Core Principles
1.  **Universality:** Support every major IDE, CLI, and desktop application that implements MCP. If it's not built-in, users can add it via `clients.yaml`.
2.  **Safety:** Never lose user data. Use robust JSON/YAML parsing (supporting comments) and always backup before applying changes.
3.  **Simplicity:** Provide a beautiful Web UI for visual management and a powerful CLI for automation.
4.  **Extensibility:** Allow users to define custom clients and registries without waiting for app updates.
5.  **Robustness:** Handle edge cases (missing files, fresh installs, race conditions) gracefully.

## Architectural Design
*   **Registry-Driven:** Client detection logic is data-driven (`internal/client/registry.go`), making it easy to add new tools.
*   **Core Logic Isolation:** Business logic (`internal/core`) is decoupled from the UI/CLI layers, enabling consistent behavior across interfaces.
*   **Web UI:** A lightweight, dependency-free (embedded) web interface for managing the entire ecosystem.
*   **Versioning:** Strict semantic versioning with a single source of truth (`internal/version/VERSION`).

## Future Roadmap
*   **Cloud Sync:** Synchronize configurations across machines via GitHub Gists or a cloud backend.
*   **Plugin System:** Allow third-party plugins to extend functionality (e.g., custom translators).
*   **MCP Server Marketplace:** Integration with a centralized registry for one-click installation of community servers.
*   **Docker Integration:** Native support for managing Docker-based MCP servers.

## Development Guidelines
*   **Autonomous Workflow:** Developers (AI and human) should work autonomously, committing frequently and ensuring high code quality.
*   **Documentation First:** All features must be documented in the UI (tooltips/help) and in markdown files (`README.md`, `Help` tab).
*   **Test-Driven:** Verify changes with tests and manual verification (Playwright for UI).
