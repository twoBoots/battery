# RFC: Herdr Awareness for Multi-Barrel Dispatch and Status

- **RFC ID**: rfc-herdr-awareness
- **Author**: Andrew Mitchell
- **Status**: In Review (Draft)
- **Created**: 2026-10-22

## 1. Summary & Motivation

Battery's multi-barrel model assumes that each target barrel has its own autonomous agent session. The root writes contracts and spec deltas, and each barrel writes its own `plan.md` and does the TDD work (AGENTS.md, Spec-Only Orchestrator Boundary). Battery has no mechanism for reaching those sessions. `battery track dispatch` writes files, then stops. A human has to open a terminal in each barrel, start an agent, and paste a kickoff prompt.

In production use this gap is visible:

- A production Battery workspace with nine barrels ran its last three multi-barrel tracks without `battery track dispatch` at all. Operators skipped it because it seeds stub deltas and writes into a barrel's main checkout when no track worktree exists. Instead, barrel agents planned directly from the spec deltas authored at the root.
- Those barrel sessions already existed, hand-arranged in Herdr (a terminal multiplexer for coding agents): one Herdr workspace per barrel, plus one for the Battery root. Cross-barrel coordination meant the operator switching workspaces and relaying findings between agents by hand.
- `battery track status` reports plan checkbox progress, but not whether a barrel agent is working, finished, or blocked waiting for a human. Today that takes visiting every workspace.

Herdr exposes a JSON CLI over a local socket. With it, a process inside a managed pane can list workspaces and agents, split panes, start a named agent of a given kind, send it a prompt, and read its lifecycle state (`idle`, `working`, `blocked`, `done`, `unknown`). It also reads Git worktrees natively, so Troop's `.worktrees/<track_id>` checkouts already show up in `herdr worktree list`.

This RFC proposes making Battery **optionally** Herdr aware. When Battery runs inside a Herdr pane, `dispatch` can start a barrel agent in the track worktree and hand it the spec deltas, and `status` can show each barrel agent's live state. Outside Herdr, Battery behaves exactly as it does today.

It also fixes the main-checkout dispatch defect first, because automating it would spread the damage faster.

## 2. Architecture & Detailed Design

### 2.1 Prerequisite: dispatch never writes to a barrel's main checkout

`seedBarrelTrack` (`internal/track/dispatch.go`) seeds `<barrel>/.worktrees/<track_id>/.cooper/active/<track_id>/` when the worktree exists. When it does not, it falls back to `<barrel>/.cooper/active/<track_id>/`. That fallback writes untracked files onto the barrel's `main` checkout, which breaks the Troop isolation protocol.

Change:

- When `.worktrees/<track_id>` is missing, dispatch fails for that barrel with an actionable error naming the command to run (`git -C <barrel> agent-start <track_id>`), unless `--create-worktree` is passed, in which case Battery runs `git agent-start <track_id>` in the barrel itself.
- The main-checkout fallback is removed. Reading from `.cooper/active/<track_id>/` in `status` stays unchanged for backward compatibility with tracks seeded before this change.

This change does not depend on Herdr and ships first.

### 2.2 Multiplexer detection (`internal/multiplexer`)

A small interface keeps Herdr out of the track packages:

```go
type Multiplexer interface {
    Name() string
    WorkspaceForPath(absPath string) (WorkspaceRef, bool, error)
    OpenWorktree(repoPath, worktreePath, label string) (WorkspaceRef, PaneRef, error)
    StartAgent(pane PaneRef, name, kind string, args []string) (AgentRef, error)
    PromptAgent(agent AgentRef, text string) error
    ListAgents() ([]AgentState, error)
}
```

`Detect()` returns the Herdr implementation when all of these hold, otherwise a no-op implementation:

