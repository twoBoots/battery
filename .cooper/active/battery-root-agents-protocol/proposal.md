# Track Proposal: Standardize Battery Root AGENTS.md Protocol for Cooper SDD Lifecycle

## Rationale & Background

When AI coding assistants or human developers execute Cooper lifecycle skills (`cooper-new-track`, `cooper-implement`, `cooper-rfc`, etc.) from the root of a Battery multi-barrel repository:
1. Cooper skills search for `.cooper/index.md` at `./`, which is absent at the Battery root when barrels reside in child directories (`packages/*`) or sibling repositories.
2. The agent reports that Cooper is uninitialized and prompts for `cooper-setup`.
3. When informed by the user that "this is a Battery root", the agent experiences boundary and protocol breakdown: it bypasses Cooper's Spec-Driven Development (SDD) guardrails completely and jumps straight into raw code editing without producing `proposal.md`, `design.md`, `spec-deltas/`, or `plan.md`.

### Architectural Boundary Principle
* **Cooper (Level 2)**: Single-barrel SDD framework that must remain **100% pure and decoupled** from multi-barrel orchestrators. Cooper skills MUST NOT introspect `.batteryrc` directly to avoid upward coupling and leaky abstractions.
* **Battery (Level 3)**: Multi-barrel orchestrator that owns the workspace root context. Battery is responsible for establishing unequivocal instructions in `AGENTS.md` (and scaffolding template `AGENTS.template.md`) that guide the agent on how to resolve child barrels and enforce SDD guardrails.

## Proposed Changes

1. **Standardize Battery Root `AGENTS.template.md` & `AGENTS.md`**:
   - Establish an explicit **Root Runtime Protocol for Cooper Skills**:
     - Step 1: Detect multi-barrel workspace topology via `.batteryrc` / `.batteryrc.local`.
     - Step 2: Identify or prompt the user for the target barrel (`packages/<barrel>` or sibling repo).
     - Step 3: Switch context into the target barrel's worktree context (`<barrel>/.worktrees/<track_id>`) or initialize an orchestrator macro track.
     - Step 4: **Strict SDD Mandate**: Under NO circumstances may an agent bypass Spec Deltas, proposals, or plans to write code directly.
   - Establish a **Track Scope & Dispatch Decision Matrix**:
     - Cross-barrel initiatives affecting multiple repositories -> Use Battery multi-barrel track dispatch (`battery track init` / `battery track dispatch`).
     - Single-barrel initiatives -> Target the specific barrel directly and execute Cooper SDD in `<barrel>/.worktrees/<track_id>`.
2. **Update Battery Framework Documentation**:
   - Synchronize `docs/BATTERY.md` and `internal/framework/templates/docs/BATTERY.md` with the standardized root AGENTS protocol and decision matrix.

## Value & Impact
- Prevents context breakdown across all AI agents (Cursor, Claude Code, Antigravity, Copilot Workspace, Windsurf).
- Ensures 100% compliance with Spec-Driven Development (SDD) in multi-barrel setups without compromising Cooper's decoupling.
- Standardizes scaffolding so newly initialized Battery workspaces automatically inherit the robust protocol.
