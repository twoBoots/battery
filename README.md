# `battery`: Multi-Repository SDD Orchestration Protocol (A Collection of Barrels)

`battery` is an open, **agent-agnostic** and **repository-agnostic** multi-repository Specification-Driven Development (SDD) orchestration protocol. Built on top of **[Cooper](https://github.com/twoBoots/cooper) (the barrel maker & Hybrid SDD Framework)** and **[Troop](https://github.com/twoBoots/troop)**, `battery` coordinates multi-repository tracks and living capability specs across a **collection of barrels** (individual repositories or packages) for human developers and autonomous AI agents alike.

See [.cooper/BATTERY.md](.cooper/BATTERY.md) for the complete pattern specification and architectural guidelines.

---

## 🛢️ Core Metaphor & Terminology

- **[Cooper](https://github.com/twoBoots/cooper)**: The barrel maker managing SDD specifications (`.cooper/`) and worktree lifecycle for human developers and autonomous AI agents.
- **[Troop](https://github.com/twoBoots/troop)**: Worktree isolation tool providing shared Git aliases for human developers and AI code monkeys working in isolated worktrees (`.worktrees/`).
- **Barrel**: An individual repository, package, or service within the system landscape.
- **Battery**: A collection of barrels orchestrated together for cross-repository feature epics.

---

## 📦 Installation

To install `battery` into any repository or orchestration workspace:

```bash
# Remote install via curl
curl -fsSL https://raw.githubusercontent.com/twoBoots/battery/main/install.sh | bash

# Or run locally from clone
./install.sh [target_directory]
```

The installer will:
1. Initialize Git and set up Troop worktree aliases (`.gitaliases`, `.worktrees/`).
2. Scaffold the Cooper Hybrid SDD workspace (`.cooper/`).
3. Install the `battery` CLI binary (via prebuilt GitHub release or native Deno compile).
4. Prompt you to select your workspace topology and discover target barrels.

---

## ⚙️ Configuration (`.batteryrc` & `.batteryrc.local`)

`battery` uses a layered configuration model:
- **[`.batteryrc`](file:///.batteryrc)**: Committed to Git, storing canonical team project topology and barrels.
- **`.batteryrc.local`**: Git-ignored, allowing individual developers to override barrel paths to match local directory layouts.

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "version": "1.0.0",
  "structure": "multi-repo",
  "barrels": [
    {
      "name": "auth-service",
      "path": "../auth-service"
    },
    {
      "name": "web-dashboard",
      "path": "../web-dashboard"
    }
  ]
}
```

### Supported Workspace Topologies
- **Multi-Repo (`"structure": "multi-repo"`)**: Barrels reside in sibling folders (`../auth-service`, `../web-client`).
- **Monorepo (`"structure": "monorepo"`)**: Barrels reside as packages in subdirectories (`./packages/ui`, `./apps/web`).
- **Custom (`"structure": "custom"`)**: Hybrid layouts, distributed paths, or Git submodules.
- **Hierarchical Sub-Batteries (`"type": "battery"`)**: A barrel that is itself a composite battery orchestrating nested barrels.

---

## 🛠️ Decoupled Barrel Tech Stacks

`battery` intentionally does **not** duplicate tech stack settings in `.batteryrc`. Instead, `battery` dynamically resolves and respects each barrel's individual Cooper definition:
- `<barrel_path>/.cooper/definition/tech-stack.md`

Each barrel is free to use its own programming language, testing tools, and styleguides.

---

## 🚀 CLI Commands

```bash
# Workspace status and barrel connectivity
battery status

# List all registered barrels and resolved Cooper tech stacks
battery barrel list

# Add a barrel to .batteryrc (or --local)
battery barrel add ../payment-service --name payment

# Remove a barrel
battery barrel remove payment

# Interactive initialization or auto-discovery
battery init [--structure <multi-repo|monorepo|custom>] [--non-interactive] [-y]

# Multi-barrel track initialization & dispatch
battery track init <track_id> [--barrels folder-a,folder-b] [--name <title>]
battery track dispatch <track_id> [--force]
battery track status [<track_id>]
battery track list

# Self-update CLI binary
battery update [--check] [--force] [--target-version <v1.3.0>]
battery self-update

# Model Context Protocol (MCP) server for AI coding assistants
battery mcp [--transport stdio]

# Configure Battery MCP server in AI coding assistants (Cursor, Antigravity, Claude, Windsurf, VS Code)
battery mcp install [--client cursor,antigravity] [--all]
```

---

## 🤖 Model Context Protocol (MCP) Support

`battery` natively implements the **Model Context Protocol (MCP)** over stdio (`battery mcp` / `battery serve`), exposing multi-repository SDD orchestration directly to AI coding assistants (e.g. Antigravity, Claude Code, Cursor, Windsurf, Copilot, Cline, Roo Code).

### 🚀 Auto-Configuration
Configure your AI coding assistants with one command:
```bash
battery mcp install
```

See **[Agent-Agnostic MCP Setup Guide](docs/mcp-setup-guide.md)** for complete configuration instructions across all editors and agents.


### Quick Setup (`mcp_config.json`, `.cursor/mcp.json`, `claude_desktop_config.json`)

```json
{
  "mcpServers": {
    "battery": {
      "command": "battery",
      "args": ["mcp"]
    }
  }
}
```

### Available MCP Primitives
* **Tools**: `battery_status`, `battery_list_barrels`, `battery_init_track`, `battery_dispatch_track`, `battery_track_status`
* **Resources**: `battery://topology`, `battery://barrels/{name}/tech-stack`, `battery://tracks/{track_id}`
* **Prompts**: `plan_multi_barrel_track`

---

## 🔗 Quick Links
- [Agent-Agnostic MCP Setup Guide](docs/mcp-setup-guide.md)
- [Multi-Barrel Track Dispatch & Decoupled Planning Guide](docs/multi-barrel-track-dispatch.md)
- [Battery Architecture](.cooper/BATTERY.md)
- [Cooper Architecture](.cooper/COOPER.md)
- [Cooper Workflow](.cooper/definition/workflow.md)
- [Troop Architecture](.cooper/TROOP.md)
- [Agent Guidelines](AGENTS.md)
- [Go Code Styleguide](.cooper/code_styleguides/go.md)
- [Cooper Repository](https://github.com/twoBoots/cooper)
- [Troop Repository](https://github.com/twoBoots/troop)
