# Spec Delta: mcp-server (rfc-herdr-awareness)

Capability: `.cooper/specs/mcp-server/spec.md`

## Requirement 2: MCP Tools for Barrel & Track Operations

### Scenario 2.3: `battery_init_track` Tool (modified)

```diff
+ - **AND** accept an optional `barrel_capabilities` object from barrel name to a list of capabilities, recorded in the track's `metadata.json` exactly as `battery track init --map` records it
+ - **AND** fail without writing any file if it names a barrel that is not in `barrels`.
```

### Scenario 2.4: `battery_dispatch_track` Tool (modified)

```diff
- - **THEN** it MUST seed target barrel worktree/active folders with spec deltas and metadata while preserving local planning autonomy.
+ - **THEN** it MUST dispatch exactly as `battery track dispatch` does without `--spawn`, seeding only target barrel worktrees with the root's spec deltas and metadata while preserving local planning autonomy
+ - **AND** accept an optional boolean `create_worktree` (default `false`) with the same behaviour as `--create-worktree`
+ - **AND** return every barrel's result, including its outcome, errors, and per-delta states (`up-to-date`, `added`, `stale`, `updated`, `orphaned`), as structured content
+ - **AND** mark the tool result as an error when any barrel failed or has outcome `stale`.
```

```diff
+ #### Scenario 2.9: No Agent Spawning Over MCP
+ - **GIVEN** any MCP tool call
+ - **WHEN** the MCP server handles it
+ - **THEN** it MUST NOT start, prompt, or control coding agents
+ - **AND** MUST NOT invoke a terminal multiplexer.
```

## Requirement 4: MCP Prompts for Multi-Barrel Workflows

```diff
+ #### Scenario 4.3: Accurate Dispatch Guidance in `plan_multi_barrel_track`
+ - **GIVEN** the `plan_multi_barrel_track` prompt
+ - **WHEN** it describes the dispatch step
+ - **THEN** it MUST state that dispatch copies the root's spec deltas into existing barrel worktrees, that each barrel receives only the deltas `barrel_capabilities` maps to it, that a missing worktree fails the barrel unless `create_worktree` is set, that a missing mapping or root delta fails the barrel, and that changed root deltas are reported as `stale` until re-dispatched with `force`.
```
