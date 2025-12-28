# 🤖 Claude Instructions

*(See `INSTRUCTIONS.md` for universal project context)*

## Specific Instructions for Claude
*   When editing `internal/client/registry.go`, ensure paths use `filepath.Join` and appropriate base constants (`BaseHome`, `BaseAppData`, `BaseUserProfile`).
*   Prioritize Windows compatibility (`%APPDATA%`) in registry definitions.
*   Use `hujson` when parsing user configuration files to avoid errors with comments.
