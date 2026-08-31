---
name: cooper-status
description: Displays a comprehensive status overview of active [Troop](https://github.com/twoBoots/troop) worktrees, track progress, Spec Deltas, and phase checkpoints across the Cooper project.
metadata:
  version: "1.0"
---

# Cooper Status Skill

You are an AI agent. Your primary function is to provide a unified status overview of the Cooper project by inspecting active [Troop](https://github.com/twoBoots/troop) worktrees (`git troop`), the Tracks Registry (`.cooper/tracks.md`), and individual track plans.

## Operational Standards

- **Interactive Question Protocol (Mandatory):** When proposing follow-up actions or presenting single-choice/multiple-choice options to the user, agents MUST invoke available interactive question tools (e.g. `ask_question`) rather than printing text choice lists in chat.
- **Native File Tools Mandate:** Always use dedicated file tools (`view_file`, `write_to_file`, `replace_file_content`) to read or inspect files. Do NOT use shell pipes, stream editors (`sed`, `awk`), or redirections.

---

## 1. Handshake & Context Initialization

1. **Verify Handshake Index:** Check for `.cooper/index.md`. If missing, notify the user.
2. **Read Tracks Registry:** Read `.cooper/tracks.md`.
3. **Inspect Active Worktrees:**
   - Run `git troop` (or `git worktree list`) to identify all active worktree workspaces in `.worktrees/`.

---

## 2. Status Parsing Protocol

### 2.1 Parse Active Tracks & Plans
For each track found in `.cooper/tracks.md` or `.cooper/active/`:
1. Check `metadata.json` for status (`new`, `in_progress`, `completed`).
2. Read the track's `plan.md`:
   - Parse major Phases.
   - Count completed tasks (`[x]`), in-progress tasks (`[~]`), and pending tasks (`[ ]`).
   - Identify recorded checkpoint commits (e.g. `[checkpoint: <sha>]`).
3. Check for Spec Deltas in `.cooper/active/<track_id>/spec-deltas/`.

---

## 3. Present Status Overview

Output the status report in the following clean format:

```markdown
# Cooper Project Status Overview

- **Timestamp**: <Current Date/Time>
- **Living Capabilities**: <Count of specs in .cooper/specs/>
- **Active Worktrees**: <Count of active worktrees via git troop>

## Active Tracks

### Track: [<Track ID>] - <Track Description>
- **Worktree**: `.worktrees/<track_id>` (Branch: `<track_id>`)
- **Status**: [In Progress / New / Ready for PR]
- **Progress**: <Completed Tasks>/<Total Tasks> (<Percentage>%)
- **Current Phase**: Phase <N>: <Phase Name>
- **Active Task**: `<Task Name>`
- **Latest Checkpoint**: `<Commit SHA>` (or 'None')
- **Spec Deltas**: `<List of modified capability specs>`

## Historical Summary
- **Completed Tracks in Archive**: <Count of directories in .cooper/archive/>
```

### Next Action Recommendation
Provide a clear, context-aware suggestion for what to do next (e.g. "Resume implementation with `cooper-implement` on `.worktrees/<track_id>`", or "Run `cooper-review` to prepare PR").
