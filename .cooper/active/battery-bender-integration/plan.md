# Implementation Plan: Battery Bender CLI Integration

Refactor Battery to import and standardize on `github.com/twoBoots/bender` for self-updating and MCP server/installer capabilities, removing redundant internal packages.

---

## Phase 1: Dependency Integration & Updater Refactoring
Refactor `battery update` to use `bender/pkg/updater` and remove `internal/updater`.

- [x] Task 1.1: Add `github.com/twoBoots/bender` dependency to `go.mod` (74a4997)
  - [x] Sub-task: Add dependency in `go.mod` and run `go mod tidy`
  - [x] Sub-task: Verify package imports succeed

- [x] Task 1.2: Refactor `cmd/update.go` to use `bender/pkg/updater` (8555ea1)
  - [x] Sub-task: Write unit tests in `cmd/update_test.go` for Bender updater integration (Red)
  - [x] Sub-task: Update `cmd/update.go` with `updater.SelfUpdate` implementation (Green)
  - [x] Sub-task: Remove `internal/updater/` package and obsolete updater tests (Refactor)
  - [x] Sub-task: Verify all `cmd/` tests pass with coverage >80%

- [x] Task 1.3: Phase 1 Verification & Checkpoint [checkpoint: 95eb265]
  - [x] Sub-task: Sync phase rules & specs (`git fetch origin main`)
  - [x] Sub-task: Run unit tests: `go test -v ./cmd/...`
  - [x] Sub-task: Push checkpoint: `git push origin battery-bender-integration`

---

## Phase 2: MCP Server Engine & Protocol Migration
Refactor `internal/mcp` to use `bender/pkg/mcp.Server` for JSON-RPC 2.0 stdio handling.

- [x] Task 2.1: Refactor `internal/mcp/server.go` to wrap `bender/pkg/mcp.Server` (cc716b6)
  - [x] Sub-task: Write/update unit tests in `internal/mcp/server_test.go` (Red)
  - [x] Sub-task: Implement `internal/mcp/server.go` using Bender server (Green)
  - [x] Sub-task: Remove redundant `internal/mcp/protocol.go` and `protocol_test.go` (Refactor)

- [x] Task 2.2: Verify domain tools, resources, and prompts registration (cc716b6)
  - [x] Sub-task: Adapt `internal/mcp/tools.go`, `resources.go`, `prompts.go` to Bender handler signatures (Red/Green)
  - [x] Sub-task: Verify all domain tools execute correctly via `internal/mcp/tools_test.go`
  - [x] Sub-task: Verify framework status and template tools via `internal/mcp/framework_tools_test.go`

- [x] Task 2.3: Phase 2 Verification & Checkpoint [checkpoint: 878808f]
  - [x] Sub-task: Sync phase rules & specs (`git fetch origin main`)
  - [x] Sub-task: Run unit tests: `go test -v ./internal/mcp/...`
  - [x] Sub-task: Push checkpoint: `git push origin battery-bender-integration`

---

## Phase 3: MCP Client Installer Migration & CLI Ergonomics
Migrate `cmd/mcp.go` to use `bender/pkg/mcp` client installer and clean up redundant installer logic.

- [~] Task 3.1: Refactor `cmd/mcp.go` to use `mcp.InstallClients` & `mcp.GetSupportedClients`
  - [ ] Sub-task: Update `cmd/mcp_test.go` for Bender installer integration (Red)
  - [ ] Sub-task: Implement installer dispatch in `cmd/mcp.go` (Green)
  - [ ] Sub-task: Remove redundant `internal/mcp/installer.go` and `installer_test.go` (Refactor)

- [ ] Task 3.2: Phase 3 Verification & Checkpoint
  - [ ] Sub-task: Sync phase rules & specs (`git fetch origin main`)
  - [ ] Sub-task: Run full test suite: `go test -v -coverprofile=coverage.out ./...`
  - [ ] Sub-task: Push checkpoint: `git push origin battery-bender-integration`

---

## Phase 4: Final Quality Verification & Spec Sync
Verify end-to-end CLI behavior, test coverage, and update living capability specs.

- [ ] Task 4.1: Comprehensive Linting & Coverage Verification
  - [ ] Sub-task: Run `gofmt -l .` and `go vet ./...`
  - [ ] Sub-task: Verify total coverage exceeds 80%
  - [ ] Sub-task: Verify binary compilation: `go build -o bin/battery .`

- [ ] Task 4.2: Merge Spec Deltas into Living Specs
  - [ ] Sub-task: Update `.cooper/specs/cli-self-update/spec.md` with Bender updater references
  - [ ] Sub-task: Update `.cooper/specs/mcp-server/spec.md` with Bender MCP references

- [ ] Task 4.3: Final Phase Verification & Checkpoint
  - [ ] Sub-task: Record completion metadata in `.cooper/active/battery-bender-integration/metadata.json`
  - [ ] Sub-task: Update `.cooper/tracks.md` registry
  - [ ] Sub-task: Push final branch: `git push origin battery-bender-integration`
