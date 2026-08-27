# Implementation Plan: Standardize Battery Root AGENTS.md Protocol

## Phase 1: Standardize AGENTS Templates & Battery Architecture Guidelines

- [x] Task 1.1: Update `AGENTS.template.md` with Battery Root Runtime Protocol & Decision Matrix
  - [x] Sub-task: Add explicit Root Runtime Protocol for Cooper Skills (Red/Draft)
  - [x] Sub-task: Add First-Class Multi-Barrel vs Single-Barrel Track Decision Matrix (Green)
  - [x] Sub-task: Ensure complete Cooper SDD Framework + Troop Workflow guidelines are included (Refactor)

- [x] Task 1.2: Synchronize `AGENTS.md` at Repository Root
  - [x] Sub-task: Update repository `AGENTS.md` with the standardized protocol (Green)
  - [x] Sub-task: Verify alignment between `AGENTS.md` and `AGENTS.template.md` (Refactor)

- [x] Task 1.3: Update Battery Architecture Reference Manuals
  - [x] Sub-task: Update `.cooper/BATTERY.md` with Root Agent Protocol & Decision Matrix (Green)
  - [x] Sub-task: Update `internal/framework/templates/docs/BATTERY.md` embedded template (Green)

- [x] Task 1.4: Phase 1 Verification & Checkpoint
  - [x] Sub-task: Run Go unit and framework tests (`go test ./...`)
  - [x] Sub-task: Run installer e2e tests (`go test -v ./cmd/... -run TestInstallScriptE2E`)
  - [x] Sub-task: Record Phase 1 Git Notes and commit checkpoint

## Phase 2: Capability Spec Integration & Final Verification

- [ ] Task 2.1: Merge Spec Delta into Living Capability Spec
  - [ ] Sub-task: Update `.cooper/specs/track-dispatch/spec.md` with Requirement 4 (Green)

- [ ] Task 2.2: Comprehensive Validation & Quality Gates
  - [ ] Sub-task: Run full test suite (`go test -v -race ./...`)
  - [ ] Sub-task: Verify clean status and formatting
  - [ ] Sub-task: Record Phase 2 Git Notes and push branch
