# Cooper 🛢️🐒

**Cooper** is a unified **Spec-Driven Development (SDD)** framework that merges **OpenSpec's Living Spec Deltas** with **Conductor's** quality governance and **Troop's** worktree isolation ([troop](https://github.com/twoBoots/troop)).

---

## Workspace Structure (`.cooper/`)

```
.cooper/
├── COOPER.md                      # Cooper SDD reference manual & cheatsheet (this file)
├── TROOP.md                       # Troop worktree reference manual
├── definition/                    # Project-wide baseline definitions
│   ├── product.md                 # Product vision & requirements
│   ├── product-guidelines.md      # UX, design, and product quality standards
│   ├── tech-stack.md              # Languages, frameworks, testing, and CI/CD
│   └── workflow.md                # TDD rules, coverage (>80%), and checkpoint protocol
├── code_styleguides/              # Language-specific conventions (e.g. typescript.md, python.md)
├── specs/                         # LIVING CAPABILITY SPECS (OpenSpec Living Spec model)
│   └── <capability>/spec.md       # Baseline system behavior & requirements
├── active/                        # ACTIVE TRACKS (Living inside .worktrees/<track_id>/)
│   └── <track_id>/
│       ├── proposal.md            # Rationale & business context
│       ├── design.md              # Technical architecture decisions
│       ├── plan.md                # TDD execution roadmap with status markers [ ]
│       ├── metadata.json          # Track metadata and timestamps
│       └── spec-deltas/           # Requirement diffs (+ added, - removed)
└── archive/                       # HISTORICAL COMPLETED TRACKS
```

---

## The SDD Track Lifecycle

### 1. Spawn Isolated Track Worktree
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

### 4. Execute Tasks with TDD & Git Notes
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

### 6. Pull Request & Teardown
1. Submit PR via GitHub CLI: `gh pr create --body-file prbody.md`.
2. Once merged, Spec Deltas integrate into `.cooper/specs/` and active tracks move to `.cooper/archive/`.
3. Teardown the isolated worktree:
   ```bash
   git agent-stop <track_id>
   ```

---

## Ecosystem Cheatsheet

| Component | File Reference | Primary Role |
| :--- | :--- | :--- |
| **Cooper** | `.cooper/COOPER.md` | Spec-Driven Development (SDD) & track lifecycle |
| **Troop** | `.cooper/TROOP.md` | Git worktree isolation (`git agent-start`, `git troop`, `git agent-stop`) |
| **Workflow** | `.cooper/definition/workflow.md` | Project-specific quality and operational governance |
