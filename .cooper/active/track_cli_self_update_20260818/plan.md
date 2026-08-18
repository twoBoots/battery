# Implementation Plan: CLI Self-Updater & Release Versioning Enhancement

## Phase 1: Core Updater Domain Logic (`internal/updater`)
- [x] Task: Platform Detection & Asset Matching (b642f6a)
  - [x] Sub-task: Write unit tests for OS/Arch mapping against release binary matrix (Red)
  - [x] Sub-task: Implement `GetPlatformBinaryName(goos, goarch)` (Green)
  - [x] Sub-task: Refactor & verify tests (Refactor)
- [x] Task: GitHub Release Fetching & Semver Logic (e03bb5d)
  - [x] Sub-task: Write unit tests for release API parsing and version comparison with httptest (Red)
  - [x] Sub-task: Implement release fetching and semver comparator (Green)
  - [x] Sub-task: Refactor & verify tests (Refactor)
- [x] Task: Binary Download & Atomic Replacement Service (740dcd3)
  - [x] Sub-task: Write unit tests for executable download and file replacement using temp directories (Red)
  - [x] Sub-task: Implement safe download and swap logic (Green)
  - [x] Sub-task: Refactor & verify coverage >80% (Refactor)
- [x] Task: Phase 1 Verification & Checkpoint (41d3640)
  - [x] Sub-task: Run unit test suite and verify coverage >80%
  - [x] Sub-task: Git fetch and notes synchronization

## Phase 2: CLI Command Integration (`cmd/update.go`)
- [x] Task: `battery update` CLI Command Implementation (8567c0c)
  - [x] Sub-task: Write CLI command tests for `battery update --check`, `--force`, `--version` (Red)
  - [x] Sub-task: Implement `updateCmd` in `cmd/update.go` and wire into `RootCmd` with aliases (Green)
  - [x] Sub-task: Refactor & verify command outputs (Refactor)
- [ ] Task: Phase 2 Verification & Checkpoint
  - [ ] Sub-task: Run full CLI test suite (`go test -v ./...`)
  - [ ] Sub-task: Phase checkpoint notes and branch synchronization

## Phase 3: Specification Promotion & Documentation
- [ ] Task: Promote Living Spec & Update Docs
  - [ ] Sub-task: Promote `spec-deltas/cli-self-update/spec.md` to `.cooper/specs/cli-self-update/spec.md`
  - [ ] Sub-task: Update `README.md`, `install.sh`, and `.cooper/index.md`
- [ ] Task: Final Track Review & Checkpoint
  - [ ] Sub-task: Run full test coverage suite and format/lint verification
  - [ ] Sub-task: Push track branch to remote
