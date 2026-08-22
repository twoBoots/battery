# Proposal: Battery Bender CLI Integration

## Intent
Adopt the standardized `twoBoots` Go CLI archetype and core package library, **`bender`** (`github.com/twoBoots/bender`), inside `battery`. This eliminates duplicated self-updater and MCP protocol logic, streamlines maintenance, and aligns `battery` with ecosystem conventions.

## Rationale & Benefits
- **Zero Maintenance Duplication**: Removes bespoke GitHub release updater and stdio JSON-RPC 2.0 protocol engine from `battery/internal/`.
- **Ecosystem Consistency**: Ensures updater mechanics (SemVer comparison, atomic binary replacement, macOS quarantine clearing and ad-hoc code signing) and MCP client auto-installers behave identically across `battery`, `cooper`, and future twoBoots CLIs.
- **Clean Architecture**: `battery` focuses purely on multi-barrel orchestration logic, domain tools, living specifications, and prompt templates, delegating generic CLI infrastructure to `bender`.

## Scope Boundaries
- **In Scope**:
  - Add `github.com/twoBoots/bender` dependency to `go.mod`.
  - Replace `battery/internal/updater` with `bender/pkg/updater` in `cmd/update.go`.
  - Replace `battery/internal/mcp/{protocol,server,installer}.go` with `bender/pkg/mcp` while preserving all domain tools, resources, and prompts.
  - Delete `battery/internal/updater/` and redundant MCP files.
  - Update unit and end-to-end CLI tests, maintaining >80% coverage.
- **Out of Scope**:
  - Changes to multi-barrel track dispatch, framework status algorithms, or barrel configuration logic.
  - Direct modifications to downstream barrels (e.g. `cooper` or `planet-express`).
