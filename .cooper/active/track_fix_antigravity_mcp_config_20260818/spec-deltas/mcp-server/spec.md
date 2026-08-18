# Spec Delta: Fix Antigravity MCP Configuration Target

## Capability: `mcp-server`

### Requirement 5: MCP Client Auto-Configuration & Installation (`battery mcp install`)
The system MUST support automatic detection, interactive selection, and safe JSON configuration insertion of Battery MCP into supported AI assistants.

#### Scenario 5.1: Non-Interactive / Flag-Based Client Installation
- **GIVEN** a list of target clients (e.g. `--client cursor,claude-desktop,antigravity` or `--all`)
- **WHEN** `battery mcp install` is executed
- **THEN** it MUST create or update each client's configuration file with the `battery mcp` server definition without corrupting existing JSON keys.
+- For Google Antigravity (`agy`), the target configuration file MUST be `~/.gemini/config/mcp_config.json` with display name referencing `mcp_config.json`.

#### Scenario 5.2: Interactive Client Discovery & Multi-Select
- **GIVEN** an interactive terminal session
- **WHEN** `battery mcp install` is executed without client flags
- **THEN** it MUST scan for known client configurations and present a multi-select prompt allowing the user to choose which clients to configure.
+- Antigravity detection MUST check for `~/.gemini/config`, `~/.gemini`, `.agents`, and `.gemini`.
