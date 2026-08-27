# Track Proposal: Lightweight Documentation & Context Metadata for Non-Cooper Barrels

## 1. Context & Motivation

In heterogeneous multi-barrel workspaces, some registered barrels do not use the Cooper SDD framework (e.g. third-party forks, firmware, micro-utilities, legacy codebases, or standalone Docker environments) and will likely never maintain a `.cooper/` directory.

Currently:
1. When a barrel lacks `<barrel>/.cooper/definition/tech-stack.md` (or `conductor/tech-stack.md` / `tech-stack.md`), commands like `battery status` and `battery barrel list` report `No Cooper tech-stack.md defined`.
2. Battery orchestrators lack a standardized, zero-intrusion mechanism to store architectural context, runtime commands, and interface contracts for non-Cooper barrels without modifying the target repository.

## 2. Proposed Solution

This track introduces a lightweight, zero-intrusion documentation and context metadata architecture for non-Cooper barrels:

1. **Orchestrator-Level Barrel Profile Directory (`docs/barrels/`)**:
   - Standardizes `docs/barrels/<barrel-name>.md` (with fallback to `.cooper/barrels/<barrel-name>.md` or custom `.batteryrc` `docs` paths) to capture domain responsibilities, runtime environments, build/test commands, data contracts, and AI agent notes directly in the Battery root repository.
   - Preserves 100% zero intrusion on the target repository.

2. **`.batteryrc` / `.batteryrc.local` Schema Enhancement**:
   - Enhances `BarrelConfig` to support typed metadata attributes (`role`, `tech`, `docs`, `jira`) and dynamic preservation of custom JSON fields during serialization and updates.
   - Enhances `battery barrel add` with flags `--role`, `--tech`, `--docs`, and `--jira`.

3. **Hybrid Fallback Resolution in CLI & Status**:
   - Fallback hierarchy for barrel technology and context summaries:
     1. Barrel Cooper living spec (`.cooper/definition/tech-stack.md`)
     2. `.batteryrc` / `.batteryrc.local` inline `tech` / `role` fields
     3. Auto-summarized orchestrator barrel profile (`docs/barrels/<name>.md`)
     4. Default fallback: `No Cooper tech-stack or profile defined`

4. **Profile Scaffolding CLI Subcommand**:
   - Adds `battery barrel doc init <name>` (aliased to `battery barrel profile init <name>`) to scaffold a standardized Markdown profile in `docs/barrels/<name>.md`.

5. **MCP Server Resource & Tool Enhancements**:
   - Exposes new MCP resource `battery://barrels/{name}/docs` returning the full markdown profile.
   - Falls back to barrel profile content in `battery://barrels/{name}/tech-stack` when Cooper `tech-stack.md` is absent.
   - Enriches `battery_list_barrels` tool and `battery://topology` resource with barrel metadata attributes.

## 3. Scope Boundaries

- **In Scope:**
  - `internal/config`: `BarrelConfig` schema updates, custom field preservation, and effective configuration merging.
  - `internal/techstack`: Context resolution fallback hierarchy and doc profile summarization.
  - `cmd`: `battery barrel add` flags, `battery barrel doc init` scaffold command, and updated `status`/`list` views.
  - `internal/mcp`: `battery://barrels/{name}/docs` resource registration, fallback tech-stack resource handling, and metadata field enrichment in tools and topology.
  - Living capability spec deltas for `barrel-config` and `mcp-server`.
  - Upstream documentation updates in `.cooper/BATTERY.md` and `docs/`.

- **Out of Scope:**
  - Modifying files inside target non-Cooper barrel directories.
  - Enforcing Cooper TDD tracking on repositories that do not opt into Cooper.

## 4. Value & Benefits

- **Zero-Intrusion Management**: Manage any Git repository or microservice seamlessly from Battery without creating dirty worktrees or PR noise on upstream repos.
- **Autonomous Agent Grounding**: Autonomous coding agents operating at the Battery multi-barrel orchestrator level can read rich architectural context for all barrels via MCP resources and CLI tools.
- **Unified Workspace Dashboard**: Clean, informative terminal views in `battery status` and `battery barrel list`.
