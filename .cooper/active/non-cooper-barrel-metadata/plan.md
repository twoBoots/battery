# Implementation Plan: Lightweight Documentation & Context Metadata for Non-Cooper Barrels

All implementation tasks strictly follow the Cooper TDD lifecycle (Red -> Green -> Refactor) and maintain code coverage >80%.

## Phase 1: Configuration Schema & Dynamic Metadata Preservation (`internal/config`)

- [x] Task: Extend `BarrelConfig` & `EffectiveBarrel` with Metadata Fields & Dynamic Property Preservation (84720cc)
  - [x] Sub-task: Write unit tests for JSON marshaling/unmarshaling of `role`, `tech`, `docs`, `jira`, and preservation of unknown custom attributes across round-trips (Red)
  - [x] Sub-task: Implement metadata fields, custom unmarshaler/marshaler, and effective configuration merging in `internal/config/config.go` (Green)
  - [x] Sub-task: Refactor & verify coverage >80% (Refactor)
- [x] Task: Update `AddBarrel` & `RemoveBarrel` APIs to Support Metadata Attributes (f4948cb)
  - [x] Sub-task: Write unit tests for adding barrels with metadata and removing barrels while preserving metadata and custom fields in `.batteryrc` / `.batteryrc.local` (Red)
  - [x] Sub-task: Implement metadata plumbing in `AddBarrel` and `RemoveBarrel` (Green)
  - [x] Sub-task: Refactor & verify coverage >80% (Refactor)
- [x] Task: Phase 1 Verification & Checkpoint [checkpoint: 78ad73c]

## Phase 2: Tech Stack & Barrel Profile Resolution Engine (`internal/techstack`)

- [x] Task: Barrel Profile Discovery & Scaffolding Engine (bc0f21a)
  - [x] Sub-task: Write unit tests for discovering profiles in `docs/barrels/`, `.cooper/barrels/`, and custom `docs` paths, plus starter template generation (Red)
  - [x] Sub-task: Implement `ResolveBarrelProfile` and `ScaffoldBarrelProfile` with embedded starter template (Green)
  - [x] Sub-task: Refactor & verify coverage >80% (Refactor)
- [x] Task: Hybrid Context & Fallback Resolver (e8afe60)
  - [x] Sub-task: Write unit tests for fallback resolution hierarchy (Cooper `tech-stack.md` -> `.batteryrc` `tech`/`role` -> profile markdown auto-summary -> default fallback) (Red)
  - [x] Sub-task: Implement `ResolveBarrelContext` and markdown summary extraction in `internal/techstack/` (Green)
  - [x] Sub-task: Refactor & verify coverage >80% (Refactor)
- [~] Task: Phase 2 Verification & Checkpoint

## Phase 3: CLI Command Enhancements (`cmd/`)

- [ ] Task: Extend `battery barrel add` with Metadata Flags
  - [ ] Sub-task: Write CLI unit tests for `battery barrel add` with `--role`, `--tech`, `--docs`, `--jira` (Red)
  - [ ] Sub-task: Implement Cobra flags and config invocation in `cmd/barrel.go` (Green)
  - [ ] Sub-task: Refactor & verify coverage >80% (Refactor)
- [ ] Task: Add `battery barrel doc init <name>` / `battery barrel profile init <name>` Subcommand
  - [ ] Sub-task: Write CLI unit tests for profile initialization and `--force` flag (Red)
  - [ ] Sub-task: Implement `battery barrel doc init` subcommand (Green)
  - [ ] Sub-task: Refactor & verify coverage >80% (Refactor)
- [ ] Task: Enhance `battery status` and `battery barrel list` Output Displays
  - [ ] Sub-task: Write CLI unit tests verifying display of metadata attributes and hybrid summaries (Red)
  - [ ] Sub-task: Update output formatting in `cmd/status.go` and `cmd/barrel.go` (Green)
  - [ ] Sub-task: Refactor & verify coverage >80% (Refactor)
- [ ] Task: Phase 3 Verification & Checkpoint

## Phase 4: MCP Server Resources & Tools (`internal/mcp`)

- [ ] Task: Expose `battery://barrels/{name}/docs` & Fallback `tech-stack` MCP Resources
  - [ ] Sub-task: Write MCP unit tests for reading barrel documentation profiles and fallback tech-stack payloads (Red)
  - [ ] Sub-task: Implement resource registration in `internal/mcp/resources.go` (Green)
  - [ ] Sub-task: Refactor & verify coverage >80% (Refactor)
- [ ] Task: Enrich `battery_list_barrels` Tool & `battery://topology` Resource with Metadata
  - [ ] Sub-task: Write MCP unit tests verifying JSON payloads include `role`, `tech`, `docs`, `jira`, `hasProfile` (Red)
  - [ ] Sub-task: Update `internal/mcp/tools.go` and topology handler (Green)
  - [ ] Sub-task: Refactor & verify coverage >80% (Refactor)
- [ ] Task: Phase 4 Verification & Checkpoint

## Phase 5: Documentation, Templates & Final Integration

- [ ] Task: Update Living Capability Specs & Upstream Framework Templates
  - [ ] Sub-task: Update `.cooper/BATTERY.md` and `internal/framework/templates/docs/BATTERY.md` documenting non-Cooper barrel profiles
  - [ ] Sub-task: Run full test suite across all packages (`go test ./...`) and lint checks (`go vet ./...`)
  - [ ] Sub-task: Verify test coverage >80% across all modified packages
- [ ] Task: Phase 5 Verification & Final Checkpoint
