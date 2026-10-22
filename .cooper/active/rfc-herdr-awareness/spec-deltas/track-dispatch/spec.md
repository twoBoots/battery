# Spec Delta: track-dispatch (rfc-herdr-awareness)

Capability: `.cooper/specs/track-dispatch/spec.md`

## Requirement 2: Contract-First Track Dispatch (`battery track dispatch`)

### Scenario 2.1: Dispatching Spec Deltas to Target Barrels

```diff
- Seed `.cooper/active/<track_id>/` (or worktree canopy) in each target barrel with localized `metadata.json`, `proposal.md`, and `spec-deltas/<capability>/spec.md`
+ Seed `<barrel>/.worktrees/<track_id>/.cooper/active/<track_id>/` in each target barrel with localized `metadata.json`, `proposal.md`, and `spec-deltas/<capability>/spec.md`
+ MUST NOT write to the barrel's main checkout (`<barrel>/.cooper/active/<track_id>/`)
```

```diff
+ #### Scenario 2.2: Missing Track Worktree
+ - **GIVEN** a target barrel with no `.worktrees/<track_id>` worktree
+ - **WHEN** `battery track dispatch <track_id>` is executed without `--create-worktree`
+ - **THEN** dispatch MUST fail for that barrel only, with an error naming `git -C <barrel> agent-start <track_id>`
+ - **AND** MUST NOT create any file in that barrel
+ - **AND** MUST continue dispatching to the remaining barrels.
+
+ #### Scenario 2.3: Creating the Track Worktree on Dispatch
+ - **GIVEN** a target barrel with no `.worktrees/<track_id>` worktree
+ - **WHEN** `battery track dispatch <track_id> --create-worktree` is executed
+ - **THEN** Battery MUST run `git agent-start <track_id>` in that barrel before seeding
+ - **AND** seed only inside the resulting worktree.
```

## Requirement 5 (new): Optional Multiplexer Integration

