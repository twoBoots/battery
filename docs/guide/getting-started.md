# Getting Started with Battery 🔋

**Battery** is an open, agent-agnostic and repository-agnostic multi-repository Specification-Driven Development (SDD) orchestration protocol. Built on top of **[Cooper](https://github.com/twoBoots/cooper) (the barrel maker & Hybrid SDD Framework)** and **[Troop](https://github.com/twoBoots/troop) (worktree isolation tool)**, Battery coordinates multi-repository tracks and living capability specs across a collection of barrels (individual repositories or packages) for human developers and autonomous AI agents alike.

---

## 🛢️ Core Metaphor & Concepts

* **[Cooper](https://github.com/twoBoots/cooper)**: The barrel maker managing SDD specifications (`.cooper/`) and worktree lifecycle for human developers and autonomous AI agents.
* **[Troop](https://github.com/twoBoots/troop)**: Worktree isolation tool providing shared Git aliases for human developers and AI code monkeys working in isolated worktrees (`.worktrees/`).
* **Barrel**: An individual repository, package, or microservice within the system landscape.
* **Battery**: A collection of barrels orchestrated together for cross-repository feature epics.

---

## ⚡ 5-Minute Quickstart

### Step 1: Install Battery

Install the `battery` CLI into your system with the official installer script:

```bash
curl -fsSL https://raw.githubusercontent.com/twoBoots/battery/main/install.sh | bash
```

The installer verifies prerequisites, sets up [Troop](https://github.com/twoBoots/troop) Git aliases, and installs the precompiled `battery` binary into `~/.local/bin` (or `/usr/local/bin`).

### Step 2: Initialize Workspace & Discover Barrels

Navigate to your multi-repo workspace root or monorepo root and initialize Battery:

```bash
# Interactive setup with automatic barrel discovery
battery init

# Or non-interactive setup for multi-repo sibling directories
battery init --structure multi-repo -y
```

This creates the canonical [`.batteryrc`](../architecture.md) topology file defining your barrels.

### Step 3: Inspect Workspace Status

Verify that all barrels and their resolved [Cooper](https://github.com/twoBoots/cooper) tech stacks are connected:

```bash
# Check overall workspace health and active tracks
battery status

# List all registered barrels and resolved language/test stacks
battery barrel list
```

### Step 4: Configure AI Coding Assistants

Connect Battery directly to your favorite AI coding assistant (Google Antigravity, Claude Code, Cursor, Windsurf, Copilot, Cline, Roo Code) via Model Context Protocol:

```bash
battery mcp install
```

---

## 📚 Next Steps

* Master the multi-barrel track lifecycle in the [Workflow Guide](workflow.md).
* Understand the multi-repo and monorepo topology models in [Architecture & Topology](../architecture.md).
* Explore the full [MCP Server Reference](../mcp.md) for AI assistant integrations.
* Review detailed build and installation options in the [Installation Guide](../installation.md).
