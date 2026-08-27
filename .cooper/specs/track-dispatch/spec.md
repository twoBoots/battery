# Capability Specification: Multi-Barrel Track Dispatch & Status Governance

## Description
Provides multi-barrel track lifecycle management, Contract-First Spec Delta dispatching, and cross-barrel progress aggregation across independent repositories (barrels) in active worktrees, main trunk directories, and completed archive locations.

## Requirements

### Requirement 1: Macro Track Initialization (`battery track init`)
The system MUST provide commands to initialize a new multi-barrel track in `battery` with target barrel mappings.

#### Scenario 1.1: Multi-Barrel Track Creation
- **GIVEN** a valid track identifier and list of target barrels
- **WHEN** `battery track init <track_id> --barrels <names>` is executed
- **THEN** it MUST create `.cooper/active/<track_id>/` in the current workspace containing `metadata.json`, `proposal.md`, `design.md`, and `plan.md` referencing the target barrels.

### Requirement 2: Contract-First Track Dispatch (`battery track dispatch`)
The system MUST dispatch track metadata, proposals, and initial Spec Delta templates to all target barrels while preserving planning autonomy.

#### Scenario 2.1: Dispatching Spec Deltas to Target Barrels
- **GIVEN** an active track in battery with configured target barrels
- **WHEN** `battery track dispatch <track_id>` is executed
- **THEN** it MUST:
  - Verify each target barrel exists and is configured in `.batteryrc`
  - Seed `.cooper/active/<track_id>/` (or worktree canopy) in each target barrel with localized `metadata.json`, `proposal.md`, and `spec-deltas/<capability>/spec.md`
  - Omit internal `plan.md` in the barrel to allow the local agent to generate its own TDD task plan.

### Requirement 3: Multi-Barrel Track Status Aggregation (`battery track status`)
The system MUST inspect and aggregate track execution progress across all participating barrels, supporting active worktrees, main trunk active tracks, and completed/archived tracks.

#### Scenario 3.1: Aggregated Multi-Barrel Progress Report Across Active and Archived States
- **GIVEN** a track ID dispatched across multiple barrels
- **WHEN** `battery track status <track_id>` or `battery track list` is executed
- **THEN** it MUST check each target barrel across three locations in order:
  1. `.worktrees/<track_id>/.cooper/active/<track_id>/` (active worktree)
  2. `.cooper/active/<track_id>/` (main trunk active)
  3. `.cooper/archive/<track_id>/` (completed and archived track)
- **AND** if found in `.cooper/archive/`, mark that barrel as `completed` / `archived` (100% complete)
- **AND** if found in active folders, parse `plan.md` to compute completed vs pending tasks and active phase.

### Requirement 4: Battery Root Agent Protocol & Cooper SDD Lifecycle Governance
The system MUST provide explicit root agent guidelines (`AGENTS.md` and `AGENTS.template.md`) enforcing Cooper Spec-Driven Development (SDD) compliance and preventing context breakdown when single-barrel skills are triggered from a Battery root.

#### Scenario 4.1: Battery Root Cooper Skill Invocation & Target Barrel Resolution
- **GIVEN** an AI coding agent or developer running Cooper lifecycle skills (`cooper-new-track`, `cooper-implement`, `cooper-rfc`, `cooper-review`) from the Battery root
- **WHEN** `.cooper/index.md` is not present at the root or barrels reside in subpackages/sibling directories
- **THEN** the agent MUST NOT bypass Cooper SDD or proceed with un-specced, un-planned code editing
- **AND** the agent MUST inspect `.batteryrc` / `.batteryrc.local` to identify registered barrels
- **AND** prompt or resolve the intended target barrel and operate within that barrel's worktree context (`<barrel_path>/.worktrees/<track_id>`)
- **AND** strictly generate `proposal.md`, `design.md`, `spec-deltas/`, and `plan.md` before implementation begins.

#### Scenario 4.2: Multi-Barrel vs Single-Barrel Track Scope Decision Matrix
- **GIVEN** a feature, bug fix, or architecture initiative being initiated from a Battery workspace
- **WHEN** determining execution scope
- **THEN** the guidelines MUST direct cross-barrel initiatives affecting multiple repositories/packages to use Battery multi-barrel track orchestration (`battery track init` -> `battery track dispatch`)
- **AND** direct single-barrel initiatives to targeted Cooper tracks within the appropriate target barrel.

