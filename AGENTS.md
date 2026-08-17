# Battery Agent Rules (Multi-Barrel Orchestrator)

1. **Isolation Protocol**: Multi-barrel tracks defined in `battery` dispatch tasks into target barrel worktrees under `.worktrees/<track_id>`.
2. **Branching Strategy**: Track IDs must be consistent across `battery` and target barrels.
3. **Execution & Cleanup**: Use `git agent-start <track_id>` in target barrels for isolated work, and `git agent-stop <track_id>` upon PR merge.
4. **Configuration & Topology**: Respect `.batteryrc` (canonical) and `.batteryrc.local` (local overrides). Do not hardcode tech stacks; resolve each barrel's `.cooper/definition/tech-stack.md`.

See [.cooper/BATTERY.md](.cooper/BATTERY.md) and [.cooper/COOPER.md](.cooper/COOPER.md) for full context.

# Agent Guidelines (Cooper Hybrid Framework + Troop Workflow)

## Rules

1. **Cooper Hybrid Mandate (.cooper/)**:
   - All feature development, bug fixes, and system changes MUST follow the **Cooper Hybrid Workflow**.
   - Refer to `.cooper/COOPER.md` for a quick reference and `.cooper/definition/workflow.md` for full track lifecycle guidelines.
   - Read living capability specs from `.cooper/specs/<capability>/spec.md` before starting new tracks.
   - All proposed feature changes MUST produce a **Spec Delta** (`.cooper/active/<track_id>/spec-deltas/<capability>/spec.md`) showing requirement additions (`+`) and deletions (`-`) before code is written.

2. **Troop Worktree Isolation Protocol**:
   - Work inside an isolated worktree under `.worktrees/<track_id>`. Do NOT modify code in the main repository trunk directly.
   - Refer to `.cooper/TROOP.md` for complete Troop worktree commands and guidelines.
   - Base track worktrees off `main` using `git agent-start <track_id>`.

3. **Phase & Remote Synchronization**:
   - At phase completion, run `git fetch origin main` to synchronize workflow rules and living capability specs across parallel worktrees.
   - Push completed phase checkpoints and Git Notes metadata to remote using `git push origin <track_id>`.

4. **Quality & Execution Control**:
   - Enforce TDD (Red -> Green -> Refactor) and maintain test coverage >80%.
   - Attach task execution summaries and phase checkpoint reports via `git notes add -m`.
   - List active tracks with `git troop`.
   - Teardown completed worktrees with `git agent-stop <track_id>` after PR approval and merge.
