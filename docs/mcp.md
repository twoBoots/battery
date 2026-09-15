# Model Context Protocol (MCP) Server 🤖

**Battery** natively implements the [Model Context Protocol (MCP)](https://modelcontextprotocol.io/) specification over standard input/output (stdio) via the `battery mcp` (or `battery serve`) command.

By exposing Battery as an MCP server, any compliant AI coding assistant or autonomous agent (Google Antigravity, Claude Code, Cursor, Windsurf, GitHub Copilot, Cline, Roo Code) can interact directly with your multi-repository SDD environment using structured JSON-RPC tools, real-time context resources, and interactive prompts.

---

## 🚀 Automatic Client Configuration

Configure all installed AI coding assistants with a single command:

```bash
# Interactive configuration menu
battery mcp install

# Non-interactive configuration for specific clients
battery mcp install --client cursor,antigravity,claude-desktop

# Configure all supported clients
battery mcp install --all
```

---

## ⚙️ Manual Configuration

Add `battery` to your assistant's MCP configuration using the JSON snippets below:

### 1. Google Antigravity & Gemini CLI
In `~/.gemini/config/mcp_config.json`:
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

### 2. Anthropic Claude Desktop & Claude Code
In macOS `~/Library/Application Support/Claude/claude_desktop_config.json` or Linux `~/.config/Claude/claude_desktop_config.json`:
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

### 3. Cursor & Windsurf
In your workspace `.cursor/mcp.json` or `.windsurf/mcp.json`:
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

## 🛠️ Available MCP Tools

| Tool Name | Description | Arguments |
| :--- | :--- | :--- |
| `battery_status` | Inspects workspace topology, barrel connectivity, and active tracks. | `verbose` (boolean, optional) |
| `battery_list_barrels` | Lists registered barrels and resolves their [Cooper](https://github.com/twoBoots/cooper) tech stacks (`.cooper/definition/tech-stack.md`). | *None* |
| `battery_init_barrel_tech_stack` | Scaffolds or updates `.cooper/definition/tech-stack.md` and code styleguides for a barrel or monorepo package. | `barrel` (string, required), `language`, `framework`, `test_runner`, `linter`, `coverage_threshold`, `force` |
| `battery_init_track` | Scaffolds a new track under `.cooper/active/<track_id>/`. | `track_id` (string, required), `barrels` (array), `name` (string), `force` (boolean) |
| `battery_dispatch_track` | Dispatches spec deltas to barrel worktrees managed by [Troop](https://github.com/twoBoots/troop) while omitting `plan.md` to preserve local planning autonomy. | `track_id` (string, required), `force` (boolean) |
| `battery_track_status` | Aggregates phase completion and task checklists across all participating barrels. | `track_id` (string, required) |

---

## 📚 Living Context Resources (`battery://`)

AI assistants can query real-time workspace state using standard MCP `resources/read`:

* **`battery://topology`** (`application/json`): Merged canonical `.batteryrc` and local `.batteryrc.local` configuration.
* **`battery://barrels/{name}/tech-stack`** (`text/markdown`): Resolved [Cooper](https://github.com/twoBoots/cooper) tech stack guidelines, language idioms, and test runner configurations for a specific barrel.
* **`battery://tracks/{track_id}`** (`application/json`): Comprehensive track status report, task completion counts, and participating barrel progress.

---

## 💡 Prompt Templates

* **`plan_multi_barrel_track`**: Interactive planning prompt that guides AI assistants through barrel discovery, contract generation, spec delta authoring, and decentralized track dispatching.

---

## 🔗 Related Resources

* [Getting Started Guide](guide/getting-started.md)
* [Multi-Barrel Workflow Guide](guide/workflow.md)
* [Architecture & Topology](architecture.md)
* [Installation Guide](installation.md)
* [Cooper SDD Framework](https://github.com/twoBoots/cooper)
* [Troop Worktree Isolation](https://github.com/twoBoots/troop)

