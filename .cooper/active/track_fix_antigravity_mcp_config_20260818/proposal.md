# Track Proposal: Fix Antigravity / agy MCP Configuration Target to mcp_config.json

## Context & Motivation
When running `battery mcp install` (or during automated setup via `battery init`), `battery` registers its MCP server for Google Antigravity / Gemini by writing to `.gemini/settings.json`.

However, Google Antigravity (`agy` CLI / Antigravity IDE) does not read MCP configuration from `.gemini/settings.json`. Antigravity follows the MCP configuration standard using `mcp_config.json`:
1. Global configuration: `~/.gemini/config/mcp_config.json`
2. Workspace configuration: `.agents/mcp_config.json`

As a result, Antigravity sessions fail to discover or spawn the Battery MCP server. GitHub Issue #4 details this defect.

## Scope & Deliverables
1. Update `internal/mcp/installer.go`:
   - Target `~/.gemini/config/mcp_config.json` for Antigravity configuration.
   - Update `displayName` to `"Google Antigravity / agy (~/.gemini/config/mcp_config.json)"`.
   - Update `detectPaths` to check `homeDir/.gemini/config`, `homeDir/.gemini`, `cwd/.agents`, and `cwd/.gemini`.
2. Update documentation and CLI help text:
   - `cmd/mcp.go`: Update `mcp install` help text.
   - `docs/mcp-setup-guide.md` and `README.md`: Update Antigravity configuration instructions.
3. Update tests:
   - `internal/mcp/installer_test.go` and `cmd/mcp_test.go` to assert target configuration at `~/.gemini/config/mcp_config.json`.
