# Implementation Plan: Bump to Go 1.27.0

## Track Information
- **Track ID**: `bump-go-1-27-0`
- **Type**: Chore / Infrastructure
- **Worktree**: `.worktrees/bump-go-1-27-0`

---

## Phase 1: Go Module & CI/Release Toolchain Upgrade
Update Go module version directive, GitHub Actions CI & Release workflows, and tech stack specification.

- [x] Task 1.1: Update Go Module & GitHub Actions Workflows (7ebe6af)
  - [x] Sub-task: Update `go.mod` to `go 1.27.0`
  - [x] Sub-task: Update `go-version: "1.27"` in `.github/workflows/ci.yml` and `.github/workflows/release.yml`
  - [x] Sub-task: Update `.cooper/definition/tech-stack.md` to `Go 1.27+`
- [x] Task 1.2: Phase 1 Verification & Checkpoint [checkpoint: 959e11f]
  - [x] Sub-task: Fetch origin rules and sync living specs (`git fetch origin main`)
  - [x] Sub-task: Run unit tests and format checks (`go vet ./...` & `go test -v ./...`)
  - [x] Sub-task: Create checkpoint commit and push: `git push origin bump-go-1-27-0`

---

## Phase 2: Tech Stack Scaffolding Inference Engine Update (TDD)
Update Barrel tech stack auto-inference logic and tests to standardise on Go 1.27+.

- [ ] Task 2.1: Write Failing Tests for Go 1.27+ Inference (Red)
  - [ ] Sub-task: Update `internal/techstack/scaffold_test.go` to assert `"Go 1.27+"`
  - [ ] Sub-task: Update `internal/techstack/techstack_test.go` to assert `"Go 1.27+"`
  - [ ] Sub-task: Run tests to confirm red failure state
- [ ] Task 2.2: Implement Scaffolding Inference Update (Green & Refactor)
  - [ ] Sub-task: Update `InferTechStack` in `internal/techstack/scaffold.go` to return `"Go 1.27+"`
  - [ ] Sub-task: Run tests to confirm green state with coverage >80%
- [ ] Task 2.3: Phase 2 Verification & Checkpoint
  - [ ] Sub-task: Fetch origin rules (`git fetch origin main`)
  - [ ] Sub-task: Run full test suite with coverage (`go test -coverprofile=coverage.out ./...`)
  - [ ] Sub-task: Create checkpoint commit and push: `git push origin bump-go-1-27-0`

---

## Phase 3: Final Quality Verification & Spec Synchronization
Verify all quality gates, synchronize living capability specs, and register final track state.

- [ ] Task 3.1: Final Quality Verification & Spec Sync
  - [ ] Sub-task: Run `gofmt -l .` and verify clean formatting
  - [ ] Sub-task: Run `go vet ./...` and `go test -race ./...`
  - [ ] Sub-task: Merge spec deltas into `.cooper/specs/ci-release/spec.md` and `.cooper/specs/barrel-config/spec.md`
  - [ ] Sub-task: Update `metadata.json` status to `completed`
  - [ ] Sub-task: Update `.cooper/tracks.md` registry
  - [ ] Sub-task: Push final branch: `git push origin bump-go-1-27-0`
