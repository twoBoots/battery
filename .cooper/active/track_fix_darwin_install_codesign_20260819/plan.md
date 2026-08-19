# Implementation Plan: Fix Darwin Codesign & Atomic Binary Replacement

## Phase 1: Installer Script Fix & Local Validation
- [ ] Task: Update `install.sh` with atomic binary replacement, quarantine removal, and ad-hoc code signing
  - [ ] Sub-task: Implement atomic download/build temp file handling in `install.sh`
  - [ ] Sub-task: Add `xattr` and `codesign` handling for Darwin
  - [ ] Sub-task: Verify with end-to-end installer tests (`go test -v ./...`)
- [ ] Task: Phase 1 Verification & Checkpoint
  - [ ] Sub-task: Workflow rule & living spec sync (`git fetch origin main`)
  - [ ] Sub-task: Run automated test suite (`go test -v ./...`)
  - [ ] Sub-task: Manual verification of `install.sh` upgrading an existing binary without `Killed: 9`
  - [ ] Sub-task: Checkpoint commit, Git notes report & push (`git push origin track_fix_darwin_install_codesign_20260819`)
