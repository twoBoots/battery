# Cooper Hybrid Project Workflow (.cooper/)

## Guiding Principles

1. **The Plan is the Source of Truth:** All active work must be tracked in `plan.md` inside `.cooper/active/<track_id>/`.
2. **Living Capability Specifications:** System behaviors are maintained in living capability specs inside `.cooper/specs/<capability>/spec.md`. All feature changes must produce a **Spec Delta** (`spec-deltas/`) before implementation.
3. **The Tech Stack is Deliberate:** Changes to the tech stack must be documented in `.cooper/definition/tech-stack.md` *before* implementation.
4. **Isolated Workflow Updates:** If changes are required for `.cooper/definition/workflow.md` itself, they MUST be performed in isolation on a separate branch, committed, and submitted via Pull Request before returning to active tracks.
5. **Test-Driven Development (TDD):** Write unit tests before implementing functionality (Red -> Green -> Refactor).
6. **High Code Coverage:** Aim for >80% code coverage for all modules.
7. **User Experience First:** Every decision should prioritize user experience and explicit requirements.
8. **Non-Interactive & CI-Aware:** Prefer non-interactive commands. Use `CI=true` for watch-mode tools (tests, linters) to ensure single execution.

---

## Workspace Structure (`.cooper/`)

```
.cooper/
├── COOPER.md                      # Cooper SDD reference manual & cheatsheet
├── TROOP.md                       # Troop worktree reference manual
├── definition/                    # Global project definitions
│   ├── product.md                 # Product vision & initial concepts
│   ├── product-guidelines.md      # UX, branding, prose standards
│   ├── tech-stack.md              # Languages, frameworks, DBs
│   └── workflow.md                # Coverage (>80%), TDD rules, commit frequency & Troop protocol
├── code_styleguides/              # Language-specific conventions (python.md, typescript.md)
├── specs/                         # LIVING CAPABILITY SPECS (OpenSpec Living Spec model)
│   ├── auth-login/spec.md
│   ├── auth-session/spec.md
│   └── checkout-cart/spec.md
├── active/                        # ACTIVE TRACKS (Living inside .worktrees/<track_id>/)
│   └── track_add_remember_me_20260813/
│       ├── proposal.md            # High-level rationale & decisions
│       ├── design.md              # Technical architecture details
│       ├── plan.md                # TDD-enforced, phase-checkpointed plan
│       ├── metadata.json          # Track metadata & status
│       └── spec-deltas/           # Requirement diffs (+ added, - removed)
│           └── auth-session/spec.md
└── archive/                       # HISTORICAL COMPLETED TRACKS
    └── track_initial_setup_20260801/
```

---

## Track Workflow

### Track Worktree Protocol (Troop)
Before starting any work on a new track:
1. **Identify Track ID:** Determine the track ID (e.g., `extension_popup` or `auth-flow`).
2. **Spawn Track Worktree:** Create an isolated worktree under `.worktrees/<track_id>` using Troop:
   ```bash
   git agent-start <track_id>
   ```
   Under the hood, this creates a Git branch `<track_id>` and checks out the isolated worktree into `.worktrees/<track_id>`.
3. **Work in Worktree:** All work related to the track (code, plan updates, spec deltas, checkpoints) MUST be performed inside `.worktrees/<track_id>`.

---

### Track Initiation & Spec Delta Creation

When proposing a new track:
1. **Check Active Worktrees & Branches:** Run `git troop` and inspect existing branches to avoid duplicating active work.
2. **Read Living Capability Specs:** Read `.cooper/specs/<capability>/spec.md` to understand current system behavior.
3. **Create Track Workspace:** Create `.cooper/active/<track_id>/` inside `.worktrees/<track_id>`.
4. **Generate Proposals & Spec Deltas:**
   - Write `proposal.md` (high-level goal & user benefit).
   - Write `design.md` (technical architecture & decisions).
   - Write **Spec Delta** (`.cooper/active/<track_id>/spec-deltas/<capability>/spec.md`) detailing requirement additions (`+`) and removals (`-`) using GIVEN/WHEN/THEN format.
   - Write `plan.md` containing TDD sub-tasks and phase completion checkpoints.