1. `HERDR_ENV=1`
2. `HERDR_SOCKET_PATH` exists and is a socket
3. the `herdr` binary resolves (`HERDR_BIN_PATH`, then `PATH`)
4. `herdr status --json` succeeds and reports `server.running` and `server.compatible` as `true`, so that a client/server version mismatch is caught before any mutation. Plain `herdr status` prints text and is not parsed.
5. `BATTERY_HERDR` is not `0`

Battery needs no flag of its own to enable this: Herdr already sets `HERDR_ENV`. The adapter shells out to the `herdr` CLI and parses its JSON. It does not speak the socket protocol directly, so Battery does not pin a Herdr wire version.

Herdr reports server-side failures as `{"error":{"code":"...","message":"..."}}` with exit status 1, and CLI syntax errors with exit status 2. The adapter branches on `error.code` (for example `agent_name_taken`, `invalid_agent_name`, `agent_not_ready`, `timeout`), not on the exit status alone.

### 2.3 Barrel to workspace mapping

Herdr IDs (`w3`, `w3:p33`) are opaque, specific to one machine and one server, and never reused. They MUST NOT appear in `.batteryrc`, track `metadata.json`, `tracks.md`, or any other committed artifact.

Resolution order for each barrel:

1. `.batteryrc.local` override: `"multiplexer": { "herdr": { "workspaceLabel": "api" } }` on the barrel entry.
2. A Herdr workspace whose checkout path equals the barrel's resolved absolute path (via `herdr worktree list --cwd <barrel>`, `source_workspace_id`).
3. None. Dispatch reports that the barrel has no open workspace and that the operator should open one. Battery does not create workspaces for barrels the operator has not opened, unless `--create-workspace` is passed.

The check matters because `herdr worktree open` creates a workspace for the source checkout as a side effect when none is open. `--spawn` therefore resolves the barrel's workspace before opening the worktree (2.4 step 2), so that the side effect only happens when the operator asked for it.

### 2.4 `battery track dispatch --spawn`

`--spawn` requires a detected multiplexer and fails clearly without one. For each target barrel, after the 2.1 seeding succeeds:

1. Ensure the worktree exists (2.1). Troop owns worktree creation. Battery MUST NOT call `herdr worktree create`, so there is only one tool creating worktrees and their layout cannot drift.
2. Resolve the barrel's workspace (2.3). If the barrel is unmapped and `--create-workspace` was not passed, skip the barrel and report it.
3. `herdr worktree open --cwd <barrel> --path <barrel>/.worktrees/<track_id> --label <barrel>:<track_id> --no-focus`, which gives a linked worktree workspace and its `root_pane`, whose working directory is the worktree. The call is idempotent: on a re-run it returns `already_open: true` and the same root pane.
4. If `herdr agent list` shows an agent already occupying that root pane, skip the barrel: report the existing agent's name and state, and print the kickoff text for the operator to send if it is still needed. Battery never starts a second agent for the same barrel and track, and never sends the kickoff twice, because it keeps no record of having sent it.
5. `herdr agent start <name> --kind <kind> --pane <root_pane>`, where:
   - `<name>` is `<barrel>-<short_track>`, normalised to `[a-z][a-z0-9_-]{0,31}`. On `agent_name_taken` Battery appends `-2`, `-3`, and records the final name in the dispatch result, not on disk.
   - `<kind>` comes from `--agent-kind`, then `.batteryrc.local` `multiplexer.herdr.agentKind`, then `claude`.
6. `agent start` blocks until the agent is ready. If it returns `agent_not_ready`, the agent is waiting at its own startup gate, typically a folder-trust or login dialog. Battery reports the barrel as waiting for the operator, then runs `herdr agent wait <name> --until idle --timeout <spawn-timeout>` (default 10 minutes, `--spawn-timeout` to change). If the operator clears the gate and the agent reaches `idle` within the timeout, Battery continues to step 7. On `timeout`, Battery reports the barrel and prints the kickoff text instead of sending it. Battery never answers the gate itself.
7. `herdr agent prompt <name> <kickoff>`, without `--wait`. The kickoff is a fixed template naming the track ID, the dispatched spec delta paths, and the instruction to run `cooper-new-track` against them in the current worktree. Battery does not wait for planning to finish.

