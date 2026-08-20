# Design: Sync Cooper RFC PR Approval Protocols & Reviewer Guidance

## Architectural Overview
This track brings upstream `twoBoots/cooper` PR #7 and PR #6 changes into Battery's project-local Cooper infrastructure:

1. **`cooper-rfc` Skill (`.agents/skills/cooper-rfc/SKILL.md`)**:
   - **Reviewer Actions in PR Body**: Step 5.3 includes clear feedback, approval, and graduation guidance for reviewers.
   - **Automated Approval Detection**: Step 6.1 adds `gh pr view --json state,reviews,reviewDecision --jq ...` inspection to detect `APPROVED` state as well as `/approve` comment triggers.
   - **Graduation & Ready for Review**: Step 6.2 transitions PR from draft to ready (`gh pr ready`) and registers child tracks in `.cooper/tracks.md`.

2. **Core Cooper Guidance (`.cooper/COOPER.md` & `.cooper/definition/workflow.md`)**:
   - Updates the Two-Tier SDD overview in `COOPER.md` to document the automated approval detection and graduation flow.
   - Updates workflow Principle 9 in `workflow.md` to reference `.cooper/tracks.md` registration upon RFC merge.

3. **Capability Specification Grounding (`.cooper/specs/barrel-config/spec.md`)**:
   - Updates Requirement 7 (Native Agent Skills & Two-Tier Planning Availability) to specify automated approval detection and PR lifecycle integration.
