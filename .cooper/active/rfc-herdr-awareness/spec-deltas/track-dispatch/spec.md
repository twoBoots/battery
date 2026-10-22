# Spec Delta: track-dispatch (rfc-herdr-awareness)

Capability: `.cooper/specs/track-dispatch/spec.md`

Requirements are listed in living-spec order. New scenarios continue each requirement's existing numbering.

## Requirement 1: Macro Track Initialization (`battery track init`)

```diff
+ #### Scenario 1.2: Mapping Capabilities to Barrels
+ - **GIVEN** a valid track identifier, target barrels, and one or more `--map <barrel>=<cap>[,<cap>...]` flags
+ - **WHEN** `battery track init <track_id> --barrels <names> --map ...` is executed
+ - **THEN** the root track's `metadata.json` MUST record `barrel_capabilities` as an object from each mapped barrel name to its list of capabilities
+ - **AND** the command MUST fail without writing any file if a `--map` names a barrel that is not in `--barrels`.
```

## Requirement 2: Contract-First Track Dispatch (`battery track dispatch`)

### Scenario 2.1: Dispatching Spec Deltas to Target Barrels (modified)

```diff
-  - Seed `.cooper/active/<track_id>/` (or worktree canopy) in each target barrel with localized `metadata.json`, `proposal.md`, and `spec-deltas/<capability>/spec.md`
+  - Seed `<barrel>/.worktrees/<track_id>/.cooper/active/<track_id>/` in each target barrel with localized `metadata.json`, `proposal.md`, and a copy of `<root>/.cooper/active/<track_id>/spec-deltas/<capability>/spec.md` for each capability that the root track's `barrel_capabilities` maps to that barrel
+  - MUST NOT copy a spec delta to a barrel it is not mapped to
+  - MUST NOT write to the barrel's main checkout (`<barrel>/.cooper/active/<track_id>/`)
+  - MUST NOT write placeholder spec delta content
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
+ - **THEN** Battery MUST run `git fetch origin main`, limited to 20 seconds and concurrently with the other barrels' fetches, and then `git agent-start <track_id>` in that barrel before seeding
+ - **AND** seed only inside the resulting worktree.
+
+ #### Scenario 2.4: Worktree Creation Failing
+ - **GIVEN** `--create-worktree` and a target barrel where `git agent-start <track_id>` exits non-zero, including when Troop is not installed
+ - **WHEN** dispatch handles that barrel
+ - **THEN** it MUST fail that barrel, report the Git output, and MUST NOT seed any file in it.
+
+ #### Scenario 2.5: Offline Worktree Creation
+ - **GIVEN** `--create-worktree` and a target barrel where `git fetch origin main` fails or exceeds 20 seconds, but `git agent-start <track_id>` succeeds
+ - **WHEN** dispatch handles that barrel
+ - **THEN** it MUST warn that the worktree is based on the local `main`
+ - **AND** continue seeding.
+
+ #### Scenario 2.6: Missing Capability Mapping
+ - **GIVEN** a root track whose `metadata.json` has no `barrel_capabilities`, or has no entry for a target barrel
+ - **WHEN** `battery track dispatch <track_id>` is executed
+ - **THEN** with no `barrel_capabilities` at all, dispatch MUST fail before touching any barrel
+ - **AND** with no entry for a barrel, dispatch MUST fail that barrel, naming the missing entry
+ - **AND** MUST NOT fall back to a capability named after the track ID.
+
+ #### Scenario 2.7: Missing Root Spec Delta
+ - **GIVEN** a capability mapped to a barrel with no `<root>/.cooper/active/<track_id>/spec-deltas/<capability>/spec.md`
+ - **WHEN** dispatch handles that barrel
+ - **THEN** it MUST fail that barrel, name the missing file, and MUST NOT seed any file in it.
+
+ #### Scenario 2.8: Re-dispatch Delta States
+ - **GIVEN** a barrel that already contains the track
+ - **WHEN** `battery track dispatch <track_id>` is executed without `--force`
+ - **THEN** for each mapped capability Battery MUST report `up-to-date` when the barrel copy equals the root delta, copy the root delta and report `added` when the barrel has no copy, and report `stale` without modifying the file when the barrel copy differs
+ - **AND** report `orphaned` for each barrel delta whose capability is no longer mapped to the barrel, without modifying or deleting it
+ - **AND** give the barrel outcome `stale` if any of its deltas is `stale`.
+
+ #### Scenario 2.9: Forced Re-dispatch
+ - **GIVEN** a barrel with `stale` spec deltas
+ - **WHEN** `battery track dispatch <track_id> --force` is executed
+ - **THEN** Battery MUST overwrite each `stale` delta with its root delta and report it as `updated`
+ - **AND** MUST NOT delete `orphaned` deltas.
+
+ #### Scenario 2.10: Dispatch Exit Status
+ - **GIVEN** a dispatch without `--spawn`
+ - **WHEN** every barrel has been handled
+ - **THEN** the command MUST exit with status 1 if any barrel failed or has outcome `stale`
+ - **AND** otherwise exit with status 0.
```

