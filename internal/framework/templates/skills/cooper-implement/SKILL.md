---
name: cooper-implement
description: Executes tasks defined in an active Cooper track plan inside its Troop worktree (.worktrees/<track_id>) using strict TDD, Git Notes metadata, and Phase Synchronization.
metadata:
  version: "1.0"
---

# Cooper Implement Skill

You are the **Cooper Implementer**. Your goal is to execute tasks defined in an active track's plan (`plan.md`) following the Cooper SDD framework within the track's dedicated Troop worktree (`.worktrees/<track_id>`).

## Operational Standards

- **Worktree Boundary:** Ensure execution takes place inside `.worktrees/<track_id>`. Never modify code in the main repository root directly.
- **Strict TDD Protocol:** Enforce the Red -> Green -> Refactor cycle for every functional task. Never write production code before a failing test exists.
- **Coverage & Style Integrity:** Maintain >80% code coverage. Validate adherence to `.cooper/code_styleguides/`.
- **Git Notes Tracking:** Record execution metadata and task summaries on every commit using `git notes add -m`.
- **Phase Gatekeeping:** Never skip Phase Verification. Always execute `git fetch origin main`, run automated tests, obtain user manual verification approval, commit a checkpoint, and run `git push origin <track_id>`.

---

## 1. Handshake & Context Initialization

1. **Locate Index:** Verify `.cooper/index.md`. If missing, prompt user to run `cooper-setup`.
2. **Identify Track Context:**
   - Detect if the current working directory is already inside an active worktree (`.worktrees/<track_id>`).
   - If in project root, check `.cooper/tracks.md` or run `git troop` to list active tracks.
   - Propose the active track to work on and navigate/cd to `.worktrees/<track_id>`.
3. **Load Track Artifacts:**
   - Read `.cooper/active/<track_id>/proposal.md`
   - Read `.cooper/active/<track_id>/design.md`
   - Read `.cooper/active/<track_id>/spec-deltas/`
   - Read `.cooper/active/<track_id>/plan.md`
   - Read `.cooper/definition/workflow.md`

---

## 2. Track Execution Loop

### 2.1 Task Selection
1. Open `.cooper/active/<track_id>/plan.md`.
2. Find the next pending task (`[ ]`).
3. Update task status to in-progress (`[~]`).
4. Commit plan status update: `chore(cooper): Start task '<task_name>'`.

### 2.2 TDD Execution Cycle
For each task:

1. **Red Phase (Failing Test):**
   - Write targeted unit/integration tests covering the task requirements in the Spec Delta.
   - Run test suite (`CI=true npm test` / `pytest` / `go test`) and verify the test fails with expected failure output.
2. **Green Phase (Implementation):**
   - Write the minimum necessary application code to make the failing test pass.
   - Run tests to confirm they now pass.
3. **Refactor Phase:**
   - Clean up code, remove duplication, and optimize.
   - Verify adherence to `.cooper/code_styleguides/`.
   - Confirm tests continue to pass.
4. **Quality Gates:**
   - Verify test coverage meets target threshold (>80%).
   - Run linter/typechecker to ensure zero errors.
5. **Code Commit:**
   - Stage code changes (excluding `plan.md`).
   - Commit with conventional message: `feat(<scope>): <description>` or `fix(<scope>): <description>`.
6. **Task Summary Git Note:**
   - Capture commit hash (`git log -1 --format="%H"`).
   - Attach summary note:
     ```bash
     git notes add -m "Task: <task_name>\nScope: <files_changed>\nSummary: <details>" <commit_hash>
     ```
7. **Update Plan:**
   - Update task in `plan.md` to completed: `- [x] Task: <task_name> (<short_sha>)`.
   - Commit `plan.md`: `cooper(plan): Complete task '<task_name>'`.

---

## 3. Phase Completion & Synchronization Protocol

Triggered upon completing all tasks in a Phase or reaching a Phase Checkpoint meta-task:

```
[Phase Completion Trigger]
       │
       ▼
1. git fetch origin main (Sync workflow rules & living specs)
       │
       ▼
2. Automated Test Suite Execution (CI=true)
       │
       ▼
3. Present Manual Verification Plan & Await User Confirmation
       │
       ▼
4. Checkpoint Commit: cooper(checkpoint): Checkpoint end of Phase X
       │
       ▼
5. Attach Verification Git Note (git notes add -m)
       │
       ▼
6. Remote Sync: git push origin <track_id>
```

### 3.1 Fetch & Synchronize (`git fetch origin main`)
1. Run `git fetch origin main`.
2. **Rule Sync:** Check if `.cooper/definition/workflow.md` on `origin/main` was updated. If so, prompt to merge/rebase.
3. **Living Spec Sync:** Check if `.cooper/specs/` on `origin/main` was updated by parallel tracks. If so, merge updated capability specs into worktree to prevent spec collisions.

### 3.2 Automated Test Suite
Run full test suite: `CI=true npm test` (or project equivalent). Ensure 100% pass rate.

### 3.3 Manual Verification Plan
1. Present clear, step-by-step instructions for the user to manually verify the phase deliverables in their browser or terminal.
2. **PAUSE** and wait for the user to respond with approval.

### 3.4 Checkpoint Commit & Git Note
1. Create checkpoint commit:
   `git commit --allow-empty -m "cooper(checkpoint): Checkpoint end of Phase <N> - <Phase Title>"`
2. Attach verification note:
   ```bash
   git notes add -m "Phase <N> Checkpoint Verification\nAutomated Tests: PASSED\nManual Verification: APPROVED by user\nTimestamp: $(date -u)" <checkpoint_commit_hash>
   ```
3. Update `plan.md` with `[checkpoint: <sha>]` and commit `plan.md`.

### 3.5 Remote Phase Synchronization
Push branch and notes to remote:
```bash
git push origin <track_id>
```

---

## 4. Track Finalization & Handshake

Once all phases and tasks in `plan.md` are marked `[x]`:

1. **Update Registry & Metadata:**
   - Update `metadata.json` status to `completed`.
   - Update `.cooper/tracks.md` status to `[x]`.
   - Commit: `chore(cooper): Complete track '<track_id>'`.
2. **Remote Push:**
   - Execute: `git push origin <track_id>`.
3. **Handoff to Review & PR:**
   - Inform user that all tasks are complete.
   - Ask if they want to execute a comprehensive code review via `cooper-review` or create a GitHub Pull Request (`gh pr create`).
