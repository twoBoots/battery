# Implementation Plan: Standardize Battery Root AGENTS.md Protocol

## Phase 1: Standardize AGENTS Templates & Battery Architecture Guidelines

- [ ] Task 1.1: Update `AGENTS.template.md` with Battery Root Runtime Protocol & Decision Matrix
  - [ ] Sub-task: Add explicit Root Runtime Protocol for Cooper Skills (Red/Draft)
  - [ ] Sub-task: Add First-Class Multi-Barrel vs Single-Barrel Track Decision Matrix (Green)
  - [ ] Sub-task: Ensure complete Cooper SDD Framework + Troop Workflow guidelines are included (Refactor)

- [ ] Task 1.2: Synchronize `AGENTS.md` at Repository Root
  - [ ] Sub-task: Update repository `AGENTS.md` with the standardized protocol (Green)
  - [ ] Sub-task: Verify alignment between `AGENTS.md` and `AGENTS.template.md` (Refactor)

- [ ] Task 1.3: Update Battery Architecture Reference Manuals
  - [ ] Sub-task: Update `.cooper/BATTERY.md` with Root Agent Protocol & Decision Matrix (Green)
  - [ ] Sub-task: Update `internal/framework/templates/docs/BATTERY.md` embedded template (Green)

- [ ] Task 1.4: Phase 1 Verification & Checkpoint
  - [ ] Sub-task: Run Go unit and framework tests (`go test ./...`)
  - [ ] Sub-task: Run installer e2e tests (`go test -v ./cmd/... -run TestInstallScriptE2E`)
  - [ ] Sub-task: Record Phase 1 Git Notes and commit checkpoint

## Phase 2: Capability Spec Integration & Final Verification

- [ ] Task 2.1: Merge Spec Delta into Living Capability Spec
  - [ ] Sub-task: Update `.cooper/specs/track-dispatch/spec.md` with Requirement 4 (Green)

- [ ] Task 2.2: Comprehensive Validation & Quality Gates
  - [ ] Sub-task: Run full test suite (`go test -v -race ./...`)
  - [ ] Sub-task: Verify clean status and formatting
  - [ ] Sub-task: Record Phase 2 Git Notes and push branch