## Requirement 3: Multi-Barrel Track Status Aggregation (`battery track status`)

```diff
+ #### Scenario 3.2: Live Agent State (`battery track status --live`)
+ - **GIVEN** a detected Herdr backend
+ - **WHEN** `battery track status <track_id> --live` is executed
+ - **THEN** Battery MUST call `herdr agent list` once
+ - **AND** attach to each barrel every live agent whose canonicalised `cwd` is under that barrel's canonicalised `.worktrees/<track_id>` path, compared by path component
+ - **AND** MUST NOT match on `foreground_cwd`
+ - **AND** print one `Agent` line per matched agent with its name, kind, and state, or `Agent      : none` when nothing matches.
+
+ #### Scenario 3.3: Live Status Without a Multiplexer
+ - **GIVEN** the no-op backend
+ - **WHEN** `battery track status <track_id> --live` is executed
+ - **THEN** Battery MUST print the file-based status and a single notice that live state is unavailable
+ - **AND** exit with status 0.
+
+ #### Scenario 3.4: Path Component Matching
+ - **GIVEN** worktrees `.worktrees/t1` and `.worktrees/t10`, and an agent whose `cwd` is under `.worktrees/t10`
+ - **WHEN** live status is computed for track `t1`
+ - **THEN** that agent MUST NOT be attached to track `t1`.
+
+ #### Scenario 3.5: Symlinked Paths
+ - **GIVEN** a barrel path configured through a symlink, and Herdr reporting the agent's `cwd` as the resolved path
+ - **WHEN** live status is computed
+ - **THEN** the agent MUST be matched, because both sides are canonicalised with `filepath.EvalSymlinks` before comparison.
```

## Requirement 4: Battery Root Agent Protocol & Cooper SDD Lifecycle Governance

```diff
+ #### Scenario 4.3: Multiplexer-Agnostic Guidance
+ - **GIVEN** `AGENTS.template.md`
+ - **WHEN** it describes multi-barrel dispatch
+ - **THEN** it MAY recommend `battery track dispatch --spawn` when Battery reports a multiplexer
+ - **AND** MUST NOT contain multiplexer-specific commands.
```

## Requirement 5 (new): Optional Multiplexer Integration

