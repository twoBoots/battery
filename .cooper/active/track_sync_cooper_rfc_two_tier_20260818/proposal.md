# Track Proposal: Sync Cooper RFC Skill & Two-Tier SDD Planning Architecture

## Context & Motivation
Following updates in `twoBoots/cooper` (PR #3), Cooper now introduces the **`cooper-rfc`** skill for collaborative upstream architectural design (epics, large refactors, Draft PR reviews) and establishes the **Two-Tier SDD Architecture** (Upstream RFC Alignment vs Downstream Tactical Track Execution) with silent scope heuristics.

In `battery`, we need to:
1. Equip the `battery` repository with the new `.agents/skills/cooper-rfc/SKILL.md` skill.
2. Update `.agents/skills/cooper-new-track/SKILL.md` and `.agents/skills/cooper-setup/SKILL.md` to support silent scope heuristics and full skill inventory discovery.
3. Update `.cooper/COOPER.md` and `.cooper/definition/workflow.md` to document the Two-Tier Planning Model and Rule 9 (Upstream Architecture vs. Track Execution).
4. Update `AGENTS.md` and `AGENTS.template.md` to document `cooper-rfc` under Project-Local Skills (Rule 5).
5. Update `install.sh` and ensure references and installation commands reflect the 6 native Cooper skills (`cooper-setup`, `cooper-rfc`, `cooper-new-track`, `cooper-implement`, `cooper-review`, `cooper-status`).
6. Update living specification `.cooper/specs/barrel-config/spec.md`.

## Scope & Deliverables
- `.agents/skills/cooper-rfc/SKILL.md` installed.
- `.agents/skills/cooper-new-track/SKILL.md` updated with silent scope check.
- `.agents/skills/cooper-setup/SKILL.md` updated with `cooper-rfc`.
- `.cooper/COOPER.md` and `.cooper/definition/workflow.md` synchronized with Cooper Two-Tier SDD model.
- `AGENTS.md` and `AGENTS.template.md` updated.
- `install.sh` updated.
- `.cooper/specs/barrel-config/spec.md` updated.
- `.cooper/tracks.md` updated and active/completed states reconciled.
