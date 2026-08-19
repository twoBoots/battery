# Implementation Plan: Fix Darwin Codesign & Atomic Binary Replacement

## Phase 1: Installer Script Fix & Local Validation
- [x] Task: Update `install.sh` with atomic binary replacement, quarantine removal, and ad-hoc code signing
  - [x] Sub-task: Implement atomic download/build temp file handling in `install.sh`
  - [x] Sub-task: Add `xattr` and `codesign` handling for Darwin
  - [x] Sub-task: Verify with end-to-end installer tests (`go test -v ./...`)
- [x] Task: Phase 1 Verification & Checkpoint
  - [x] Sub-task: Workflow rule & living spec sync (`git fetch origin main`)
  - [x] Sub-task: Run automated test suite (`go test -v ./...`)
  - [x] Sub-task: Manual verification of `install.sh` upgrading an existing binary without `Killed: 9`
  - [x] Sub-task: Checkpoint commit, Git notes report & push (`git push origin track_fix_darwin_install_codesign_20260819`)
