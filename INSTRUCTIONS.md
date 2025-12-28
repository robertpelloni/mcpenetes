# 🤖 Universal LLM Instructions

**Project Name:** mcpenetes
**Current Version:** 1.1.0

## Core Context
`mcpenetes` is a CLI tool and Web UI for managing Model Context Protocol (MCP) server configurations across multiple AI clients (IDEs, CLIs, Desktop Apps).

## Key Features
1.  **Registry:** Supports ~33 clients via `internal/client/registry.go`.
2.  **UI:** Web dashboard at `http://localhost:3000` via `mcpenetes ui`.
3.  **Safety:** Uses `hujson` for JSONC parsing (safe for VS Code).
4.  **Extensibility:** User-defined `clients.yaml` support.

## Development Guidelines
*   **Versioning:** Always update `VERSION` file and `CHANGELOG.md` on release.
*   **Testing:** Run `go test ./...` before committing.
*   **Frontend:** Simple HTML/JS in `internal/ui/static`. No complex build steps.
*   **Documentation:** Keep `README.md` and `HANDOFF.md` updated.

## Submodules
*   None.
