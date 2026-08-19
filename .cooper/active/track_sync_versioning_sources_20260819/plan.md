# Implementation Plan: Synchronize Versioning Sources

## Phase 1: Dynamic MCP Server Version & Battery Status Tool
- [x] Task: Make MCP Server version dynamic in `internal/mcp/server.go` and `cmd/mcp.go`
  - [x] Sub-task: Update `NewServer` constructor to accept version parameter
  - [x] Sub-task: Pass `cmd.Version` in `cmd/mcp.go`
  - [x] Sub-task: Update `battery_status` tool to report `cli_version` and `config_version`
  - [x] Sub-task: Update MCP unit tests in `internal/mcp/server_test.go` and `internal/mcp/tools_test.go`
- [x] Task: Update GitHub Actions release workflow (`.github/workflows/release.yml`)
  - [x] Sub-task: Dynamically parse `cmd/root.go` version for untagged `main` builds
- [x] Task: Phase 1 Verification & Checkpoint
  - [x] Sub-task: Run full test suite (`go test -v ./...`)
  - [x] Sub-task: Commit, git notes checkpoint, and push branch
