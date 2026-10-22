# RFC: Herdr Awareness for Multi-Barrel Dispatch and Status

- **RFC ID**: rfc-herdr-awareness
- **Author**: Andrew Mitchell
- **Status**: In Review (Draft)
- **Created**: 2026-10-22
- **Revised**: 2026-10-22

## 1. Summary & Motivation

Battery's multi-barrel model assumes that each target barrel has its own autonomous agent session. The root writes contracts and spec deltas, and each barrel writes its own `plan.md` and does the TDD work (AGENTS.md, Spec-Only Orchestrator Boundary). Each barrel agent works in an isolated Git worktree at `<barrel>/.worktrees/<track_id>`, created by [Troop](https://github.com/twoBoots/troop) (`git agent-start <track_id>`), and plans with the [Cooper](https://github.com/twoBoots/cooper) `cooper-new-track` skill.

Battery has no way to reach those sessions. `battery track dispatch` writes files, then stops. A human has to open a terminal in each barrel, start an agent, and paste a kickoff prompt.

In production use this gap is visible:

- A production Battery workspace with nine barrels ran its last three multi-barrel tracks without `battery track dispatch` at all. Operators skipped it because it writes placeholder spec deltas instead of the ones authored at the root, and writes into a barrel's main checkout when no track worktree exists. Instead, barrel agents planned directly from the root's spec deltas.
- Those barrel sessions already existed, hand-arranged in Herdr (a terminal multiplexer for coding agents): one Herdr workspace per barrel, plus one for the Battery root. Cross-barrel coordination meant the operator switching workspaces and relaying findings between agents by hand.
- `battery track status` reports plan checkbox progress, but not whether a barrel agent is working, finished, or blocked waiting for a human. Today that takes visiting every workspace.

Herdr exposes a JSON CLI over a local socket. With it, a process inside a managed pane can list workspaces and agents, split panes, start a named agent of a given kind, send it a prompt, and read its lifecycle state (`idle`, `working`, `blocked`, `done`, `unknown`). It also reads Git worktrees natively, so Troop's `.worktrees/<track_id>` checkouts already show up in `herdr worktree list`.

This RFC proposes:

1. **Fixing dispatch first** (2.1), independent of Herdr: never write to a barrel's main checkout, send each barrel only the root spec deltas mapped to it, report drift on re-dispatch, and exit non-zero when a barrel fails. Automating dispatch as it stands would spread placeholder specs to agents faster.
2. **Making the Battery CLI optionally Herdr aware** (2.2 to 2.6). When the CLI runs inside a Herdr pane, `dispatch --spawn` can start a barrel agent in each track worktree and hand it the spec deltas, and `status --live` can show each barrel agent's live state. Outside Herdr, Battery behaves exactly as it does today.

The MCP server (`battery mcp`) gains only what the dispatch fix requires (2.1). It never spawns agents or calls Herdr.

## 2. Architecture & Detailed Design

### 2.1 Prerequisite: fix `battery track dispatch`

These changes do not depend on Herdr and ship first, as one track.

**a. Never write to a barrel's main checkout.** `seedBarrelTrack` (`internal/track/dispatch.go`) seeds `<barrel>/.worktrees/<track_id>/.cooper/active/<track_id>/` when the worktree exists, and otherwise falls back to `<barrel>/.cooper/active/<track_id>/`. That fallback writes untracked files onto the barrel's `main` checkout, which breaks worktree isolation. The fallback is removed:

- When `.worktrees/<track_id>` is missing, dispatch fails for that barrel with an error naming the command to run (`git -C <barrel> agent-start <track_id>`).
- With `--create-worktree`, Battery runs `git fetch origin main` in the barrel, then `git agent-start <track_id>`. The fetches for all barrels that need a worktree run concurrently, each limited to 20 seconds; a fetch that times out is treated as a failed fetch. If `agent-start` exits non-zero (including when Troop is not installed, whether as an alias or an executable), that barrel fails with the Git output, and nothing is seeded. Troop's `agent-start` ignores fetch failures and falls back to the local `main`, so Battery runs the fetch itself first: if it fails, Battery warns that the worktree is based on the local `main`, and continues.
- Seeded files are untracked in the new worktree. The barrel agent commits them as part of its track, as it does today.
- `status` keeps reading `<barrel>/.cooper/active/<track_id>/`, for tracks seeded before this change.

