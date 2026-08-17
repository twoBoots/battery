# Troop Architecture (Monkeys & Trees)

## Overview

**Troop** is a worktree-based workflow for concurrent development where human developers and autonomous AI agents collaborate as a **troop of code monkeys working in trees**.

In the Troop paradigm:
* **The Troop**: The collective group of Code Monkeys (human and AI developers) working in parallel.
* **Code Monkeys**: All developers—human or AI—are code monkeys.
* **Trees (Worktrees)**: Isolated Git worktrees (`.worktrees/`) where monkeys climb up to work on tasks without swinging into each other.
* **The Main Trunk**: The root repository workspace (`main` branch) shared by the troop.

---

## Workspace Structure

Trees (worktrees) for code monkeys are isolated in a hidden `.worktrees/` canopy at the project root.

```text
my-project/
├── .git/                 # Core Git metadata
├── .worktrees/           # Hidden tree canopy for troop monkeys
│   ├── monkey-1-auth/    # Isolated tree for monkey 1
│   └── monkey-2-billing/ # Isolated tree for monkey 2
├── src/                  # Main application source code
├── .gitaliases           # Shared Git aliases for Troop
├── .gitignore            # Ignores .worktrees/
├── README.md             # Project quick start
├── AGENTS.md             # Guidelines for AI agents
└── TROOP.md              # Architecture specification
```

---

## The Troop Lifecycle

### 1. Initializing Aliases
To initialize Troop in your repository:
```bash
curl -fsSL https://raw.githubusercontent.com/twoboots/troop/main/install.sh | bash
```

### 2. Spawning a Monkey Worktree (`agent-start`)
When a new task is assigned to a code monkey:
```bash
git agent-start <task-name>
```
What this does under the hood:
1. Fetches `origin/main` (if a remote exists).
2. Spawns a new tree (worktree) at `.worktrees/<task-name>`.
3. Creates branch `<task-name>` off `origin/main` (falling back to local `main` if offline).

### 3. Listing Monkeys in Trees (`troop`)
To see all monkeys actively working in trees across the troop:
```bash
git troop
```

### 4. Execution & Isolation
* A code monkey navigates up to its assigned tree (`.worktrees/<task-name>`).
* Other troop members (human or AI code monkeys) continue working in the main trunk or their own trees without interference.

### 5. Teardown (`agent-stop`)
Once a monkey's work is committed, pushed, and merged:
```bash
git agent-stop <task-name>
```
What this does:
1. Removes the `.worktrees/<task-name>` directory.
2. Deletes the local `<task-name>` branch.
