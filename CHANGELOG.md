## [1.3.2](https://github.com/tuannvm/mcpenetes/compare/v1.3.1...v1.3.2) (2025-05-27)

### Features

* **ui:** add ability to delete configuration backups.
* **ui:** add ability to clear and filter application logs.
* **ui:** add ability to edit existing registry URLs.
* **ui:** consolidate system health ("Doctor") into the System tab.
* **core:** implement backend logic for removing backups and clearing logs.
* **core:** implement registry update logic.

## [1.3.1](https://github.com/tuannvm/mcpenetes/compare/v1.3.0...v1.3.1) (2025-05-27)

### Features

* **ui:** add "System" dashboard tab with version info and project structure.
* **ui:** add "Config Key" support to Custom Client wizard for advanced configurations.
* **docs:** add comprehensive `VISION.md` and standardize developer documentation.

## [1.3.0](https://github.com/tuannvm/mcpenetes/compare/v1.2.0...v1.3.0) (2025-05-27)

### Features

* **ui:** add custom client management interface (add, delete, list).
* **ui:** add configuration backup manager with restore capabilities.
* **ui:** add real-time application log viewer.
* **ui:** add comprehensive documentation and help section.
* **ui:** add client configuration inspection (view raw config file).
* **ui:** add import configuration modal with clipboard support.
* **core:** implement safe custom client registry management (`clients.yaml`).
* **core:** fix data race in logging system.

## [1.2.0](https://github.com/tuannvm/mcpenetes/compare/v1.1.0...v1.2.0) (2025-05-27)

### Features

* **ui:** major refactor to introduce advanced web-based dashboard.
* **ui:** add toast notifications and modal interactions.
* **core:** refactor backup and restore logic.

## [1.1.0](https://github.com/tuannvm/mcpenetes/compare/v1.0.1...v1.1.0) (2025-05-27)

### Features

* **client:** add extensive client support including VS Code, Claude, Cursor, Windsurf, Zed, Trae, Goose, Mistral Vibe, and more.
* **client:** implement user-defined registry for custom tool support via `clients.yaml`.
* **ui:** add web-based user interface (`mcpenetes ui`) for managing configurations.
* **ui:** implement install wizard for customizing server commands (npx/uvx) before installation.
* **search:** refactor search workflow to automatically add default configurations to `mcp.json`.
* **translator:** add support for JSONC (VS Code settings with comments), YAML, and TOML formats.
* **translator:** add safe parsing to prevent data loss in existing configuration files.
* **core:** refactor application logic into reusable `Manager` for CLI and UI consistency.
* **integration:** add support for JetBrains IDEs via Junie agent.
* **integration:** add support for Continue extension with nested configuration format.

### Bug Fixes

* **load:** fix clipboard loading crashing on JSON with comments.
* **search:** fix search command not persisting usable configurations.

## [1.0.1](https://github.com/tuannvm/mcpenetes/compare/v1.0.0...v1.0.1) (2025-04-25)

# 1.0.0 (2025-04-18)


### Features

* **cache:** implement server caching and add refresh flag to search cmd ([636ee98](https://github.com/tuannvm/mcpenetes/commit/636ee98d7e3eff3cac99a6ee17d76d13b5b8646e))
* **ci:** add GitHub Actions for build, release, and dependency updates ([55e01e8](https://github.com/tuannvm/mcpenetes/commit/55e01e8f319290f69d4ec519026eca3b25432c9d))
* **cli:** add `load` command for clipboard config loading ([90db793](https://github.com/tuannvm/mcpenetes/commit/90db7933ed32f2d8fc85672f7d2e0b178be2e075))
* **cmd:** add interactive fuzzy search for MCP versions ([72fc02d](https://github.com/tuannvm/mcpenetes/commit/72fc02db0f5825aaf2a21c27ac5db7fb04d1a100))
* **cmd:** implement CLI commands for resource management ([7ddaf19](https://github.com/tuannvm/mcpenetes/commit/7ddaf19a64a01781b7d917b91bc7402c1e53ddb4))
* **config:** add configuration management and caching system ([d5871b7](https://github.com/tuannvm/mcpenetes/commit/d5871b741127a8481ed92dc9f0dc4ab6764025b0))
* **config:** enhance MCPServer struct and add client detection ([23c41cf](https://github.com/tuannvm/mcpenetes/commit/23c41cfdb83dd243ddb757a6eec2bfaf48bef852))