**b. Map capabilities to barrels.** `DispatchTrackOptions.BarrelCapMap` exists but nothing sets it, so every barrel would receive every capability's delta, including deltas written for other barrels. If a track has no capabilities, dispatch falls back to a capability named after the track ID. Instead:

- The root track's `metadata.json` gains `barrel_capabilities`, an object from barrel name to the list of capabilities that barrel implements, for example `{"api": ["track-dispatch"], "cli": ["cli-track-commands"]}`.
- `battery track init` sets it with a repeatable `--map <barrel>=<cap>[,<cap>...]` flag, and the `battery_init_track` MCP tool with a `barrel_capabilities` object. It can also be edited by hand before dispatch.
- Dispatch reads it. A barrel listed in `barrels` with no `barrel_capabilities` entry fails, with an error naming the missing entry. A track without `barrel_capabilities` fails before any barrel is touched. The track-ID fallback is removed.

**c. Copy the root's spec deltas.** Today, `seedBarrelTrack` writes a placeholder `spec.md` for each capability (a template with "Description of requirement additions and changes" and a sample GIVEN/WHEN/THEN), and never reads the deltas authored at the root. Dispatch instead copies `<root>/.cooper/active/<track_id>/spec-deltas/<capability>/spec.md` for each capability mapped to the barrel. If a mapped capability has no root delta, that barrel fails and names the missing file. No placeholder is ever written.

**d. Re-dispatch reports drift instead of hiding it.** Today, dispatch returns early when the barrel's track directory exists, and `--force` rewrites placeholders. Instead, when the barrel already has the track, dispatch compares each mapped capability's barrel delta with its root delta and reports one state per delta:

| State | Meaning | Action |
| :--- | :--- | :--- |
| `up-to-date` | Barrel copy equals the root delta. | None. |
| `added` | Root delta exists, barrel has no copy (a capability mapped after the first dispatch). | Copied. Adding cannot overwrite barrel work. |
| `stale` | Barrel copy differs from the root delta. | Left alone, unless `--force` is passed, in which case it is overwritten and reported as `updated`. |
| `orphaned` | Barrel has a delta for a capability no longer mapped to it. | Left alone and reported. Never deleted. |

A barrel with any `stale` delta (without `--force`) has outcome `stale`, which counts as a failure for the exit status and blocks `--spawn` for that barrel (2.4).

**e. Exit status.** `runTrackDispatch` (`cmd/track.go`) prints failed barrels and returns `nil`, so a partial dispatch exits 0. Dispatch exits 1 if any barrel fails or is `stale` (see 2.4 for the full exit-status table with `--spawn`).

**f. MCP dispatch.** The `battery_dispatch_track` tool (`internal/mcp/tools.go`) takes `track_id` and `force`. Without a way to create worktrees, 2.1a would leave MCP callers unable to dispatch to a new track. The tool gains `create_worktree` (boolean, default false), with the same behaviour as the flag. It also returns the per-barrel results, including failures and delta states, as structured content. `battery_init_track` gains `barrel_capabilities` (2.1b). The `plan_multi_barrel_track` MCP prompt (`internal/mcp/prompts.go`) currently says dispatch "seeds spec deltas into barrel worktrees". Its wording is corrected to match 2.1a to 2.1e.

The MCP server does not gain `spawn` or `live`. An agent starting agents through a tool call needs its own authority design (section 6, question 6).

### 2.2 Multiplexer detection (`internal/multiplexer`)

A small interface keeps Herdr out of the track packages:

