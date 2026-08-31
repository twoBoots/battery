# Battery Agent Rules (Multi-Barrel Orchestrator)

1. **Isolation Protocol**: Multi-barrel tracks defined in `battery` dispatch tasks into target barrel worktrees under `.worktrees/<track_id>`.
2. **Branching Strategy**: Track IDs must be consistent across `battery` and target barrels.
3. **Execution & Cleanup**: Use `git agent-start <track_id>` in target barrels for isolated work, and `git agent-stop <track_id>` upon PR merge.
4. **Configuration & Topology**: Respect `.batteryrc` (canonical) and `.batteryrc.local` (local overrides). Do not hardcode tech stacks; resolve each barrel's `.cooper/definition/tech-stack.md`.
5. **Spec-Only Orchestrator Boundary**: When operating from the Battery multi-barrel root, agents MUST strictly limit cross-barrel track actions to authoring macro contracts, scaffolding living capability specs (`.cooper/specs/`), and dispatching spec deltas (`.cooper/active/<track_id>/spec-deltas/`). Orchestrator agents MUST NOT author target barrel `plan.md` files or execute code implementation from the Battery root—delegating local planning and TDD execution to the target barrel's autonomous session.

## Battery Root Runtime Protocol for Cooper Skills

When an agent or developer invokes single-barrel Cooper lifecycle skills (`cooper-new-track`, `cooper-implement`, `cooper-rfc`, `cooper-review`) from the Battery workspace root:

1. **Non-Bypassable SDD Mandate**:
   - If `.cooper/index.md` is absent at the root (standard for monorepos where `.cooper/` lives inside individual packages, or when operating from an orchestrator root), agents **MUST NOT** bypass Cooper's Spec-Driven Development (SDD) lifecycle or jump into un-specced, un-planned code editing.
   - Do not prompt to run `cooper-setup` at the Battery root if the intent is to work on a child barrel.
2. **Target Barrel Resolution**:
   - Inspect `.batteryrc` (and `.batteryrc.local`) to list available barrels.
   - If the target barrel was not specified in the prompt, prompt the user or identify the relevant target barrel (e.g., `packages/<barrel>` or sibling repository).
   - Navigate into the target barrel directory or execute the Cooper skill targeted within that barrel's worktree context (`<barrel>/.worktrees/<track_id>`).
   - Ground all planning in the target barrel's living specs (`<barrel>/.cooper/specs/`) and tech stack (`<barrel>/.cooper/definition/tech-stack.md`).
3. **Strict Artifact Ordering**:
   - Always generate `proposal.md`, `design.md`, `spec-deltas/<capability>/spec.md`, and `plan.md` in the target barrel's active track directory before authoring any source or test code.

## Track Scope & Dispatch Decision Matrix

| Scope | Mechanism | Workflow & Artifacts |
| :--- | :--- | :--- |
| **Multi-Barrel / Cross-Cutting Epic** | `battery track init` + `battery track dispatch` | 1. Initialize macro track in Battery root (`.cooper/active/<track_id>/`).<br>2. Author macro contracts, shared schemas, and initial Spec Deltas in Battery.<br>3. Dispatch track spec deltas to target barrels via `battery track dispatch <track_id>`.<br>4. Barrel agents spawn worktrees (`git agent-start <track_id>`) and autonomously author local `plan.md` + execute TDD. |
| **Single-Barrel Feature / Bug Fix** | Targeted Cooper Track (`cooper-new-track` in target barrel) | 1. Identify target barrel from `.batteryrc`.<br>2. Spawn worktree in target barrel (`<barrel>/.worktrees/<track_id>`).<br>3. Author local `proposal.md`, `design.md`, `spec-deltas/`, and TDD `plan.md`.<br>4. Execute TDD Red -> Green -> Refactor and submit barrel PR. |

See [.cooper/BATTERY.md](.cooper/BATTERY.md) and [.cooper/COOPER.md](.cooper/COOPER.md) for full context.

# Agent Guidelines (Cooper SDD Framework + [Troop](https://github.com/twoBoots/troop) Workflow)

## Operational Rules

1. **Cooper Framework Mandate (.cooper/)**:
   - All feature development, bug fixes, and system changes MUST follow the **Cooper Spec-Driven Development (SDD)** lifecycle.
   - Refer to `.cooper/COOPER.md` for framework reference and `.cooper/definition/workflow.md` for track lifecycle rules.
   - Ground all planning in living capability specifications (`.cooper/specs/<capability>/spec.md`).
   - Every feature/change proposal MUST produce a **Spec Delta** (`.cooper/active/<track_id>/spec-deltas/<capability>/spec.md`) documenting requirement additions (`+`) and deletions (`-`) before code is written.

2. **[Troop](https://github.com/twoBoots/troop) Worktree Isolation Protocol**:
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

6. **Interaction & Native Tool Protocols**:
   - **Interactive Question Tools:** When presenting single-choice or multiple-choice questions, options, or confirmations, agents MUST invoke available interactive question tools (e.g. `ask_question`) rather than printing text choice lists in chat. Plain text formatting is strictly a fallback when no interactive question tool is available.
   - **Native File Tools:** Use dedicated file tools (`view_file`, `write_to_file`, `replace_file_content`) for file operations. Do NOT use shell pipes, stream editors (`sed`, `awk`), heredocs (`cat << 'EOF'`), or stream redirections to create or modify files.
