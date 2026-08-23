# Design: Battery Go Module Casing Standardization

## Architecture Overview
Standardize the Go module path to `github.com/twoBoots/battery` across the codebase to ensure consistency with the GitHub organization name (`twoBoots`) and eliminate casing discrepancies across Go tooling, proxy caches, and sibling dependencies (such as `github.com/twoBoots/bender`).

## Changes Required

### 1. `go.mod` Definition
- Update module path:
  ```go
  module github.com/twoBoots/battery
  ```

### 2. Internal Package Import Rewrites
Replace all instances of `github.com/twoboots/battery/...` with `github.com/twoBoots/battery/...` across:
- `main.go`
- `cmd/`:
  - `cmd/barrel.go`, `cmd/barrel_test.go`
  - `cmd/init.go`, `cmd/init_test.go`
  - `cmd/installer_e2e_test.go`
  - `cmd/mcp.go`, `cmd/mcp_test.go`
  - `cmd/root_test.go`
  - `cmd/status.go`, `cmd/status_test.go`
  - `cmd/track.go`, `cmd/track_test.go`
- `internal/`:
  - `internal/config/config_test.go`
  - `internal/discovery/discovery.go`, `internal/discovery/discovery_test.go`
  - `internal/mcp/framework_tools_test.go`
  - `internal/mcp/resources.go`, `internal/mcp/resources_prompts_test.go`
  - `internal/mcp/tools.go`, `internal/mcp/tools_test.go`
  - `internal/techstack/techstack_test.go`
  - `internal/track/dispatch.go`, `internal/track/dispatch_test.go`
  - `internal/track/locate_test.go`, `internal/track/models_test.go`
  - `internal/track/status.go`, `internal/track/status_test.go`

### 3. CI / Release Workflow Linker Flags
In `.github/workflows/release.yml`:
- Update line 96:
  ```yaml
  go build -ldflags="-s -w -X github.com/twoBoots/battery/cmd.Version=${VERSION}" -o dist/${{ matrix.binary }} .
  ```

### 4. Verification & Testing
- Run `go test -v ./...` to verify all package tests compile and execute cleanly.
- Run `go vet ./...` to confirm static analysis passes with zero warnings.
- Build local test binary via `go build -o bin/battery .` and execute `bin/battery version`.
