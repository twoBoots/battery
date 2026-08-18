# Battery Agent Rules (Multi-Barrel Orchestrator)

1. **Isolation Protocol**: Multi-barrel tracks defined in `battery` dispatch tasks into target barrel worktrees under `.worktrees/<track_id>`.
2. **Branching Strategy**: Track IDs must be consistent across `battery` and target barrels.
3. **Execution & Cleanup**: Use `git agent-start <track_id>` in target barrels for isolated work, and `git agent-stop <track_id>` upon PR merge.
4. **Configuration & Topology**: Respect `.batteryrc` (canonical) and `.batteryrc.local` (local overrides). Do not hardcode tech stacks; resolve each barrel's `.cooper/definition/tech-stack.md`.

See [.cooper/BATTERY.md](.cooper/BATTERY.md) and [.cooper/COOPER.md](.cooper/COOPER.md) for full context.

# Agent Guidelines (Cooper SDD Framework + Troop Workflow)

## Operational Rules

1. **Cooper Framework Mandate (.cooper/)**:
   - All feature development, bug fixes, and system changes MUST follow the **Cooper Spec-Driven Development (SDD)** lifecycle.
   - Refer to `.cooper/COOPER.md` for framework reference and `.cooper/definition/workflow.md` for track lifecycle rules.
   - Ground all planning in living capability specifications (`.cooper/specs/<capability>/spec.md`).
   - Every feature/change proposal MUST produce a **Spec Delta** (`.cooper/active/<track_id>/spec-deltas/<capability>/spec.md`) documenting requirement additions (`+`) and deletions (`-`) before code is written.

2. **Troop Worktree Isolation Protocol**:
   - Work inside an isolated worktree under `.worktrees/<track_id>`. Do NOT write feature code directly on the main repository trunk.
   - Base track worktrees off `main` using `git agent-start <track_id>`.
   - List active worktrees with `git troop`.
   - Teardown completed worktrees with `git agent-stop <track_id>` after PR approval and merge.

3. **Phase & Remote Synchronization**:
   - At phase completion, run `git fetch origin main` to synchronize workflow rules and living capability specs across parallel worktrees.
   - Push completed phase checkpoints and Git Notes metadata to remote using `git push origin <track_id>`.

4. **Quality & Execution Control**:
   - Enforce strict TDD (Red -> Green -> Refactor) and maintain test coverage >80%.
   - Attach task execution summaries and phase checkpoint reports via `git notes add -m`.

5. **Project-Local Skills (.agents/skills/)**:
   - When available, activate Cooper's dedicated project skills for structured workflows:
     - `cooper-setup`: Audit, scaffold, or reconfigure `.cooper/` infrastructure.
     - `cooper-rfc`: Plan collaborative architectural initiatives, draft RFCs & spec deltas, open Draft PRs, and decompose into tracks.
     - `cooper-new-track`: Spawn worktree, analyze living specs, and create proposal/design/spec-deltas/plan.
     - `cooper-implement`: Execute TDD tasks, record Git Notes, run phase checkpoints and syncs.
     - `cooper-review`: Conduct Principal Engineer code review against spec deltas, styleguides, and tests.
     - `cooper-status`: Inspect active worktrees, track progress, and phase checkpoints.
