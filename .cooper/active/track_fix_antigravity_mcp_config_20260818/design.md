# Technical Design: Fix Antigravity / agy MCP Configuration Target

## Architecture & Changes

### 1. Installer Target Definition (`internal/mcp/installer.go`)
Update the `antigravity` target in `GetSupportedClients(cwd, homeDir)`:
- `id`: `"antigravity"`
- `displayName`: `"Google Antigravity / agy (~/.gemini/config/mcp_config.json)"`
- `configPath`: `filepath.Join(homeDir, ".gemini", "config", "mcp_config.json")`
- `detectPaths`:
  - `filepath.Join(homeDir, ".gemini", "config")`
  - `filepath.Join(homeDir, ".gemini")`
  - `filepath.Join(cwd, ".agents")`
  - `filepath.Join(cwd, ".gemini")`

### 2. CLI Help Strings (`cmd/mcp.go`)
Update `mcpInstallCmd.Long`:
- Replace `* antigravity     - Google Antigravity / Gemini (.gemini/settings.json)` with `* antigravity     - Google Antigravity / agy (~/.gemini/config/mcp_config.json)`

### 3. Documentation (`docs/mcp-setup-guide.md`, `README.md`)
- Update Antigravity manual setup instructions to reference `~/.gemini/config/mcp_config.json` (or `.agents/mcp_config.json`).
- Update quick setup snippets.

### 4. Tests (`internal/mcp/installer_test.go`, `cmd/mcp_test.go`)
- Update `TestGetSupportedClients_Detections` to check detection with `homeDir/.gemini/config` or `workspaceDir/.agents`.
- Update `TestMCPInstall_WithClientsFlag` to assert creation of `homeDir/.gemini/config/mcp_config.json`.
