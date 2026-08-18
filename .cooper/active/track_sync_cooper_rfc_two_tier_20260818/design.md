# Technical Design: Sync Cooper RFC Skill & Two-Tier SDD Planning Architecture

## Architecture & Integration

### 1. Project-Local Cooper Skills (`.agents/skills/`)
Package and synchronize the complete 6 native Cooper skills under `.agents/skills/`:
- `cooper-setup/SKILL.md`: Audits and scaffolds definitions and skills (including `cooper-rfc`).
- `cooper-rfc/SKILL.md`: Upstream collaborative RFC drafting, Draft PR creation, comment synthesis, and track decomposition.
- `cooper-new-track/SKILL.md`: Worktree spawning, silent scope check, and tactical spec delta/plan authoring.
- `cooper-implement/SKILL.md`: Strict TDD execution, Git Notes, and phase sync.
- `cooper-review/SKILL.md`: Principal engineer code review and PR preparation.
- `cooper-status/SKILL.md`: Active worktree and track progress inspection.

### 2. Two-Tier Planning Architecture
- Document the two-tier planning model across `.cooper/COOPER.md` and `.cooper/definition/workflow.md`:
  - **Tier 1 (Upstream RFC Alignment)**: `cooper-rfc` for epics, large refactors, cross-capability initiatives, draft PRs, and team review.
  - **Tier 2 (Downstream Track Execution)**: `cooper-new-track` & `cooper-implement` for tactical single-capability features, bug fixes, or child tracks from approved RFCs.

### 3. Agent Guidelines & Rules (`AGENTS.md` & `AGENTS.template.md`)
- Update Rule 5 to include `cooper-rfc` alongside existing skills.

### 4. Installer Protocol (`install.sh`)
- Ensure `install.sh` displays and handles `cooper-rfc`.
