# Technical Design: Synchronized Dynamic Versioning

## Technical Specifications

### 1. `internal/mcp/server.go`
- Update `NewServer(cwd string, version ...string) *Server` or `SetVersion(version string)`:
  - If version is provided, strip any leading `v` or format as `v<version>`.
  - Default to `v1.4.0` if empty.
- In `cmd/mcp.go`:
  - Pass `Version` from `cmd.Version` to `mcp.NewServer(cwd, Version)`.

### 2. `internal/mcp/tools.go`
- In `battery_status` response:
  - Add `CliVersion: s.version` and rename/clarify `ConfigVersion: effCfg.Version`.
  - Maintain backwards compatibility by populating `Version: effCfg.Version` while providing explicit `CliVersion`.

### 3. `.github/workflows/release.yml`
- In `build-and-release` job:
  ```bash
  VERSION=$(grep 'Version = ' cmd/root.go | sed -E 's/.*"([^"]+)".*/\1/')
  if [ "$GITHUB_REF_TYPE" = "tag" ]; then
    VERSION="${GITHUB_REF_NAME#v}"
  fi
  ```
