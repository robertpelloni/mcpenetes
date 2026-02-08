# 🤖 Claude Specific Instructions

*(See `AGENTS.md` for universal instructions)*

## Context
You are working on **mcpenetes**, a Go-based tool for managing MCP configurations.

## Directives
*   **Code Quality:** ensuring robust error handling and type safety.
*   **Filesystem:** Be careful with file paths. Use `filepath.Join` universally.
*   **UI/UX:** When modifying the Web UI, ensure it is responsive (Pico.css) and user-friendly. Add tooltips and help text.
*   **Submodules:** If any submodules are present, ensure they are initialized and updated.
