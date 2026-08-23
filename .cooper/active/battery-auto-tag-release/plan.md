# Implementation Plan: Battery Auto-Tag Release Pipeline

## Track Information
- **Track ID**: `battery-auto-tag-release`
- **Type**: Feature / Infrastructure
- **Worktree**: `.worktrees/battery-auto-tag-release`

---

## Phase 1: Auto-Tag Job Implementation & Workflow Refactoring
Implement `auto-tag` job and dual-release publishing in `.github/workflows/release.yml`.

- [~] Task 1.1: Add `auto-tag` Job to `release.yml`
  - [ ] Sub-task: Add `auto-tag` job with `fetch-depth: 0` and `git ls-remote` tag detection
  - [ ] Sub-task: Configure Git author and automated tag push
  - [ ] Sub-task: Update `build-and-release` job dependencies to `[ci, auto-tag]`

- [ ] Task 1.2: Implement Dual Release Publishing (`v<Version>` and `latest`)
  - [ ] Sub-task: Update `publish-release` to publish/clobber both `v<Version>` and `latest` releases
  - [ ] Sub-task: Verify YAML syntax and action versions (`@v6`, `@v7`, `@v8`)

- [ ] Task 1.3: Phase 1 Verification & Checkpoint
  - [ ] Sub-task: Sync phase rules & specs (`git fetch origin main`)
  - [ ] Sub-task: Run unit tests across all packages (`go test -v ./...`)
  - [ ] Sub-task: Push checkpoint: `git push origin battery-auto-tag-release`

---

## Phase 2: Spec Sync & Final Quality Verification
Verify repository consistency, sync living specs, and finalize track.

- [ ] Task 2.1: Final Quality Verification & Spec Sync
  - [ ] Sub-task: Run `gofmt -l .` and `go vet ./...`
  - [ ] Sub-task: Run full test suite with coverage (`go test -coverprofile=coverage.out ./...`)
  - [ ] Sub-task: Merge spec delta into `.cooper/specs/ci-release/spec.md`
  - [ ] Sub-task: Record completion metadata in `.cooper/active/battery-auto-tag-release/metadata.json`
  - [ ] Sub-task: Update `.cooper/tracks.md` registry
  - [ ] Sub-task: Push final branch: `git push origin battery-auto-tag-release`
