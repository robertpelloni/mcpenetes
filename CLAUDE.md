# 🧠 CLAUDE.md

## 🛠️ Build & Test Commands
- **Build**: `go build -o bin/mcpenetes main.go`
- **Run**: `go run main.go`
- **Test**: `go test ./...`
- **Run UI**: `go run main.go ui` (starts on localhost:3000)
- **Lint**: `golangci-lint run` (if available)

## 📂 Project Structure
- `cmd/`: Cobra CLI commands (`apply`, `search`, `ui`, `load`, `restore`).
- `internal/client/`: Registry of supported MCP clients (IDEs, CLIs). **Add new tools here.**
- `internal/core/`: Application logic (`ApplyToClient`) shared by CLI and UI.
- `internal/translator/`: Logic for reading/writing config files (JSON/YAML/TOML). **Handles safety.**
- `internal/ui/`: Web server (`server.go`) and static assets (`static/`).
- `internal/util/`: Helper functions (path expansion, detection wrappers).

## 🧩 Code Style & Conventions
- **JSON Handling**: ALWAYS use `github.com/tailscale/hujson` for parsing config files to support comments (JSONC).
- **Versioning**: Update `internal/version/version.go` and `CHANGELOG.md` on every feature change.
- **Safety**: Never overwrite a user config file if the read parse fails. Return an error instead.
- **Paths**: Prioritize Windows paths (`APPDATA`, `USERPROFILE`) in the registry.
- **Error Handling**: Wrap errors with context (`fmt.Errorf("failed to ...: %w", err)`).

## 🚀 Key Features
- **Registry**: Data-driven client detection. Supports user extensions via `clients.yaml`.
- **Web UI**: Dashboard for managing MCP configs visually.
- **Search**: Interactive registry search + install wizard.
- **Backups**: Automatically backs up config files before modification.