```diff
+ ### Requirement 5: Optional Multiplexer Integration
+ The Battery CLI MUST detect a supported terminal multiplexer only for commands that use one, and MUST behave identically to a multiplexer-free run when none is detected. The MCP server MUST NOT use a multiplexer.
+
+ #### Scenario 5.1: Herdr Detection
+ - **GIVEN** `BATTERY_HERDR` unset or not `0`, `HERDR_ENV=1`, a socket at `HERDR_SOCKET_PATH`, a resolvable `herdr` binary, and `herdr status --json` reporting `server.running` and `server.compatible` as `true`
+ - **WHEN** `battery track dispatch --spawn`, `battery track status --live`, or `battery status` starts
+ - **THEN** Battery MUST select the Herdr backend.
+
+ #### Scenario 5.2: No Multiplexer Detected
+ - **GIVEN** any condition in Scenario 5.1 does not hold
+ - **WHEN** one of those commands starts
+ - **THEN** Battery MUST select the no-op backend and MUST NOT invoke `herdr` again in that run.
+
+ #### Scenario 5.3: Detection Scope
+ - **GIVEN** any Battery command other than those in Scenario 5.1, or any MCP tool call
+ - **WHEN** it runs
+ - **THEN** Battery MUST NOT invoke `herdr`.
+
+ #### Scenario 5.4: No Persisted Multiplexer Identifiers
+ - **GIVEN** any Herdr workspace, tab, pane, or agent identifier
+ - **WHEN** Battery writes `.batteryrc`, `.batteryrc.local`, track `metadata.json`, `tracks.md`, or any seeded barrel file
+ - **THEN** Battery MUST NOT write any of those identifiers.
+
+ #### Scenario 5.5: Unmapped Barrel
+ - **GIVEN** a detected Herdr backend and a barrel for which `herdr worktree list --cwd <barrel>` reports no `source_workspace_id`
+ - **WHEN** `battery track dispatch <track_id> --spawn` handles that barrel without `--create-workspace`
+ - **THEN** Battery MUST NOT call `herdr worktree open` for it
+ - **AND** MUST record outcome `unmapped`.
+
+ #### Scenario 5.6: Unmapped Barrel With `--create-workspace`
+ - **GIVEN** the barrel in Scenario 5.5
+ - **WHEN** `battery track dispatch <track_id> --spawn --create-workspace` handles it
+ - **THEN** Battery MUST call `herdr worktree open --cwd <barrel>` for it
+ - **AND** MUST NOT call `herdr workspace create`.
+
+ #### Scenario 5.7: Spawning a Barrel Agent
+ - **GIVEN** a detected Herdr backend, a seeded barrel worktree with no `stale` deltas, a resolved source workspace, and no live agent whose canonicalised `cwd` is under the worktree
+ - **WHEN** `battery track dispatch <track_id> --spawn` handles that barrel
+ - **THEN** Battery MUST open the worktree with `herdr worktree open --workspace <source_workspace_id>` and MUST NOT call `herdr worktree create` or `herdr workspace create`
+ - **AND** start one agent in the returned root pane with no agent arguments, using the agent kind from `--agent-kind`, then `BATTERY_AGENT_KIND`, then `claude`
+ - **AND** send the fixed kickoff template once, without waiting for the agent's turn to finish
+ - **AND** record outcome `prompted`.
+
+ #### Scenario 5.8: Agent Name Normalisation
+ - **GIVEN** barrel `Web.App` and track `example-track`
+ - **WHEN** Battery names the barrel's agent
+ - **THEN** the name MUST be `web-app-example-track`
+ - **AND** every generated name MUST match `[a-z][a-z0-9_-]{0,31}`, applying the lowercase, replace, collapse, prefix and truncate rule from the RFC to the whole name
+ - **AND** on `agent_name_taken` Battery MUST truncate to 30 characters and append `-2` to `-9`, stopping with outcome `failed` after `-9`.
+
+ #### Scenario 5.9: Stale Deltas Block Spawn
+ - **GIVEN** a barrel with outcome `stale` after seeding
+ - **WHEN** `battery track dispatch <track_id> --spawn` is executed without `--force`
+ - **THEN** Battery MUST NOT start an agent or send a kickoff for that barrel.
+
+ #### Scenario 5.10: Spawn Without a Multiplexer
+ - **GIVEN** the no-op backend
+ - **WHEN** `battery track dispatch <track_id> --spawn` is executed
+ - **THEN** Battery MUST exit non-zero before seeding any barrel, stating that `--spawn` requires a supported multiplexer.
+
+ #### Scenario 5.11: Existing Agent in the Worktree
+ - **GIVEN** a live agent whose canonicalised `cwd` is under the barrel's worktree, in any pane
+ - **WHEN** `battery track dispatch <track_id> --spawn` handles that barrel
+ - **THEN** Battery MUST NOT start another agent or send the kickoff
+ - **AND** MUST record outcome `already-running` with the existing agent's name and state.
+
+ #### Scenario 5.12: Concurrent Agent Start
+ - **GIVEN** several barrels ready to start an agent
+ - **WHEN** `battery track dispatch <track_id> --spawn` reaches agent start
+ - **THEN** Battery MUST start all of their agents concurrently, each with a start timeout of 60 seconds
+ - **AND** from then on MUST target each agent by its root pane ID, not its name, whenever it prompts or waits on it.
+
+ #### Scenario 5.13: Agent Gated at Startup
+ - **GIVEN** `herdr agent start` returns `agent_not_ready`, or returns `timeout` and `herdr agent get <root_pane>` then detects an agent in the pane, or `herdr agent prompt` returns `agent_blocked`
+ - **WHEN** Battery handles that barrel
+ - **THEN** it MUST record outcome `waiting-for-operator`
+ - **AND** MUST NOT send keys or input to that agent.
+
+ #### Scenario 5.14: Agent Never Started
+ - **GIVEN** `herdr agent start` returns `timeout` and `herdr agent get <root_pane>` detects no agent in the pane
+ - **WHEN** Battery handles that barrel
+ - **THEN** it MUST record outcome `failed` with code `agent_not_detected`.
+
+ #### Scenario 5.15: Waiting on Gated Agents (`--wait-gates`)
+ - **GIVEN** one or more barrels with outcome `waiting-for-operator` after kickoffs have been sent to every ready agent
+ - **WHEN** `--wait-gates` was passed
+ - **THEN** Battery MUST wait on all of them together for at most `--gate-timeout` (default 10 minutes), for `idle`, `done`, or `working`
+ - **AND** send the kickoff and record `prompted` for each agent that reaches `idle` or `done` first
+ - **AND** send nothing and record `operator-active` for each agent that reaches `working` first
+ - **AND** record `failed` with code `agent_exited` for each agent that exits during the wait
+ - **AND** leave the rest as `waiting-for-operator`.
+
+ #### Scenario 5.16: Kickoff Text for the Operator
+ - **GIVEN** a barrel whose final outcome is `waiting-for-operator` or `operator-active`
+ - **WHEN** Battery prints its results
+ - **THEN** it MUST print the exact kickoff text for that barrel
+ - **AND** for `operator-active`, ask the operator to check whether the agent still needs it
+ - **AND** under `--wait-gates`, MUST print it only after the wait in Scenario 5.15 has ended.
+
+ #### Scenario 5.17: Herdr Error Handling
+ - **GIVEN** a `herdr` command exits non-zero
+ - **WHEN** the adapter handles the result
+ - **THEN** it MUST decide on the JSON `error.code`, not on the exit status alone
+ - **AND** record outcome `failed` with the code for any code not handled by Scenarios 5.8 and 5.13 to 5.15
+ - **AND** MUST NOT retry `agent start` or `agent prompt`, except for renaming on `agent_name_taken`.
+
+ #### Scenario 5.18: Spawn Exit Status
+ - **GIVEN** `battery track dispatch <track_id> --spawn` has handled every barrel
+ - **WHEN** it exits
+ - **THEN** it MUST exit 1 if any barrel is `failed`, `stale`, or `unmapped`
+ - **AND** otherwise exit 3 if any barrel is `waiting-for-operator` or `operator-active`
+ - **AND** otherwise exit 0.
+
+ #### Scenario 5.19: Operator Authority
+ - **GIVEN** any spawned or observed barrel agent
+ - **WHEN** Battery runs any command
+ - **THEN** Battery MUST NOT answer approvals or dialogs, close workspaces, panes, or agents, split panes, interrupt a foreground process, stop the Herdr server, or pass `--trust-repository`.
```