```diff
+ ### Requirement 5: Optional Multiplexer Integration
+ The system MUST detect a supported terminal multiplexer at runtime and MUST behave identically to a multiplexer-free run when none is detected.
+
+ #### Scenario 5.1: Herdr Detection
+ - **GIVEN** `HERDR_ENV=1`, a socket at `HERDR_SOCKET_PATH`, a resolvable `herdr` binary, and `herdr status --json` reporting `server.running` and `server.compatible` as `true`
+ - **AND** `BATTERY_HERDR` is unset or not `0`
+ - **WHEN** any Battery command starts
+ - **THEN** Battery MUST select the Herdr multiplexer backend
+ - **OTHERWISE** Battery MUST select the no-op backend and MUST NOT invoke `herdr`.
+
+ #### Scenario 5.2: No Persisted Multiplexer Identifiers
+ - **GIVEN** any Herdr workspace, tab, pane, or agent identifier
+ - **WHEN** Battery writes `.batteryrc`, track `metadata.json`, `tracks.md`, or any seeded barrel file
+ - **THEN** none of those identifiers MUST be written
+ - **AND** barrel-to-workspace overrides MAY appear only in `.batteryrc.local` and only as workspace labels.
+
+ #### Scenario 5.3: Barrel to Workspace Resolution
+ - **GIVEN** a detected Herdr backend
+ - **WHEN** Battery resolves a barrel's workspace
+ - **THEN** it MUST use, in order: the `.batteryrc.local` `multiplexer.herdr.workspaceLabel` override, then the workspace whose source checkout path equals the barrel's resolved absolute path
+ - **AND** if neither matches, report the barrel as unmapped and MUST NOT create a workspace unless `--create-workspace` is passed
+ - **AND** because `herdr worktree open` creates a source workspace as a side effect, MUST NOT call it for an unmapped barrel unless `--create-workspace` is passed.
+
+ #### Scenario 5.4: Spawning Barrel Agents (`battery track dispatch --spawn`)
+ - **GIVEN** a detected Herdr backend and a successfully seeded barrel worktree
+ - **WHEN** `battery track dispatch <track_id> --spawn` is executed
+ - **THEN** Battery MUST open the worktree with `herdr worktree open` and MUST NOT call `herdr worktree create`
+ - **AND** start one agent per barrel in that workspace's root pane, named `<barrel>-<short_track>` normalised to `[a-z][a-z0-9_-]{0,31}` and suffixed on `agent_name_taken`
+ - **AND** use the agent kind from `--agent-kind`, then `.batteryrc.local` `multiplexer.herdr.agentKind`, then the built-in default
+ - **AND** send the fixed kickoff template once, without waiting for the agent's turn to finish.
+
+ #### Scenario 5.4a: Agent Waiting at a Startup Gate
+ - **GIVEN** `herdr agent start` returns `agent_not_ready`
+ - **WHEN** Battery handles that barrel
+ - **THEN** it MUST report the barrel as waiting for the operator
+ - **AND** wait for the agent to reach `idle` for at most the spawn timeout (default 10 minutes, `--spawn-timeout`)
+ - **AND** send the kickoff only if `idle` is reached within the timeout, otherwise print the kickoff text without sending it
+ - **AND** MUST NOT send any keys or input to the agent while it is gated.
+
+ #### Scenario 5.4b: Re-running Spawn
+ - **GIVEN** `herdr worktree open` returns a root pane already occupied by an agent
+ - **WHEN** `battery track dispatch <track_id> --spawn` is executed again
+ - **THEN** Battery MUST NOT start another agent or send the kickoff for that barrel
+ - **AND** MUST report the existing agent's name and state and print the kickoff text.
+
+ #### Scenario 5.4c: Herdr Error Handling
+ - **GIVEN** a `herdr` command exits non-zero
+ - **WHEN** the adapter handles the result
+ - **THEN** it MUST decide on the JSON `error.code`, not on the exit status alone
+ - **AND** MUST NOT retry `agent start` or `agent prompt` except on `agent_name_taken`.
+
+ #### Scenario 5.5: Spawn Without a Multiplexer
+ - **GIVEN** the no-op backend
+ - **WHEN** `battery track dispatch <track_id> --spawn` is executed
+ - **THEN** Battery MUST exit non-zero before seeding any barrel, stating that `--spawn` requires a supported multiplexer.
+
+ #### Scenario 5.6: Operator Authority
+ - **GIVEN** any spawned or observed barrel agent
+ - **WHEN** that agent is `blocked` on an approval or question
+ - **THEN** Battery MUST NOT send keys, prompts, or approvals to it
+ - **AND** Battery MUST NOT close workspaces, panes, or agents, MUST NOT stop the Herdr server, and MUST NOT pass `--trust-repository`.
```

## Requirement 3: Multi-Barrel Track Status Aggregation (`battery track status`)

```diff
+ #### Scenario 3.2: Live Agent State (`battery track status --live`)
+ - **GIVEN** a detected Herdr backend
+ - **WHEN** `battery track status <track_id> --live` is executed, or the MCP `battery_track_status` tool is called with `live: true`
+ - **THEN** Battery MUST call `herdr agent list` once
+ - **AND** for each barrel, attach the agent whose `cwd` or `foreground_cwd` lies under `<barrel>/.worktrees/<track_id>`, reporting its name, kind, and state (`idle`, `working`, `blocked`, `done`, `unknown`)
+ - **AND** report `no agent` when none matches.
+ - **GIVEN** the no-op backend
+ - **WHEN** `--live` is passed
+ - **THEN** Battery MUST print the file-based status and a single notice that live state is unavailable, and exit zero.
```

## Requirement 4: Battery Root Agent Protocol

```diff
+ #### Scenario 4.3: Multiplexer-Agnostic Guidance
+ - **GIVEN** `AGENTS.template.md`
+ - **WHEN** it describes multi-barrel dispatch
+ - **THEN** it MAY recommend `battery track dispatch --spawn` when Battery reports a multiplexer
+ - **AND** MUST NOT contain multiplexer-specific commands.
```