Barrels are processed one at a time, so a barrel waiting in step 6 delays the ones after it. Parallel spawning is left for later (section 6).

Dispatch output adds one line per barrel: workspace label, agent name, kind, and outcome (`prompted`, `waiting-timed-out`, `already-running`, `unmapped`).

### 2.5 `battery track status --live`

When a multiplexer is detected, `--live` (and the MCP `battery_track_status` tool with `live: true`) joins `herdr agent list` against each barrel's worktree path (agent `cwd` or `foreground_cwd` under `<barrel>/.worktrees/<track_id>`) and adds the agent's name, kind, and state to each `BarrelTrackSummary`:

```
Track: example-track
  api  planning   3/14 tasks   agent api-example-track  claude  working
  cli  planning   0/6 tasks    agent cli-example-track  agy     blocked
```

`blocked` is shown and never acted on. Battery MUST NOT send keys, approve, or answer on an agent's behalf. Those decisions belong to the operator.

### 2.6 Validation against Herdr

The steps in 2.2 to 2.5 were exercised against Herdr 0.9.1 in an isolated named session (`herdr --session <name> server`) using a throwaway barrel repository with a Troop-style worktree:

| Step | Observed |
| :--- | :--- |
| `herdr status --json` | JSON with `server.running` and `server.compatible`. Without `--json` the output is text. |
| `herdr worktree list --cwd <barrel>` | Lists the Troop worktree. `source_workspace_id` is absent until a workspace is open for the barrel. |
| `herdr worktree open ... --no-focus` | Returns `root_pane` with its working directory set to the worktree. Creates a workspace for the source checkout if none was open. A second call returns `already_open: true` and the same pane. |
| `herdr agent start` (`claude`), worktree under a trusted repository | `agent_started` in about 4 seconds, state `idle`. |
| `herdr agent start` (`claude`), untrusted path | `agent_not_ready`, exit 1. The pane shows the agent's folder-trust dialog. Worktrees under a repository the agent already trusts did not show it. |
| `herdr agent start`, duplicate or invalid name | `agent_name_taken` or `invalid_agent_name`, exit 1, nothing launched. |
| `herdr agent prompt` without `--wait` | Delivered. The agent answered the prompt. |
| `herdr agent list`, cwd prefix match | Matched the spawned agent to its worktree and nothing else. |
| `herdr agent wait --until idle` on a blocked agent | `timeout`, exit 1, no input sent. |

### 2.7 Agent-facing guidance

`AGENTS.template.md` gains one multiplexer-agnostic clause in the dispatch row: "If Battery reports a multiplexer, prefer `battery track dispatch --spawn` so barrel sessions start in the track worktree." No Herdr commands go into the template. A teammate without Herdr reads the same template and still has a working process.

## 3. Alternatives Considered

| Option | Summary | Why not (or not alone) |
| :--- | :--- | :--- |
| A. Do nothing | Operators keep starting barrel agents by hand | The production gap in section 1 stays. Dispatch stays unused. |
| B. Skill only | Ship a `cooper-battery-dispatch` agent skill that teaches the root agent to drive `herdr` after `battery track dispatch` | Zero Go code and available now, but every agent kind needs the skill, behaviour drifts between agents, and `status` gains nothing. Kept as an **interim step** (see section 5), not the end state. |
| C. Battery talks to the Herdr socket directly | Native client in Go | Ties Battery to a wire protocol that Herdr does not promise to keep stable. The CLI is the documented interface. |
| D. Herdr owns worktrees (`herdr worktree create`) | Let Herdr create worktrees during dispatch | Two tools creating worktrees under different path conventions. Troop already owns `.worktrees/<track_id>`, and Herdr reads it. |
| E. Generic tmux/zellij support | Abstract over all multiplexers now | tmux cannot report agent lifecycle state, which is most of the value. The interface in 2.2 leaves room for a second backend later without building one now. |

