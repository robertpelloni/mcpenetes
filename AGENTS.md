# 🤖 AI Agent Instructions

*(See `INSTRUCTIONS.md` for universal project context)*

## Workflow
1.  **Plan:** Use `set_plan` to outline steps.
2.  **Verify:** Use `ls`, `read_file` to check state.
3.  **Implement:** Edit code.
4.  **Test:** Run `go test ./...` and `go run main.go doctor`.
5.  **Submit:** Commit with descriptive messages matching the changelog.

## Key Files
*   `internal/client/registry.go`: Client definitions.
*   `internal/translator/translator.go`: Logic for applying configs.
*   `internal/ui/server.go`: Web UI handler.
