# Cooper 🛢️🐒

**Cooper** is a unified **Spec-Driven Development (SDD)** framework that merges **OpenSpec's Living Spec Deltas** with **Conductor's** quality governance and **Troop's** worktree isolation ([troop](https://github.com/twoBoots/troop)).

---

## Workspace Structure (`.cooper/` & `.agents/skills/`)

```
your-project/
├── .agents/
│   └── skills/                        # Packaged Project-Local Cooper Skills
│       ├── cooper-setup/SKILL.md      # Project initialization & scaffolding
│       ├── cooper-rfc/SKILL.md        # Upstream collaborative RFC & architectural design
│       ├── cooper-new-track/SKILL.md  # Worktree spawning & spec delta planning
│       ├── cooper-implement/SKILL.md  # TDD execution & phase sync
│       ├── cooper-review/SKILL.md     # Code & spec delta review
│       └── cooper-status/SKILL.md     # Worktrees & track overview
├── .cooper/
│   ├── index.md                       # Handshake index (Single Source of Truth)
│   ├── COOPER.md                      # Cooper SDD reference manual & cheatsheet
│   ├── TROOP.md                       # Troop worktree reference manual
│   ├── tracks.md                      # Tracks Registry
│   ├── definition/                    # Project-wide baseline definitions
│   │   ├── product.md                 # Product vision & requirements
│   │   ├── product-guidelines.md      # UX, design, and product quality standards
│   │   ├── tech-stack.md              # Languages, frameworks, testing, and CI/CD
│   │   └── workflow.md                # TDD rules, coverage (>80%), and checkpoint protocol
│   ├── code_styleguides/              # Language-specific conventions (typescript.md, python.md)
│   ├── specs/                         # LIVING CAPABILITY SPECS (OpenSpec Living Spec model)
│   │   └── <capability>/spec.md       # Baseline system behavior & requirements
│   ├── active/                        # ACTIVE TRACKS (Living inside .worktrees/<track_id>/)
│   │   └── <track_id>/
│   │       ├── proposal.md            # Rationale & business context
│   │       ├── design.md              # Technical architecture decisions
│   │       ├── plan.md                # TDD execution roadmap with status markers [ ]
│   │       ├── metadata.json          # Track metadata and timestamps
│   │       └── spec-deltas/           # Requirement diffs (+ added, - removed)
│   └── archive/                       # HISTORICAL COMPLETED TRACKS
├── .worktrees/                        # Isolated Git worktrees for active tracks
└── AGENTS.md                          # Universal agent guidelines
```

---

---

## Planning Architecture: The Two-Tier Model

Cooper cleanly isolates architectural design from tactical code execution:

1. **Upstream Alignment (`cooper-rfc`)**:
   - For epics, major refactors, multi-system changes, or initiatives requiring team consensus.
   - Spawns `.worktrees/rfc-<name>`, drafts `rfc.md` + cross-capability `spec-deltas/`, and opens a **Draft Pull Request**.
   - Review loop synthesizes PR comments, registers decomposed tracks in `.cooper/tracks.md`, and awaits user PR merge to `main`.
2. **Downstream Execution (`cooper-new-track` & `workflow.md`)**:
   - For single-capability features, bug fixes, or child tracks from an approved RFC.
   - Spawns `.worktrees/<track_id>`, drafts tactical `plan.md`, runs strict TDD (Red -> Green -> Refactor), records Git Notes, and syncs checkpoints.

---

## The SDD Track Lifecycle

### 1. Spawn Isolated Track Worktree (`cooper-new-track`)
Never write feature code directly on `main`. Start an isolated worktree with Troop:
```bash
git agent-start <track_id>
```
This checks out a dedicated worktree at `.worktrees/<track_id>`.

### 2. Inspect Living Capability Specs
Before designing new features or bug fixes, read existing capability specifications:
```
.cooper/specs/<capability>/spec.md
```

### 3. Create Track Proposal, Design & Spec Deltas
Inside `.worktrees/<track_id>/.cooper/active/<track_id>/`:
* `proposal.md`: Summary of changes, intent, and value.
* `design.md`: Technical architecture and implementation details.
* `spec-deltas/<capability>/spec.md`: Requirement diffs using `+` (additions) and `-` (removals) in GIVEN/WHEN/THEN format.
* `plan.md`: Step-by-step TDD task checklist organized into Phases.

### 4. Execute Tasks with TDD & Git Notes (`cooper-implement`)
Follow the strict TDD cycle:
1. **Red**: Write failing unit tests.
2. **Green**: Write minimal code to make tests pass.
3. **Refactor**: Clean up and ensure coverage meets threshold (>80%).
4. **Git Note**: Attach task execution notes to the commit:
   ```bash
   git notes add -m "Task summary: <details>" <commit_hash>
   ```
5. **Update Plan**: Update task in `plan.md` to `[x] Task (commit_hash)`.

### 5. Phase Completion & Checkpoint
At the end of each Phase in `plan.md`:
1. **Sync**: Run `git fetch origin main` to pull latest workflow rules and living specs.
2. **Verify**: Run the full test suite (`CI=true npm test`).
3. **Checkpoint**: Commit and attach verification notes, then push:
   ```bash
   git commit -m "cooper(checkpoint): Checkpoint end of Phase X"
   git notes add -m "<verification_report>" <checkpoint_hash>
   git push origin <track_id>
   ```

### 6. Code Review, Pull Request & Teardown (`cooper-review`)
1. Review implementation against Spec Deltas and styleguides.
2. Submit PR via GitHub CLI: `gh pr create --body-file prbody.md`.
3. Once merged, Spec Deltas integrate into `.cooper/specs/` and active tracks move to `.cooper/archive/`.
4. Teardown the isolated worktree:
   ```bash
   git agent-stop <track_id>
   ```

---

## Ecosystem Cheatsheet

| Component | File Reference | Primary Role |
| :--- | :--- | :--- |
| **Cooper Skills** | `.agents/skills/cooper-*` | Executable agent skills for setup, RFCs, planning, TDD, review, and status |
| **Cooper Reference** | `.cooper/COOPER.md` | Spec-Driven Development (SDD) & track lifecycle cheatsheet |
| **Troop** | `.cooper/TROOP.md` | Git worktree isolation (`git agent-start`, `git troop`, `git agent-stop`) |
| **Workflow** | `.cooper/definition/workflow.md` | Project-specific quality and operational governance |
| **Handshake** | `.cooper/index.md` | Single source of truth index linking project context |
