# Implementation Plan: Battery Go Module Casing Standardization

Standardize Go module path to `github.com/twoBoots/battery` across `go.mod`, internal package imports, and CI release workflow linker flags.

## Phase 1: Module Declaration & CI Workflow Configuration
- [~] Task 1.1: Update `go.mod` declaration and CI workflow ldflags
  - [ ] Sub-task: Update `go.mod` line 1 to `module github.com/twoBoots/battery`
  - [ ] Sub-task: Update `.github/workflows/release.yml` ldflags target to `-X github.com/twoBoots/battery/cmd.Version=${VERSION}`
- [ ] Task 1.2: Phase 1 Verification & Checkpoint
  - [ ] Sub-task: Verify `go.mod` syntax and `.github/workflows/release.yml` formatting

## Phase 2: Internal Package Import Migration
- [ ] Task 2.1: Update CLI package imports (`cmd/` and `main.go`)
  - [ ] Sub-task: Update `main.go` and `cmd/*.go` imports to `github.com/twoBoots/battery/...`
  - [ ] Sub-task: Update `cmd/*_test.go` imports to `github.com/twoBoots/battery/...`
- [ ] Task 2.2: Update Internal package imports (`internal/` packages and unit tests)
  - [ ] Sub-task: Update `internal/config`, `internal/discovery`, `internal/mcp`, `internal/techstack`, `internal/track` imports to `github.com/twoBoots/battery/...`
  - [ ] Sub-task: Update all unit tests in `internal/**/` to import `github.com/twoBoots/battery/...`
- [ ] Task 2.3: Phase 2 Verification & Checkpoint
  - [ ] Sub-task: Run `go vet ./...` and `go test ./...`
  - [ ] Sub-task: Build local binary `go build -o bin/battery .` and verify execution

## Phase 3: Final Quality Gate & Spec Sync Checkpoint
- [ ] Task 3.1: Full Test Suite & Coverage Verification
  - [ ] Sub-task: Execute `CI=true go test -v -cover ./...` ensuring all packages pass and test coverage remains >80%
  - [ ] Sub-task: Verify zero linting or compiler warnings with `go vet ./...`
- [ ] Task 3.2: Phase 3 Verification & Checkpoint
  - [ ] Sub-task: Synchronize living specs and workflow rules (`git fetch origin main`)
  - [ ] Sub-task: Prepare PR and handoff for merge