```go
type Multiplexer interface {
    Name() string
    SourceWorkspace(barrelPath string) (WorkspaceRef, bool, error)
    OpenWorktree(ws *WorkspaceRef, barrelPath, worktreePath, label string) (PaneRef, error)
    StartAgent(pane PaneRef, name, kind string) (AgentRef, error)
    WaitReady(agents []AgentRef, timeout time.Duration) (map[AgentRef]AgentState, error)
    PromptAgent(agent AgentRef, text string) error
    ListAgents() ([]AgentState, error)
}
```

`StartAgent` takes no agent arguments. Battery never passes flags to the agent it starts, so it cannot be used to start an agent with its permission prompts disabled.

Detection is lazy. It runs only for the commands that use it: `battery track dispatch --spawn`, `battery track status --live`, and `battery status` (which prints the detected multiplexer and whether each barrel has a source workspace open). `Detect()` returns the Herdr implementation when all of these hold, otherwise a no-op implementation:

1. `BATTERY_HERDR` is not `0`
2. `HERDR_ENV=1`
3. `HERDR_SOCKET_PATH` exists and is a socket
4. the `herdr` binary resolves (`HERDR_BIN_PATH`, then `PATH`)
5. `herdr status --json` succeeds and reports `server.running` and `server.compatible` as `true`, so that a client/server version mismatch is caught before any mutation. Plain `herdr status` prints text and is not parsed.

Battery needs no flag of its own to enable this: Herdr already sets `HERDR_ENV`. The adapter shells out to the `herdr` CLI and parses its JSON. It does not speak the socket protocol directly, so Battery does not pin a Herdr wire version.

Herdr reports server-side failures as `{"error":{"code":"...","message":"..."}}` with exit status 1, and CLI syntax errors with exit status 2. The adapter decides on `error.code`, not on the exit status alone. Codes Battery handles explicitly: `agent_name_taken`, `invalid_agent_name`, `agent_not_ready`, `agent_blocked`, `timeout`. Any other code becomes outcome `failed` with the code attached. Battery adds two codes of its own for `failed` outcomes: `agent_not_detected` (a start timed out and no agent is in the pane) and `agent_exited` (a gated agent disappeared while Battery waited on it).

### 2.3 Barrel to workspace mapping and path matching

Herdr IDs (for example `w3` or `w3:p33`) identify workspaces, tabs and panes on one Herdr server on one machine. They MUST NOT appear in `.batteryrc`, `.batteryrc.local`, track `metadata.json`, `tracks.md`, or any other file Battery writes. Battery resolves them at runtime and uses them only within one command.

**Path canonicalisation.** Herdr reports some paths as given and others resolved (for example `/tmp/...` versus `/private/tmp/...` on macOS). Every path comparison in this RFC canonicalises both sides with `filepath.EvalSymlinks` and `filepath.Clean`, and compares path components, not string prefixes. "Under the worktree" means the path equals the worktree root or has it as a leading sequence of components, so `.worktrees/t10` is never under `.worktrees/t1`.

**Source workspace.** A barrel's source workspace is the Herdr workspace open on the barrel's main checkout: `herdr worktree list --cwd <barrel>` reports it as `source.source_workspace_id` when one is open.

There is no label-based override. Herdr allows several workspaces with the same label, so a label cannot identify one reliably, and the path match covers the normal case.

If the barrel has no source workspace, it is **unmapped**. The reason this matters: `herdr worktree open` creates a source workspace as a side effect when none is open. `--spawn` therefore skips unmapped barrels (outcome `unmapped`) unless `--create-workspace` is passed. `--create-workspace` does exactly one thing: it allows `herdr worktree open` to run for an unmapped barrel, accepting the source workspace it creates (labelled by Herdr with the repository name). Battery never calls `herdr workspace create`.

### 2.4 `battery track dispatch --spawn`

`--spawn` requires a detected multiplexer. Without one, dispatch exits non-zero before seeding any barrel.