5. **Push Initial Branch:** Push the new branch and initial track files to remote: `git push -u origin <track_id>`.

---

### Task Execution Protocol (TDD + Git Notes)

All tasks in `plan.md` follow a strict TDD lifecycle:

1. **Select Task:** Choose the next available task from `plan.md` in sequential order.
2. **Mark In Progress:** Change task status from `[ ]` to `[~]`.
3. **Red Phase (Failing Test):**
   - Write unit tests clearly defining expected behavior.
   - Run test suite and confirm that tests fail as expected.
4. **Green Phase (Implementation):**
   - Write the minimum application code necessary to make failing tests pass.
   - Run tests again to verify all pass.
5. **Refactor Phase:**
   - Clean up implementation and test code for legibility and performance.
   - Re-run tests to confirm safety.
6. **Coverage & Quality Gates:**
   - Verify test coverage meets target (>80%).
   - Verify code adheres to `.cooper/code_styleguides/` and passes lint checks.
7. **Commit Code Changes:**
   - Commit code with conventional commit message (e.g. `feat(auth): Add 30-day session validation`).
8. **Attach Task Summary Git Note:**
   - Get hash of commit (`git log -1 --format="%H"`).
   - Draft summary note (task name, changed files, rationale).
   - Attach note: `git notes add -m "<summary>" <commit_hash>`.
9. **Record Commit SHA in Plan:**
   - Update task status in `plan.md` from `[~]` to `[x]` and append the short commit SHA (e.g. `[x] Task... (a1b2c3d)`).
10. **Commit Plan Update:**
    - Commit `plan.md` change (e.g. `cooper(plan): Mark task 'Session Handler' complete`).

---

### Phase Completion & Synchronization Protocol

Triggered immediately after completing a task that concludes a Phase in `plan.md`.

#### 1. Workflow & Living Spec Synchronization (`git fetch origin main`)
- Execute `git fetch origin main`.
- **Workflow Rule Sync:** Check if `.cooper/definition/workflow.md` on `origin/main` has changed. If updated, merge or rebase `origin/main` into `.worktrees/<track_id>`.
- **Living Spec Sync:** Check if `.cooper/specs/` on `origin/main` has changed from parallel track merges. Merge updated capability specs into `.worktrees/<track_id>` to prevent spec collisions.

#### 2. Automated Testing & Coverage Verification
- Run full test suite: `CI=true npm test` (or language equivalent).
- Verify all modified code files have corresponding passing unit tests.

#### 3. Manual Verification Plan
- Present a step-by-step manual verification guide to the user with exact setup commands and expected outcomes.
- **PAUSE** and await explicit user confirmation.

#### 4. Checkpoint Commit & Git Notes
- Create checkpoint commit: `git commit -m "cooper(checkpoint): Checkpoint end of Phase X"`.
- Attach full verification report (test command, manual steps, user confirmation) as a Git Note: `git notes add -m "<report>" <checkpoint_commit_hash>`.
- Update `plan.md` with `[checkpoint: <sha>]` and commit `plan.md`.

#### 5. Remote Phase Synchronization (`git push origin <track_id>`)
- Push branch and notes to remote: `git push origin <track_id>`.

---

### Track Finalization (PR & Teardown)

Once all phases are complete:
1. **Push Branch:** Ensure local branch is fully pushed (`git push -u origin <track_id>`).
2. **Version Bump:** Prompt user for version increment (`major`, `minor`, `patch`, or `none`). Update `package.json` / `manifest.json` if selected.
3. **Open Pull Request (`gh pr create`)**:
   - Create temporary `prbody.md` detailing Track ID, work summary, requirement mappings, Spec Delta references, and verification steps.
   - Run `gh pr create --title "feat(cooper): <track_description>" --body-file prbody.md`.
   - Remove `prbody.md`.
4. **Merge & Spec Delta Integration**:
   - Upon PR approval and merge to `main`, Cooper automatically merges Spec Deltas (`.cooper/active/<track_id>/spec-deltas/`) into living capability specs (`.cooper/specs/`).
   - Move active track to `.cooper/archive/`.
5. **Teardown Worktree (`git agent-stop <track_id>`)**:
   - Execute `git agent-stop <track_id>` to delete the `.worktrees/<track_id>` directory and local branch.
