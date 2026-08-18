# Track Proposal: Sync Cooper Skills, Handshake Index & Installer Refinement

## Context & Motivation
Following updates in `twoBoots/cooper` (PR #1 & PR #2), Cooper now packages 5 native project-local skills in `.agents/skills/cooper-*`, establishes `.cooper/index.md` as the unified project Handshake index, and maintains `.cooper/tracks.md`.

In `battery`, we need to:
1. Equip the `battery` repository with project-local Cooper skills (`.agents/skills/cooper-*`) and the `.cooper/index.md` / `.cooper/tracks.md` structure.
2. Refine `battery/install.sh` to cleanly rely on Cooper's installer for `.agents/skills/`, `.cooper/index.md`, and Troop setup, while ensuring Battery's orchestrator definitions (`.cooper/BATTERY.md`) and rules are registered properly.
3. Update `AGENTS.md` and `AGENTS.template.md` to reflect Rule 5 for project-local skills and the latest Cooper SDD framework rules.
4. Sync reference manuals (`.cooper/COOPER.md` and `.cooper/definition/workflow.md`).

## Scope & Deliverables
- `.agents/skills/cooper-{setup,new-track,implement,review,status}/SKILL.md` installed.
- `.cooper/index.md` created with Battery & Cooper references.
- `.cooper/tracks.md` initialized with historical and active tracks.
- `.cooper/COOPER.md` and `.cooper/definition/workflow.md` synchronized with Cooper.
- `AGENTS.md` and `AGENTS.template.md` updated with Rule 5.
- `install.sh` streamlined and verified.