**Agent names.** A barrel's agent is named by one rule, applied to the whole name: take `<barrel>-<track_id>`, lowercase it, replace each character outside `[a-z0-9_-]` with `-`, collapse repeated `-`, and prefix `a` if it does not start with a letter. Truncate to 32 characters. On `agent_name_taken`, truncate to 30 characters and append `-2`, then `-3`, up to `-9`. The final name is shown in the output and not written to disk. Example: barrel `Web.App`, track `example-track` gives `web-app-example-track`.

**Phase 1: prepare each barrel in turn.** These steps are fast local and Herdr calls, with no waiting on agents.

1. **Seed** (2.1). If seeding fails, the outcome is `failed`. If any mapped delta is `stale`, the outcome is `stale`: Battery does not start an agent on outdated deltas, and the operator re-runs with `--force` once the change is intended.
2. **Resolve the source workspace** (2.3). If the barrel is unmapped and `--create-workspace` was not passed, the outcome is `unmapped`.
3. **Open the worktree:** `herdr worktree open --workspace <source_workspace_id> --path <barrel>/.worktrees/<track_id> --label <barrel>:<track_id> --no-focus`. Under `--create-workspace`, an unmapped barrel uses `--cwd <barrel>` in place of `--workspace`, which lets Herdr create the source workspace. The result includes the linked workspace's `root_pane`, whose working directory is the worktree. The call is idempotent (`already_open: true`), but the pane it returns is not stable: if the operator closed the original pane, Herdr returns a different one.
4. **Check for an existing agent:** if `herdr agent list` shows any live agent whose canonical `cwd` is under the worktree (2.3), in any pane, the outcome is `already-running`. Battery reports the agent and does not start another. The guarantee is "no second agent while one is live in this worktree". Battery keeps no record of past kickoffs, so if the earlier agent has exited, a re-run starts a new one and sends the kickoff again.

**Phase 2: start agents concurrently.** For every barrel still pending, Battery runs `herdr agent start <name> --kind <kind> --pane <root_pane> --timeout 60000` at the same time, so phase 2 takes about one start timeout in total, not one per barrel.

