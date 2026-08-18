# Implementation Plan: Preserve Existing Configuration on Battery Install & Init

## Phase 1: Core CLI Initialization Logic & Flags
- [~] Task: Add CLI flags and configuration detection logic to `battery init`
  - [ ] Sub-task: Write unit tests for non-interactive init with existing config (Red)
  - [ ] Sub-task: Implement existing config detection, preservation logic, and `--force`/`--overwrite` flag in `cmd/init.go` (Green)
  - [ ] Sub-task: Refactor and verify test coverage >80% (Refactor)
- [ ] Task: Interactive prompt for preserving or overwriting existing config
  - [ ] Sub-task: Write unit tests for interactive prompt flow / config preservation handler (Red)
  - [ ] Sub-task: Implement interactive prompt options in `cmd/init.go` (Green)
  - [ ] Sub-task: Refactor and verify test coverage >80% (Refactor)
- [ ] Task: Phase 1 Verification & Checkpoint
  - [ ] Sub-task: Workflow rule & living spec sync (`git fetch origin main`)
  - [ ] Sub-task: Run automated test suite (`go test -v -cover ./...`)
  - [ ] Sub-task: Manual verification of interactive & non-interactive init behavior
  - [ ] Sub-task: Checkpoint commit, Git notes report & push (`git push origin track_preserve_config_on_install_20260818`)

## Phase 2: Installer Script (`install.sh`) Resilience & End-to-End Upgrade Validation
- [ ] Task: Update `install.sh` fallback and documentation
  - [ ] Sub-task: Update `install.sh` fallback logic when `.batteryrc` exists
  - [ ] Sub-task: Test `install.sh` on fresh and existing setups
- [ ] Task: Phase 2 Verification & Checkpoint
  - [ ] Sub-task: Workflow rule & living spec sync (`git fetch origin main`)
  - [ ] Sub-task: Run full test suite & coverage validation
  - [ ] Sub-task: Manual end-to-end upgrade verification with `install.sh`
  - [ ] Sub-task: Checkpoint commit, Git notes report & push (`git push origin track_preserve_config_on_install_20260818`)
