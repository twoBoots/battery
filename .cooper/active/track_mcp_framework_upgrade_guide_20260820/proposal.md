# Track Proposal: MCP Framework & Standards Upgrade Guide

## Context & Motivation
When Battery is updated via the CLI (`battery update`), it updates the global executable binary in place. It intentionally does not modify project-local Cooper standards (`.cooper/`), workflow definitions (`.cooper/definition/workflow.md`), or agent skills (`.agents/skills/cooper-*`) because teams frequently customize these files with project-specific rules, linters, and conventions.

However, developers and teams need visibility into when new framework capabilities, skills, or guidelines are available (e.g. new Cooper RFC skills, updated TDD review protocols, or refined workflow definitions), and they need an agent-guided way to adopt them safely.

Rather than implementing a destructive or blunt CLI overwrite command, this track establishes an **MCP-first Framework Upgrade Guide**:
1. Battery's MCP server inspects local workspace standards against the embedded upstream baseline.
2. The agent is provided with tools, resources, and prompt templates to evaluate divergence, understand team customizations, and present an actionable migration path.
3. The agent guides the user to adopt upstream enhancements as a dedicated Cooper track (e.g. `track_upgrade_cooper_1_3_0`) executed in an isolated Troop worktree with full git history and PR review.

## Scope & Deliverables
- **Embedded Upstream Templates**: Embed canonical Cooper skills (`.agents/skills/cooper-*`) and framework documents (`COOPER.md`, `BATTERY.md`, etc.) in the Battery Go binary.
- **Framework Status Inspection Tool & Resource**:
  - MCP Tool `battery_framework_status`: Scans workspace `.cooper/` and `.agents/skills/` to detect versions, checksum divergence, and local customizations.
  - MCP Resource `battery://framework-status`: Exposes structured JSON status for assistant context.
  - Integration with `battery_status` to include summary framework alignment metadata.
- **Upstream Template Retrieval**:
  - MCP Tool `battery_get_template`: Allows agents to retrieve the exact upstream markdown content of any canonical skill or specification.
  - MCP Resource `battery://templates/{template_name}`.
- **MCP Prompt Template**:
  - `guide_framework_upgrade_track`: Guided prompt for AI agents to evaluate divergence, draft track proposals/spec deltas, and execute non-destructive upgrades.
- **Tests & Documentation**: Full unit and integration test coverage for MCP tools, resources, and prompt templates.
