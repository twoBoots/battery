# Installation & Setup Guide 📦

This guide covers installing the `battery` CLI, configuring [Troop](https://twoboots.github.io/troop) Git aliases, scaffolding [Cooper](https://twoboots.github.io/cooper) SDD infrastructure, and troubleshooting.

## ⚡ Quickstart (One-Line Installer)

Install Battery into your environment with the official curl installer:

```bash
curl -fsSL https://raw.githubusercontent.com/twoBoots/battery/main/install.sh | bash
```

The installer will:
1. Detect your operating system (macOS / Linux) and architecture (arm64 / x86_64).
2. Download and verify the latest precompiled `battery` binary.
3. Install the binary into `~/.local/bin` (or `/usr/local/bin`).
4. Configure [Troop](https://twoboots.github.io/troop) Git aliases in your global Git configuration (`~/.gitconfig`).
5. Ensure your `PATH` includes the binary destination directory.

## 📋 System Prerequisites

* **Git**: Version `2.30.0` or higher (required for `git worktree` features).
* **cURL**: For remote script and binary retrieval.
* **Bash / Zsh**: Compatible POSIX shell.
* **Go** (optional): Go `1.22+` if compiling from source.

## 🔧 Git Aliases Configuration ([Troop](https://twoboots.github.io/troop))

Battery relies on [Troop](https://twoboots.github.io/troop) for Git worktree isolation. The installer automatically registers the following aliases in your Git configuration:

```ini
[alias]
    troop = "!troop"
    agent-start = "!f() { git worktree add .worktrees/$1 -b $1 main; }; f"
    agent-stop = "!f() { git worktree remove .worktrees/$1; git branch -d $1; }; f"
```

Verify your Troop aliases:
```bash
git troop
```

## 🔨 Building from Source

You can build `battery` directly from source using Go:

```bash
# Clone the repository
git clone https://github.com/twoBoots/battery.git
cd battery

# Build the binary
go build -o battery ./cmd/battery

# Install into local path
mv battery ~/.local/bin/
```

## 🔄 Self-Updating Battery

Keep Battery up to date using the built-in update command:

```bash
# Check if a new version is available
battery update --check

# Upgrade to the latest release
battery update

# Upgrade or downgrade to a specific release
battery update --target-version v1.3.0
```

## 🛠️ Troubleshooting

### Command Not Found: `battery`
If running `battery` fails with `command not found`, ensure `~/.local/bin` is in your shell's `PATH`:

```bash
# In ~/.zshrc or ~/.bashrc
export PATH="$HOME/.local/bin:$PATH"
```

Then reload your shell:
```bash
source ~/.zshrc # or source ~/.bashrc
```

### Worktree Creation Issues
If `git agent-start` fails with `fatal: already exists`, check for stale worktrees with:
```bash
git worktree list
git worktree prune
```

## 🔗 Related Resources

* [Getting Started Guide](guide/getting-started.md)
* [Multi-Barrel Workflow Guide](guide/workflow.md)
* [Architecture & Topology Model](architecture.md)
* [Model Context Protocol (MCP) Server](mcp.md)
* [Cooper SDD Framework](https://twoboots.github.io/cooper)
* [Troop Worktree Isolation](https://twoboots.github.io/troop)