## 4. Security, Performance & Scalability

- **Authority boundary.** `--spawn` starts agents that edit code in barrel worktrees. It is opt-in for each invocation and never the default. The kickoff prompt is a fixed template with no operator free text, so spec content cannot inject instructions into the prompt itself. The agent still reads the spec deltas, as it does today.
- **No auto-approval.** Battery never answers `blocked` agents or startup gates (2.4 step 6, 2.5), and never passes `--trust-repository` to Herdr.
- **Process safety.** Battery never closes workspaces, panes, or agents it did not create, and never stops the Herdr server. It does not close the ones it created either: their lifetime belongs to the operator.
- **Local scope only.** Battery does not use `herdr --machine`. Remote barrels are out of scope (section 6).
- **Cost.** Detection is one `herdr status` call per Battery invocation, made only when `HERDR_ENV=1`. `--live` adds one `herdr agent list` call per `status` invocation, not one per barrel.
- **Failure isolation.** A Herdr failure on one barrel is recorded in that barrel's dispatch result and does not undo file seeding or affect other barrels. Because a timeout does not prove that a mutation was not applied, Battery never retries `agent start` or `agent prompt` blindly. It reports the failure and the agent name.

## 5. Rollout & Migration Strategy

1. **Dispatch worktree guard** (2.1). This is a behaviour change: dispatch against a barrel with no worktree now fails instead of writing to `main`. Called out in the release notes with the one-line fix.
2. **Interim skill** (the skill-only alternative in section 3), optional, for operators who want spawn behaviour before step 4 ships.
3. **Detection and mapping** (2.2, 2.3), with `battery status` printing the detected multiplexer and mapped workspaces. Read-only.
4. **`dispatch --spawn`** (2.4).
5. **`status --live`** and the MCP `live` field (2.5).
6. **Template wording** (2.7) through the normal Cooper skill sync.

Steps 3 to 5 are additive and inactive outside Herdr. Rolling back means removing the flag. No data migration is needed: nothing Herdr-specific is ever persisted.

## 6. Open Questions & Discussion Topics

1. **Default agent kind.** Is `claude` the right default, or should there be no default and `--spawn` require `--agent-kind` until `.batteryrc.local` sets one?
2. **Kickoff template ownership.** Should the kickoff prompt live in Battery, or be a Cooper template (`battery_get_template`) so barrels can override it? *Authors' leaning: fixed and built in for the first cut. An overridable template invites operator free text into the prompt (section 4).*
3. **Workspace topology.** Section 2.4 opens a linked worktree workspace per barrel per track. Is a tab inside the barrel's existing workspace a better default for operators running many tracks? *Authors' leaning: linked worktree workspace. It maps one workspace to one checkout, which keeps the `status --live` cwd match unambiguous.*
4. **Teardown.** Should `git agent-stop <track_id>` or a future `battery track close` offer to close the linked Herdr workspace, given the rule that Battery never closes what it did not create?
5. **Remote barrels.** Herdr can control saved SSH machines (`--machine`). Is there demand for barrels whose worktrees live on another host?
6. **Root agent messaging.** Should Battery offer `battery track ask <track_id> <barrel> "<question>"` (a wrapper over `herdr agent prompt`) so the root can query barrel agents? Or does that let the root steer barrel planning, which the Spec-Only Orchestrator Boundary forbids? *Authors' leaning: leave it out. Cross-barrel findings go back through the spec deltas, so the macro contract stays the single source of truth.*
7. **Parallel spawning.** Section 2.4 spawns barrels one at a time, so one barrel waiting at a startup gate holds up the rest for up to the spawn timeout. Should Battery start all barrels first and then wait on the gated ones together, at the cost of a more complex output and failure model?
