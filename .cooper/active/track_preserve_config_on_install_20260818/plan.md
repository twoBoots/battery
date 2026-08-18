# Implementation Plan: Preserve Existing Configuration on Battery Install & Init

## Phase 1: Core CLI Initialization Logic & Flags
- [x] Task: Add CLI flags and configuration detection logic to `battery init` (070777c)
  - [x] Sub-task: Write unit tests for non-interactive init with existing config (Red)
  - [x] Sub-task: Implement existing config detection, preservation logic, and `--force`/`--overwrite` flag in `cmd/init.go` (Green)
  - [x] Sub-task: Refactor and verify test coverage >80% (Refactor)
- [x] Task: Interactive prompt for preserving or overwriting existing config (070777c)
  - [x] Sub-task: Write unit tests for interactive prompt flow / config preservation handler (Red)
  - [x] Sub-task: Implement interactive prompt options in `cmd/init.go` (Green)
  - [x] Sub-task: Refactor and verify test coverage >80% (Refactor)
- [x] Task: Phase 1 Verification & Checkpoint (accdfee)
  - [x] Sub-task: Workflow rule & living spec sync (`git fetch origin main`)
  - [x] Sub-task: Run automated test suite (`go test -v -cover ./...`)
  - [x] Sub-task: Manual verification of interactive & non-interactive init behavior
  - [x] Sub-task: Checkpoint commit, Git notes report & push (`git push origin track_preserve_config_on_install_20260818`)

## Phase 2: Installer Script (`install.sh`) Resilience & End-to-End Upgrade Validation
- [x] Task: Update `install.sh` fallback and documentation (5886b52)
  - [x] Sub-task: Update `install.sh` fallback logic when `.batteryrc` exists
  - [x] Sub-task: Test `install.sh` on fresh and existing setups
- [x] Task: Phase 2 Verification & Checkpoint (accdfee)
  - [x] Sub-task: Workflow rule & living spec sync (`git fetch origin main`)
  - [x] Sub-task: Run full test suite & coverage validation
  - [x] Sub-task: Manual end-to-end upgrade verification with `install.sh`
  - [x] Sub-task: Checkpoint commit, Git notes report & push (`git push origin track_preserve_config_on_install_20260818`)
