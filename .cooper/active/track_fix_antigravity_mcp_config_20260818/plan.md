# Execution Plan: Fix Antigravity / agy MCP Configuration Target

## Phase 1: Update MCP Installer Target & Tests (TDD)
- [x] Task: Update Installer Logic & Tests (0ec6751)
  - [x] Sub-task: Update `cmd/mcp_test.go` and `internal/mcp/installer_test.go` for Antigravity target (Red)
  - [x] Sub-task: Update `internal/mcp/installer.go` and `cmd/mcp.go` with `~/.gemini/config/mcp_config.json` (Green)
  - [x] Sub-task: Verify tests pass with coverage >80% (Refactor)
- [x] Task: Phase 1 Checkpoint (5644a61)

## Phase 2: Documentation & Living Spec Updates
- [x] Task: Update Documentation (`README.md`, `docs/mcp-setup-guide.md`) (3ed45d8)
- [x] Task: Update Living Spec (`.cooper/specs/mcp-server/spec.md`) (34246d6)
- [ ] Task: Phase 2 Checkpoint & Verification
