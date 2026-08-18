# Execution Plan: Fix Antigravity / agy MCP Configuration Target

## Phase 1: Update MCP Installer Target & Tests (TDD)
- [ ] Task: Update Installer Logic & Tests
  - [ ] Sub-task: Update `cmd/mcp_test.go` and `internal/mcp/installer_test.go` for Antigravity target (Red)
  - [ ] Sub-task: Update `internal/mcp/installer.go` and `cmd/mcp.go` with `~/.gemini/config/mcp_config.json` (Green)
  - [ ] Sub-task: Verify tests pass with coverage >80% (Refactor)
- [ ] Task: Phase 1 Checkpoint

## Phase 2: Documentation & Living Spec Updates
- [ ] Task: Update Documentation (`README.md`, `docs/mcp-setup-guide.md`)
- [ ] Task: Update Living Spec (`.cooper/specs/mcp-server/spec.md`)
- [ ] Task: Phase 2 Checkpoint & Verification
