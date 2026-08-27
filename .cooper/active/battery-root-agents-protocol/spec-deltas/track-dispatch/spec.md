# Spec Delta: Multi-Barrel Track Dispatch & Status Governance

## Living Specification
- Target Capability: `track-dispatch`
- Target Spec File: `.cooper/specs/track-dispatch/spec.md`

## Spec Diffs

```diff
--- a/.cooper/specs/track-dispatch/spec.md
+++ b/.cooper/specs/track-dispatch/spec.md
@@ -37,3 +37,21 @@
 - **AND** if found in `.cooper/archive/`, mark that barrel as `completed` / `archived` (100% complete)
 - **AND** if found in active folders, parse `plan.md` to compute completed vs pending tasks and active phase.
 
+### Requirement 4: Battery Root Agent Protocol & Cooper SDD Lifecycle Governance
+The system MUST provide explicit root agent guidelines (`AGENTS.md` and `AGENTS.template.md`) enforcing Cooper Spec-Driven Development (SDD) compliance and preventing context breakdown when single-barrel skills are triggered from a Battery root.
+
+#### Scenario 4.1: Battery Root Cooper Skill Invocation & Target Barrel Resolution
+- **GIVEN** an AI coding agent or developer running Cooper lifecycle skills (`cooper-new-track`, `cooper-implement`, `cooper-rfc`, `cooper-review`) from the Battery root
+- **WHEN** `.cooper/index.md` is not present at the root or barrels reside in subpackages/sibling directories
+- **THEN** the agent MUST NOT bypass Cooper SDD or proceed with un-specced, un-planned code editing
+- **AND** the agent MUST inspect `.batteryrc` / `.batteryrc.local` to identify registered barrels
+- **AND** prompt or resolve the intended target barrel and operate within that barrel's worktree context (`<barrel_path>/.worktrees/<track_id>`)
+- **AND** strictly generate `proposal.md`, `design.md`, `spec-deltas/`, and `plan.md` before implementation begins.
+
+#### Scenario 4.2: Multi-Barrel vs Single-Barrel Track Scope Decision Matrix
+- **GIVEN** a feature, bug fix, or architecture initiative being initiated from a Battery workspace
+- **WHEN** determining execution scope
+- **THEN** the guidelines MUST direct cross-barrel initiatives affecting multiple repositories/packages to use Battery multi-barrel track orchestration (`battery track init` -> `battery track dispatch`)
+- **AND** direct single-barrel initiatives to targeted Cooper tracks within the appropriate target barrel.
```
