# Agent-Agnostic Model Context Protocol (MCP) Setup Guide

## Overview

`battery` natively implements the **Model Context Protocol (MCP)** specification over standard input/output (stdio) via the `battery mcp` (or `battery serve`) command.

Exposing Battery as an MCP server enables **any** compliant AI coding assistant or agent (Antigravity, Claude Code, Cursor, Windsurf, Copilot, Cline, Roo Code) to interact with your multi-repository SDD environment using structured JSON-RPC tools, real-time context resources, and prompt templates—without needing fragile shell scripting or terminal output parsing.

---

## 🚀 Quick Setup (Automatic Configuration)

Battery can automatically detect and configure your installed AI coding assistants with a single command:

```bash
# Interactive multi-select menu with auto-detected clients
battery mcp install

# Or configure specific clients non-interactively
battery mcp install --client cursor,antigravity,claude-desktop

# Or configure all supported clients
battery mcp install --all
```

`battery init` also offers to configure your AI assistant automatically during initial workspace setup!

---

## ⚙️ Client Configuration Reference (Manual)

Add `battery` to your AI client's MCP configuration using the examples below.

### 1. Google Antigravity & Gemini CLI

In your global `~/.gemini/config/mcp_config.json` (or workspace `.agents/mcp_config.json`):

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

---

### 2. Anthropic Claude Desktop & Claude Code

In macOS `~/Library/Application Support/Claude/claude_desktop_config.json` or Windows `%APPDATA%\Claude\claude_desktop_config.json`:

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

For **Claude Code** CLI (`~/.claude.json` or run via CLI):
```bash
claude mcp add battery -- battery mcp
```

---

### 3. Cursor IDE

In your workspace root `.cursor/mcp.json` or in **Cursor Settings > Features > MCP**:

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

---

### 4. Windsurf IDE (Codeium)

In `~/.codeium/windsurf/mcp_config.json`:

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

---

### 5. VS Code Extensions (Roo Code, Cline, GitHub Copilot)

In `mcpSettings.json` or extension MCP settings:

```json
{
  "mcpServers": {
    "battery": {
      "command": "battery",
      "args": ["mcp"],
      "disabled": false,
      "autoApprove": [
        "battery_status",
        "battery_list_barrels",
        "battery_track_status"
      ]
    }
  }
}
```

---

### 6. Generic / Custom Stdio MCP Client

Any custom AI agent runner can launch `battery` as a child process using:

* **Executable**: `battery` (or absolute path to compiled binary `/path/to/battery`)
* **Arguments**: `["mcp"]` or `["serve"]`
* **Transport**: `stdio` (JSON-RPC 2.0 delimited by newlines)
* **Working Directory**: The root of the Battery orchestrator workspace containing `.batteryrc` / `.cooper/`.

---

## 🛠️ Available MCP Tools

| Tool Name | Description | Required Arguments | Optional Arguments |
| :--- | :--- | :--- | :--- |
| `battery_status` | Inspects workspace topology, barrel connectivity, and active tracks. | *None* | `verbose` (boolean) |
| `battery_list_barrels` | Lists registered barrels and resolves their Cooper tech stacks (`.cooper/definition/tech-stack.md`). | *None* | *None* |
| `battery_init_track` | Scaffolds a new track under `.cooper/active/<track_id>/`. | `track_id` (string) | `barrels` (array), `name` (string), `force` (boolean) |
| `battery_dispatch_track` | Dispatches spec deltas to barrel worktrees while omitting `plan.md` to preserve local planning autonomy. | `track_id` (string) | `force` (boolean) |
| `battery_track_status` | Aggregates phase completion and task checklists across all participating barrels. | `track_id` (string) | *None* |

---

## 📚 Living Context Resources (`battery://`)

AI agents can query real-time workspace state using standard MCP `resources/read`:

* **`battery://topology`** (`application/json`): Merged canonical `.batteryrc` and local `.batteryrc.local` configuration.
* **`battery://barrels/{name}/tech-stack`** (`text/markdown`): Resolved Cooper tech stack guidelines, language idioms, and test runner configurations for a specific barrel.
* **`battery://tracks/{track_id}`** (`application/json`): Comprehensive track status report, task completion counts, and participating barrel progress.

---

## 💡 Prompt Templates

* **`plan_multi_barrel_track`**: Interactive planning prompt that guides AI assistants through barrel discovery, contract generation, spec delta authoring, and decentralized track dispatching.

---

## 🧪 Testing & Verification

You can verify that the Battery MCP server is responsive by sending a JSON-RPC ping over stdin:

```bash
echo '{"jsonrpc":"2.0","id":1,"method":"ping"}' | battery mcp
```

Expected output:
```json
{"jsonrpc":"2.0","id":1,"result":{}}
```
