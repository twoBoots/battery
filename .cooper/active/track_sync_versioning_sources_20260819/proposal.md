# Proposal: Synchronize Battery Versioning Sources

## Problem Statement
Battery versioning was fragmented across three disconnected locations:
1. **MCP Server `Initialize` Response**: Hardcoded as `"v1.2.1"` in `internal/mcp/server.go`.
2. **MCP `battery_status` Output**: Reported `effCfg.Version` (`"1.0.0"`), which represents the `.batteryrc` schema version, rather than the CLI binary version.
3. **GitHub Actions Release Pipeline**: `.github/workflows/release.yml` defaulted to `VERSION="1.0.0"` on pushes to `main` without tag prefix, causing builds to report `v1.0.0`.

## Proposed Changes
1. **Dynamic MCP Server Version**: Pass `cmd.Version` (or server version argument) to `mcp.NewServer` so MCP `initialize` returns the exact CLI version (`1.4.0`).
2. **Clearer `battery_status` Tool Output**: Include both `config_version` and `cli_version` in `battery_status` tool output.
3. **CI Release Workflow Dynamic Version Extraction**: In `.github/workflows/release.yml`, dynamically extract the version from `cmd/root.go` when `GITHUB_REF_TYPE != "tag"`.
