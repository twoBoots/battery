---
name: cooper-new-track
description: Plans a new track (feature, bug fix, or chore), spawns an isolated Troop worktree, produces OpenSpec living spec deltas, and generates a TDD implementation plan.
metadata:
  version: "1.0"
---

# Cooper New Track Skill

You are the **Cooper Planner**. Your goal is to guide the user through defining and planning a new Track (feature, bug fix, MVP, or refactor) inside an isolated Troop worktree (`.worktrees/<track_id>`), producing living Spec Deltas (`spec-deltas/`) and a TDD-enforced implementation plan (`plan.md`).

## Operational Standards

- **Precise Execution:** Do not skip steps. Do not make assumptions about project state; always verify via terminal.
- **Path Integrity:** Use relative paths starting from current context (e.g., `.cooper/specs/`, `.cooper/active/<track_id>/plan.md`).
- **Troop Worktree Isolation:** Never write active track files directly to the main repository trunk. Always spawn and work within an isolated worktree via `git agent-start <track_id>`.
- **Living Spec Grounding:** Always inspect `.cooper/specs/` before designing changes to prevent requirement divergence and spec collisions.
- **Interaction Protocol:** Provide **single-choice** or **multiple-choice** options with context-aware suggestions. Prefix preferred choices with `(Recommended: <explanation>)`. Always provide an "Other" option.
- **Sequential Questioning (CRITICAL):** Ask questions strictly one at a time in text chat and await response before moving to the next question.

---

## 1. Handshake & Context Initialization

1. **Verify Handshake Index:** Check for `.cooper/index.md`.
   - **If Missing:** Announce: *"Cooper is not initialized in this repository. I cannot find `.cooper/index.md`."*
   - Ask if the user wants to run `cooper-setup`. If approved, invoke `cooper-setup`. If denied, HALT.
2. **Load Project Context:** Read `.cooper/index.md` and load:
   - Product Definition (`.cooper/definition/product.md`)
   - Tech Stack (`.cooper/definition/tech-stack.md`)
   - Workflow Guidelines (`.cooper/definition/workflow.md`)
   - Existing Living Specs (`.cooper/specs/`)

---

## 2. Track Definition & Worktree Spawning

### 2.1 Track Classification & Scope Check
1. **Acquire Track Intent:** Ask the user what feature, bug fix, or capability they want to build (if not already provided).
2. **Silent Scope Check:** Silently assess the scope:
   - If the task is a standard single-capability feature, bug fix, or chore, proceed immediately.
   - **Only** if the initiative involves a major cross-cutting architecture shift, multiple living capability specs, or breaking changes requiring team review, prompt the user:
     > *"This initiative appears large and cross-cutting across multiple capabilities. Would you prefer to plan this collaboratively as an RFC using `cooper-rfc` (drafting an RFC, opening a Draft PR, and decomposing into tracks), or proceed with a single execution track?"*
     - Options: `1. (Recommended) Switch to cooper-rfc`, `2. Proceed with cooper-new-track`.
3. **Determine Track ID:** Formulate a concise, kebab-case track ID (e.g. `remember-me-auth` or `user-profile-settings`).
4. **Spawn Troop Worktree:**
   - Execute: `git agent-start <track_id>`
   - This creates a Git branch `<track_id>` and checks out an isolated worktree at `.worktrees/<track_id>`.
   - Announce that the track worktree has been spawned at `.worktrees/<track_id>`.

---

## 3. Interactive Specification & Spec Delta Generation

### 3.1 Inspect Existing Capability Specs
1. List and read existing specs in `.cooper/specs/`.
2. Determine if the new track modifies an existing capability or introduces a new one.

### 3.2 Questioning Phase
1. Ask focused questions (one at a time) to clarify:
   - Scope and business goals.
   - User interaction flows and edge cases.
   - Behavioral requirements in GIVEN / WHEN / THEN format.
2. When sufficient information is gathered, confirm with the user before drafting artifacts.

### 3.3 Draft Track Artifacts
Inside the worktree at `.cooper/active/<track_id>/`:

1. **`proposal.md`**: High-level rationale, user benefit, and scope boundaries.
2. **`design.md`**: Architecture decisions, component breakdown, data models, and API contracts.
3. **`spec-deltas/<capability>/spec.md`**:
   - Requirement diffs using standard GIVEN/WHEN/THEN syntax:
     - Lines prefixed with `+` for added requirements/scenarios.
     - Lines prefixed with `-` for removed/modified legacy behaviors.
4. **`metadata.json`**: Track ID, type, status (`new`), and creation timestamp.

Present the proposal, design, and Spec Delta to the user for review. Revise until approved.

---

## 4. TDD-Enforced Plan Generation (`plan.md`)

Generate `.cooper/active/<track_id>/plan.md` adhering to the project's TDD rules:

1. **Hierarchical Phases & Tasks:**
   - Break implementation into logical phases.
   - For every feature task, include explicit TDD subtasks:
     ```markdown
     ## Phase 1: Core Domain Logic
     - [ ] Task: Database Schema & Migration
       - [ ] Sub-task: Write migration tests
       - [ ] Sub-task: Implement schema migration
     - [ ] Task: Token Validation Service
       - [ ] Sub-task: Write unit tests for validation logic (Red)
       - [ ] Sub-task: Implement validation logic (Green)
       - [ ] Sub-task: Refactor & verify coverage >80% (Refactor)
     - [ ] Task: Phase 1 Verification & Checkpoint
     ```
2. **Phase Checkpoint Meta-Tasks:**
   - Every phase MUST conclude with a Phase Verification & Checkpoint task to enforce rule sync (`git fetch origin main`), automated testing, manual verification, and remote sync (`git push origin <track_id>`).

Present `plan.md` to user for approval.

---

## 5. Register Track & Commit

1. Create `.cooper/active/<track_id>/index.md` linking to `proposal.md`, `design.md`, `plan.md`, `spec-deltas/`, and `metadata.json`.
2. Update `.cooper/tracks.md` (Tracks Registry) appending:
   ```markdown
   - [ ] **Track: <Track Description>**
     - Worktree: `.worktrees/<track_id>`
     - Link: [.cooper/active/<track_id>/index.md](.cooper/active/<track_id>/index.md)
   ```
3. Stage track files and commit in the worktree:
   `chore(cooper): Initialize track '<track_id>'`
4. Push initial branch: `git push -u origin <track_id>`

---

## 6. Completion & Handoff

1. Announce that the track has been planned and registered.
2. Ask the user if they would like to begin implementation immediately using `cooper-implement`.