- `<kind>` comes from `--agent-kind`, then the `BATTERY_AGENT_KIND` environment variable, then `claude`. No config file key is added.
- `agent_started`: the agent is ready. It goes to phase 3.
- `agent_not_ready`: Herdr returns this immediately when the agent is held at its own startup dialog, typically folder trust or login. The outcome is `waiting-for-operator`. Battery never answers the dialog.
- `timeout`: Battery runs `herdr agent get <root_pane>`, looking the agent up by pane rather than by name, because Herdr may not have bound the name yet. If Herdr detects an agent in the pane, it is still starting or gated, and the outcome is `waiting-for-operator`. If no agent is detected (the kind's binary is missing, crashed, or never started), the outcome is `failed` with code `agent_not_detected`.
- If the root pane is not at a shell prompt (an editor or a test run is in the foreground), `agent start` fails and the outcome is `failed` with Herdr's code. Battery never splits a pane or interrupts the foreground process.

**Agent targets.** After `agent start`, Battery addresses each agent by its root pane ID, never by name, in phases 3 and 4. Herdr agent commands accept a pane ID as the target, and an agent that came through the `timeout` path may not have its name bound. The name is for display only.

**Phase 3: send kickoffs.** For each ready agent, `herdr agent prompt <root_pane> <kickoff>`, without `--wait`. The kickoff is a fixed template naming the track ID, the copied spec delta paths, and the instruction to run `cooper-new-track` against them in the current worktree. If Herdr rejects it with `agent_blocked` (the agent hit a dialog after becoming ready), the outcome is `waiting-for-operator`. Otherwise the outcome is `prompted`. Battery does not wait for planning to finish.

**Phase 4, only with `--wait-gates`:** Battery waits on all `waiting-for-operator` agents together, for one overall timeout (`--gate-timeout`, default 10 minutes), using `herdr agent wait <root_pane> --until idle --until done --until working`:

- `idle` or `done` first: the operator cleared the dialog and the agent is ready. Battery sends the kickoff (phase 3) and the outcome becomes `prompted`.
- `working` first: something other than Battery gave the agent work. The operator has not seen the kickoff text yet, so this is either their own message or the agent doing work on startup. Battery sends nothing, so it can never send a second kickoff, and the outcome becomes `operator-active`.
- The agent exits during the wait: the outcome becomes `failed` with code `agent_exited`.
- Timeout: the outcome stays `waiting-for-operator`.

**Output** is one block per barrel: workspace label, agent name, kind, outcome, and the Herdr error code when there is one. For every barrel that ends `waiting-for-operator` or `operator-active`, Battery prints the exact kickoff text. For `waiting-for-operator`, the operator sends it after clearing the dialog. For `operator-active`, the output asks the operator to check whether the agent still needs it. Under `--wait-gates`, the text is printed only after phase 4 ends, so the operator is never invited to send a kickoff that phase 4 might send too.

**Duration.** Without `--wait-gates`, a run takes at most one 20 second fetch timeout (only with `--create-worktree`), plus local seeding, plus a few Herdr calls per barrel, plus at most one 60 second start timeout. That fits inside the shell command timeout of a root agent calling Battery, which is usually a few minutes. With `--wait-gates`, add up to `--gate-timeout`, which is why it is opt-in.

**Exit status:**

| Exit | Meaning |
| :--- | :--- |
| 0 | Every barrel is `prompted` or `already-running`. |
| 3 | No barrel `failed`, `stale` or `unmapped`, but at least one is `waiting-for-operator` or `operator-active`. |
| 1 | At least one barrel is `failed`, `stale` or `unmapped`. |

Without `--spawn`, dispatch uses 0 and 1 only (2.1e).

### 2.5 `battery track status --live`

CLI only. With a detected multiplexer, `--live` calls `herdr agent list` once and attaches to each barrel every live agent whose canonical `cwd` is under that barrel's worktree (2.3). Matching uses `cwd` only. `foreground_cwd` is ignored, because a foreground command can change directory into another barrel's worktree. If several agents match one barrel, all are listed.

The existing multi-line status format gains one line per matched agent:

```
  🟡 api [worktree]
     Status     : planning
     Progress   : 21% (3/14 tasks completed)
     Spec Deltas: track-dispatch
     Agent      : api-example-track (claude) working
```

A barrel with no matching agent shows `Agent      : none`. Without a multiplexer, `--live` prints the normal status plus one notice that live state is unavailable, and exits 0.

`blocked` is shown and never acted on. Battery MUST NOT send keys, approve, or answer on an agent's behalf. Those decisions belong to the operator.

### 2.6 Validation against Herdr

The Herdr steps in 2.2 to 2.5 were exercised against Herdr 0.9.1 in isolated named sessions (`herdr --session <name> server`), using throwaway repositories with Troop-style worktrees:

| Step | Observed |
| :--- | :--- |
| `herdr status --json` | JSON with `server.running` and `server.compatible`. Without `--json`, the output is text. |
| `herdr worktree list --cwd <barrel>` | Lists the Troop worktree. `source_workspace_id` is absent until a workspace is open for the barrel. |
| `herdr worktree open ... --no-focus` | Returns `root_pane` with its working directory set to the worktree. Creates a workspace for the source checkout if none was open. A second call returns `already_open: true`. |
| `herdr worktree open` after the root pane is closed | `already_open: true` with a different pane. The root pane is not a stable handle. |
| Two workspaces with the same label | Both created without error. Labels are not unique. |
| Paths given via a symlink (`/tmp` on macOS) | Pane `cwd` and `checkout_path` come back resolved (`/private/tmp/...`). `repo_root` and `source_checkout_path` echo the path as given. |
| `herdr agent start` (`claude`), worktree under a trusted repository | `agent_started` in about 4 seconds, state `idle`. |
| `herdr agent start` (`claude`), untrusted path | `agent_not_ready`, exit 1. The pane shows the agent's folder-trust dialog. Worktrees under a repository the agent already trusts did not show it. |
| `herdr agent start`, duplicate or invalid name | `agent_name_taken` or `invalid_agent_name`, exit 1, nothing launched. |
| `herdr agent start --timeout` | Default 30000 ms, maximum 300000 ms (from `--help`). |
| `herdr agent prompt` without `--wait` | Delivered. The agent answered the prompt. Help text: a prompt to an already-blocked agent is rejected with `agent_blocked` before any input is sent. |
| `herdr agent list`, `cwd` match | Matched the spawned agent to its worktree and nothing else. |
| `herdr agent wait --until idle` on a blocked agent | `timeout`, exit 1, no input sent. |

### 2.7 Guidance and documentation

- `AGENTS.template.md` gains one multiplexer-agnostic clause in the dispatch row: "If Battery reports a multiplexer, prefer `battery track dispatch --spawn` so barrel sessions start in the track worktree." No Herdr commands go into the template. A teammate without Herdr reads the same template and still has a working process.
- The documentation site describes `battery track init --map`, `--create-worktree`, `--spawn`, `--create-workspace`, `--wait-gates`, `--gate-timeout`, `--agent-kind`, `--live`, the `BATTERY_HERDR` and `BATTERY_AGENT_KIND` variables, the delta states, the outcomes, and the exit-status table.

### 2.8 Out of scope

- **`.batteryrc.local` barrel merge.** A local barrel entry currently replaces the canonical entry's `path` and `type` outright (`internal/config/config.go`), so a local entry without `path` leaves it empty. This RFC adds no `.batteryrc.local` keys, so it does not depend on that behaviour, but it is a latent defect worth its own track.
- Remote barrels (`herdr --machine`), multiplexers other than Herdr, and closing workspaces or agents.

## 3. Alternatives Considered

| Option | Summary | Why not (or not alone) |
| :--- | :--- | :--- |
| A. Do nothing | Operators keep starting barrel agents by hand | The production gap in section 1 stays. Dispatch stays unused. |
| B. Skill only | Ship a `cooper-battery-dispatch` agent skill that teaches the root agent to drive `herdr` after `battery track dispatch` | No Herdr code in Battery, and available as soon as the dispatch fixes (2.1) ship, so it never kicks agents off on placeholder deltas. But every agent kind needs the skill, behaviour drifts between agents, and `status` gains nothing. Kept as an **interim step** (section 5), not the end state. |
| C. Battery talks to the Herdr socket directly | Native client in Go | Ties Battery to a wire protocol that Herdr does not promise to keep stable. The CLI is the documented interface. |
| D. Herdr owns worktrees (`herdr worktree create`) | Let Herdr create worktrees during dispatch | Two tools creating worktrees under different path conventions. Troop already owns `.worktrees/<track_id>`, and Herdr reads it. |
| E. Generic tmux/zellij support | Abstract over all multiplexers now | tmux cannot report agent lifecycle state, which is most of the value. The interface in 2.2 leaves room for a second backend later without building one now. |
| F. Spawn through MCP | Add `spawn` to `battery_dispatch_track` | An agent starting agents through a tool call is a different authority model from an operator running a command. Deferred (section 6, question 6). |
| G. Wait on gated barrels one at a time | Block on each gated barrel before moving to the next | Nine gated barrels could take more than an hour, far beyond a calling agent's shell timeout. Replaced by concurrent starts and the opt-in phase 4 (2.4). |

## 4. Security, Performance & Scalability

- **Authority boundary.** `--spawn` starts agents that edit code in barrel worktrees. It is CLI only, opt-in for each invocation, and never the default. The MCP server never spawns agents or calls Herdr (2.1f). It can create branches and worktrees in barrels through `create_worktree`, which runs the Troop `agent-start` command from the user's Git setup. That stays within the authority dispatch already has: writing track files into barrels.
- **Fixed inputs.** The kickoff prompt is a fixed template with no operator free text, so spec content cannot inject instructions into the prompt itself. The agent still reads the spec deltas, as it does today. Battery passes no arguments to the agent it starts (2.2).
- **No auto-approval.** Battery never answers `blocked` agents or startup dialogs (2.4 phases 2 to 4, 2.5), and never passes `--trust-repository` to Herdr.
- **Process safety.** Battery never closes workspaces, panes, or agents, never splits panes, never interrupts a foreground process, and never stops the Herdr server.
- **Local scope only.** Battery does not use `herdr --machine`.
- **Cost.** Detection runs only for the three commands that need it (2.2): one `herdr status --json` call. `--live` adds one `herdr agent list` call per run, not one per barrel. `--spawn` adds a bounded number of calls per barrel.
- **Failure isolation.** A Herdr failure on one barrel is recorded in that barrel's outcome and does not undo seeding or affect other barrels. Because a timeout does not prove that a mutation was not applied, Battery never retries `agent start` or `agent prompt`, except for the name retry on `agent_name_taken` (which launches nothing).

## 5. Rollout & Migration Strategy

1. **Dispatch fixes** (2.1). Behaviour changes, called out in the release notes:
   - Dispatch to a barrel without a worktree now fails instead of writing to `main`. The fix is one command, or `--create-worktree`.
   - Tracks need `barrel_capabilities` (2.1b). Barrels receive only the root spec deltas mapped to them. A missing mapping or root delta fails the barrel.
   - Re-dispatch reports delta states, and a `stale` delta fails the barrel until `--force`.
   - Partial dispatch exits 1.
   - MCP callers gain `create_worktree` and `barrel_capabilities`.
2. **Interim skill** (alternative B in section 3), optional, for operators who want spawn behaviour before step 4 ships. Removed when step 4 ships.
3. **Detection and mapping** (2.2, 2.3), with `battery status` printing the detected multiplexer. Read-only.
4. **`dispatch --spawn`** (2.4).
5. **`status --live`** (2.5).
6. **Guidance and documentation** (2.7), through the normal Cooper skill sync.

Steps 3 to 5 are additive and inactive outside Herdr. Rolling them back means removing the flags. No data migration is needed: nothing Herdr-specific is ever persisted.

## 6. Open Questions & Discussion Topics

1. **Default agent kind.** Is `claude` the right default, or should `--spawn` require `--agent-kind` (or `BATTERY_AGENT_KIND`) with no default?
2. **Kickoff template ownership.** Should the kickoff prompt live in Battery, or be a Cooper template (`battery_get_template`) so barrels can override it? *Current leaning: fixed and built in for the first cut. An overridable template invites operator free text into the prompt (section 4).*
3. **Workspace topology.** Section 2.4 opens a linked worktree workspace per barrel per track. Is a tab inside the barrel's existing workspace a better default for operators running many tracks? *Current leaning: linked worktree workspace. One workspace per checkout keeps the `status --live` match simple.*
4. **Teardown.** Should `git agent-stop <track_id>` or a future `battery track close` offer to close the linked Herdr workspace, given the rule that Battery never closes workspaces?
5. **Remote barrels.** Herdr can control saved SSH machines (`--machine`). Is there demand for barrels whose worktrees live on another host?
6. **Spawn through MCP.** If root agents should spawn barrel agents through `battery_dispatch_track`, what authority model applies: a per-workspace opt-in, a confirmation step, or never?
7. **Root agent messaging.** Should Battery offer `battery track ask <track_id> <barrel> "<question>"` (a wrapper over `herdr agent prompt`) so the root can query barrel agents? Or does that let the root steer barrel planning, which the Spec-Only Orchestrator Boundary forbids? *Current leaning: leave it out. Cross-barrel findings go back through the spec deltas, so the macro contract stays the single source of truth.*
8. **Resuming a gated barrel.** Without `--wait-gates`, the operator sends the printed kickoff by hand after clearing a dialog. Is a `battery track kickoff <track_id> <barrel>` command worth adding, which finds the barrel's agent by `cwd`, checks it is ready, and sends the kickoff?
9. **Exit status 3.** Is a distinct exit status for "waiting for the operator" useful to callers, or should it fold into 0 with the outcome in the output?
