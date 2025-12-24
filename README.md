# 🧙‍♂️ mcpenetes

![mcpenetes](https://img.shields.io/badge/mcpenetes-MCP%20Configuration%20Manager-blue)
![License](https://img.shields.io/badge/license-MIT-green)

> *"One CLI to rule them all, one CLI to find them, one CLI to bring them all, and in the configurations bind them."*

![mcpenetes in action](https://github.com/user-attachments/assets/bf84188d-e043-4436-bebb-ef206ff10d8f)

## 🌟 What is mcpenetes?

**mcpenetes** is a magical CLI tool that helps you manage multiple Model Context Protocol (MCP) server configurations with ease! If you're tired of manually editing config files for different MCP-compatible clients whenever you want to switch servers, mcpenetes is here to save your day.

Think of mcpenetes as your friendly neighborhood wizard who can:

- 🔍 Search for available MCP servers from configured registries
- 🔄 Switch between different MCP server configurations
- 🧠 Apply configurations across all your MCP clients automatically
- 🖥️ **New!** Manage everything via a beautiful Web UI
- 💾 Backup your configurations before making any changes
- 🛡️ Restore configurations if something goes wrong
- 🏥 Diagnose system health with the `doctor` command

## 🚀 Installation

### From Source

```bash
git clone https://github.com/tuannvm/mcpenetes.git
cd mcpenetes
make build
# The binary will be available at ./bin/mcpenetes
```

### Using Go

```bash
go install github.com/tuannvm/mcpenetes@latest
```

## 🏄‍♂️ Quick Start

### Option 1: The Web UI (Recommended)

Start the dashboard to view your clients, search for servers, and apply configurations visually:

```bash
mcpenetes ui
```

This will open `http://localhost:3000` in your default browser.

### Option 2: The CLI Way

1. **Search for available MCP servers**:

```bash
mcpenetes search
```

2. **Apply selected configuration** to all your clients:

```bash
mcpenetes apply
```

That's it! Your MCP configurations are now synced across all clients. Magic! ✨

## 📚 Usage Guide

### 🛠️ Available Commands

```
ui             Start the Web UI dashboard
search         Interactive fuzzy search for MCP versions and apply them
apply          Applies MCP configuration to all clients
load           Load MCP server configuration from clipboard
restore        Restores client configurations from the latest backups
doctor         Run system health checks and client detection verification
```

### 📋 Searching for MCP Servers

The `search` command lets you interactively find and select MCP servers from configured registries. It will present you with a list of available servers that you can select from.

```bash
mcpenetes search
```

By default, search results are cached to improve performance. Use the `--refresh` flag to force a refresh:

```bash
mcpenetes search --refresh
```

### 📥 Loading Configuration from Clipboard

If you've copied an MCP configuration to your clipboard, you can load it directly:

```bash
mcpenetes load
```

### 🗑️ Removing Resources

To remove a registry:

```bash
mcpenetes remove registry my-registry
```

### ⏪ Restoring Configurations

If something goes wrong, you can restore your clients' configurations from backups:

```bash
mcpenetes restore
```

## 🧩 Supported Clients

mcpenetes automatically detects and configures over 30 MCP-compatible clients, including:

**IDEs & Editors:**
*   VS Code, VS Code Insiders
*   Cursor, Windsurf, Zed, Trae, PearAI, Void
*   **JetBrains IDEs** (IntelliJ, PyCharm, etc.) via Junie
*   **Melty** (VS Code Fork)
*   **CodeBuddy** (VS Code Fork)
*   **Kiro**

**Extensions:**
*   **Cline**
*   **Roo Code**
*   **Continue**
*   **Cody (Sourcegraph)** (Configures `openctx.providers` in VS Code settings)

**Desktop Apps:**
*   Claude Desktop
*   LM Studio
*   AnythingLLM
*   Tabby
*   LibreChat
*   Jan
*   BoltAI

**CLIs & Terminals:**
*   **Amazon Q (CodeWhisperer)**
*   **Claude Code CLI**
*   **LLM CLI** (Simon Willison)
*   Goose CLI
*   Mistral Vibe
*   Code CLI (Codex)
*   Grok CLI
*   Open Interpreter
*   Factory CLI
*   Aider
*   Warp Terminal

### Adding Custom Clients
You can support additional tools by creating a `clients.yaml` file in your config directory (e.g., `~/.config/mcpetes/clients.yaml`).

## 📁 Configuration Files

mcpenetes uses the following configuration files:

- `~/.config/mcpetes/config.yaml`: Stores global configuration, including registered registries and selected MCP servers
- `~/.config/mcpetes/mcp.json`: Stores the MCP server configurations
- `~/.config/mcpetes/cache/`: Caches registry responses for faster access

## 🤝 Contributing

Contributions are welcome! Feel free to:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📜 License

Licensed under the MIT License. See the LICENSE file for details.

## 🌐 Related Projects

- [mcp-trino](https://github.com/tuannvm/mcp-trino): Trino MCP server implementation

---

Made with ❤️ by humans (and occasionally with the help of AI)
